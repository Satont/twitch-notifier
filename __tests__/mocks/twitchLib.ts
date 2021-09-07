import { HelixBroadcasterType, HelixStream, HelixStreamType, HelixUser, HelixUserType } from 'twitch'

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

const userStream = new HelixStream({
    id: '123456',
    user_id: '12345',
    user_login: 'sadisnamenya',
    user_name: 'sadisnamenya',
    game_id: '12345',
    game_name: 'Dota 2',
    community_ids: [],
    type: HelixStreamType.Live,
    title: 'This is title',
    viewer_count: 5,
    started_at: new Date().toISOString(),
    language: 'ru',
    thumbnail_url: '',
    tag_ids: null,
    is_mature: true,
}, {} as any)

jest.mock('../../src/libs/twitch', () => ({
  Twitch: {
    getUser: jest.fn().mockImplementation(() => Promise.resolve(user)),
    getUsers: jest.fn().mockImplementation(() => Promise.resolve([user])),
    getStreams: jest.fn().mockImplementation(() => Promise.resolve([userStream]))
  }
}))