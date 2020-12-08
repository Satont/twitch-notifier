import { Chat } from '../../entities/Chat'
import { I18n } from '../../libs/i18n'

declare module 'telegraf' {
  interface Context {
    ChatEntity: Chat
    isAction?: boolean
    i18n: I18n
  }
}
