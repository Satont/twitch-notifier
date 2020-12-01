import { Chat } from '../entities/Chat'
import { Twitch } from '../libs/twitch'

export async function followsCommand({ chat }: { chat: Chat }) {
  const streams: Array<{ channel: string, category: string, title: string }> = []
  for (const stream of await Twitch.getStreams(chat.follows.map(f => f.channel.id))) {
    const category = (await stream.getGame())?.name
    streams.push({ channel: (await stream.getUser()).name, category, title: stream.title })
  }

  const names = streams.map(s => `https://twitch.tv/${s.channel} | Title: ${s.title} | Category: ${s.category}`)
  return `Currently live:\n${names.join('\n')}`
}
