import { getConnection } from 'typeorm'
import { Stream } from '../entities/Stream'

export function getLatestStream(channelId: string) {
  return getConnection().getRepository(Stream).findOne({
    where: {
      channel: {
        id: channelId,
      },
    },
    order: {
      startedAt: 'DESC',
    },
  })
}