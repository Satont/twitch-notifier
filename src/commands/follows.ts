import { Chat } from '../entities/Chat'
import { Twitch } from '../libs/twitch'

export async function followsCommand({ chat }: { chat: Chat }) {
  const streamers = await Twitch.getUsers({ ids: chat.follows.map(f => f.channel.id) })
  const names = streamers.map(s => `https://twitch.tv/${s.name}`)

  return `You are followed to:\n${names.join('\n')}`
}
