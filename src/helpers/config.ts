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
  }
}

export const db = require('../../database.js')[process.env.NODE_ENV]
