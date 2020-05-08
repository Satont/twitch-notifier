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
      const request = await this.kraken.get(`users?login=${channelName}`)
      const response = request.data
      if (!response.users.length) throw new Error(`Channel ${channelName} not found.`)
      else return { id: Number(response.users[0]._id), login: response.users[0].name, displayName: response.users[0].display_name }
    } catch (e) {
      error(e.message)
      throw new Error(e.message)
    }
  }

  public async getChannelsById(channels: number[]): Promise<[{ id: number, displayName: string, login: string }]> {
    try {
      const request = await this.kraken.get(`users?id=${channels.join('&id=')}`)
      return request.data.users.map(o => { return { id: Number(o._id), displayName: o.display_name, login: o.name } })
    } catch (e) {
      error(e.message)
      throw new Error(e.message)
    }
  }

  public async checkOnline (channels: number[]) {
    try {
      const request = await this.kraken.get(`streams?first=100&channel=${channels.join('&channel=')}`)
      return request.data.streams
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
