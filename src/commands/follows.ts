import { Chat } from '../entities/Chat'
import { I18n } from '../libs/i18n'
import { Twitch } from '../libs/twitch'

export async function followsCommand({ chat, i18n }: { chat: Chat, i18n: I18n }) {
  if (!chat.follows?.length) {
    return i18n.translate('commands.follows.emptyList')
  }

  const streamers = await Twitch.getUsers({ ids: chat.follows.map(f => f.channel.id) })
  const names = streamers.map(s => `https://twitch.tv/${s.name}`)

  return i18n.translate('commands.follows.list', { list: names.join('\n') })
}
