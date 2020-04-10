import "reflect-metadata"
import Sentry from '@sentry/node'

if (process.env.SENTRY_DSN && process.env.SENTRY_DSN !== '') Sentry.init({ dsn: process.env.SENTRY_DSN })

require('dotenv').config()


import { connected } from './libs/db'
import { info } from './helpers/logs'

function init () {
  if (!connected) return setTimeout(() => init(), 500)
  require('./messengers/vk')
  require('./messengers/telegram')
  require('./libs/checker')
  require('./libs/http')
  info('Application works now.')
}
init()