import { SendMessageOpts, ServiceInterface } from './_interface'
import { VK as VKIO, MessageContext, Keyboard  } from 'vk-io'
import { error, info, warning } from '../libs/logger'
import { Chat, Services } from '../entities/Chat'
import { getConnection } from 'typeorm'
import { command } from '../decorators/command'
import { chunk } from 'lodash'
import { Languages } from '../entities/ChatSettings'
import { HearManager } from '@vk-io/hear'
import { i18n } from '../libs/i18n'
import { followCommand } from '../commands/follow'
import { followsCommand } from '../commands/follows'
import { liveCommand } from '../commands/live'
import { unFollowCommand } from '../commands/unfollow'
import { vkAction } from '../decorators/vkAction'

class VK extends ServiceInterface {
  service = Services.VK
  bot: VKIO = null
  hearManager: HearManager<MessageContext> = null

  async init() {
    const token = process.env.VK_GROUP_TOKEN
    if (!token) {
      warning('VK: group token not setuped, library will not works.')
      return
    }

    try {
      this.bot = new VKIO({ token })
      this.hearManager = new HearManager<MessageContext>()

      this.bot.updates.on('message', async (ctx, next) => {
        if (!ctx.isUser) return
        await this.ensureUser(ctx)
        ctx.i18n = i18n.clone(ctx.ChatEntity.settings.language)
        await this.listener(ctx)
        next()
      })
      this.bot.updates.on('message_new', this.hearManager.middleware)
      await this.bot.updates.start()
      info('VK Service initialized.')
      this.inited = true
    } catch (e) {
      error(e)
    }
  }

  async ensureUser(ctx: MessageContext) {
    const repository = getConnection().getRepository(Chat)
    const data = { chatId: String(ctx.chatId || ctx.peerId || ctx.senderId), service: Services.VK }
    const chat = await repository.findOne(data) || repository.create({ ...data, settings: { language: Languages.RUSSIAN } })
    chat.save()

    ctx.ChatEntity = chat
  }

  async listener(msg: MessageContext) {
    if (!msg.hasText || !msg.text?.startsWith('/')) return
    const commandName = msg.text.substring(1).split(' ')[0]
    const args = msg.text.split(' ').slice(1)
    const arg = msg.text.substring(1).replace(commandName, '')

    const command = this.commands.find(c => c.name === commandName)
    if (!command) return

    await this[command.fnc](msg, args, arg)
    return true
  }

  @command('start', { description: 'Start command' })
  async startCommand(ctx: MessageContext) {
    const description = ctx.i18n.translate('bot.description')
    ctx.send(`${description}\n\n${this.commands.map(c => `/${c.name}`).join('\n')}`)
  }

  private getInlineKeyboard = (ctx: MessageContext) => Keyboard.builder()
    .oneTime()
    .inline()
    .textButton({
      label: `${!ctx.ChatEntity.settings.game_change_notification ? '◻︎' : '☑︎'} ${ctx.i18n.translate('settings.game_change_notification_setting.button')}`,
      payload: {
        command: 'game_change_notification_setting',
      },
    })
    .textButton({
      label: ctx.i18n.translate('settings.language.button'),
      payload: {
        command: 'language_setting',
      },
    })

  @command('settings', { description: 'Settings command' })
  async settings(ctx: MessageContext) {
    ctx.send({
      message: ctx.i18n.translate('bot.description'),
      keyboard: this.getInlineKeyboard(ctx),
    })
  }

  @vkAction('language_setting')
  async languageMenu(ctx: MessageContext) {
    const keyboard = Keyboard.builder().oneTime().inline()
    Object.values(Languages).forEach(v => {
      const name = v.charAt(0).toUpperCase() + v.slice(1)
      const emoji = ctx.i18n.translate('language.emoji')
      keyboard.textButton({ label: `${emoji} ${name}`, payload: { command: `language_set_${v}_setting` } })
    })

    await ctx.send({
      message: '🌍',
      keyboard,
    })
  }

  @vkAction(Object.values(Languages).map(v => `language_set_${v}_setting`))
  async setLang(ctx: MessageContext) {
    const lang = ctx.messagePayload.command.split('_')[2] as Languages
    ctx.ChatEntity.settings.language = lang
    ctx.i18n = ctx.i18n.clone(lang)
    await ctx.ChatEntity.save()
    await this.settings(ctx)
  }

  @vkAction('game_change_notification_setting')
  async setGameChangeNotification(ctx: MessageContext) {
    const currentState = ctx.ChatEntity.settings.game_change_notification
    ctx.ChatEntity.settings.game_change_notification = !currentState
    await ctx.ChatEntity.save()
    await this.settings(ctx)
  }

  @command('follow', { description: 'Follow to some user.' })
  async follow(ctx: MessageContext, _args: string[], arg: string) {
    this.sendMessage({
      target: String(ctx.chat.id),
      message: await followCommand({ chat: ctx.ChatEntity, channelName: arg, i18n: ctx.i18n }),
    })
  }

  @command('unfollow', { description: 'Unfollow from some user.' })
  async unfollow(ctx: MessageContext, _args: string[], arg: string) {
    this.sendMessage({
      target: String(ctx.chat.id),
      message: await unFollowCommand({ chat: ctx.ChatEntity, channelName: arg, i18n: ctx.i18n }),
    })
  }

  @command('follows', { description: 'Shows list of your follows.' })
  async follows(ctx: MessageContext) {
    this.sendMessage({
      target: String(ctx.chat.id),
      message: await followsCommand({ chat: ctx.ChatEntity, i18n: ctx.i18n }),
    })
  }

  @command('live', { description: 'Check currently live streams from your follow list.' })
  async live(ctx: MessageContext) {
    this.sendMessage({
      target: String(ctx.chat.id),
      message: await liveCommand({ chat: ctx.ChatEntity, i18n: ctx.i18n }),
    })
  }

  async sendMessage(opts: SendMessageOpts) {
    const targets = Array.isArray(opts.target) ? opts.target : [opts.target]
    const chunks = chunk(targets.map(t => Number(t)), 100)
    const attachment = opts.image ? await this.uploadPhoto(opts.image) : undefined
    for (const chunk of chunks) {
      await this.bot.api.messages.send({
        random_id: Math.random() * (1000000000 - 9) + 10,
        user_ids: chunk,
        message: opts.message,
        dont_parse_links: true,
        attachment,
      })
    }
  }

  public async uploadPhoto(source: string) {
    return await this.bot.upload.messagePhoto({
      source: {
        value: source,
      },
    })
  }
}

export default new VK()
