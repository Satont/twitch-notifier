export const config = { 
  vk: {
    token: process.env.VKTOKEN
  },
  telegram: {
    token: process.env.TGTOKEN
  },
  twitch: {
    clientId: process.env.TWITCH_CLIENTID
  },
  panel: {
    port: Number(process.env.PANEL_PORT) || 3000
  },
  proxy: {
    host: process.env.PROXY_HOST,
    port: Number(process.env.PROXY_PORT),
    username: process.env.PROXY_USERNAME,
    password: process.env.PROXY_PASSWORD
  }
}

export const db = require('../../database.js')[process.env.NODE_ENV]
