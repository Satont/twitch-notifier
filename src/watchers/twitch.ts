import { getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'
import { info } from '../libs/logger'
import Twitch from '../libs/twitch'
import { services } from '../services/_interface'
import { ITwitchStreamChangedPayload } from '../typings/twitch'

class TwitchWatcherClass {
  private readonly channelsRepository = getConnection().getRepository(Channel)


  async init() {
    const channelsRepository = getConnection().getRepository(Channel)
    const channels = await channelsRepository.find()

    for (const channel of channels) {
      await this.addChannelToWebhooks(channel.id)
    }

    info(`TWITCH: webhook subscribed to ${channels.length} channels`)
    setTimeout((() => this.init()), 864000 * 1000)
  }

  async addChannelToWebhooks(channelId: string) {
    const siteUrl = process.env.SITE_URL
    if (!siteUrl) return
    const options = {
      callbackUrl: `${siteUrl}/twitch/webhooks/callback`,
      validityInSeconds: 864000,
    }

    await Twitch.apiClient.helix.webHooks.unsubscribeFromStreamChanges(channelId, options)
    await Twitch.apiClient.helix.webHooks.subscribeToStreamChanges(channelId, options)
  }

  async processPayload(data: ITwitchStreamChangedPayload['data']) {
    for (const stream of data) {
      const category = stream.game_name
      const channel = await this.channelsRepository.findOne(stream.user_id, { relations: ['followers', 'followers.chat' ] })
        || this.channelsRepository.create({
          id: stream.user_id,
        })

      const messageOpts = {
        image: `${stream.thumbnail_url?.replace('{width}', '1920').replace('{height}', '1080')}?timestamp=${Date.now()}`,
      }

      if (stream.type === 'live') {
        if (!channel.online) {
          for (const service of services) {
            service.makeAnnounce({
              message: `${stream.user_name} online!\nCategory: ${category}\nTitle: ${stream.title}\nhttps://twitch.tv/${stream.user_name}`,
              target: channel.followers?.map(f => f.chat.chatId),
              ...messageOpts,
            })
          }
        }

        if (channel.online && channel.category !== category) {
          for (const service of services) {
            service.makeAnnounce({
              message: `${stream.user_name} now streaming ${category}\nPrevious category: ${channel.category}\nhttps://twitch.tv/${stream.user_name}`,
              target: channel.followers?.filter(f => f.chat.settings.game_change_notification).map(f => f.chat.chatId),
              ...messageOpts,
            })
          }
        }

        channel.category = category
        channel.title = stream.title
        channel.username = stream.user_name
        channel.online = true
      } else if (stream.type === 'offline') {
        channel.online = false
      }

      await channel.save()
    }
  }
}

const TwitchWatcher = new TwitchWatcherClass()
export default TwitchWatcher
