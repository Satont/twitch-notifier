import Twitch from '../libs/twitch'
import { User } from '../models/User'

export default async (
  { userId, service}: { userId: number, service: 'vk' | 'telegram'}
): Promise<boolean | string[]> => {
  const user = await User.findOne({ where: { id: userId, service } })
  if (!user.follows.length) {
    return false
  } else {
    const channels = await Twitch.getChannelsById(user.follows)
    return channels.map(o => o.login)
  }
}
