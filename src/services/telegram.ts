import { SendMessageOpts, ServiceInterface } from './_interface'
import Telegraf, { Context, Markup } from 'telegraf'
import { chatIn, chatOut, error, info, warning } from '../libs/logger'
import { Chat, Services } from '../entities/Chat'
import { getConnection } from 'typeorm'
import { command } from '../decorators/command'
import { followCommand } from '../commands/follow'
import { followsCommand } from '../commands/follows'
import { liveCommand } from '../commands/live'
import { telegramAction } from '../decorators/telegramAction'
import { ChatSettings } from '../entities/ChatSettings'

class Telegram extends ServiceInterface {
  readonly service = Services.TELEGRAM
  bot: Telegraf<any> = null
  private readonly chatRepository = getConnection().getRepository(Chat)

  constructor() {
    super()

    const accessToken = process.env.TELEGRAM_BOT_TOKEN
    if (!accessToken) {
      warning('TELEGRAM: bot token not setuped, telegram library will not works.')
      return
    }
    this.bot = new Telegraf(accessToken)
  }

  async init() {
    try {
      await this.bot.launch()
      await this.bot.telegram.setMyCommands(this.commands.map(c => ({ command: c.name, description: c.description })))

      this.bot.use(async (ctx: Context, next) => {
        if (ctx.message?.text) chatIn(`TG [${ctx.from.username}]: ${ctx.message?.text}`)
        await this.ensureUser(ctx)
        next()
      })
      this.bot.on('message', (msg) => this.listener(msg))

      info('Telegram Service initialized.')
      this.inited = true
    } catch (e) {
      error(e)
    }
  }

  async ensureUser(ctx: Context) {
    const data = { chatId: String(ctx.chat.id), service: Services.TELEGRAM }
    const chat = await this.chatRepository.findOne(data, { relations: ['follows', 'follows.channel'] })
      ?? this.chatRepository.create({ ...data, settings: new ChatSettings() })
    chat.save()

    ctx.ChatEntity = chat
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

  @command('follow', { description: 'Follow to some user.' })
  async follow(msg: Context, args?: string[], arg?: string) {
    if (!arg) return msg.reply('arg is empty')
    this.sendMessage({
      target: String(msg.chat.id),
      message: await followCommand({ chat: msg.ChatEntity, channelName: arg }),
    })
  }

  @command('follows', { description: 'Shows list of your follows.' })
  async follows(msg: Context) {
    this.sendMessage({
      target: String(msg.chat.id),
      message: await followsCommand({ chat: msg.ChatEntity }),
    })
  }

  @command('live', { description: 'Check currently live streams from your follow list.' })
  async live(msg: Context) {
    this.sendMessage({
      target: String(msg.chat.id),
      message: await liveCommand({ chat: msg.ChatEntity }),
    })
  }



  @command('settings', { description: 'Settings menu.' })
  @telegramAction('get_settings')
  async settings(msg: Context) {
    const getInlineKeyboard = () => Markup.inlineKeyboard([
      Markup.callbackButton('Game change notification', 'game_change_notification'),
    ])

    if (msg.message?.text) {
      await msg.reply('This is list of bot settings', getInlineKeyboard().extra())
    } else if (msg.isAction) {
      msg.editMessageReplyMarkup(getInlineKeyboard())
    } else {
      return getInlineKeyboard()
    }
  }

  @telegramAction('game_change_notification')
  async gameChangeNotificationAction(msg: Context) {
    await msg.editMessageReplyMarkup(Markup.inlineKeyboard([
      Markup.urlButton('Test', 'http://vk.com'),
      Markup.callbackButton('«', 'get_settings'),
    ]))
  }

  async sendMessage(opts: SendMessageOpts) {
    const targets = Array.isArray(opts.target) ? opts.target : [opts.target]
    for (const target of targets) {
      const log = () => chatOut(`TG [${target}]: ${opts.message}`.replace(/(\r\n|\n|\r)/gm, ' '))
      if (opts.image) {
        this.bot?.telegram.sendPhoto(target, opts.image, {
          caption: opts.message,
          parse_mode: 'HTML',
        })
          .then(() => log())
          .catch(console.error)
      } else {
        this.bot?.telegram.sendMessage(target, opts.message, {
          disable_web_page_preview: true,
          parse_mode: 'HTML',
        })
          .then(() => log())
          .catch(console.error)
      }
    }
  }
}

export default new Telegram()
