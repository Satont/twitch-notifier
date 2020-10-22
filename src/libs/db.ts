import { MikroORM } from '@mikro-orm/core'

export let orm: MikroORM
export const start = async () => {
  orm = await MikroORM.init()
}
