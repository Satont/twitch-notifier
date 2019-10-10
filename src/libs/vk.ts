import { VkBot } from 'nodejs-vk-bot'
 
const bot = new VkBot(process.env.VKTOKEN)
 
bot.command(['!подписка', '!follow'], (ctx) => {
  ctx.reply('Hello!')
})

bot.command(['!отписка', '!unfollow'], (ctx) => {
  ctx.reply('Hello!')
})

export function say(userId: number | number[], message: string, attachment: string) {
  bot.sendMessage(userId, message, attachment)
} 

bot.startPolling().then(() => {
  console.log('VK bot connected.')
})
