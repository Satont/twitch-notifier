import { SendMessageOpts, ServiceInterface } from './_interface'
import Telegraf, { Context } from 'telegraf'
import { chatIn, chatOut, error, info, warning } from '../libs/logger'
import { Chat } from '../entities/Chat'
import { getConnection } from 'typeorm'
import { command } from '../decorators/command'
import { followCommand } from '../commands/follow'
import { followsCommand } from '../commands/follows'
import { liveCommand } from '../commands/live'

class Telegram extends ServiceInterface {
  service = 'telegram'
  bot: Telegraf<any> = null

  async init() {
    const accessToken = process.env.TELEGRAM_BOT_TOKEN
    if (!accessToken) {
      warning('TELEGRAM: bot token not setuped, telegram library will not works.')
      return
    }

    try {
      this.bot = new Telegraf(accessToken)

      await this.bot.launch()
      this.bot.use(async (ctx: Context, next) => {
        chatIn(`TG [${ctx.from.username}]: ${ctx.message.text}`)
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
    const repository = getConnection().getRepository(Chat)
    const data = { id: String(ctx.chat.id) }
    const user = await repository.findOne(data, { relations: ['follows', 'follows.channel'] }) || await repository.create(data).save()

    ctx.ChatEntity = user
  }

  async listener(msg: Context) {
    if (!msg.chat?.id || !msg.message?.text) return
    const commandName = msg.message.text.substring(1).split(' ')[0]
    const args = msg.message.text.split(' ').slice(1)
    const arg = msg.message.text.substring(1).replace(commandName, '')

    const command = this.commands.find(c => c.name === commandName)
    if (!command) return

    await this[command.fnc](msg, args, arg)
    return true
  }

  @command('follow')
  async follow(msg: Context, args?: string[], arg?: string) {
    if (!arg) return msg.reply('arg is empty')
    this.sendMessage({
      target: String(msg.chat.id),
      message: await followCommand({ chat: msg.ChatEntity, channelName: arg }),
    })
  }

  @command('follows')
  async follows(msg: Context) {
    this.sendMessage({
      target: String(msg.chat.id),
      message: await followsCommand({ chat: msg.ChatEntity }),
    })
  }

  @command('live')
  async live(msg: Context) {
    this.sendMessage({
      target: String(msg.chat.id),
      message: await liveCommand({ chat: msg.ChatEntity }),
    })
  }

  async sendMessage(opts: SendMessageOpts) {
    const targets = Array.isArray(opts.target) ? opts.target : [opts.target]
    for (const target of targets) {
      const log = () => chatOut(`TG [${target}]: ${opts.message}`.replace(/(\r\n|\n|\r)/gm, ' '))
      if (opts.image) {
        const image = opts.image.replace('{width}', '1080').replace('{height}', '1920')
        this.bot?.telegram.sendPhoto(target, image, {
          caption: opts.message,
          parse_mode: 'HTML',
        }).then(() => {
          log()
        }).catch(console.error)
      } else {
        this.bot?.telegram.sendMessage(target, opts.message, {
          disable_web_page_preview: true,
          parse_mode: 'HTML',
        }).then(() => {
          log()
        }).catch(console.error)
      }
    }
  }
}

export default new Telegram()
