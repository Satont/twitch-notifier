import { Twitch } from './twitch'
import { Channel } from '../models/Channel'
import { User } from '../models/User'
import { chunk } from 'lodash'
import { say } from './vk'
import { config } from '../helpers/config'
import { Op } from 'sequelize'
const twitch = new Twitch(config.twitch.clientId)

async function check () {
  setTimeout(() => check(), 5 * 60 * 1000)
  const channels = await Channel.findAll({ raw: true })
  let chunks = chunk(channels, 100)
  for (let chunk of chunks) {
    const checkChannels = await twitch.checkOnline(chunk.map(o => Number(o.id)))
    for (let channel of checkChannels) {
      const dbChannel = await Channel.findOne({ where: { id: Number(channel.user_id) }})
      if (dbChannel.online) return false
      else {
        await Channel.update({ online: true }, { where: { id: dbChannel.id }})
        const users = await User.findAll({ 
          where: { 
            follows: { [Op.contains]: [dbChannel.id] }  
          },
          raw: true
        })
        say(users.map(o => o.id), `${channel.user_name} онлайн!`)
      }
    }
  }
}
check()
