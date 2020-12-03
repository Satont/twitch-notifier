import { Chat } from '../entities/Chat'
import { I18n } from '../libs/i18n'

declare module 'telegraf' {
  interface Context {
    public ChatEntity: Chat
    public isAction?: boolean
    public i18n: I18n
  }
}
