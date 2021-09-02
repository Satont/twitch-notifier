import { Follow } from './Follow'
import { Entity, Column, BaseEntity, OneToMany, CreateDateColumn, UpdateDateColumn, Unique, PrimaryGeneratedColumn, OneToOne, JoinColumn } from 'typeorm'
import { ChatSettings } from './ChatSettings'

export type DiscordType = 'discord_user' | 'discord_server'


export enum Services {
  VK = 'vk',
  TELEGRAM = 'tg',
  DISCORD_SERVER = 'discord_server'
}

@Entity('chats')
@Unique(['chatId', 'service'])
export class Chat extends BaseEntity {
  @PrimaryGeneratedColumn()
  id: string

  @Column()
  chatId: string

  @Column()
  service: Services

  @CreateDateColumn()
  createdAt!: Date

  @UpdateDateColumn()
  updatedAt!: Date

  @OneToMany(() => Follow, category => category.chat)
  follows: Follow[]

  @OneToOne(() => ChatSettings, settings => settings.chat, { cascade: true, eager: true })
  @JoinColumn()
  settings: ChatSettings
}
