import Telegraf, { Context } from 'telegraf'
import { ServiceInterface } from '../services/_interface'

export function telegramAction(name: string | string[]): MethodDecorator {
  const names = Array.isArray(name) ? name : [name]

  return (_service: ServiceInterface & { bot: Telegraf<any> }, methodName: string): void => {
    import('../services/telegram').then((v) => {
      names.forEach(n => {
        v.default.bot.action(n, async (ctx: Context, next) => {
          ctx.isAction = true
          await v.default[methodName](ctx, next)
          ctx.answerCbQuery()
          next()
        })
      })
    })
  }
}
