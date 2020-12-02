import { Chat } from '../entities/Chat'

declare module 'telegraf' {
  interface Context {
    public ChatEntity: Chat
    public isAction?: boolean
  }
}
