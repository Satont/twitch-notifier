import { orm, start as connect } from './libs/db'
import { error } from './libs/logger'

const start = async () => {
  if (!orm?.isConnected()) {
    await connect()
    return start()
  }
}
start()

process.on('unhandledRejection', (reason, promise) => {
  error(reason)
  error(promise)
})
process.on('uncaughtException', (err: Error) => {
  const date = new Date().toISOString()

  process.report?.writeReport(`uncaughtException-${date}`, err)
  error(err)

  process.exit(1)
})
