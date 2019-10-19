import { Twitch } from './twitch'
import { Channel } from '../models/Channel'
import { User } from '../models/User'
import { chunk, flattenDeep } from 'lodash'
import { say } from './vk'
import { config } from '../helpers/config'
import { Op } from 'sequelize'
const twitch = new Twitch(config.twitch.clientId)

async function check () {
  setTimeout(() => check(), 5 * 60 * 1000)
  const dbChannels = await Channel.findAll()
  const onlineChannels = flattenDeep(await getOnlineStreams(dbChannels.map(o => o.id)))

  for (let dbChannel of dbChannels) {
    const channel = onlineChannels.find(o => Number(o.user_id) === dbChannel.id)

    if (channel && !dbChannel.online) { // channel online, do notify
      await dbChannel.update({ online: true })
      const users = await User.findAll({ 
        where: { follows: { [Op.contains]: [dbChannel.id] } },
        raw: true
      })
      say(users.map(o => o.id), `${channel.user_name} онлайн!`)
    } else if (!channel && dbChannel.online) { // if channel offline but online in db, then set channel as offline in db
      await dbChannel.update({ online: false })
    } else if (channel && dbChannel.online) { // skip if channel online and online in db
      continue
    } else await dbChannel.update({ online: false })
  }
}
check()

async function getOnlineStreams(channels: number[]) {
  let onlineChannels: any[] = []
  const chunks = chunk(channels, 100)
  for (const chunk of chunks) {
    onlineChannels.push((await twitch.checkOnline(chunk)))
  }
  return onlineChannels
}
