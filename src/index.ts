import 'reflect-metadata'
import 'source-map-support/register'
import { createConnection, getConnection } from 'typeorm'
import { error } from './libs/logger'
import * as Sentry from '@sentry/node'

if (process.env.SENTRY_DSN && process.env.NODE_ENV === 'production') {
  Sentry.init({
    dsn: process.env.SENTRY_DSN,
  })
}

const start = async () => {
  await createConnection()
  if (!getConnection().isConnected) {
    return setTimeout(() => start(), 1000)
  }
  import('./libs/loader')
  import('./web')
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
