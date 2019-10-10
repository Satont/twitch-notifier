const VkBot = require('node-vk-bot-api')
 
const bot = new VkBot(process.env.VKTOKEN)
 
bot.command(['!подписка', '!follow'], (ctx) => {
  ctx.reply('Hello!')
})

bot.command(['!отписка', '!unfollow'], (ctx) => {
  ctx.reply('Hello!')
})

module.exports.say = (userId, message, attachment) => bot.sendMesage(userId, message, attachment)

bot.startPolling()