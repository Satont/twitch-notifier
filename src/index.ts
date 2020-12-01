import 'reflect-metadata'
import 'source-map-support/register'
import { createConnection, getConnection } from 'typeorm'
import { error } from './libs/logger'

const start = async () => {
  await createConnection()
  if (!getConnection().isConnected) {
    return setTimeout(() => start(), 1000)
  }
  import('./libs/loader')
}
start()

process.on('unhandledRejection', (reason) => {
  error(reason)
})
process.on('uncaughtException', (err: Error) => {
  const date = new Date().toISOString()

  process.report?.writeReport(`uncaughtException-${date}`, err)
  error(err)

  process.exit(1)
})
