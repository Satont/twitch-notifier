import { join } from 'path'
import { createConnection } from 'typeorm'

export default async () => {
  return await createConnection({
    name: 'default',
    type: 'better-sqlite3',
    database: ':memory:',
    entities: [join(process.cwd(), 'dist', 'entities', '*.js')],
    synchronize: true,
  })
}
