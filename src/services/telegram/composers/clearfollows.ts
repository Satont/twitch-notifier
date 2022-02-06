import { Composer } from 'grammy'
import { unFollowAllCommand } from '../../../commands/unFollowAll'
import { TelegramMessageSender } from '../MessageSender'
import { Context } from '../types'

export const composer = new Composer<Context>()

composer.command('clearfollows', async (ctx) => {
  const { message } = await unFollowAllCommand({ chat: ctx.session.entity, i18n: ctx.session.i18n })

  TelegramMessageSender.sendMessage({
    target: String(ctx.from.id),
    message,
  })
})