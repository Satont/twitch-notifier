import { Follow } from './Follow'
import { Entity, Column, BaseEntity, OneToMany, CreateDateColumn, UpdateDateColumn, Unique, PrimaryGeneratedColumn, OneToOne, JoinColumn } from 'typeorm'
import { ChatSettings } from './ChatSettings'
import { IgnoredCategory } from './IgnoredCategory'

export enum Services {
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

  @OneToMany(() => Follow, follow => follow.chat, { onDelete: 'CASCADE' })
  follows: Follow[]

  @OneToOne(() => ChatSettings, settings => settings.chat, { cascade: true, eager: true, onDelete: 'CASCADE' })
  @JoinColumn()
  settings: ChatSettings

  @OneToMany(() => IgnoredCategory, category => category.chat, { onDelete: 'CASCADE' })
  ingoredCategories: IgnoredCategory[]
}
