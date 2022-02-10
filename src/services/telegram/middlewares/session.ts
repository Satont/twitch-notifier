import { session } from 'grammy'

export const middleware = session({
  getSessionKey: (ctx) => (ctx.from?.id ?? ctx.chat?.id).toString(),
  initial: () => ({
    entity: null,
    i18n: null,
    menus: {
      unfollow: {
        currentPage: 0,
        totalPages: 0,
      },
      removeCategory: {
        currentPage: 0,
        totalPages: 0,
      },
      categories: {
        ignoreId: '',
      },
    },
  }),
})