import { Entity, ManyToOne, PrimaryKey } from '@mikro-orm/core'
import { Channel } from './Channel'
import { Chat } from './Chat'

@Entity({
  tableName: 'chats',
})
export class Follow {
  @PrimaryKey()
  id!: number

  @ManyToOne({ fieldName: 'chatId' })
  chat!: Chat

  @ManyToOne({ fieldName: 'channelId' })
  channel!: Channel
}
