import { Composer } from 'grammy'
import { TelegramMessageSender } from '../MessageSender'
import { Context } from '../types'
import { Menu } from '@grammyjs/menu'
import { getRepository } from 'typeorm'
import { Follow } from '../../../entities/Follow'
import Twitch from '../../../libs/twitch'
import { unFollowCommand } from '../../../commands/unfollow'

export const composer = new Composer<Context>()

const KEYBOARD_ITEMS_MAX_SIZE = 20

const menu = new Menu<Context>('unfollow', {
  onMenuOutdated: async (ctx) => {
    await ctx.answerCallbackQuery()
    await ctx.reply(ctx.session.i18n.translate('scenes.unfollow.enter'), { reply_markup: menu })
  },
})

menu.dynamic(async (ctx, range) => {
  const skip = KEYBOARD_ITEMS_MAX_SIZE * (ctx.session.menus.unfollow.currentPage - 1)
  const follows = (await getRepository(Follow).find({ 
    where: { 
      chat: { chatId: String(ctx.chat.id) },
    },
    skip,
    relations: ['chat', 'channel'],
  }))
  const users = await Twitch.getUsers({ ids: follows.map(f => f.channel.id) })


  let i = 1
  for (const user of users.sort((a, b) => a.name.localeCompare(b.name))) {
    const text = user.displayName.toLowerCase() === user.name ? user.displayName : `${user.displayName} (${user.name})`

    range.text(text, async (ctx) => {
      const result = await unFollowCommand({ chat: ctx.session.entity, channelName: user.name, i18n: ctx.session.i18n })

      TelegramMessageSender.sendMessage({
        target: String(ctx.chat.id),
        message: result.message,
      })
    })

    if (i > 1 && i % 2 === 0) range.row()

    if (i === KEYBOARD_ITEMS_MAX_SIZE) {
      break
    }
    i++
  }
  range.row()
  if (ctx.session.menus.unfollow.currentPage > 1) {
    range.text('«', (ctx) => {
      ctx.session.menus.unfollow.currentPage = ctx.session.menus.unfollow.currentPage - 1
      ctx.menu.update()
    })
  }
  range.text(`Page ${ctx.session.menus.unfollow.currentPage} / ${ctx.session.menus.unfollow.totalPages}`)
  if (ctx.session.menus.unfollow.currentPage + 1 <= ctx.session.menus.unfollow.totalPages) {
    range.text('»', (ctx) => {
      ctx.session.menus.unfollow.currentPage = ctx.session.menus.unfollow.currentPage + 1
      ctx.menu.update()
    })
  }
})

composer.use(menu)

composer.command('unfollow', async (ctx) => {
  const follows = await getRepository(Follow).count({ 
    where: { chat: { chatId: String(ctx.chat.id) } }, 
    relations: ['chat'],
  })
  ctx.session.menus.unfollow.currentPage = 1
  ctx.session.menus.unfollow.totalPages = Math.ceil(follows / KEYBOARD_ITEMS_MAX_SIZE)

  ctx.reply(ctx.session.i18n.translate('scenes.unfollow.enter'), { 
    reply_markup: menu,
  })
})
