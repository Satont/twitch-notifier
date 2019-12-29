import { Twitch } from '../libs/twitch'
import { config } from '../helpers/config'
import { User } from '../models/User'
import { Channel } from '../models/Channel'
import { Op } from 'sequelize'

const twitch = new Twitch(config.twitch.clientId)

export default async ({ userId, service}: { userId: number, service: 'vk' | 'telegram'}): Promise<boolean | string[]> => {
  const user = await User.findOne({ where: { id: userId, service } })
  if (!user.follows.length) {
    return false
  } else {
    const liveChannels = await Channel.findAll({
      where: {
        id: { [Op.in]: user.follows },
        online: true,
      },
      raw: true,
    })

    const channels = await twitch.getChannelsById(liveChannels.map((o) => o.id))
    return channels.map(o => o.displayName)
  }
}