import { Intents } from 'discord.js'
import { Client } from 'discordx'
import { resolve } from 'path'

const currentPath = resolve(__dirname)

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const client = new Client({ 
  intents: [Intents.FLAGS.GUILDS],
  classes: [
    `${currentPath}/commands/**/*.js`,
    `${currentPath}/events/**/*.js`,
  ],
  silent: true,
})

client.login(process.env.DISCORD_BOT_TOKEN)