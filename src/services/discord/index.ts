import { Services } from '../../entities/Chat'
import { ServiceInterface } from '../_interface'
import { ShardingManager } from 'discord.js'
import { info, warning } from '../../libs/logger'
import { resolve } from 'path'

export class DiscordService extends ServiceInterface {
  manager: ShardingManager

  constructor() {
    super({ service: Services.DISCORD_SERVER })

    const token = process.env.DISCORD_BOT_TOKEN
    if (!token) {
      warning('DISCORD: bot token not setuped, discord service will not works.')
      return
    }

    this.manager = new ShardingManager(resolve(__dirname, 'bot.js'), { token })
  }

  async init() {
    if (!this.manager) {
      return
    }

    this.initListeners()
    await this.manager.spawn()
  }

  async initListeners() {
    this.manager.on('shardCreate', shard => info(`DISCORD: Launched shard ${shard.id}`))
  }
}

export default new DiscordService()