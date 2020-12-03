import Telegraf, { Context } from 'telegraf'
import { ServiceInterface } from '../services/_interface'

export function telegramAction(name: string): MethodDecorator {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  return (service: ServiceInterface & { bot: Telegraf<any> }, methodName: string): void => {
    import('../services/telegram').then((v) => {
      v.default.bot.action(name, async (ctx: Context, next) => {
        ctx.isAction = true
        await v.default[methodName](ctx, next)
        ctx.answerCbQuery()
        next()
      })
    })
  }
}
