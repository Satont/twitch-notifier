import { Twitch } from './twitch'
import { Channel } from '../models/Channel'
import { User } from '../models/User'
import { chunk, flattenDeep } from 'lodash'
import { say as sayVK } from './vk'
import { say as sayTG } from './telegram'
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
      notifyVk(channel.user_name, dbChannel.id)
      notifyTg(channel.user_name, dbChannel.id)
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

async function notifyVk (streamerName: string, streamerId: number) {
  const users = await User.findAll({ 
    where: { follows: { [Op.contains]: [streamerId], service: 'vk' } },
    raw: true
  })
  sayVK(users.map(o => o.id), `${streamerName} онлайн!\nhttps://twitch.tv/${streamerName}`)
}

async function notifyTg (streamerName: string, streamerId: number) {
  const users = await User.findAll({ 
    where: { follows: { [Op.contains]: [streamerId], service: 'telegram' } },
    raw: true
  })
  sayTG(users.map(o => o.id), `${streamerName} online!\nhttps://twitch.tv/${streamerName}`)
}