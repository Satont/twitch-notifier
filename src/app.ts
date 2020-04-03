import "reflect-metadata"
require('dotenv').config()

import { connected } from './libs/db'
import { info } from './helpers/logs'

function init () {
  if (!connected) return setTimeout(() => init(), 500)
  require('./libs/vk')
  require('./libs/telegram')
  require('./libs/checker')
  require('./libs/http')
  info('Application works now.')
}
init()