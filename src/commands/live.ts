import { getConnection } from 'typeorm'
import { Chat } from '../entities/Chat'
import { Follow } from '../entities/Follow'

const followRepository = getConnection().getRepository(Follow)

export async function liveCommand({ chat }: { chat: Chat }) {
  const streams = (await followRepository.find({
    where: { chat },
    relations: ['channel'],
  })).map(f => f.channel).filter(c => c.online)

  const names = streams.map(s => `https://twitch.tv/${s.username} | Title: ${s.title} | Category: ${s.category}`)
  return `Currently live:\n${names.join('\n')}`
}
