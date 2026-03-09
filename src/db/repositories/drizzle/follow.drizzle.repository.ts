import type { DrizzleD1Database } from 'drizzle-orm/d1';
import { eq, and, count } from 'drizzle-orm';
import { randomUUID } from 'node:crypto';
import { follows } from '../../schema';
import { Follow, FollowAlreadyExistsError, FollowNotFoundError } from '../../../domain/models';
import { DomainMapper } from '../../../domain/mapper';
import type { IFollowRepository } from '../interfaces';

export class FollowDrizzleRepository implements IFollowRepository {
  constructor(private db: DrizzleD1Database) {}

  async findByChatAndChannel(chatId: string, channelId: string): Promise<Follow | undefined> {
    const result = await this.db
      .select()
      .from(follows)
      .where(and(eq(follows.chatId, chatId), eq(follows.channelId, channelId)))
      .limit(1);
    
    return result[0] ? DomainMapper.toDomainFollow(result[0]) : undefined;
  }

  async findByChatId(chatId: string): Promise<Follow[]> {
    const results = await this.db
      .select()
      .from(follows)
      .where(eq(follows.chatId, chatId));
    
    return results.map(r => DomainMapper.toDomainFollow(r));
  }

  async create(chatId: string, channelId: string): Promise<string> {
    // Check if already exists
    const existing = await this.findByChatAndChannel(chatId, channelId);
    if (existing) {
      throw new FollowAlreadyExistsError();
    }
    
    const id = randomUUID();
    await this.db.insert(follows).values({ id, chatId, channelId });
    return id;
  }

  async delete(id: string): Promise<void> {
    const result = await this.db
      .delete(follows)
      .where(eq(follows.id, id))
      .returning();
    
    if (result.length === 0) {
      throw new FollowNotFoundError();
    }
  }

  async findByChannelId(channelId: string): Promise<Follow[]> {
    const results = await this.db
      .select()
      .from(follows)
      .where(eq(follows.channelId, channelId));
    
    return results.map(r => DomainMapper.toDomainFollow(r));
  }

  async findByChatIdPaginated(chatId: string, limit: number, offset: number): Promise<Follow[]> {
    const results = await this.db
      .select()
      .from(follows)
      .where(eq(follows.chatId, chatId))
      .limit(limit)
      .offset(offset);
    
    return results.map(r => DomainMapper.toDomainFollow(r));
  }

  async countByChatId(chatId: string): Promise<number> {
    const result = await this.db
      .select({ count: count() })
      .from(follows)
      .where(eq(follows.chatId, chatId));
    
    return result[0]?.count ?? 0;
  }
}
