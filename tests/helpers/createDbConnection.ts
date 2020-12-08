import { resolve } from 'path'
import { createConnection } from 'typeorm'

export default async () => {
  await createConnection({
    type: 'better-sqlite3',
    database: ':memory:',
    entities: [resolve(process.cwd(), 'src', 'entities', '*{.js,.ts}')],
    synchronize: true,
  })
}
