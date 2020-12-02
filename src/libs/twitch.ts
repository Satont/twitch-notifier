import { ApiClient, ClientCredentialsAuthProvider, HelixUser, HelixStream } from 'twitch'
import { chunk } from 'lodash'
import { info, warning } from './logger'

class TwitchClient {
  authProvider: ClientCredentialsAuthProvider = null
  apiClient: ApiClient = null

  async init() {
    const [client_id, client_secret] = [process.env.TWITCH_CLIENT_ID, process.env.TWITCH_CLIENT_SECRET]
    if (!client_id || !client_secret) {
      warning('TWITCH: client_id or client_secret not setuped, twitch library will not works.')
      return
    }
    this.authProvider = new ClientCredentialsAuthProvider(client_id, client_secret)
    this.apiClient = new ApiClient({ authProvider: this.authProvider })

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
}

export const Twitch = new TwitchClient()
export default Twitch
