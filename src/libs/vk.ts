import { VkBot } from 'nodejs-vk-bot'
import { info, error } from '../helpers/logs'
import { Twitch, Methods } from './twitch'

const bot = new VkBot(process.env.VKTOKEN)
const twitch = new Twitch(process.env.TWITCH_CLIENTID)

bot.command(['!подписка', '!follow'], async (ctx) => {
  const streamer: string = ctx.message.text.split(' ').slice().join(' ')
  try {
    const request = await twitch.request({ method: Methods.GET, endpoint: 'users', data: { login: streamer} })
    info(request)
    ctx.reply(request.data)
  } catch (e) {
    error(e.message)
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
