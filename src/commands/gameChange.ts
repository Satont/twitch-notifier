import { User } from '../models/User'

export default async (
  { userId, service}: { userId: number, service: 'vk' | 'telegram'}
): Promise<boolean> => {
  const user = await User.findOne({ where: { id: userId, service } })
  const current = user.follow_game_change
  await user.update({ follow_game_change: !current })

  return !current
}
