import { getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'
import { info } from '../libs/logger'
import Twitch from '../libs/twitch'
import * as TwitchEventSub from '@twurple/eventsub'
import { getAppLication, listened } from '../web'
import { Stream } from '../entities/Stream'
import { listenedChannels } from '../cache/listenedChannels'
import { Announcer } from '../libs/announcer'
import { Express } from 'express' 

class TwitchWatcherClass {
  private readonly channelsRepository = getConnection().getRepository(Channel)
  private readonly streamsRepository = getConnection().getRepository(Stream)
  private adapter: TwitchEventSub.EventSubMiddleware

  async init() {
    if (!Twitch.apiClient) {
      return setTimeout(() => this.init(), 1000)
    }
    listenedChannels.clear()

    this.adapter = new TwitchEventSub.EventSubMiddleware({
      apiClient: Twitch.apiClient,
      secret: process.env.TWITCH_CLIENT_ID,
      hostName: await this.getAdapterHostname(),
      pathPrefix: 'twitch/eventsub',
      logger: {
        minLevel: 'debug',
      },
    })

    await this.adapter.apply(getAppLication().getHttpAdapter() as unknown as Express)
    
    await this.initChannels()

    // Add channels to watcher on start
    info(`TWITCH: EventSub watcher started.`)
  }

  private async initChannels() {
    if (!listened) {
      return setTimeout(() => this.initChannels(), 1000)
    }

    await this.adapter.markAsReady()

    const channels = await getConnection().getRepository(Channel).find()
    const currentSubscriptions = await Twitch.apiClient.eventSub.getSubscriptionsPaginated().getAll()

    const hostname = await this.getAdapterHostname()
    const forDelete = currentSubscriptions.filter((sub) => !sub._transport.callback.includes(hostname))

    for (const sub of forDelete) {
      const callback = sub._transport.callback
      await sub.unsubscribe()
      const message = `Deleting redutant subscription ${sub.type}${sub.condition.broadcaster_user_id}, domain: ${callback}`
      info(`TWITCH: EventSub ${message}.`)
    }

    const forCreate = channels
      .filter(channel => currentSubscriptions.find(sub => sub.condition.broadcaster_user_id == channel.id))

    for (const channel of forCreate) {
      await this.addChannelToWatch(channel.id)
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
    const channel = await this.channelsRepository.findOne(channelId)
    || await this.channelsRepository.create({
      id: channelId,
    }).save()
    
    const stream = await Twitch.apiClient.streams.getStreamByUserId(channelId) 

    if (stream) {
      channel.online = true
      await channel.save()
    }

    if (stream && stream.id !== (await this.getLatestStream(channelId))?.id) {
      await this.streamsRepository.create({ 
        id: stream.id, 
        startedAt: stream.startDate, 
        channel,
        category: stream.gameName,
        title: stream.title,
      }).save()
    }

    const announcer = new Announcer(channelId)
    await announcer.init()
    
    this.adapter.subscribeToStreamOnlineEvents(channelId, (event) => announcer.announceLive(event))
    this.adapter.subscribeToStreamOfflineEvents(channelId, (event) => announcer.announceOffline(event))
    this.adapter.subscribeToChannelUpdateEvents(channelId, (event) => announcer.announceUpdate(event))
    listenedChannels.add(channelId)
  }

  private async getAdapterHostname() {
    const hostname = process.env.SITE_URL

    return hostname.replace('http://', '').replace('https://', '')
  }
}

const TwitchWatcher = new TwitchWatcherClass()
export default TwitchWatcher
