import { getConnection } from 'typeorm'
import { Chat } from '../entities/Chat'
import { Follow } from '../entities/Follow'
import { I18n } from '../libs/i18n'
import { Twitch } from '../libs/twitch'

const followRepository = getConnection().getRepository(Follow)

export async function unFollowCommand({ chat, channelName, i18n }: { chat: Chat, channelName: string, i18n: I18n }) {
  channelName = channelName.replace(/\s/g, '')
  if (/[^a-zA-Z0-9_]/gmu.test(channelName) || !channelName.length) {
    return {
      success: false,
      message: i18n.translate('commands.follow.errors.badUsername'),
    }
  }

  const streamer = await Twitch.getUser({ name: channelName.toLowerCase() })
  if (!streamer) {
    return {
      success: false,
      message: i18n.translate('commands.follow.errors.streamerNotFound', { streamer: streamer.displayName }),
    }
  }

  const follow = await followRepository.findOne({
    chat,
    channel: { id: streamer.id },
  })

  if (!follow) {
    return {
      success: true,
      message: i18n.translate('commands.unfollow.notFollowed', { streamer: streamer.name }),
    }
  } else {
    await follow.remove()
    return {
      success: true,
      message: i18n.translate('commands.unfollow.success', { streamer: streamer.name }),
    }
  }
}
