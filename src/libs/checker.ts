import { Twitch } from './twitch'
import { Channel } from '../models/Channel'
import { User } from '../models/User'
import { chunk } from 'lodash'
import { say } from './vk'
const twitch = new Twitch(process.env.TWITCH_CLIENTID)

async function check () {
  setTimeout(() => check(), 5 * 60 * 1000)
  const channels = await Channel.findAll()
  let chunks = chunk(channels, 100)

  for (let chunk of chunks) {
    const checkChannels = await twitch.checkOnline(chunk.map(o => Number(o.id)))
    for (let channel of checkChannels) {
      const dbChannel = await Channel.findOne({ where: { id: Number(channel.user_id) }})
      if (dbChannel.online) return false
      else {
        await Channel.update({ online: true }, { where: { id: dbChannel.id }})
        const users = await User.findAll({ where: { follows: [dbChannel.id] }})
        say(users.map(o => o.id), `${channel.user_name} just wen't online!`, null)
      }
    }
  }
}
check()
