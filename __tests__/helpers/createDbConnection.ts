import { createConnection } from 'typeorm'

export default async () => {
  return await createConnection({
    name: 'default',
    type: 'better-sqlite3',
    database: ':memory:',
    entities: ['src/entities/*.ts'],
    synchronize: true,
  })
}
