import { HelixUser, HelixUserType, HelixBroadcasterType } from 'twitch'

export const createHelixUser = (opts: { name?: string, id?: string} = {}) => {
  return new HelixUser({
    id: opts.id ?? '12345',
    login: opts.name ?? 'sadisnamenya',
    display_name: opts.name ?? 'sadisnamenya',
    description: '',
    type: HelixUserType.None,
    broadcaster_type: HelixBroadcasterType.Affiliate,
    profile_image_url: '',
    offline_image_url: '',
    view_count: 5,
    created_at: new Date().toISOString(),
  }, {} as any)
}