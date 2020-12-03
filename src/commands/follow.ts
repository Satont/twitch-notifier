import { getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'
import { Chat } from '../entities/Chat'
import { Follow } from '../entities/Follow'
import { I18n } from '../libs/i18n'
import { Twitch } from '../libs/twitch'
import TwitchWatcher from '../watchers/twitch'

const channelRepository = getConnection().getRepository(Channel)
const followRepository = getConnection().getRepository(Follow)

export async function followCommand({ chat, channelName, i18n }: { chat: Chat, channelName: string, i18n: I18n }) {
  channelName = channelName.replace(/\s/g, '')
  if (/[^a-zA-Z0-9_]/gmu.test(channelName) || !channelName.length) {
    return i18n.translate('commands.follow.errors.badUsername')
  }

  const streamer = await Twitch.getUser({ name: channelName.toLowerCase() })
  if (!streamer) {
    return i18n.translate('commands.follow.errors.streamerNotFound', { streamer: streamer.displayName })
  }

  const channel = await channelRepository.findOne({ id: streamer.id }) || await channelRepository.create({
    id: streamer.id,
    username: streamer.name,
  }).save()
  TwitchWatcher.addChannelToWebhooks(channel.id)

  if (chat.follows?.find(f => f.channel.id === streamer.id)) {
    return i18n.translate('commands.follow.alreadyFollowed', { streamer: streamer.displayName })
  } else {
    const follow = await followRepository.create({ chat, channel }).save()
    chat.follows.push(follow)
    await chat.save()
    return i18n.translate('commands.follow.success', { streamer: streamer.displayName })
  }
}
