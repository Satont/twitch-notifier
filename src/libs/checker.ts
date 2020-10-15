import Twitch from './twitch'
import { Channel } from '../models/Channel'
import { chunk, flattenDeep } from 'lodash'
import { notify as notifyUsers, notifyGameChange } from './sender'

async function check () {
  if (!Twitch.inited) return setTimeout(() => check(), 2 * 1000)

  setTimeout(() => check(), 2 * 60 * 1000)
  const dbChannels = await Channel.findAll()
  const onlineChannels = flattenDeep(await getOnlineStreams(dbChannels.map(o => o.id)))

  for (let dbChannel of dbChannels) {
    const channel = onlineChannels.find(o => Number(o.user_id) === dbChannel.id)
    const metadata = channel ? await Twitch.getStreamMetaData(Number(channel?.user_id)) : null

    if (channel && metadata) {
      if (!dbChannel.online) {
        await dbChannel.update({ online: true, game: metadata?.game })
        notifyUsers(metadata)
      } else {
        await checkGame(channel, { old: dbChannel.game, new: metadata?.game })
        await dbChannel.update({ game: metadata?.game })
      }
    } else {
      await dbChannel.update({ online: false })
    }
  }
}
check()

async function getOnlineStreams(channels: number[]): Promise<Array<{
  user_id: string,
  user_name: string,
  game_id: string
}>> {
  let onlineChannels: any[] = []
  const chunks = chunk(channels, 100)
  for (const chunk of chunks) {
    onlineChannels.push((await Twitch.checkOnline(chunk)))
  }

  return onlineChannels
}

async function checkGame(channel: {
  user_id: string,
  user_name: string,
  game_id: string
}, game: { old: string, new: string }) {

  if (game.old === game.new) return;

  await notifyGameChange(
    { name: channel.user_name, id: Number(channel.user_id) }, game.old, game.new
  )
}
