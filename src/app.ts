import "reflect-metadata"

import { connected } from './libs/db'
import { info } from './helpers/logs'

function init () {
  if (!connected) return setTimeout(() => init(), 500)
  require('./libs/vk')
  require('./libs/telegram')
  require('./libs/checker')
  info('Application works now.')
}
init()