import { getConnection } from 'typeorm'
import { Chat } from '../entities/Chat'
import { Follow } from '../entities/Follow'
import { I18n } from '../libs/i18n'

const followRepository = getConnection().getRepository(Follow)

export async function liveCommand({ chat, i18n }: { chat: Chat, i18n: I18n }) {
  const streams = (await followRepository.find({
    where: { chat },
    relations: ['channel'],
  })).filter(f => f.channel.online).map(f => f.channel)

  if (!streams.length) {
    return i18n.translate('commands.live.empty')
  } else {
    const names = streams.map(s => `https://twitch.tv/${s.username} | Title: ${s.title} | Category: ${s.category}`)
    return i18n.translate('commands.live.list', { list: names.join('\n') })
  }
}
