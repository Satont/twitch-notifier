import { Chat } from '../entities/Chat'

declare module 'vk-io' {
  interface MessageContext {
    public ChatEntity: Chat
  } 
}