import type { Stream } from '../../../domain/models';

export interface IStreamRepository {
  findLatestByChannelId(channelId: string): Promise<Stream | undefined>;
  create(id: string, channelId: string, category: string, title: string): Promise<string>;
  update(id: string, data: { 
    isLive?: boolean; 
    category?: string; 
    title?: string; 
    endedAt?: string;
    categories?: string[];
    titles?: string[];
  }): Promise<Stream>;
  findById(id: string): Promise<Stream | undefined>;
}
