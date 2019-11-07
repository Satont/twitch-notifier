import { Twitch } from './twitch'
import { Channel } from '../models/Channel'
import { User } from '../models/User'
import { chunk, flattenDeep } from 'lodash'
import { say as sayVK, bot as vk } from './vk'
import { say as sayTG } from './telegram'
import { config } from '../helpers/config'
import { error } from '../helpers/logs'
import { Op } from 'sequelize'
import axios from 'axios'

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
  const streamMetaData = await getStreamMetaData(streamerId)
  const game = streamMetaData.game ? `Игра: ${streamMetaData.game}\n` : ''
  const title = streamMetaData.channel.status ? `Название стрима: ${streamMetaData.channel.status}\n` : ''
  const users = await User.findAll({ 
    where: { follows: { [Op.contains]: [streamerId] }, service: 'vk' },
    raw: true
  })
  const photo = await vk.upload.messagePhoto({
    source: streamMetaData.preview.template.replace('{width}', '1280').replace('{height}', '720')
  })
  sayVK(users.map(o => o.id), `${streamerName} онлайн!\n${game}${title}https://twitch.tv/${streamerName}`, photo.toString())
}

async function notifyTg (streamerName: string, streamerId: number) {
  const streamMetaData = await getStreamMetaData(streamerId)
  const game = streamMetaData.game ? `Game: ${streamMetaData.game}\n` : ''
  const title = streamMetaData.channel.status ? `Title: ${streamMetaData.channel.status}\n` : ''
  const preview = streamMetaData.preview.template.replace('{width}', '1280').replace('{height}', '720')
  const users = await User.findAll({ 
    where: { follows: { [Op.contains]: [streamerId] }, service: 'telegram' },
    raw: true
  })
  sayTG(users.map(o => o.id), `${streamerName} online!\n${game}${title}https://twitch.tv/${streamerName}`, `${preview}?timestamp=${Date.now()}`)
}

async function getStreamMetaData (id: number) {
  let request: any
  try {
    request = await axios.get(`https://api.twitch.tv/kraken/streams/${id}`, { 
      headers: { 'Accept': 'application/vnd.twitchtv.v5+json', 'Client-ID': config.twitch.clientId }
    })
    return request.data.stream
  } catch (e) {
    error(e)
    throw new Error(e)
  }
}