import { Composer } from 'grammy'
import { followsCommand } from '../../../commands/follows'
import { TelegramMessageSender } from '../MessageSender'
import { Context } from '../types'

export const composer = new Composer<Context>()

composer.command('follows', async (ctx) => {
  TelegramMessageSender.sendMessage({
    target: ctx.from.id,
    message: await followsCommand({ chat: ctx.session.entity, i18n: ctx.session.i18n }),
  })
})