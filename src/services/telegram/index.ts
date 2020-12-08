import { SendMessageOpts, ServiceInterface } from '../_interface'
import Telegraf, {  Markup, session, Stage } from 'telegraf'
import { chatIn, chatOut, error, info, warning } from '../../libs/logger'
import { Chat, Services } from '../../entities/Chat'
import { getConnection } from 'typeorm'
import { command } from '../../decorators/command'
import { followCommand } from '../../commands/follow'
import { followsCommand } from '../../commands/follows'
import { liveCommand } from '../../commands/live'
import { telegramAction } from '../../decorators/telegramAction'
import { ChatSettings } from '../../entities/ChatSettings'
import { unFollowCommand } from '../../commands/unfollow'
import { i18n } from '../../libs/i18n'
import { followScene } from './scenes/follow'
import { SceneContextMessageUpdate } from 'telegraf/typings/stage'
import { unFollowScene } from './scenes/unfollow'
import { unFollowAllCommand } from '../../commands/unFollowAll'

class Telegram extends ServiceInterface {
  bot: Telegraf<any> = null
  stage = new Stage([
    followScene,
    unFollowScene,
  ])
  private readonly chatRepository = getConnection().getRepository(Chat)

  constructor(token?: string) {
    super({
      service: Services.TELEGRAM,
    })

    const accessToken = token || process.env.TELEGRAM_BOT_TOKEN
    if (!accessToken) {
      warning('TELEGRAM: bot token not setuped, telegram library will not works.')
      return
    }

    this.bot = new Telegraf(accessToken)
    this.bot.use(async (ctx: Context, next) => {
      if (ctx.message?.text) chatIn(`TG [${ctx.from?.username || ctx.from?.id}]: ${ctx.message?.text}`)

      ctx.ChatEntity = await this.ensureUser(ctx)
      ctx.i18n = i18n.clone(ctx.ChatEntity.settings.language)
      next()
    })
    this.bot.use(session())
    this.bot.use(this.stage.middleware())
    this.bot.on('message', (msg) => this.listener(msg))
    this.bot.catch((err, ctx) => {
      error(err)
      error(ctx)
    })
  }

  async init() {
    try {
      if (!this.bot) return
      await this.bot.launch()
      this.bot.telegram.getMe()
      const commands = this.commands
        .filter(c => c.isVisible ?? true)
        .map(c => ({ command: c.name, description: c.description }))

      await this.bot.telegram.setMyCommands([
        ...commands,
        {
          command: 'cancel',
          description: 'Cancel current action.',
        },
      ])

      info('Telegram Service initialized.')
      this.inited = true
    } catch (e) {
      error(e)
    }
  }

  async ensureUser(ctx: SceneContextMessageUpdate) {
    if (!ctx.chat?.id) return null

    const data = { chatId: String(ctx.chat?.id), service: Services.TELEGRAM }
    const chat = await this.chatRepository.findOne(data, { relations: ['follows', 'follows.channel'] })
      ?? this.chatRepository.create({ ...data, settings: new ChatSettings() })
    await chat.save()

    return chat
  }

  async listener(ctx: SceneContextMessageUpdate) {
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
  async follow(ctx: SceneContextMessageUpdate, args: string[], arg: string) {
    if (!arg.length) {
      ctx.scene.enter('followScene')
    } else {
      const { message } = await followCommand({ chat: ctx.ChatEntity, channelName: arg, i18n: ctx.i18n })
      this.sendMessage({
        target: String(ctx.chat.id),
        message,
      })
    }
  }

  @command('unfollow', { description: 'Unfollow from some user.' })
  async unfollow(ctx: SceneContextMessageUpdate, args: string[], arg: string) {
    if (!arg.length) {
      ctx.scene.enter('unfollowScene')
    } else {
      const { message } = await unFollowCommand({ chat: ctx.ChatEntity, channelName: arg, i18n: ctx.i18n })
      this.sendMessage({
        target: String(ctx.chat.id),
        message,
      })
    }
  }

  @command('clearfollows', { description: 'Unfollow from all users.' })
  async unfollowAll(ctx: SceneContextMessageUpdate) {
    const { message } = await unFollowAllCommand({ chat: ctx.ChatEntity, i18n: ctx.i18n })
    this.sendMessage({
      target: String(ctx.chat.id),
      message,
    })
  }

  @command('follows', { description: 'Shows list of your follows.' })
  async follows(ctx: SceneContextMessageUpdate) {
    this.sendMessage({
      target: String(ctx.chat.id),
      message: await followsCommand({ chat: ctx.ChatEntity, i18n: ctx.i18n }),
    })
  }

  @command('live', { description: 'Check currently live streams from your follow list.' })
  async live(ctx: SceneContextMessageUpdate) {
    this.sendMessage({
      target: String(ctx.chat.id),
      message: await liveCommand({ chat: ctx.ChatEntity, i18n: ctx.i18n }),
    })
  }

  @command('start', { description: 'Start command' })
  @command('settings', { description: 'Settings menu.' })
  @telegramAction('get_settings')
  async settings(ctx: SceneContextMessageUpdate) {
    const getMarkEmoji = (state: boolean) => !state ? '◻︎' : '☑︎'

    const getInlineKeyboard = () => Markup.inlineKeyboard([
      Markup.callbackButton(
        `${getMarkEmoji(ctx.ChatEntity.settings.game_change_notification)} ${ctx.i18n.translate('settings.game_change_notification_setting.button')}`, 'game_change_notification_setting'
      ),
      Markup.callbackButton(
        `${getMarkEmoji(ctx.ChatEntity.settings.offline_notification)} ${ctx.i18n.translate('settings.offline_notification.button')}`,
        'offline_notification_setting'
      ),
      Markup.callbackButton(ctx.i18n.translate('settings.language.button'), 'language_setting'),
      Markup.urlButton('GitHub', 'https://github.com/Satont/twitch-notifier'),
      Markup.urlButton('Patreon', 'https://www.patreon.com/satont'),
    ], { columns: 1 })

    if (ctx.message?.text) {
      await ctx.reply(ctx.i18n.translate('bot.description'), getInlineKeyboard().extra())
    } else if (ctx.isAction) {
      await ctx.editMessageText(ctx.i18n.translate('bot.description'), getInlineKeyboard().extra())
    } else {
      return getInlineKeyboard()
    }
  }

  @telegramAction('game_change_notification_setting')
  async gameChangeNotificationAction(ctx: SceneContextMessageUpdate) {
    const currentState = ctx.ChatEntity.settings.game_change_notification
    ctx.ChatEntity.settings.game_change_notification = !currentState
    await ctx.ChatEntity.save()
    await this.settings(ctx)
  }

  @telegramAction('offline_notification_setting')
  async offLineNotificationAction(ctx: SceneContextMessageUpdate) {
    const currentState = ctx.ChatEntity.settings.offline_notification
    ctx.ChatEntity.settings.offline_notification = !currentState
    await ctx.ChatEntity.save()
    await this.settings(ctx)
  }

  @telegramAction('language_setting')
  async language(ctx: SceneContextMessageUpdate) {
    const buttons = Object.keys(i18n.translations).map(key => {
      const name = i18n.translations[key].language.name
      const emoji = i18n.translations[key].language.emoji
      return Markup.callbackButton(`${name} ${emoji}`, `language_set_${key}_setting`)
    })

    await ctx.editMessageReplyMarkup(Markup.inlineKeyboard([
      ...buttons,
      Markup.callbackButton('«', 'get_settings'),
    ], { columns: 1 }))
  }

  @telegramAction(Object.keys(i18n.translations).map(key => `language_set_${key}_setting`))
  async languageSet(ctx: SceneContextMessageUpdate) {
    const lang = ctx.callbackQuery.data.split('_')[2] as string
    ctx.ChatEntity.settings.language = lang
    ctx.i18n = ctx.i18n.clone(lang)
    await ctx.ChatEntity.save()
    await this.settings(ctx)
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
export {
  Telegram,
}
