import { Entity, PrimaryKey, Property } from '@mikro-orm/core'


@Entity({
  tableName: 'channels',
})
export class Channel {
  @PrimaryKey()
  id!: number

  @Property()
  username!: string

  @Property()
  online: boolean = false

  @Property()
  game?: string

  @Property()
  createdAt = new Date()

  @Property({ onUpdate: () => new Date() })
  updatedAt = new Date()
}
