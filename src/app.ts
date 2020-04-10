import 'source-map-support/register'
import 'reflect-metadata'

import { init as SentryInit } from '@sentry/node'

require('dotenv').config()
if (process.env.SENTRY_DSN && process.env.SENTRY_DSN !== '') SentryInit({ dsn: process.env.SENTRY_DSN })


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