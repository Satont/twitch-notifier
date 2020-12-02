import { Chat } from '../entities/Chat'
import { i18n } from 'i18next'

declare module 'telegraf' {
  interface Context {
    public ChatEntity: Chat
    public isAction?: boolean
    public i18n: i18n
  }
}
