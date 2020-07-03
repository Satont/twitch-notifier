import VkService from '../messengers/vk'
import TelegramService from '../messengers/telegram'
import { Op } from 'sequelize'
import { User } from '../models/User'
import { StreamMetadata } from './twitch'
import { error } from '../helpers/logs'

const spaceRegexp = /^\s*$/

export const notify = async (metadata: StreamMetadata) => {
  if (!metadata) return
  try {
    await sendVk(metadata)
    await sendTelegram(metadata)
  } catch (e) {
    console.debug(e)
  }
}

const sendVk = async (streamMetaData: StreamMetadata) => {
  const photo = await VkService.uploadPhoto(streamMetaData.preview.template.replace('{width}', '1280').replace('{height}', '720'))
  const game = streamMetaData.game ? `Игра: ${streamMetaData.game}\n` : ''
  const title = spaceRegexp.test(streamMetaData.channel.status) ? '' : `Название стрима: ${streamMetaData.channel.status}\n`
  const users = (await User.findAll({ 
    where: { follows: { [Op.contains]: [Number(streamMetaData.channel._id)] }, service: 'vk' },
    raw: true
  })).map(o => o.id)
  const message = `${streamMetaData.channel.display_name} онлайн!\n${game}${title}https://twitch.tv/${streamMetaData.channel.name}`
  await VkService.sendMessage({ target: users, message, image: photo.toString() })
}

const sendTelegram = async (streamMetaData: StreamMetadata) => {
  const game = streamMetaData.game ? `Game: ${streamMetaData.game}\n` : ''
  const title = spaceRegexp.test(streamMetaData.channel.status) ? '' : `Title: ${streamMetaData.channel.status}\n`
  const preview = streamMetaData.preview.template.replace('{width}', '1280').replace('{height}', '720')
  const users = (await User.findAll({ 
    where: { follows: { [Op.contains]: [Number(streamMetaData.channel._id)] }, service: 'telegram' },
    raw: true
  })).map(o => o.id)
  const message = `${streamMetaData.channel.display_name} online!\n${game}${title}https://twitch.tv/${streamMetaData.channel.name}`
  TelegramService.sendMessage({ target: users, message, image: `${preview}?timestamp=${Date.now()}` })
}

export const notifyGameChange = async (streamer: { name: string, id: number }, oldGame: string, newGame: string) => {
  try {
    sendVkGameChange(streamer, oldGame, newGame)
    sendTelegramGameChange(streamer, oldGame, newGame)
  } catch (e) {
    error(e)
  }
}

const sendVkGameChange = async (streamer: { name: string, id: number }, oldGame: string, newGame: string) => {
  const users = (await User.findAll({ 
    where: { follows: { [Op.contains]: [streamer.id] }, service: 'vk', follow_game_change: true },
    raw: true
  })).map(o => o.id)
  const message = `${streamer.name} теперь стримит ${newGame}! Предыдущая категория: ${oldGame}\nhttps://twitch.tv/${streamer.name}`
  await VkService.sendMessage({ target: users, message })
}

const sendTelegramGameChange = async (streamer: { name: string, id: number }, oldGame: string, newGame: string) => {
  const users = (await User.findAll({
    where: { follows: { [Op.contains]: [streamer.id] }, service: 'telegram', follow_game_change: true },
    raw: true
  })).map(o => o.id)
  const message = `${streamer.name} now streaming ${newGame}! Previous category: ${oldGame}\nhttps://twitch.tv/${streamer.name}`
  TelegramService.sendMessage({ target: users, message,  })
}
