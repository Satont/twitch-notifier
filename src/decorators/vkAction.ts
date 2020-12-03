import { ServiceInterface } from '../services/_interface'

export function vkAction(name: string | string[]): MethodDecorator {
  const names = Array.isArray(name) ? name : [name]

  return (_service: ServiceInterface, methodName: string): void => {
    import('../services/vk').then((v) => {
      names.forEach(n => {
        v.default.bot.updates.on('message_new', async (ctx, next) => {
          if (ctx.messagePayload?.command === n) await v.default[methodName](ctx)
          next()
        })
      })
    })
  }
}
