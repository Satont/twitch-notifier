import axios from 'axios'
import { error } from '../helpers/logs'


export enum Methods {
  GET = 'GET',
  POST = 'POST',
  PUT = 'PUT',
  DELETE = 'DELETE'
}

interface TwitchError {
  status: string,
  message: string,
}

interface Request {
  method: Methods,
  endpoint: string,
  data?: object,
}

export class Twitch {
  clientId: string = ''
  baseUrl: string = 'https://api.twitch.tv/helix/'

  constructor(clientId: string) {
    this.clientId = clientId
  }

  public async request (options: Request) {
    let query: any;
    const url: string = `${this.baseUrl}${options.endpoint}`
    try {
      query = await axios({
        method: options.method,
        url,
        data: options.data || {},
        headers: {
          'Client-ID': this.clientId
        }
      })
      return query.data
    } catch (e) {
      const twitchError: boolean = Boolean(e.response.data)
      const twitchData: TwitchError = e.response.data

      const errorMessage: string = twitchError ? `${twitchData.status} — ${twitchData.message}` : `${e.response.status} — ${e.response.statusText}`
      error(`Ошибка при запросе ${options.method} ${url}, тело ошибки: ${errorMessage}`)
      throw new Error('Произошла ошибка при запросе к twitch.')
    }
  }

  public async getChannel(channelName: string): Promise<{id: number, login: string, displayName: string}> {
    try {
      const request = await this.request({ method: Methods.GET, endpoint: `users?login=${channelName}` })
      const response = request.data[0]
      if (!request.data.length) throw new Error(`Канал ${channelName} не найден.`)
      else return { id: Number(response.id), login: response.login, displayName: response.display_name }
    } catch (e) {
      error(e.message)
      throw new Error(e.message)
    }
  }

  public async getChannelsById(channels: number[]): Promise<[{ id: number, displayName: string, login: string }]> {
    try {
      const request = await this.request({ method: Methods.GET, endpoint: `users?id=${channels.join('&id=')}` })
      return request.data.map(o => { return { id: Number(o.id), displayName: o.display_name, login: o.login } })
    } catch (e) {
      error(e.message)
      throw new Error(e.message)
    }
  }

  public async checkOnline (channels: number[]) {
    try {
      const request = await this.request({ method: Methods.GET, endpoint: `streams?first=100&user_id=${channels.join('&user_id')}` })
      return request.data
    } catch (e) {
      throw new Error(e.message)
    }
  }
}