import { EventSubStreamOfflineEvent, EventSubStreamOnlineEvent, EventSubChannelUpdateEvent } from '@twurple/eventsub'
import { getRepository, getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'
import { Follow } from '../entities/Follow'
import { Stream } from '../entities/Stream'
import { SendMessageOpts, services } from '../services/_interface'
import Twitch from './twitch'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import { info } from 'console'
import { TelegramMessageSender } from '../services/telegram/MessageSender'

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

  private getLatestStream(channelId: string) {
    return this.streamsRepository.findOne({
      where: {
        channel: {
          id: channelId,
        },
      },
      order: {
        startedAt: 'DESC',
      },
    })
  }

  async announceLive(event: EventSubStreamOnlineEvent) {
    if (event.streamType !== 'live') return
    const latestStream = await this.getLatestStream(event.broadcasterId)

    const stream = await Twitch.apiClient.streams.getStreamByUserId(event.broadcasterId)

    if (stream?.id !== latestStream?.id) {
      this.announce({
        message: `
          ${event.broadcasterDisplayName} online!
          Category: ${stream.gameName}
          Title: ${stream.title}
          https://twitch.tv/${event.broadcasterName}
        `,
        target: (await this.getChannelFollowers(this.channel.id)).map(f => f.chat.chatId),
        image: this.getThumnailUrl(stream.thumbnailUrl),
      })
      info(`EventSub: Sended notification about ${event.broadcasterName}[${event.broadcasterId}].`)
    } else {
      info(`EventSub: Stream ${stream?.id} of ${event.broadcasterName}[${event.broadcasterId}] already in database, skipping announce.`)
    }
    
    this.channel.username = event.broadcasterName
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

  async announceOffline(event: EventSubStreamOfflineEvent) {
    const latestStream = await this.getLatestStream(event.broadcasterId)
    const streamDuration = dayjs().from(dayjs(latestStream?.startedAt), true)

    this.announce({
      message: `
        ${event.broadcasterDisplayName} now offline
        ${latestStream ? 'Stream duration was: ' + streamDuration : ''}
      `,
      target: (await this.getChannelFollowers(this.channel.id)).filter(f => f.chat.settings.offline_notification).map(f => f.chat.chatId),
    })

    this.channel.online = false
    this.channel.save()
  }

  async announceUpdate(event: EventSubChannelUpdateEvent) {
    const stream = await Twitch.apiClient.helix.streams.getStreamByUserId(event.broadcasterId)
    if (stream?.type !== 'live') return

    const latestStream = await this.getLatestStream(event.broadcasterId)
    if (this.channel.online && latestStream?.category !== event.categoryName) {
      this.announce({
        message: `
          ${event.broadcasterDisplayName} now streaming ${event.categoryName}
          Previous category: ${latestStream?.category}
          https://twitch.tv/${event.broadcasterName}
        `,
        target: (await this.getChannelFollowers(this.channel.id)).filter(f => f.chat.settings.game_change_notification).map(f => f.chat.chatId),
        image: this.getThumnailUrl(stream.thumbnailUrl),
      })
    }

    if (latestStream) {
      latestStream.updatedAt = new Date()
      latestStream.category = event.categoryName
      await this.streamsRepository.save(latestStream)
    }

    this.channel.save()
  }

  private announce(opts: SendMessageOpts) {
    TelegramMessageSender.sendMessage({
      ...opts,
      message: opts.message.replace(/  +/g, ''),
    })
  }

  private getThumnailUrl(url: string) {
    return `${url.replace('{width}', '1920').replace('{height}', '1080')}?timestamp=${Date.now()}`
  }
}
