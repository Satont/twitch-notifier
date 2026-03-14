import type { DrizzleD1Database } from 'drizzle-orm/d1';
import { eq } from 'drizzle-orm';
import { randomUUID } from 'node:crypto';
import { chats, chatSettings } from '../../schema';
import { Chat, ChatSettings } from '../../../domain/models';
import { DomainMapper } from '../../../domain/mapper';
import type { IChatRepository } from '../interfaces';

export class ChatDrizzleRepository implements IChatRepository {
  constructor(private db: DrizzleD1Database) {}

  async findByChatId(chatId: number, service: 'telegram' = 'telegram'): Promise<Chat | undefined> {
    const chatIdStr = chatId.toString();
    
    const chatResult = await this.db
      .select()
      .from(chats)
      .where(eq(chats.chatId, chatIdStr))
      .limit(1);
    
    if (!chatResult[0]) return undefined;

    const settingsResult = await this.db
      .select()
      .from(chatSettings)
      .where(eq(chatSettings.chatId, chatResult[0].id))
      .limit(1);
    
    return DomainMapper.toDomainChat({
      ...chatResult[0],
      settings: settingsResult[0] || null
    });
  }

  async findById(id: string): Promise<Chat | undefined> {
    const chatResult = await this.db
      .select()
      .from(chats)
      .where(eq(chats.id, id))
      .limit(1);
    
    if (!chatResult[0]) return undefined;

    const settingsResult = await this.db
      .select()
      .from(chatSettings)
      .where(eq(chatSettings.chatId, chatResult[0].id))
      .limit(1);
    
    return DomainMapper.toDomainChat({
      ...chatResult[0],
      settings: settingsResult[0] || null
    });
  }

  async findAllByService(service: 'telegram' = 'telegram'): Promise<Chat[]> {
    const chatResults = await this.db
      .select()
      .from(chats)
      .where(eq(chats.service, service));
    
    const chatsWithSettings: Chat[] = [];
    
    for (const chat of chatResults) {
      const settingsResult = await this.db
        .select()
        .from(chatSettings)
        .where(eq(chatSettings.chatId, chat.id))
        .limit(1);
      
      chatsWithSettings.push(DomainMapper.toDomainChat({
        ...chat,
        settings: settingsResult[0] || null
      }));
    }
    
    return chatsWithSettings;
  }

  async create(chatId: string, service: 'telegram' = 'telegram'): Promise<string> {
    const id = randomUUID();
    await this.db.insert(chats).values({ id, chatId, service });
    
    // Create default settings
    await this.db.insert(chatSettings).values({
      chatId: id,
      language: 'en',
      offlineNotification: true,
      gameChangeNotification: false,
      titleChangeNotification: false,
      gameAndTitleChangeNotification: false,
      imageInNotification: true,
    });
    
    return id;
  }

  async updateSettings(chatId: string, settings: Partial<ChatSettings>): Promise<void> {
    await this.db
      .update(chatSettings)
      .set(settings)
      .where(eq(chatSettings.chatId, chatId));
  }
}
