import "@babel/polyfill"
import "reflect-metadata"

import { connected } from './libs/db'
import { info } from './helpers/logs'

function init () {
  if (!connected) return setTimeout(() => init(), 500)
  require('../dest/libs/vk')

  info('Whole app initiated.')
}
init()