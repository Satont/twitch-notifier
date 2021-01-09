import { getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'
import { info } from '../libs/logger'
import Twitch from '../libs/twitch'
import { services } from '../services/_interface'
import { ITwitchStreamChangedPayload } from '../typings/twitch'

class TwitchWatcherClass {
  private readonly channelsRepository = getConnection().getRepository(Channel)
  subscriptions = new Set()

  async init() {
    const subscriptions = await (await Twitch.apiClient.helix.webHooks.getSubscriptions()).getAll()
    for (const subsciption of subscriptions) {
      await subsciption.unsubscribe()
    }

    const channelsRepository = getConnection().getRepository(Channel)
    const channels = await channelsRepository.find()

    for (const channel of channels) {
      await this.addChannelToWebhooks(channel.id)
    }

    info(`TWITCH: webhooks subscribed to ${channels.length} channels`)
    setTimeout(() => this.init(), 864000)
  }

  async addChannelToWebhooks(channelId: string) {
    const siteUrl = process.env.SITE_URL
    if (!siteUrl) return
    const options = {
      callbackUrl: `${siteUrl}/twitch/webhooks/callback`,
      validityInSeconds: 864000,
    }

    if (this.subscriptions.has(channelId)) return

    await Twitch.apiClient.helix.webHooks.subscribeToStreamChanges(channelId, options)
    this.subscriptions.add(channelId)
  }

  async processPayload(data: ITwitchStreamChangedPayload['data']) {
    for (const stream of data) {
      const user = await Twitch.getUser({ id: stream.user_id })
      const category = stream.game_name
      const channel = await this.channelsRepository.findOne(stream.user_id, { relations: ['followers', 'followers.chat' ] })
        || this.channelsRepository.create({
          id: stream.user_id,
        })

      channel.username = user.name
      const messageOpts = {
        image: `${stream.thumbnail_url?.replace('{width}', '1920').replace('{height}', '1080')}?timestamp=${Date.now()}`,
      }

      if (stream.type === 'live') {
        if (!channel.online) {
          for (const service of services) {
            service.makeAnnounce({
              message: `${stream.user_name} online!\nCategory: ${category}\nTitle: ${stream.title}\nhttps://twitch.tv/${user.name}`,
              target: channel.followers?.map(f => f.chat.chatId),
              ...messageOpts,
            })
          }
        }

        if (channel.online && channel.category !== category) {
          for (const service of services) {
            service.makeAnnounce({
              message: `${stream.user_name} now streaming ${category}\nPrevious category: ${channel.category}\nhttps://twitch.tv/${user.name}`,
              target: channel.followers?.filter(f => f.chat.settings.game_change_notification).map(f => f.chat.chatId),
              ...messageOpts,
            })
          }
        }

        channel.category = category
        channel.title = stream.title
        channel.online = true
      } else if (stream.type === 'offline') {
        channel.online = false

        for (const service of services) {
          service.makeAnnounce({
            message: `${channel.username} now offline`,
            target: channel.followers?.filter(f => f.chat.settings.offline_notification).map(f => f.chat.chatId),
          })
        }
      }

      await channel.save()
    }
  }
}

const TwitchWatcher = new TwitchWatcherClass()
export default TwitchWatcher
