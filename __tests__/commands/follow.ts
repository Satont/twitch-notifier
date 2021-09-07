import createDbConnection from '../helpers/createDbConnection'
/* import Twitch from '../../src/libs/twitch' */
import { HelixBroadcasterType, HelixUser, HelixUserType } from 'twitch'

const user = new HelixUser({
  id: '12345',
  login: 'sadisnamenya',
  display_name: 'sadisnamenya',
  description: 'sadisnamenya',
  type: HelixUserType.None,
  broadcaster_type: HelixBroadcasterType.Affiliate,
  profile_image_url: '',
  offline_image_url: '',
  view_count: 5,
  created_at: new Date().toISOString(),
}, {} as any)

jest.mock('../../src/watchers/twitch', () => ({
  TwitchWatcher: {
    addChannelToWatch: jest.fn().mockImplementation(() => true)
  },
}))

jest.mock('../../src/libs/twitch', () => ({
  Twitch: {
    getUser: jest.fn().mockImplementation(() => user),
    getUsers: jest.fn().mockImplementation(() => [user]),
  }
}))

beforeAll(async () => {
  await createDbConnection()
})

import { getConnection } from 'typeorm'
import { Chat, Services } from '../../src/entities/Chat'
import { followCommand } from '../../src/commands/follow'
import { I18n } from '../../src/libs/i18n'
import { Channel } from '../../src/entities/Channel'

it('Should create Chat follow', async () => {
  const chatRepository = getConnection().getRepository(Chat)
  const channelRepository = getConnection().getRepository(Channel)
  const chat = await chatRepository.create({ chatId: '123456', service: Services.TELEGRAM }).save()

  const follow = await followCommand({ chat, channelName: 'sadisnamenya', i18n: new I18n() })
  expect(follow.success).toEqual(true)
})