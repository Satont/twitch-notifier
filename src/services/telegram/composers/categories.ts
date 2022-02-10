import { Composer } from 'grammy'
import { Context } from '../types'
import { StatelessQuestion } from '@grammyjs/stateless-question'
import { Menu } from '@grammyjs/menu'
import Twitch from '../../../libs/twitch'
import { getRepository } from 'typeorm'
import { IgnoredCategory } from '../../../entities/IgnoredCategory'

const question = new StatelessQuestion<Context>('categoryQuestion', async (ctx) => {
  const categoryName = ctx.message.text
  const category = await Twitch.apiClient.games.getGameByName(categoryName)
  
  if (!category) {
    ctx.reply(ctx.session.i18n.translate('composers.categories.errors.gameNotFound', { category: categoryName }))
    return
  }

  const repository = getRepository(IgnoredCategory)

  if (await repository.findOne({ chatId: ctx.session.entity.id, categoryId: category.id })) {
    ctx.reply(ctx.session.i18n.translate('composers.categories.errors.alreadyIgnored', { category: category.name }))
    return
  }

  await repository.create({ 
    chatId: ctx.session.entity.id, 
    categoryId: category.id,
  }).save()

  ctx.reply(ctx.session.i18n.translate('composers.categories.success.added', { category: category.name }))
})

export const composer = new Composer<Context>()
const menu = new Menu<Context>('categoriesMenu')
  .text(
    (ctx) => ctx.session.i18n.translate('composers.categories.menu.buttons.add'),
    (ctx) => question.replyWithMarkdown(ctx, ctx.session.i18n.translate('composers.categories.questions.add'))
  )
  .row()
  .submenu((ctx) => ctx.session.i18n.translate('composers.categories.menu.buttons.remove'), 'removeMenu')


const KEYBOARD_ITEMS_MAX_SIZE = 20
const removeMenu = new Menu<Context>('removeMenu')

removeMenu.dynamic(async (ctx, range) => {
  const skip = KEYBOARD_ITEMS_MAX_SIZE * (ctx.session.menus.removeCategory.currentPage - 1)
  const repository = getRepository(IgnoredCategory)
  const categories = (await repository.find({ 
    where: { 
      chat: { chatId: String(ctx.chat.id) },
    },
    skip,
    relations: ['chat'],
  }))

  const categoriesList = await Twitch.apiClient.games.getGamesByIds(categories.map(c => c.categoryId))

  let i = 1
  for (const category of categoriesList.sort((a, b) => a.name.localeCompare(b.name))) {
    range.text(category.name, async (ctx) => {
      await repository.delete({ 
        chatId: ctx.session.entity.id,
        categoryId: category.id,
      })

      ctx.reply(ctx.session.i18n.translate('composers.categories.success.deleted'))
      ctx.menu.update()
    })

    if (i > 1 && i % 2 === 0) range.row()

    if (i === KEYBOARD_ITEMS_MAX_SIZE) {
      break
    }
    i++
  }
  range.row()
  if (ctx.session.menus.removeCategory.currentPage > 1) {
    range.text('«', (ctx) => {
      ctx.session.menus.removeCategory.currentPage = ctx.session.menus.removeCategory.currentPage - 1
      ctx.menu.update()
    })
  }
  range.text(`Page ${ctx.session.menus.removeCategory.currentPage} / ${ctx.session.menus.removeCategory.totalPages}`)
  if (ctx.session.menus.removeCategory.currentPage + 1 <= ctx.session.menus.removeCategory.totalPages) {
    range.text('»', (ctx) => {
      ctx.session.menus.removeCategory.currentPage = ctx.session.menus.removeCategory.currentPage + 1
      ctx.menu.update()
    })
  }
})

removeMenu.row()
removeMenu.back((ctx) => ctx.session.i18n.translate('composers.categories.menu.buttons.backToMainMenu'))

menu.register(removeMenu)
composer.use(menu)
composer.use(question.middleware())

composer.command('categories', async (ctx) => {
  const categories = await getRepository(IgnoredCategory).count({ 
    where: { chat: { chatId: String(ctx.chat.id) } }, 
    relations: ['chat'],
  })
  ctx.session.menus.removeCategory.currentPage = 1
  ctx.session.menus.removeCategory.totalPages = Math.ceil(categories / KEYBOARD_ITEMS_MAX_SIZE) || 1

  ctx.reply(ctx.session.i18n.translate('composers.categories.command.description'), { reply_markup: menu })
})