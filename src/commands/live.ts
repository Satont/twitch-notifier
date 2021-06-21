import { getConnection } from 'typeorm'
import { Chat } from '../entities/Chat'
import { Follow } from '../entities/Follow'
import { I18n } from '../libs/i18n'
import Twitch from '../libs/twitch'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import { HelixStream } from 'twitch/lib'

dayjs.extend(relativeTime)

const followRepository = getConnection().getRepository(Follow)

export async function liveCommand({ chat, i18n }: { chat: Chat, i18n: I18n }) {
  if (!chat.follows?.length) {
    return i18n.translate('commands.follows.emptyList')
  }

  const channels = (await followRepository.find({
    where: { chat },
    relations: ['channel'],
  })).filter(f => f.channel.online).map(f => f.channel)

  const streams: HelixStream[] = await Twitch.getStreams(channels.map(c => c.id)).catch(() => [])

  if (!streams.length) {
    return i18n.translate('commands.live.empty')
  } else {
    const names = streams.map(s => `https://twitch.tv/${s.userName} | Title: ${s.title} | Category: ${s.gameName} | Uptime: ${dayjs().from(dayjs(s.startDate), true)}`)
    return i18n.translate('commands.live.list', { list: names.join('\n') })
  }
}
