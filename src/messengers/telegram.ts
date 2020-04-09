import Telegraf, { SceneContextMessageUpdate, Stage, session, BaseScene, TelegramOptions } from 'telegraf'
import { config } from '../helpers/config'
import { info, error } from '../helpers/logs'
import { User } from '../models/User'
import follow from '../commands/follow'
import unfollow from '../commands/unfollow'
import live from '../commands/live'
import follows from '../commands/follows'
import { isBoolean } from 'util'
import { IService, SendMessageOpts } from './interface'
import { agent as ProxyAgent } from '../helpers/tgProxy'

const service = 'telegram'
const telegramOpts: TelegramOptions = {
  agent: ProxyAgent as any
} 

class Telegram extends IService {
  scenes: { [x: string]: BaseScene<SceneContextMessageUpdate> } = {
    follow: new BaseScene('follow', {
      ttl: 60,
      enterHandlers: [
        (ctx) => ctx.reply('Enter streamer username you want to follow')
      ],
      handlers: [
        async (ctx) => {
          const channel = ctx.message.text
          try {
            const followed = await follow({ userId: ctx.from.id, service, channel })
            if (followed) {
              ctx.reply(`You successfuly followed to ${channel}`)
            }
          } catch (e) {
            error(e)
            ctx.reply(e.message)
          }
          ctx.scene.leave()
        }
      ]
    }),
    unfollow: new BaseScene('unfollow', {
      ttl: 60,
      enterHandlers: [
        (ctx) => ctx.reply('Enter streamer username you want to unfollow')
      ],
      handlers: [
        async (ctx) => {
          const channel = ctx.message.text
          try {
            const unfollowed = await unfollow({ service, userId: ctx.from.id, channel })
            if (!unfollowed) {
              ctx.reply(`You aren't followed to ${channel}.`)
            } else {
              ctx.reply(`You was successfuly unfollowed from ${channel}.`)
            }
          } catch (e) {
            error(e)
            ctx.reply(e.message)
          }
          ctx.scene.leave()
        }
      ]
    }),
  }

  public bot: Telegraf<any>
  constructor() {
    super()
    this.init()
  }
  protected init(): void {
    this.bot = new Telegraf(config.telegram.token, { telegram: telegramOpts })
    this.bot.launch().then(() => info('Telegram bot started.')).catch(e => error(e))
    this.loadMiddlewares()
    this.registerScenes()
    this.loadCommands()
  }
  public async sendMessage(opts: SendMessageOpts) {
    try {
      const targets = Array.isArray(opts.target) ? opts.target : [opts.target]
      for (const target of targets) {
        this.bot.telegram.sendPhoto(target, opts.image, {
          caption: opts.message
        })
      }
      return true
    } catch (e) {
      error(e)
      return false
    }
  }
  protected async loadMiddlewares(): Promise<void> {
    super.loadMiddlewares()
    this.bot.use(session())
    this.bot.use(async (ctx, next) => {
      if (!ctx.message) return
      info(`Telegram | New message from ${ctx.from.username} [${ctx.from.id}], message: ${ctx.message.text}`)
      const [user] = await User.findOrCreate({ where: { id: ctx.from.id, service }, defaults: { follows: [], service } })
      ctx.userDb = user
      next()
    })
  }
  protected async registerScenes() {
    super.registerScenes()
    const stage = new Stage([this.scenes.follow, this.scenes.unfollow])
    stage.command('cancel', Stage.leave())
    this.bot.use(stage.middleware())
  }
  protected async loadCommands() {
    super.loadCommands()
    this.bot.command(['start', 'help'], ({ reply }) => reply(`Hi! I will notify you about the beginning of the broadcasts on Twitch.`))
    this.bot.command('follow', Stage.enter('follow'))
    this.bot.command('unfollow', Stage.enter('unfollow'))
    this.bot.command('follows', async (ctx) => {
      const channels = await follows({ userId: ctx.from.id, service })
      if (isBoolean(channels)) ctx.reply(`You aren't followed to anyone`)
      else ctx.reply(`You are followed to ${channels.join(', ')}`)
    })
    this.bot.command('live', async (ctx) => {
      const channels = await live({ userId: ctx.from.id, service })
      if (isBoolean(channels)) ctx.reply(`There is no channels currently online`) 
      else ctx.reply(`Currently online: \n${channels.map((o) => 'https://twitch.tv/' + o).join('\n')}`)
    })
  }
}

export default new Telegram()
