import { Bot } from 'grammy'
import { Services } from '../../entities/Chat'
import { error, warning } from '../../libs/logger'
import { ServiceInterface } from '../_interface'
import { clearFollowsComposer, followComposer, followsComposer, liveComposer, settingsComposer, unfollowComposer } from './composers'
import { i18nMiddleware, sessionMiddleware, userEntityMiddleware } from './middlewares'
import { Context } from './types'

class TelegramService extends ServiceInterface {
  bot: Bot<Context>
  
  constructor() {
    super({
      service: Services.TELEGRAM,
    })

    const accessToken = process.env.TELEGRAM_BOT_TOKEN
    if (!accessToken) {
      warning('TELEGRAM: bot token not setuped, telegram library will not works.')
      return
    }

    this.bot = new Bot<Context>(accessToken)
    this.bot.use(sessionMiddleware)
    this.bot.use(userEntityMiddleware)
    this.bot.use(i18nMiddleware)
    this.bot.use(followComposer)
    this.bot.use(unfollowComposer)
    this.bot.use(settingsComposer)
    this.bot.use(clearFollowsComposer)
    this.bot.use(liveComposer)
    this.bot.use(followsComposer)

    this.bot.catch(err => {
      error(err)
    })
  }

  async init() {
    await this.bot.init()
    await this.bot.api.setMyCommands([
      { command: 'follow', description: 'Follow to some user.' },
      { command: 'unfollow', description: 'Unfollow from some users.' },
      { command: 'follows', description: 'Shows list of your follows.' },
      { command: 'live', description: 'Shows list of your follows.' },
      { command: 'start', description: 'Start command' },
      { command: 'settings', description: 'Settings menu.' },
      { command: 'clearfollows', description: 'Unfollow from all users.' },
    ])
    this.bot.start()
  }
}

export default new TelegramService()
