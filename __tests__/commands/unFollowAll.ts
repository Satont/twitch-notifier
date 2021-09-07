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
import { unFollowAllCommand } from '../../src/commands/unFollowAll'

it('Should unfollow from all channels', async () => {
  const chatRepository = getConnection().getRepository(Chat)
  const chat = await chatRepository.create({ chatId: '123456', service: Services.TELEGRAM }).save()
  const followRepository = getConnection().getRepository(Follow)

  await followCommand({ chat, channelName: 'sadisnamenya', i18n: new I18n() })
  const unFollowAll = await unFollowAllCommand({ chat, i18n: new I18n() })

  expect(await followRepository.count({ chat: { id: '123456' } })).toEqual(0)
  expect(unFollowAll.success).toEqual(true)
})