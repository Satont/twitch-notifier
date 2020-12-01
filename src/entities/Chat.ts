import { Follow } from './Follow'
import { Entity, PrimaryColumn, Column, BaseEntity, OneToMany, CreateDateColumn, UpdateDateColumn } from 'typeorm'

@Entity('chats')
export class Chat extends BaseEntity {
  @PrimaryColumn()
  id: string;

  @Column()
  followGameChange: boolean = false

  @CreateDateColumn()
  createdAt!: Date

  @UpdateDateColumn()
  updatedAt!: Date

  @OneToMany(() => Follow, category => category.chat)
  follows: Follow[]
}
