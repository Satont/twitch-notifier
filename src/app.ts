import "reflect-metadata"

import { connected } from './libs/db'
import { info } from './helpers/logs'

function init () {
  if (!connected) return setTimeout(() => init(), 500)
  require('./libs/vk')
  require('./libs/checker')
  info('Whole app initiated.')
}
init()