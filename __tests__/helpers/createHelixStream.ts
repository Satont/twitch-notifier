import { HelixStream, HelixStreamType, HelixUser } from 'twitch/lib'

export const createHelixStream = (user?: Partial<HelixUser>) => {
  return new HelixStream({
    id: '123456',
    user_id: user?.id ?? '12345',
    user_login: user?.name ?? 'sadisnamenya',
    user_name: user?.name ?? 'sadisnamenya',
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
}