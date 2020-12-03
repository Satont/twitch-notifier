import { Entity, PrimaryColumn, Column, BaseEntity, OneToMany, CreateDateColumn, UpdateDateColumn } from 'typeorm'
import { Follow } from './Follow'

@Entity('channels')
export class Channel extends BaseEntity {
  @PrimaryColumn()
  id!: string

  @Column()
  username!: string

  @Column()
  online: boolean = false

  @Column({ nullable: true })
  category?: string

  @Column({ nullable: true })
  title?: string

  @CreateDateColumn()
  createdAt!: Date

  @UpdateDateColumn()
  updatedAt!: Date

  @OneToMany(() => Follow, category => category.channel)
  followers: Follow[]
}
