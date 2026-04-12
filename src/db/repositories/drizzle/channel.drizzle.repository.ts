import type { DrizzleD1Database } from 'drizzle-orm/d1';
import { eq, and } from 'drizzle-orm';
import { randomUUID } from 'node:crypto';
import { channels } from '../../schema';
import type { NewChannel } from '../../schema';
import { Channel, ChannelNotFoundError } from '../../../domain/models';
import { DomainMapper } from '../../../domain/mapper';
import type { IChannelRepository } from '../interfaces';

export class ChannelDrizzleRepository implements IChannelRepository {
  constructor(private db: DrizzleD1Database) {}

  async findByChannelId(channelId: string, service: 'twitch' = 'twitch'): Promise<Channel | undefined> {
    const result = await this.db
      .select()
      .from(channels)
      .where(and(eq(channels.channelId, channelId), eq(channels.service, service)))
      .limit(1);
    
    return result[0] ? DomainMapper.toDomainChannel(result[0]) : undefined;
  }

  async findById(id: string): Promise<Channel | undefined> {
    const result = await this.db
      .select()
      .from(channels)
      .where(eq(channels.id, id))
      .limit(1);
    
    return result[0] ? DomainMapper.toDomainChannel(result[0]) : undefined;
  }

  async findAll(): Promise<Channel[]> {
    const result = await this.db
      .select()
      .from(channels);
    
    return result.map(DomainMapper.toDomainChannel);
  }

  async create(channelId: string, service: 'twitch' = 'twitch'): Promise<Channel> {
    const id = randomUUID();
    const result = await this.db.insert(channels).values({
      id,
      channelId,
      service,
      isLive: false,
    }).returning();
    
    return DomainMapper.toDomainChannel(result[0]);
  }

  async update(id: string, data: Partial<Omit<NewChannel, 'id'>>): Promise<Channel> {
    const result = await this.db
      .update(channels)
      .set({ ...data, updatedAt: new Date().toISOString() })
      .where(eq(channels.id, id))
      .returning();
    
    if (!result[0]) {
      throw new ChannelNotFoundError();
    }
    
    return DomainMapper.toDomainChannel(result[0]);
  }

  async updateChannelId(oldChannelId: string, newChannelId: string, service: 'twitch' = 'twitch'): Promise<void> {
    await this.db
      .update(channels)
      .set({ channelId: newChannelId, updatedAt: new Date().toISOString() })
      .where(and(eq(channels.channelId, oldChannelId), eq(channels.service, service)));
  }
}
