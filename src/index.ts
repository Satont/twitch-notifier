import 'reflect-metadata'
import 'source-map-support/register'
import { createConnection, getConnection } from 'typeorm'
import { error } from './libs/logger'
import * as Sentry from '@sentry/node'
import loader from './libs/loader'

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

  const web = await import('./web')

  await web.bootstrap()
  await loader()
  await web.getAppLication().listen(process.env.PORT || 3000, '0.0.0.0').then(() => web.listened = true)
}
start().catch(console.error)

function stopListen() {
  import('./web').then(web => web.getAppLication()?.close())
  import('./services/telegram').then(service => service.default.bot.stop())
  process.exit(0)
}

process
  .on('SIGINT', () => stopListen())
  .on('SIGTERM', () => stopListen())
  .on('unhandledRejection', (reason) => {
    error(reason)
  })
  .on('uncaughtException', (err: Error) => {
    const date = new Date().toISOString()

    process.report?.writeReport(`uncaughtException-${date}`, err)
    error(err)

    process.exit(1)
  })
