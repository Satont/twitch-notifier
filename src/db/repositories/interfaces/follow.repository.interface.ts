import type { Follow } from '../../../domain/models';

export interface IFollowRepository {
  findByChatAndChannel(chatId: string, channelId: string): Promise<Follow | undefined>;
  findByChatId(chatId: string): Promise<Follow[]>;
  create(chatId: string, channelId: string): Promise<string>;
  delete(id: string): Promise<void>;
  findByChannelId(channelId: string): Promise<Follow[]>;
  findByChatIdPaginated(chatId: string, limit: number, offset: number): Promise<Follow[]>;
  countByChatId(chatId: string): Promise<number>;
}
