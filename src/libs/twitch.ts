import axios, { AxiosInstance } from 'axios'
import { error } from '../helpers/logs'

export class Twitch {
  helix: AxiosInstance
  kraken: AxiosInstance

  constructor(clientId: string) {
    this.helix = axios.create({
      baseURL: 'https://api.twitch.tv/helix/',
      headers: {
        'Client-ID': clientId,
      },
    })
    this.kraken = axios.create({
      baseURL: 'https://api.twitch.tv/kraken/',
      headers: {
        'Accept': 'application/vnd.twitchtv.v5+json',
        'Client-ID': clientId,
      },
    })
  }

  public async getChannel(channelName: string): Promise<{id: number, login: string, displayName: string}> {
    try {
      const request = await this.helix.get(`users?login=${channelName}`)
      const response = request.data.data[0]
      if (!request.data.data.length) throw new Error(`Channel ${channelName} not found.`)
      else return { id: Number(response.id), login: response.login, displayName: response.display_name }
    } catch (e) {
      error(e.message)
      throw new Error(e.message)
    }
  }

  public async getChannelsById(channels: number[]): Promise<[{ id: number, displayName: string, login: string }]> {
    try {
      const request = await this.helix.get(`users?id=${channels.join('&id=')}`)
      return request.data.data.map(o => { return { id: Number(o.id), displayName: o.display_name, login: o.login } })
    } catch (e) {
      error(e.message)
      throw new Error(e.message)
    }
  }

  public async checkOnline (channels: number[]) {
    try {
      const request = await this.helix.get(`streams?first=100&user_id=${channels.join('&user_id=')}`)
      return request.data.data
    } catch (e) {
      throw new Error(e.message)
    }
  }

  public async getStreamMetaData(id: number): Promise<StreamMetadata> {
    try {
      const { data } = await this.kraken.get(`streams/${id}`)
      return data.stream
    } catch (e) {
      error(e)
      throw new Error(e)
    }
  }
}

export type StreamMetadata = {
  game: null | string,
  channel: {
    display_name: string,
    name: string,
    status: string,
    _id: number,
  }
  preview: {
    template: string,
  }
}
