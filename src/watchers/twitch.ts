import { getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'
import { info } from '../libs/logger'
import Twitch from '../libs/twitch'
import * as TwitchEventSub from 'twitch-eventsub'
import { getAppLication, listened } from '../web'
import { Stream } from '../entities/Stream'
import localtunnel from 'localtunnel'
import { listenedChannels } from '../cache/listenedChannels'
import { Announcer } from '../libs/announcer'

class TwitchWatcherClass {
  private readonly channelsRepository = getConnection().getRepository(Channel)
  private readonly streamsRepository = getConnection().getRepository(Stream)
  private adapter: TwitchEventSub.MiddlewareAdapter
  private listener: TwitchEventSub.EventSubListener

  async init() {
    if (!Twitch.apiClient) {
      return setTimeout(() => this.init(), 1000)
    }
    listenedChannels.clear()

    this.adapter = new TwitchEventSub.MiddlewareAdapter({
      hostName: await this.getAdapterHostname(),
      pathPrefix: 'twitch/eventsub',
    })

    this.listener = new TwitchEventSub.EventSubListener(Twitch.apiClient, this.adapter, process.env.TWITCH_EVENTSUB_SECRET || '0123456789')
    await this.listener.applyMiddleware(getAppLication())

    info(`TWITCH: EventSub starting unsubscribe from all channels.`)
    // We need delete all subscriptions because our app URL can be changed.
    await Twitch.apiClient.helix.eventSub.deleteAllSubscriptions()
    
    // Add channels to watcher on start
    this.initChannels().then(() => {
      info(`TWITCH: EventSub watcher started.`)
      info(`TWITCH: EventSub watcher: ${listenedChannels.size} channels.`)
    })
  }

  async initChannels() {
    if (!listened) {
      return setTimeout(() => this.initChannels(), 1000)
    }

    for (const channel of await getConnection().getRepository(Channel).find()) {
      info(`Adding channel ${channel.username}[${channel.id}] to watcher`)
      this.addChannelToWatch(channel.id)
    }
  }
  
  private getLatestStream(channelId: string) {
    return this.streamsRepository.findOne({
      where: {
        channel: {
          id: channelId,
        },
      },
      order: {
        startedAt: 'DESC',
      },
    })
  }

  async addChannelToWatch(channelId: string) {
    if (listenedChannels.has(channelId)) return
    const announcer = new Announcer(channelId)
    await announcer.init()
    const channel = await this.channelsRepository.findOne(channelId)
    || await this.channelsRepository.create({
      id: channelId,
    }).save()
    
    const stream = await Twitch.apiClient.helix.streams.getStreamByUserId(channelId) 
    if (stream && stream.id !== (await this.getLatestStream(channelId))?.id) {
      await this.streamsRepository.create({ 
        id: stream.id, 
        startedAt: stream.startDate, 
        channel,
        category: stream.gameName,
        title: stream.title,
      }).save()
    }
    
    this.listener.subscribeToStreamOnlineEvents(channelId, async (event) => announcer.announceLive(event))
    this.listener.subscribeToStreamOfflineEvents(channelId, async (event) => announcer.announceOffline(event))
    this.listener.subscribeToChannelUpdateEvents(channelId, async (event) => announcer.announceUpdate(event))
  }

  private async getAdapterHostname() {
    let hostname: string
    if (process.env.NODE_ENV === 'production') hostname = process.env.SITE_URL
    else {
      hostname = (await localtunnel(Number(process.env.PORT))).url
      info(`EventSub: working with localtunnel. Current link is: ${hostname}`)
    }

    return hostname.replace('http://', '').replace('https://', '')
  }
}

const TwitchWatcher = new TwitchWatcherClass()
export default TwitchWatcher
