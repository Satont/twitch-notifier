import { Menu } from '@grammyjs/menu'
import { Composer } from 'grammy'
import i18n from '../../../libs/i18n'
import { Context } from '../types'

export const composer = new Composer<Context>()

const lanugageMenu = new Menu<Context>('language', { 
  onMenuOutdated: async (ctx) => {
    await ctx.answerCallbackQuery()
    await ctx.reply(ctx.session.i18n.translate('bot.description'), { reply_markup: lanugageMenu })
  },
})

lanugageMenu.dynamic((_ctx, range) => {
  for (const key of Object.keys(i18n.translations)) {
    const name = i18n.translations[key].language.name
    const emoji = i18n.translations[key].language.emoji
    range.text(`${emoji} ${name}`, async (ctx) => {
      ctx.session.entity.settings.language = key
      await ctx.session.entity.save()
      ctx.session.i18n = i18n.clone(key)
      ctx.editMessageText(ctx.session.i18n.translate('bot.description'))
    })
    range.row()
  }
})

Object.keys(i18n.translations).forEach(key => {
  const name = i18n.translations[key].language.name
  const emoji = i18n.translations[key].language.emoji
  lanugageMenu.text(`${emoji} ${name}`, async (ctx) => {
    ctx.session.entity.settings.language = key
    await ctx.session.entity.save()
    ctx.session.i18n = i18n.clone(key)
    ctx.editMessageText(ctx.session.i18n.translate('bot.description'))
  })
  lanugageMenu.row()
})

lanugageMenu.back('«')

const menu = new Menu<Context>('settings', { onMenuOutdated: false })
  .text(
    (ctx) => `${ctx.session.entity.settings.game_change_notification ? '☑︎' : '◻︎'} ${ctx.session.i18n.translate('settings.game_change_notification_setting.button')}`,
    (ctx) => {
      ctx.session.entity.settings.game_change_notification = !ctx.session.entity.settings.game_change_notification
      ctx.session.entity.save()
      ctx.menu.update()
    },
  )
  .row()
  .text(
    (ctx) => `${ctx.session.entity.settings.offline_notification ? '☑︎' : '◻︎'} ${ctx.session.i18n.translate('settings.offline_notification.button')}`,
    (ctx) => {
      ctx.session.entity.settings.offline_notification = !ctx.session.entity.settings.offline_notification
      ctx.session.entity.save()
      ctx.menu.update()
    },
  )
  .row()
  .submenu(ctx => ctx.session.i18n.translate('settings.language.button'), 'language')
  .row()
  .url('GitHub', 'https://github.com/Satont/twitch-notifier')
  .row()
  .url('Support', 'https://www.buymeacoffee.com/satont')
  

composer.use(menu)
menu.register(lanugageMenu)

composer.command(['settings', 'start'], async (ctx) => {
  ctx.reply(ctx.session.i18n.translate('bot.description'), { reply_markup: menu })
})