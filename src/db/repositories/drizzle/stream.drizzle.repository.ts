import type { DrizzleD1Database } from 'drizzle-orm/d1';
import { eq, desc } from 'drizzle-orm';
import { streams } from '../../schema';
import { Stream } from '../../../domain/models';
import { DomainMapper } from '../../../domain/mapper';
import type { IStreamRepository } from '../interfaces';

export class StreamDrizzleRepository implements IStreamRepository {
  constructor(private db: DrizzleD1Database) {}

  async findLatestByChannelId(channelId: string): Promise<Stream | undefined> {
    const result = await this.db
      .select()
      .from(streams)
      .where(eq(streams.channelId, channelId))
      .orderBy(desc(streams.startedAt))
      .limit(1);
    
    return result[0] ? DomainMapper.toDomainStream(result[0]) : undefined;
  }

  async create(id: string, channelId: string, category: string, title: string): Promise<string> {
    try {
      await this.db.insert(streams).values({
        id,
        channelId,
        isLive: true,
        category,
        title,
        startedAt: new Date().toISOString(),
        titles: [title] as any,
        categories: [category] as any,
      });
    } catch (error: any) {
      // If the stream already exists (duplicate webhook), just return the id
      if (error?.message?.includes('UNIQUE constraint failed') || 
          error?.message?.includes('already exists')) {
        console.log(`Stream ${id} already exists, skipping insert`);
        return id;
      }
      throw error;
    }
    
    return id;
  }

  async update(id: string, data: { isLive?: boolean; category?: string; title?: string; endedAt?: string }): Promise<Stream> {
    const result = await this.db
      .update(streams)
      .set(data)
      .where(eq(streams.id, id))
      .returning();
    
    if (!result[0]) {
      throw new Error('Stream not found');
    }
    
    return DomainMapper.toDomainStream(result[0]);
  }

  async findById(id: string): Promise<Stream | undefined> {
    const result = await this.db
      .select()
      .from(streams)
      .where(eq(streams.id, id))
      .limit(1);
    
    return result[0] ? DomainMapper.toDomainStream(result[0]) : undefined;
  }
}
