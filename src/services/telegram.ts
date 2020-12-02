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
import { ChatSettings, Languages } from '../entities/ChatSettings'
import { unFollowCommand } from '../commands/unfollow'
import { i18n } from '../libs/i18n'

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
    this.bot.use(async (ctx: Context, next) => {
      if (ctx.message?.text) chatIn(`TG [${ctx.from.username}]: ${ctx.message?.text}`)

      ctx.ChatEntity = await this.ensureUser(ctx)
      ctx.i18n = i18n.clone(ctx.ChatEntity.settings.language)
      next()
    })
    this.bot.on('message', (msg) => this.listener(msg))
  }

  async init() {
    try {
      await this.bot.launch()
      await this.bot.telegram.setMyCommands(this.commands.map(c => ({ command: c.name, description: c.description })))

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

    return chat
  }

  async listener(ctx: Context) {
    if (!ctx.chat?.id || !ctx.message?.text) return
    const commandName = ctx.message.text.substring(1).split(' ')[0]
    const args = ctx.message.text.split(' ').slice(1)
    const arg = ctx.message.text.substring(1).replace(commandName, '')

    const command = this.commands.find(c => c.name === commandName)
    if (!command) return

    await this[command.fnc](ctx, args, arg)
    return true
  }

  @command('follow', { description: 'Follow to some user.' })
  async follow(ctx: Context, args: string[], arg: string) {
    this.sendMessage({
      target: String(ctx.chat.id),
      message: await followCommand({ chat: ctx.ChatEntity, channelName: arg, i18n: ctx.i18n }),
    })
  }

  @command('unfollow', { description: 'Unfollow from some user.' })
  async unfollow(ctx: Context, args: string[], arg: string) {
    this.sendMessage({
      target: String(ctx.chat.id),
      message: await unFollowCommand({ chat: ctx.ChatEntity, channelName: arg, i18n: ctx.i18n }),
    })
  }

  @command('test', { description: 'test' })
  async test(ctx: Context) {
    ctx.reply(ctx.i18n.translate(ctx.message.text.replace('/test ', ''), { t: ctx.from.username }))
  }

  @command('follows', { description: 'Shows list of your follows.' })
  async follows(ctx: Context) {
    this.sendMessage({
      target: String(ctx.chat.id),
      message: await followsCommand({ chat: ctx.ChatEntity, i18n: ctx.i18n }),
    })
  }

  @command('live', { description: 'Check currently live streams from your follow list.' })
  async live(ctx: Context) {
    this.sendMessage({
      target: String(ctx.chat.id),
      message: await liveCommand({ chat: ctx.ChatEntity, i18n: ctx.i18n }),
    })
  }

  @command('start', { description: 'Start command' })
  @command('settings', { description: 'Settings menu.' })
  @telegramAction('get_settings')
  async settings(ctx: Context) {
    const getInlineKeyboard = () => Markup.inlineKeyboard([
      Markup.callbackButton(ctx.i18n.translate('telegram.settings.game_change_notification_setting.button'), 'game_change_notification_setting'),
      Markup.callbackButton(ctx.i18n.translate('telegram.settings.language.button'), 'language_setting'),
    ])

    if (ctx.message?.text) {
      await ctx.reply(ctx.i18n.translate('bot.description'), getInlineKeyboard().extra())
    } else if (ctx.isAction) {
      ctx.editMessageReplyMarkup(getInlineKeyboard())
    } else {
      return getInlineKeyboard()
    }
  }

  @telegramAction('game_change_notification_setting')
  async gameChangeNotificationAction(ctx: Context) {
    await ctx.editMessageReplyMarkup(Markup.inlineKeyboard([
      Markup.callbackButton('Test', 'http://vk.com'),
      Markup.callbackButton('«', 'get_settings'),
    ]))
  }

  @telegramAction('language_setting')
  async language(ctx: Context) {
    await ctx.editMessageReplyMarkup(Markup.inlineKeyboard([
      Markup.callbackButton('English', 'language_set_english_setting'),
      Markup.callbackButton('Russian', 'language_set_russian_setting'),
      Markup.callbackButton('«', 'get_settings'),
    ]))
  }

  @telegramAction('language_set_english_setting')
  @telegramAction('language_set_russian_setting')
  async languageSetEnglish(ctx: Context) {
    const lang = ctx.callbackQuery.data.split('_')[2] as Languages
    ctx.ChatEntity.settings.language = lang
    ctx.i18n = ctx.i18n.clone(lang)
    await ctx.ChatEntity.save()
    ctx.reply(ctx.i18n.translate('telegram.settings.language.changed'))
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
