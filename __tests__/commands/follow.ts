import createDbConnection from '../helpers/createDbConnection'
import '../mocks/twitchLib'
import '../mocks/twitchWatcher'

beforeAll(async () => {
  await createDbConnection()
})

import { getConnection } from 'typeorm'
import { Chat, Services } from '../../src/entities/Chat'
import { followCommand } from '../../src/commands/follow'
import { I18n } from '../../src/libs/i18n'
import { Follow } from '../../src/entities/Follow'
import TwitchWatcher from '../../src/watchers/twitch'

beforeAll(() => {
  jest.spyOn(TwitchWatcher, 'addChannelToWatch')
})

it('Should create Chat follow', async () => {
  const chatRepository = getConnection().getRepository(Chat)
  const chat = await chatRepository.create({ chatId: '123456', service: Services.TELEGRAM }).save()
  const followRepository = getConnection().getRepository(Follow)
  
  const follow = await followCommand({ chat, channelName: 'sadisnamenya', i18n: new I18n() })
  expect(follow.success).toEqual(true)
  expect(TwitchWatcher.addChannelToWatch).toHaveBeenCalled()
  expect(await followRepository.findOne({ chat: { id: '123456' }, channel: { id: '12345' } })).not.toBeNull()
})