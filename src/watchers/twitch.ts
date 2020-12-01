import { getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'
import Twitch from '../libs/twitch'
import { services } from '../services/_interface'
import { ITwitchStreamChangedPayload } from '../typings/twitch'

class TwitchWatcherClass {
  private readonly channelsRepository = getConnection().getRepository(Channel)

  async processPayload(data: ITwitchStreamChangedPayload['data']) {
    for (const stream of data) {
      const category = (await Twitch.apiClient.helix.games.getGameById(stream.game_id))?.name
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
            message: `${stream.user_name} online!\nCategory:${category}\nTitle:${stream.title}\nhttps://twitch.tv/${stream.user_name}`,
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
