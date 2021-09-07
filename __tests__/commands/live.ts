import createDbConnection from '../helpers/createDbConnection'
import '../mocks/twitchLib'
import '../mocks/twitchWatcher'

beforeAll(async () => {
  await createDbConnection()
})

import { getConnection } from 'typeorm'
import { Chat, Services } from '../../src/entities/Chat'
import { I18n } from '../../src/libs/i18n'
import { Channel } from '../../src/entities/Channel'
import { liveCommand } from '../../src/commands/live'
import { Follow } from '../../src/entities/Follow'

it('Should return list contained streamer name', async () => {
  const chatRepository = getConnection().getRepository(Chat)
  const channelRepository = getConnection().getRepository(Channel)
  const followRepository = getConnection().getRepository(Follow)

  const chat = await chatRepository.create({ chatId: '123456', service: Services.TELEGRAM }).save()
  const channel = await channelRepository.create({ id: '12345', username: 'sadisnamenya', online: true }).save()
  
  const i18n = new I18n()
  await i18n.init()
  const follow = await followRepository.save({ channel, chat })
  chat.follows = [follow]

  const command = await liveCommand({ chat, i18n })
  expect(command).toContain('twitch.tv/sadisnamenya')
})