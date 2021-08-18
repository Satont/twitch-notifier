import { BaseScene } from 'telegraf'
import telegram from '..'
import { followCommand } from '../../../commands/follow'

export const followScene = new BaseScene('followScene', { ttl: 30 })
  .enter((ctx) => ctx.reply(ctx.i18n.translate('scenes.follow.enter')))
  .on('message', async (ctx) => {
    const result = await followCommand({ chat: ctx.ChatEntity, channelName: ctx.message.text, i18n: ctx.i18n })

    await telegram.sendMessage({
      target: String(ctx.chat.id),
      message: result.message,
    })

    if (result.success) {
      ctx.scene.leave()
    }
  })
