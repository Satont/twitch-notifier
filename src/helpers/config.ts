const config = {
  db: {
    host: process.env.DB_HOST,
    port: Number(process.env.DB_PORT),
    name: process.env.DB_NAME,
    user: process.env.DB_USER,
    password: process.env.DB_PASSWORD,
  },
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
    port: process.env.PANEL_PORT || 3000
  }
}

export default { config }
export { config }