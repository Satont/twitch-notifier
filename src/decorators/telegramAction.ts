import Telegraf, { Context } from 'telegraf'
import { i18n } from '../libs/i18n'
import { ServiceInterface } from '../services/_interface'

export function telegramAction(name: string): MethodDecorator {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  return (service: ServiceInterface & { bot: Telegraf<any> }, methodName: string): void => {
    import('../services/telegram').then((v) => {
      v.default.bot.action(name, async (ctx: Context, next) => {
        ctx.isAction = true
        ctx = await v.default.ensureUser(ctx)
        ctx.i18n = i18n.cloneInstance({ lng: ctx.ChatEntity.settings.language })
        await v.default[methodName](ctx, next)
        next()
      })
    })
  }
}
