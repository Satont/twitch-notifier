import { Follow } from './Follow'
import { Entity, PrimaryColumn, Column, BaseEntity, OneToMany, CreateDateColumn, UpdateDateColumn, Unique } from 'typeorm'

export enum Services {
  VK = 'vk',
  TELEGRAM = 'tg'
}

@Entity('chats')
@Unique(['id', 'service'])
export class Chat extends BaseEntity {
  @PrimaryColumn()
  id: string

  @Column({ enum: Services })
  service: Services

  @CreateDateColumn()
  createdAt!: Date

  @UpdateDateColumn()
  updatedAt!: Date

  @OneToMany(() => Follow, category => category.chat)
  follows: Follow[]
}
