import { Channel } from './Channel'
import { Chat } from './Chat'
import { Entity, PrimaryGeneratedColumn, BaseEntity, ManyToOne, CreateDateColumn, Unique } from 'typeorm'

@Entity('follows')
@Unique(['chat', 'channel'])
export class Follow extends BaseEntity {
  @PrimaryGeneratedColumn()
  id!: number

  @CreateDateColumn()
  createdAt!: Date

  @ManyToOne(() => Chat, category => category.follows, { nullable: false })
  chat!: Chat

  @ManyToOne(() => Channel, category => category.followers, { nullable: false })
  channel!: Channel
}
