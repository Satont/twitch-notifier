import { getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'
import { info } from '../libs/logger'
import Twitch from '../libs/twitch'
import { services } from '../services/_interface'
import * as TwitchEventSub from 'twitch-eventsub'
import { getAppLication } from '../web'

class TwitchWatcherClass {
  private readonly channelsRepository = getConnection().getRepository(Channel)
  private adapter: TwitchEventSub.MiddlewareAdapter
  private listener: TwitchEventSub.EventSubListener

  private async getAdapterHostname() {
    let hostname: string
    if (process.env.NODE_ENV === 'production') hostname = process.env.SITE_URL
    else {
      hostname = await (await import('ngrok')).connect(Number(process.env.PORT))
      info(`EventSub: working with ngrok. Current link is: ${hostname}`)
    }

    return hostname.replace('http://', '').replace('https://', '')
  }

  async init() {
    if (!Twitch.apiClient) {
      return setTimeout(() => this.init(), 1000)
    }
    this.adapter = new TwitchEventSub.MiddlewareAdapter({
      hostName: await this.getAdapterHostname(),
      pathPrefix: 'twitch/eventsub',
    })

    this.listener = new TwitchEventSub.EventSubListener(Twitch.apiClient, this.adapter, process.env.TWITCH_EVENTSUB_SECRET || '0123456789')

    await this.listener.applyMiddleware(getAppLication())
    await Twitch.apiClient.helix.eventSub.deleteAllSubscriptions()
    
    const channelsRepository = getConnection().getRepository(Channel)
    const channels = await channelsRepository.find()

    for (const channel of channels) {
      await this.addChannelToWatch(channel.id)
    }

    info(`TWITCH: EventSub watcher started`)
    console.log('Eventsub subscriptions:', (await Twitch.apiClient.helix.eventSub.getSubscriptions()).data.length)
  }

  async addChannelToWatch(channelId: string) {
    const channel = await this.channelsRepository.findOne(channelId, { relations: ['followers', 'followers.chat' ] })
        || this.channelsRepository.create({
          id: channelId,
        })


    await this.listener.subscribeToStreamOnlineEvents(channelId, async (event) => {
      if (event.streamType !== 'live') return

      channel.username = event.broadcasterName
      const stream = await Twitch.apiClient.helix.streams.getStreamByUserId(channelId)

      for (const service of services) {
        service.makeAnnounce({
          message: `${event.broadcasterDisplayName} online!\nCategory: ${stream.gameName}\nTitle: ${stream.title}\nhttps://twitch.tv/${event.broadcasterName}`,
          target: channel.followers?.map(f => f.chat.chatId),
          image: this.getThumnailUrl(stream.thumbnailUrl),
        })
      }

      channel.title = stream.title
      channel.online = true
      channel.save()
    })

    await this.listener.subscribeToStreamOfflineEvents(channelId, async (event) => {
      for (const service of services) {
        service.makeAnnounce({
          message: `${event.broadcasterDisplayName} now offline`,
          target: channel.followers?.filter(f => f.chat.settings.offline_notification).map(f => f.chat.chatId),
        })
      }

      channel.online = false
      channel.save()
    })

    await this.listener.subscribeToChannelUpdateEvents(channelId, async (event) => {
      const stream = await Twitch.apiClient.helix.streams.getStreamByUserId(channelId)
      if (stream.type === 'live') return

      if (channel.online && channel.category !== event.categoryName) {
        for (const service of services) {
          service.makeAnnounce({
            message: `${event.broadcasterDisplayName} now streaming ${event.categoryName}\nPrevious category: ${channel.category}\nhttps://twitch.tv/${event.broadcasterName}`,
            target: channel.followers?.filter(f => f.chat.settings.game_change_notification).map(f => f.chat.chatId),
            image: this.getThumnailUrl(stream.thumbnailUrl),
          })
        }
      }
    })

    return
  }

  private getThumnailUrl(url: string) {
    return `${url.replace('{width}', '1920').replace('{height}', '1080')}?timestamp=${Date.now()}`
  }
}

const TwitchWatcher = new TwitchWatcherClass()
export default TwitchWatcher
