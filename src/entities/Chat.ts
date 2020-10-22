import { Collection, Entity, OneToMany, PrimaryKey, Property } from '@mikro-orm/core'
import { Follow } from './Follow'

@Entity({
  tableName: 'chats',
})
export class Chat {
  @PrimaryKey()
  id!: number

  @Property()
  chatId!: number

  @Property()
  followGameChange: boolean = false

  @OneToMany(() => Follow, follow => follow.chat)
  follows = new Collection<Follow>(this)
}
