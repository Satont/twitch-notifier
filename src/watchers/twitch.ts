import { getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'
import { info, warning } from '../libs/logger'
import Twitch from '../libs/twitch'
import { services } from '../services/_interface'
import { ITwitchStreamChangedPayload } from '../typings/twitch'

class TwitchWatcherClass {
  private readonly channelsRepository = getConnection().getRepository(Channel)

  constructor() {
    this.initWebhooks()
  }

  private async initWebhooks() {
    const siteUrl = process.env.SITE_URL
    if (!siteUrl) {
      warning('TWITCH: siteUrl not setuped, streams udpates will not be recieved.')
      return
    }

    const options = {
      callbackUrl: `${siteUrl}/twitch/webhooks/callback`,
      validityInSeconds: 864000,
    }

    const channelsRepository = getConnection().getRepository(Channel)
    const channels = await channelsRepository.find()

    for (const channel of channels) {
      await Twitch.apiClient.helix.webHooks.unsubscribeFromStreamChanges(channel.id, options)
      await Twitch.apiClient.helix.webHooks.subscribeToStreamChanges(channel.id, options)
    }

    info(`TWITCH: webhook subscribed to ${channels.length} channels`)

    setTimeout((() => this.initWebhooks()), options.validityInSeconds * 1000)
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
        image: stream.thumbnail_url,
        target: channel.followers.map(f => f.chat.id),
      }
      if (!channel.online) {
        for (const service of services) {
          service.sendMessage({
            message: `${stream.user_name} online!\nCategory: ${category}\nTitle: ${stream.title}\nhttps://twitch.tv/${stream.user_name}`,
            ...messageOpts,
          })
        }
      } else if (channel.category !== category && channel.online) {
        for (const service of services) {
          service.sendMessage({
            message: `${stream.user_name} now streaming ${category}\nPrevious category: ${channel.category}\nhttps://twitch.tv/${stream.user_name}`,
            ...messageOpts,
          })
        }
      }

      channel.category = category
      channel.title = stream.title
      channel.username = stream.user_name

      if (stream.type === 'live') {
        channel.online = true
      } else {
        channel.online = false
      }

      await channel.save()
    }
  }
}

const TwitchWatcher = new TwitchWatcherClass()
export default TwitchWatcher
