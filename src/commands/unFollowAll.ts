import { getConnection } from 'typeorm'
import { Chat } from '../entities/Chat'
import { Follow } from '../entities/Follow'
import { I18n } from '../libs/i18n'


export async function unFollowAllCommand({ chat, i18n }: { chat: Chat, i18n: I18n }) {
  const followRepository = getConnection().getRepository(Follow)
  const follows = await followRepository.find({
    chat,
  })

  if (!follows) {
    return {
      success: true,
      message: i18n.translate('commands.follows.emptyList'),
    }
  } else {
    for (const follow of follows) {
      await follow.remove()
    }

    return {
      success: true,
      message: i18n.translate('commands.unfollowall.success', { count: follows.length }),
    }
  }
}
