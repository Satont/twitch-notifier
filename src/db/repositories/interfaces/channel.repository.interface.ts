import type { Channel } from '../../../domain/models';
import type { NewChannel } from '../../schema';

export interface IChannelRepository {
  findByChannelId(channelId: string, service: 'twitch'): Promise<Channel | undefined>;
  findById(id: string): Promise<Channel | undefined>;
  create(channelId: string, service: 'twitch'): Promise<Channel>;
  update(id: string, data: Partial<Omit<NewChannel, 'id'>>): Promise<Channel>;
  updateChannelId(oldChannelId: string, newChannelId: string, service: 'twitch'): Promise<void>;
}
