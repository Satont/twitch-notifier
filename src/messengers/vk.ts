import { VK as VkIO } from 'vk-io'
import { info, error } from '../helpers/logs'
import { User } from '../models/User'
import follow from '../commands/follow'
import unfollow from '../commands/unfollow'
import follows from '../commands/follows'
import live from '../commands/live'
import { config } from '../helpers/config'
import { chunk, isBoolean } from 'lodash'
import { IService, SendMessageOpts } from './interface'

const service = 'vk'

class Vk extends IService {
  bot: VkIO
  constructor() {
    super()
  }
  protected init() {
    this.bot = new VkIO({ token: config.vk.token })
    this.bot.updates.start().then(() => info('VK bot connected.')).catch(e => error(e))
    this.loadMiddlewares()
    this.loadCommands()
  }
  public async sendMessage(opts: SendMessageOpts) {
    try {
      const targets = Array.isArray(opts.target) ? opts.target : [opts.target]
      const chunks = chunk(targets, 100)
      for (const chunk of chunks) {
        await this.bot.api.messages.send({
          random_id: Math.random() * (1000000000 - 9) + 10,
          user_ids: chunk,
          message: opts.message,
          dont_parse_links: true,
          attachment: opts.image
        })
      }
      return true
    } catch (e) {
      error(e)
      return false
    }
  }
  protected async loadMiddlewares() {
    this.bot.updates.on('message', async (ctx, next) => {
      if (ctx.senderId !== -187752469) info(`Vk | New message from: ${ctx.senderId}, message: ${ctx.text}`)
      const [user] = await User.findOrCreate({ where: { id: ctx.senderId, service }, defaults: { follows: [], service: 'vk' } })
      ctx.dbUser = user
      return next()
    })
  }
  protected async loadCommands() {
    this.bot.updates.hear(value => (value.startsWith('!подписка')), async (ctx) => {
      const channel: string = ctx.text.split(' ')[1]
      if (!channel) return ctx.reply('Вы должны указать на кого подписаться.')
      try {
        const followed = await follow({ service, userId: ctx.senderId, channel })
        if (followed) ctx.reply(`Вы успешно подписались на ${channel}.`)
      } catch (e) {
        error(e)
        ctx.reply(e.message)
      }
    })
    this.bot.updates.hear(value => (value.startsWith('!подписка')), async (ctx) => {
      const channel: string = ctx.text.split(' ')[1]
      if (!channel) return ctx.reply('Вы должны указать на кого подписаться.')
      try {
        const followed = await follow({ service, userId: ctx.senderId, channel })
        if (followed) ctx.reply(`Вы успешно подписались на ${channel}.`)
      } catch (e) {
        error(e)
        ctx.reply(e.message)
      }
    })
    this.bot.updates.hear(value => (value.startsWith('!отписка')), async (ctx) => {
      const channel: string = ctx.text.split(' ')[1]
      if (!channel) return ctx.reply('Вы должны указать от кого отписаться.')
      try {
        const unfollowed = await unfollow({ service, userId: ctx.senderId, channel })
        if (!unfollowed) ctx.reply(`Вы не подписаны на ${channel}.`)
        else ctx.reply(`Вы успешно отписались от ${channel}.`)
      } catch (e) {
        error(e)
        ctx.reply(e.message)
      }
    })
    this.bot.updates.hear(value => (value.startsWith('!подписки')), async (ctx) => {
      const followed = await follows({ userId: ctx.senderId, service })
      if (isBoolean(followed)) ctx.reply(`Вы ни на кого не подписаны.`) 
      else ctx.reply(`Вы подписаны на ${followed.join(', ')}.`)
    })
    this.bot.updates.hear(value => (value.startsWith('!онлайн')), async (ctx) => {
      const channels = await live({ userId: ctx.senderId, service })
      if (isBoolean(channels)) ctx.reply(`Сейчас нет ни одного канала онлайн из ваших подписок.`)
      else ctx.reply(`Сейчас онлайн: \n${channels.map((o) => 'https://twitch.tv/' + o).join('\n')}`)
    })
    this.bot.updates.hear(value => (value.startsWith('!команды')), async (ctx) => {
      ctx.reply(`
      На данный момент доступны следующие команды: 
      !подписка username
      !отписка username
      !подписки 
      !команды
      !онлайн
      `)
    })
  }
  public async uploadPhoto(source: string) {
    return await this.bot.upload.messagePhoto({ source })
  }
}

export default new Vk()
