import { getConnection } from 'typeorm'
import { Channel } from '../../entities/Channel'
import { getLatestStream } from '../../helpers/getLatestStream'
import { Announcer } from '../../libs/announcer'
import Twitch from '../../libs/twitch'
import * as locks from './locks'


class TwitchWatcherPolling {
  private timeout: NodeJS.Timeout

  async init() {
    this.poll()
  }

  async poll() {
    clearTimeout(this.timeout)
    this.timeout = setTimeout(() => this.poll(), 5 * 60 * 1000)

    const channels = await getConnection().getRepository(Channel).find()

    for (const channelIndex in channels) {
      const channel = channels[channelIndex]

      const [stream, latestStream] = await Promise.all([
        Twitch.apiClient.streams.getStreamByUserId(channel.id),
        getLatestStream(channel.id),
      ])

      const announcer = new Announcer(channel.id)
      await announcer.init()

      if (stream && stream.id !== latestStream.id) {
        return announcer.announceLive({ displayName: stream.userDisplayName, userId: stream.userId, stream })
      }

      if (!stream && channel.online) {
        return announcer.announceOffline({ displayName: stream.userDisplayName, userId: stream.userId })
      }

      if (stream.id === latestStream.id) {
        if (stream.gameName === latestStream.category) return

        announcer.announceUpdate({ displayName: stream.userDisplayName, newCategory: stream.gameName, userId: stream.userId })
      }
    }
  }
}

const instance = new TwitchWatcherPolling()
export default instance