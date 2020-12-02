import { Follow } from './Follow'
import { Entity, Column, BaseEntity, OneToMany, CreateDateColumn, UpdateDateColumn, Unique, PrimaryGeneratedColumn, OneToOne } from 'typeorm'
import { ChatSettings } from './ChatSettings'

export enum Services {
  VK = 'vk',
  TELEGRAM = 'tg'
}

@Entity('chats')
@Unique(['chatId', 'service'])
export class Chat extends BaseEntity {
  @PrimaryGeneratedColumn()
  id: string

  @Column()
  chatId: string

  @Column({ enum: Services })
  service: Services

  @CreateDateColumn()
  createdAt!: Date

  @UpdateDateColumn()
  updatedAt!: Date

  @OneToMany(() => Follow, category => category.chat)
  follows: Follow[]

  @OneToOne(() => ChatSettings, settings => settings.chat)
  settings: ChatSettings
}
