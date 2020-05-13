import Twitch from '../libs/twitch'
import { User } from '../models/User'
import { Channel } from '../models/Channel'
import { remove } from 'lodash'

export default async ({ service, userId, channel }: { service: 'telegram' | 'vk', userId: number, channel: string }): Promise<boolean> => {
  if (/[^a-zA-Z0-9_]/gmu.test(channel)) {
    throw new Error('Username can cointain only "a-z", "0-9" and "_" symbols.')
  }
  const streamer = await Twitch.getChannel(channel)
  const [user] = await User.findOrCreate({ where: { id: userId, service }, defaults: { follows: [], service } })
  await Channel.findOrCreate({ where: { id: streamer.id }, defaults: { username: streamer.login, online: false } })
  if (!user.follows.includes(streamer.id)) {
    return false
  } else {
    remove(user.follows, (o: number) => o === streamer.id)
    await user.update({ follows: user.follows })
    return true
  }
}
