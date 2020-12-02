import { getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'
import { Chat } from '../entities/Chat'
import { Follow } from '../entities/Follow'
import { Twitch } from '../libs/twitch'
import TwitchWatcher from '../watchers/twitch'

const channelRepository = getConnection().getRepository(Channel)
const followRepository = getConnection().getRepository(Follow)

export async function followCommand({ chat, channelName }: { chat: Chat, channelName: string }) {
  channelName = channelName.replace(/\s/g, '')
  if (/[^a-zA-Z0-9_]/gmu.test(channelName) || !channelName.length) {
    return 'Username can cointain only "a-z", "0-9" and "_" symbols.'
  }

  const streamer = await Twitch.getUser({ name: channelName.toLowerCase() })
  const channel = await channelRepository.findOne({ id: streamer.id }) || await channelRepository.create({
    id: streamer.id,
    username: streamer.name,
  }).save()
  TwitchWatcher.addChannelToWebhooks(channel.id)

  if (chat.follows.find(f => f.channel.id === streamer.id)) {
    return `You already followed to ${streamer.displayName}.`
  } else {
    const follow = await followRepository.create({ chat, channel }).save()
    chat.follows.push(follow)
    await chat.save()
    return `You successfuly followed to ${streamer.displayName}`
  }
}
