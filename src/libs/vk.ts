import { VkBot } from 'nodejs-vk-bot'
import { info, error } from '../helpers/logs'
import { Twitch } from './twitch'
import { User } from '../models/User'
import { config } from '../helpers/config'

const bot = new VkBot(config.vk.token)
const twitch = new Twitch(config.twitch.clientId)

bot.command(['!подписка', '!follow'], async (ctx) => {
  const user = await User.findOne({ where: { id: 123 } })
  console.log(user)
  let streamer: string = ctx.message.text.split(' ')[1]
  try {
    streamer = await twitch.getChannel(streamer)
  } catch (e) {
    error(e)
    ctx.reply(e.message)
  }
})

bot.command(['!отписка', '!unfollow'], (ctx) => {
  ctx.reply('Hello!')
})

export function say(userId: number | number[], message: string, attachment: string) {
  bot.sendMessage(userId, message, attachment)
} 

bot.startPolling().then(() => {
  info('VK bot connected.')
})
