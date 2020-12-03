import { ServiceInterface } from '../services/_interface'

export function vkAction(name: string): MethodDecorator {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  return (service: ServiceInterface, methodName: string): void => {
    import('../services/vk').then((v) => {
      v.default.bot.updates.on('message_new', async (ctx, next) => {
        if (ctx.messagePayload?.command === name) await v.default[methodName](ctx)
        next()
      })
    })
  }
}
