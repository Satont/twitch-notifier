import { Injectable } from '@nestjs/common'
import { HelixUser } from '@twurple/api'
import { getConnection } from 'typeorm'
import { Channel } from '../entities/Channel'
import { Chat } from '../entities/Chat'
import { Follow } from '../entities/Follow'
import Twitch from '../libs/twitch'

const connection = getConnection()

@Injectable()
export class AppService {
  private readonly chatsRepository = connection.getRepository(Chat)
  private readonly channelsRepository = connection.getRepository(Channel)
  private readonly followRepository = connection.getRepository(Follow)

  async counts() {
    const chats = await this.chatsRepository.count()
    const channels = await this.channelsRepository.count()

    return { chats, channels }
  }

  async top(count: number): Promise<Array<HelixUser & { count: number }>> {
    const channels = await this.channelsRepository.createQueryBuilder('channels')
      .innerJoin('follows', 'follows', '"follows"."channelId" = "channels"."id"')
      .addSelect('follows.count', 'count')
      .groupBy('channels.id')
      .orderBy('follows.count', 'DESC')
      .limit(count)
      .execute()

    const twitchChannels = await Twitch.getUsers({ ids: channels.map(c => c.channels_id) })

    return twitchChannels
      .map(channel => ({
        ...(channel as any)._data,
        count: channels.find(c => c.channels_id === channel.id).count,
      }))
      .sort((a, b) => b.count - a.count)
  }

}
