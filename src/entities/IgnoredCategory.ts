import { Chat } from './Chat'
import { Entity, PrimaryGeneratedColumn, BaseEntity, ManyToOne, CreateDateColumn, Unique, Column, JoinColumn } from 'typeorm'

@Entity('ignored_categories')
@Unique(['chat', 'categoryId'])
export class IgnoredCategory extends BaseEntity {
  @PrimaryGeneratedColumn()
  id!: number

  @CreateDateColumn()
  createdAt!: Date

  @ManyToOne(() => Chat, category => category.follows, { nullable: false })
  @JoinColumn({ name: 'chatId' })
  chat!: Chat

  @Column()
  chatId: string

  @Column()
  categoryId: string
}
