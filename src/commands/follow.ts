import { Twitch } from '../libs/twitch'
import { config } from '../helpers/config'
import { User } from '../models/User'
import { Channel } from '../models/Channel'

const twitch = new Twitch(config.twitch.clientId)

export default async ({ service, userId, channel }: { service: 'telegram' | 'vk', userId: number, channel: string }) => {
  if (/[^a-zA-Z0-9_]/gmu.test(channel)) {
    throw new Error('Username can cointain only "a-z", "0-9" and "_" symbols.')
  }
  const streamer = await twitch.getChannel(channel)
  const [user] = await User.findOrCreate({ where: { id: userId, service }, defaults: { follows: [], service } })
  await Channel.findOrCreate({ where: { id: streamer.id }, defaults: { username: streamer.login, online: false } })
  if (user.follows.includes(streamer.id)) {
     throw new Error(`You already followed to ${streamer.displayName}.`)
  } else {
    user.follows.push(streamer.id)
    await user.update({ follows: user.follows })
    return true
  }
}
