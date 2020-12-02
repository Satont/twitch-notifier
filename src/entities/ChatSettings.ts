import { Entity, BaseEntity, PrimaryGeneratedColumn, OneToOne, Column } from 'typeorm'
import { Chat } from './Chat'

export enum Languages {
  RUSSIAN = 'russian',
  ENGLISH = 'english',
}

@Entity('chats_settings')
export class ChatSettings extends BaseEntity {
  @PrimaryGeneratedColumn()
  id: string

  @OneToOne(() => Chat, chat => chat.settings)
  chat: Chat

  @Column({ default: false, nullable: false })
  game_change_notification: boolean = false

  @Column({ enum: Languages, default: Languages.ENGLISH })
  language: Languages = Languages.ENGLISH
}
