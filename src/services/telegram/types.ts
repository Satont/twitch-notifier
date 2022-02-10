import { Context as GrammyContext, SessionFlavor } from 'grammy'
import { Chat } from '../../entities/Chat'
import { I18n } from '../../libs/i18n'

interface SessionData {
  entity: Chat;
  i18n: I18n;
  menus: {
    unfollow: {
      currentPage: number,
      totalPages: number
    },
    removeCategory: {
      currentPage: number,
      totalPages: number
    },
    categories: {
      ignoreId: string,
    },
  }
}

// Flavor the context type to include sessions.
export type Context = GrammyContext & SessionFlavor<SessionData>;