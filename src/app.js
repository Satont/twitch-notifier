import "@babel/polyfill"
import "reflect-metadata"

import { connected } from './libs/db'

function init () {
  if (!connected) setTimeout(() => init(), 500)
  require('../dest/libs/vk')
}
init()