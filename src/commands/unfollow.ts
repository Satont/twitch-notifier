import { getConnection } from 'typeorm'
import { Chat } from '../entities/Chat'
import { Follow } from '../entities/Follow'
import { Twitch } from '../libs/twitch'

const followRepository = getConnection().getRepository(Follow)

export async function unFollowCommand({ chat, channelName }: { chat: Chat, channelName: string }) {
  channelName = channelName.replace(/\s/g, '')
  if (/[^a-zA-Z0-9_]/gmu.test(channelName) || !channelName.length) {
    return 'Username can cointain only "a-z", "0-9" and "_" symbols.'
  }

  const streamer = await Twitch.getUser({ name: channelName.toLowerCase() })
  const follow = await followRepository.findOne({
    chat,
    channel: { id: streamer.id },
  })

  if (!follow) {
    return `You are not followed to ${streamer.name}`
  } else {
    await follow.remove()
    return `Successuly unfollowed from ${streamer.displayName}`
  }
}
