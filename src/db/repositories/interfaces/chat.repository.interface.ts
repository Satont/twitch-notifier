import type { Chat, ChatSettings } from '../../../domain/models';

export interface IChatRepository {
  findByChatId(chatId: number, service: 'telegram'): Promise<Chat | undefined>;
  findById(id: string): Promise<Chat | undefined>;
  findAllByService(service: 'telegram'): Promise<Chat[]>;
  create(chatId: string, service: 'telegram'): Promise<string>;
  updateSettings(chatId: string, settings: Partial<ChatSettings>): Promise<void>;
}
