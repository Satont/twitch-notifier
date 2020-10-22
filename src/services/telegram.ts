import { ServiceInterface } from './_interface'
import Telegraf, { Context } from 'telegraf'
import { error, info } from '../libs/logger'
import { orm } from '../libs/db'
import { Chat } from '../entities/Chat'

export default new class Telegram extends ServiceInterface {
  service = 'telegram'
  bot: Telegraf<any> = null

  async init() {
    const accessToken = process.env.TELEGRAM_BOT_TOKEN
    if (!accessToken) return false

    try {
      this.bot = new Telegraf(accessToken)
      
      await this.bot.launch()
      this.bot.use(async (ctx, next) => {
        await this.ensureUser(ctx)
        next()
      })
      this.bot.on('message', (msg) => this.listener(msg))
      info('Telegram Service initialized.')
    } catch (e) {
      error(e)
    }
  }

  async ensureUser(ctx: Context) {
    const repository = orm.em.fork().getRepository(Chat)
    const data = { chatId: ctx.chat.id }
    const user = await repository.findOne(data) || repository.assign(new Chat(), data)
    await repository.persistAndFlush(user)
    
    ctx.ChatEntity = user
  }

  async listener(msg: Context) {
    if (!msg.chat?.id || !msg.message?.text) return
    const commandName = msg.message.text.substring(1).split(' ')[0]
    const args = msg.message.text.split(' ').slice(1)
    const arg = msg.message.text.substring(1).replace(commandName, '')

    const command = this.commands.find(c => c.name === commandName)
    if (!command) return

    command['fnc'].call(Telegram, msg, args, arg)
    return true
  }

  async follow(msg: Context, args?: string[], arg?: string) {
    return true
  }
}