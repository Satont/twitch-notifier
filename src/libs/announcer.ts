import { getRepository, getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'
import { Follow } from '../entities/Follow'
import { Stream } from '../entities/Stream'
import { SendMessageOpts, services } from '../services/_interface'
import Twitch from './twitch'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import { info } from 'console'
import { HelixStream } from '@twurple/api'
import { getLatestStream } from '../helpers/getLatestStream'

dayjs.extend(relativeTime)

export class Announcer {
  private readonly channelsRepository = getConnection().getRepository(Channel)
  private readonly streamsRepository = getConnection().getRepository(Stream)
  channelId: string = null
  private channel: Channel = null

  constructor(channelId: string) {
    this.channelId = channelId
  }

  async init() {
    this.channel = await this.channelsRepository.findOne(this.channelId)
  }

  private getChannelFollowers(channelId: string) {
    return getRepository(Follow).find({
      where: {
        channel: { id: channelId },
      },
      relations: ['chat'],
    })
  }

  async announceLive(data: { displayName: string, userId: string, stream?: HelixStream }) {
    const latestStream = await getLatestStream(data.userId)

    const stream = data.stream ?? await Twitch.apiClient.streams.getStreamByUserId(data.userId)
    if (stream.type !== 'live') return

    if (stream?.id !== latestStream?.id) {
      this.announce({
        message: `
          ${data.displayName} online!
          Category: ${stream.gameName}
          Title: ${stream.title}
          https://twitch.tv/${stream.userName}
        `,
        target: (await this.getChannelFollowers(this.channel.id)).map(f => f.chat.chatId),
        image: this.getThumnailUrl(stream.thumbnailUrl),
      })
    } else {
      info(`EventSub: Stream ${stream?.id} of ${data.displayName}[${data.userId}] already in database, skipping announce.`)
    }
    
    this.channel.username = stream.userName
    this.channel.online = true
    await this.streamsRepository.create({ 
      id: stream.id, 
      startedAt: stream.startDate, 
      channel: this.channel,
      title: stream.title,
      category: stream.gameName,
    }).save()
    this.channel.save()
  }

  async announceOffline(data: { displayName: string, userId: string }) {
    const latestStream = await getLatestStream(data.userId)
    const streamDuration = dayjs().from(dayjs(latestStream?.startedAt), true)

    this.announce({
      message: `
        ${data.displayName} now offline
        ${latestStream ? 'Stream duration was: ' + streamDuration : ''}
      `,
      target: (await this.getChannelFollowers(this.channel.id)).filter(f => f.chat.settings.offline_notification).map(f => f.chat.chatId),
    })

    this.channel.online = false
    this.channel.save()
  }

  async announceUpdate(data: { displayName: string, userId: string, newCategory: string }) {
    const stream = await Twitch.apiClient.streams.getStreamByUserId(data.userId)
    if (stream?.type !== 'live') return

    const latestStream = await getLatestStream(data.userId)
    if (this.channel.online && latestStream?.category !== data.newCategory) {
      this.announce({
        message: `
          ${data.displayName} now streaming ${data.newCategory}
          Previous category: ${latestStream?.category}
          https://twitch.tv/${stream.userName}
        `,
        target: (await this.getChannelFollowers(this.channel.id)).filter(f => f.chat.settings.game_change_notification).map(f => f.chat.chatId),
        image: this.getThumnailUrl(stream.thumbnailUrl),
      })
    }

    if (latestStream) {
      latestStream.updatedAt = new Date()
      latestStream.category = data.newCategory
      await this.streamsRepository.save(latestStream)
    }

    this.channel.save()
  }

  private announce(opts: SendMessageOpts) {
    for (const service of services) {
      service.makeAnnounce({
        ...opts,
        message: opts.message.replace(/  +/g, ''),
      })
    }
  }

  private getThumnailUrl(url: string) {
    return `${url.replace('{width}', '1920').replace('{height}', '1080')}?timestamp=${Date.now()}`
  }
}
