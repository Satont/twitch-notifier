import { Composer } from 'grammy'
import { liveCommand } from '../../../commands/live'
import { TelegramMessageSender } from '../MessageSender'
import { Context } from '../types'

export const composer = new Composer<Context>()

composer.command('live', async (ctx) => {
  TelegramMessageSender.sendMessage({ 
    target: ctx.from.id,
    message: await liveCommand({ chat: ctx.session.entity, i18n: ctx.session.i18n }),
  })
})