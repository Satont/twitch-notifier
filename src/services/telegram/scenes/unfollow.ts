import { BaseScene, Markup } from 'telegraf'
import { SceneContextMessageUpdate } from 'telegraf/typings/stage'
import { getRepository } from 'typeorm'
import telegram from '..'
import { unFollowCommand } from '../../../commands/unfollow'
import { Follow } from '../../../entities/Follow'
import Twitch from '../../../libs/twitch'
import { chunk } from 'lodash'

const KEYBOARD_ITEMS_MAX_SIZE = 30

export interface UnfollowScene extends SceneContextMessageUpdate {
  state: {
    follows: Follow[]
    total: number
    current: number
  }
}

export const unFollowScene = new BaseScene<UnfollowScene>('unfollowScene')
  .enter(async (ctx) => {
    const follows = await getRepository(Follow).find({ where: { chat: { chatId: String(ctx.chat.id) } }, relations: ['channel'] })
    const users = await Twitch.getUsers({ ids: follows.map(f => f.channel.id) })
    const keyboard = Markup.inlineKeyboard(chunk(
      users.map(u => Markup.callbackButton(u.name, u.name)),
      KEYBOARD_ITEMS_MAX_SIZE
    ), 
    { 
      columns: 3,
    })

    ctx.reply(ctx.i18n.translate('scenes.unfollow.enter'), {
      reply_markup: keyboard,
    })
  })
  .on('message', async (ctx) => {
    const result = await unFollowCommand({ chat: ctx.ChatEntity, channelName: ctx.message.text, i18n: ctx.i18n })

    await telegram.sendMessage({
      target: String(ctx.chat.id),
      message: result.message,
    })

    if (result.success) {
      ctx.scene.leave()
    }
  })
  .action(/.+/, (ctx) => {
    return
  })