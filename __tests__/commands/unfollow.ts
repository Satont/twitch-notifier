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
import { unFollowCommand } from '../../src/commands/unfollow'

it('Should unfollow from channel', async () => {
  const chatRepository = getConnection().getRepository(Chat)
  const chat = await chatRepository.create({ chatId: '123456', service: Services.TELEGRAM }).save()
  const followRepository = getConnection().getRepository(Follow)

  await followCommand({ chat, channelName: 'sadisnamenya', i18n: new I18n() })
  const unFollow = await unFollowCommand({ chat, channelName: 'sadisnamenya', i18n: new I18n() })

  expect(await followRepository.findOne({ chat: { id: '123456' }, channel: { id: '12345' } })).not.toBeNull()
  expect(unFollow.success).toEqual(true)
})