import { BaseScene, Markup } from 'telegraf'
import { SceneContextMessageUpdate, SceneContext } from 'telegraf/typings/stage'
import { getRepository } from 'typeorm'
import telegram from '..'
import { unFollowCommand } from '../../../commands/unfollow'
import { Follow } from '../../../entities/Follow'
import Twitch from '../../../libs/twitch'
import { chunk } from 'lodash'

const KEYBOARD_ITEMS_MAX_SIZE = 20

export interface UnfollowScene extends SceneContextMessageUpdate {
  scene: SceneContext<this> & {
    state: {
      follows: Follow[]
      totalItems: number
      currentPage: number
      totalPages: number
    }
  }
}

const getKeyboard = (ctx: UnfollowScene) => {
  return [
    ctx.scene.state.currentPage > 1 ? Markup.callbackButton('«', 'PREV_PAGE') : undefined,
    Markup.callbackButton(`Page ${ctx.scene.state.currentPage} / ${ctx.scene.state.totalPages}`, 'void'),
    ctx.scene.state.currentPage + 1 <= ctx.scene.state.totalPages ? Markup.callbackButton('»', 'NEXT_PAGE') : undefined,
  ].filter(Boolean)
}

const changePage = async (ctx: UnfollowScene) => {
  const startIndex = KEYBOARD_ITEMS_MAX_SIZE * (ctx.scene.state.currentPage - 1)
  const neededFollows = ctx.scene.state.follows.slice(startIndex, startIndex + 20)
  const users = await Twitch.getUsers({ ids: neededFollows.map(f => f.channel.id) })
  const keyboard = Markup.inlineKeyboard([
    ...chunk(users.map(u => Markup.callbackButton(u.name, `channel-${u.name}`)), 2),
    [...getKeyboard(ctx)],
  ], 
  { 
    columns: 3,
  })

  return keyboard
}

export const unFollowScene = new BaseScene<UnfollowScene>('unfollowScene')
  .enter(async (ctx) => {
    const follows = (await getRepository(Follow).find({ where: { chat: { chatId: String(ctx.chat.id) } }, relations: ['chat', 'channel'] }))
      .sort((a, b) => Number(a.channel.id) + Number(b.channel.id))

    ctx.scene.state = {
      follows,
      totalItems: follows.length,
      currentPage: 1,
      totalPages: Math.ceil(follows.length / KEYBOARD_ITEMS_MAX_SIZE),
    }

    ctx.reply(ctx.i18n.translate('scenes.unfollow.enter'), {
      reply_markup: await changePage(ctx),
    })
  })
  .action('NEXT_PAGE', async (ctx) => {
    await ctx.answerCbQuery()
    ctx.scene.state.currentPage = ctx.scene.state.currentPage + 1
    await ctx.editMessageReplyMarkup(await changePage(ctx))
  })
  .action('PREV_PAGE', async (ctx) => {
    await ctx.answerCbQuery()
    ctx.scene.state.currentPage = ctx.scene.state.currentPage - 1
    await ctx.editMessageReplyMarkup(await changePage(ctx))
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
  .action(/channel-.+/, async (ctx) => {
    const channelName = ctx.match.input.replace('channel-', '')
    const { message } = await unFollowCommand({ chat: ctx.ChatEntity, channelName, i18n: ctx.i18n })

    await telegram.sendMessage({
      target: String(ctx.chat.id),
      message,
    })

    await ctx.answerCbQuery()
  })
  .action('void', (ctx) => ctx.answerCbQuery())