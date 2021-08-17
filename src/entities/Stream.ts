import { Entity, PrimaryColumn, Column, BaseEntity, UpdateDateColumn, ManyToOne } from 'typeorm'
import { Channel } from './Channel'

@Entity('streams')
export class Stream extends BaseEntity {
  @PrimaryColumn()
  id!: string

  @Column()
  startedAt: Date

  @UpdateDateColumn()
  updatedAt: Date

  @ManyToOne(() => Channel, channel => channel.streams)
  channel: Channel
}
