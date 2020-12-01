import { ApiClient, ClientCredentialsAuthProvider, HelixUser, HelixStream } from 'twitch'
import { chunk } from 'lodash'
import { info, warning } from './logger'
import { getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'

class TwitchClient {
  apiClient: ApiClient = null

  async init() {
    const [client_id, client_secret] = [process.env.TWITCH_CLIENT_ID, process.env.TWITCH_CLIENT_SECRET]
    if (!client_id || !client_secret) {
      warning('TWITCH: client_id or client_secret not setuped, twitch library will not works.')
      return
    }
    const authProvider = new ClientCredentialsAuthProvider(client_id, client_secret)
    this.apiClient = new ApiClient({ authProvider })
    await this.initWebhooks()

    info('Twitch library initialized.')
  }

  async getUsers({ ids, names }: { ids?: string[], names?: string[] }) {
    const chunks = chunk(ids || names, 100)
    const result: HelixUser[] = []

    for (const chunk of chunks) {
      const data: HelixUser[] = await this.apiClient.helix.users[names ? 'getUsersByNames' : 'getUsersByIds'](chunk)
      result.push(...data)
    }

    return result
  }

  async getUser({ name, id }: { name?: string, id?: string }) {
    const query = {
      [name ? 'names' : 'ids']: name ? [name] : [id],
    }

    const [user] = await this.getUsers(query)
    return user
  }

  async getStreams(userIds: string[]) {
    const chunks = chunk(userIds, 100)
    const result: HelixStream[] = []

    for (const chunk of chunks) {
      const streams = await this.apiClient.helix.streams.getStreamsPaginated({ userId: chunk }).getAll()
      result.push(...streams)
    }

    return result
  }

  async getStream(userId: string) {
    const [stream] = await this.getStreams([userId])
    return stream
  }

  private async initWebhooks() {
    const siteUrl = process.env.SITE_URL
    if (!siteUrl) {
      warning('TWITCH: siteUrl not setuped, streams udpates will not be recieved.')
      return
    }

    const options = {
      callbackUrl: `${siteUrl}/twitch/webhooks/callback`,
      validityInSeconds: 864000,
    }

    const channelsRepository = getConnection().getRepository(Channel)
    const channels = await channelsRepository.find()

    for (const channel of channels) {
      await this.apiClient.helix.webHooks.unsubscribeFromStreamChanges(channel.id, options)
      await this.apiClient.helix.webHooks.subscribeToStreamChanges(channel.id, options)
    }

    info(`TWITCH: webhook subscribed to ${channels.length} channels`)

    setTimeout((() => this.initWebhooks()), options.validityInSeconds * 1000)
  }
}

export const Twitch = new TwitchClient()
export default Twitch
