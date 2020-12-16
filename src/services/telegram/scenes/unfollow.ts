import { BaseScene } from 'telegraf'
import telegram from '..'
import { unFollowCommand } from '../../../commands/unfollow'

export const unFollowScene = new BaseScene('unFollowScene')

unFollowScene.enter((ctx) => ctx.reply(ctx.i18n.translate('scenes.unfollow.enter')))
unFollowScene.on('message', async (ctx) => {
  const result = await unFollowCommand({ chat: ctx.ChatEntity, channelName: ctx.message.text, i18n: ctx.i18n })

  await telegram.sendMessage({
    target: String(ctx.chat.id),
    message: result.message,
  })

  if (result.success) {
    ctx.scene.leave()
  }
})
