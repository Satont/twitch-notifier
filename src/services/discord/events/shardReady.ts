import { ArgsOf, Client, Discord, On } from 'discordx'
import { info } from '../../../libs/logger'

@Discord()
export abstract class AppDiscord {
  @On('shardReady')
  async ready([shardId]: ArgsOf<'shardReady'>, client: Client) {
    await client.initApplicationCommands()
    info(`DISCORD: Shard #${shardId} is now ready.`)
  }
}