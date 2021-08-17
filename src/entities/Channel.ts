import { Entity, PrimaryColumn, Column, BaseEntity, OneToMany, CreateDateColumn, UpdateDateColumn } from 'typeorm'
import { Follow } from './Follow'
import { Stream } from './Stream'

@Entity('channels')
export class Channel extends BaseEntity {
  @PrimaryColumn()
  id!: string

  @Column()
  username!: string

  @Column()
  online: boolean = false

  @CreateDateColumn()
  createdAt!: Date

  @UpdateDateColumn()
  updatedAt!: Date

  @OneToMany(() => Follow, category => category.channel)
  followers: Follow[]

  @OneToMany(() => Stream, stream => stream.channel)
  streams: Stream[]
}
