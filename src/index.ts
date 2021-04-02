import 'reflect-metadata'
import 'source-map-support/register'
import { createConnection, getConnection } from 'typeorm'
import { error } from './libs/logger'
import * as Sentry from '@sentry/node'
import loader from './libs/loader'
import { getAppLication } from './web'

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

  const { bootstrap: webBootstrap, getAppLication } = await import('./web')
  await webBootstrap()
  await loader()
  await getAppLication().listen(process.env.PORT || 3000, '0.0.0.0')
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
