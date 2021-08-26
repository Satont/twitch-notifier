import { getConnection, getRepository } from 'typeorm'
import { Channel } from '../entities/Channel'
import { info } from '../libs/logger'
import Twitch from '../libs/twitch'
import { services } from '../services/_interface'
import * as TwitchEventSub from 'twitch-eventsub'
import { getAppLication, listened } from '../web'
import { Follow } from '../entities/Follow'
import { Stream } from '../entities/Stream'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import localtunnel from 'localtunnel'

dayjs.extend(relativeTime)

class TwitchWatcherClass {
  private readonly channelsRepository = getConnection().getRepository(Channel)
  private readonly streamsRepository = getConnection().getRepository(Stream)
  private adapter: TwitchEventSub.MiddlewareAdapter
  private listener: TwitchEventSub.EventSubListener
  private listenedChannels: string[] = []

  async init() {
    if (!Twitch.apiClient) {
      return setTimeout(() => this.init(), 1000)
    }
    this.listenedChannels = []

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
      info(`TWITCH: EventSub watcher: ${this.listenedChannels.length} channels.`)
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

  private getChannelFollowers(channelId: string) {
    return getRepository(Follow).find({
      where: {
        channel: { id: channelId },
      },
      relations: ['chat'],
    })
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
    if (this.listenedChannels.includes(channelId)) return
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
    
    await this.listener.subscribeToStreamOnlineEvents(channelId, async (event) => {
      if (event.streamType !== 'live') return
      const latestStream = await this.getLatestStream(event.broadcasterId)

      const stream = await Twitch.apiClient.helix.streams.getStreamByUserId(channelId)
      
      if (stream.id !== latestStream?.id) {
        for (const service of services) {
          service.makeAnnounce({
            message: `
              ${event.broadcasterDisplayName} online!
              Category: ${stream.gameName}
              Title: ${stream.title}
              https://twitch.tv/${event.broadcasterName}
            `.replace(/  +/g, ''),
            target: (await this.getChannelFollowers(channel.id)).map(f => f.chat.chatId),
            image: this.getThumnailUrl(stream.thumbnailUrl),
          })
        }
      }
      
      channel.username = event.broadcasterName
      channel.online = true
      await this.streamsRepository.create({ 
        id: stream.id, 
        startedAt: stream.startDate, 
        channel,
        title: stream.title,
        category: stream.gameName,
      }).save()
      channel.save()
    })

    await this.listener.subscribeToStreamOfflineEvents(channelId, async (event) => {
      const latestStream = await this.getLatestStream(event.broadcasterId)
      const streamDuration = dayjs().from(dayjs(latestStream?.startedAt), true)

      for (const service of services) {
        service.makeAnnounce({
          message: `
            ${event.broadcasterDisplayName} now offline
            ${latestStream ? 'Stream duration was: ' + streamDuration : ''}
          `.replace(/  +/g, ''),
          target: (await this.getChannelFollowers(channel.id)).filter(f => f.chat.settings.offline_notification).map(f => f.chat.chatId),
        })
      }

      channel.online = false
      channel.save()
    })

    await this.listener.subscribeToChannelUpdateEvents(channelId, async (event) => {
      const stream = await Twitch.apiClient.helix.streams.getStreamByUserId(channelId)
      if (stream?.type !== 'live') return
      const latestStream = await this.getLatestStream(channelId)

      if (channel.online && latestStream?.category !== event.categoryName) {
        for (const service of services) {
          service.makeAnnounce({
            message: `
              ${event.broadcasterDisplayName} now streaming ${event.categoryName}
              Previous category: ${latestStream?.category}
              https://twitch.tv/${event.broadcasterName}
            `.replace(/  +/g, ''),
            target: (await this.getChannelFollowers(channel.id)).filter(f => f.chat.settings.game_change_notification).map(f => f.chat.chatId),
            image: this.getThumnailUrl(stream.thumbnailUrl),
          })
        }
      }

      if (latestStream) {
        latestStream.updatedAt = new Date()
        latestStream.category = event.categoryName
        await this.streamsRepository.save(latestStream)
      }

      channel.save()
    })

    return
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

  private getThumnailUrl(url: string) {
    return `${url.replace('{width}', '1920').replace('{height}', '1080')}?timestamp=${Date.now()}`
  }
}

const TwitchWatcher = new TwitchWatcherClass()
export default TwitchWatcher
