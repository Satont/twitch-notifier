// Mappers to convert between database schema and domain models
import type { Chat as DbChat, ChatSettings as DbChatSettings, Channel as DbChannel, Follow as DbFollow, Stream as DbStream } from '../db/schema';
import { Chat, ChatSettings, Channel, Follow, Stream } from './models';
import type { SupportedLanguage } from './models';

export class DomainMapper {
  static toDomainChat(dbChat: DbChat & { settings: DbChatSettings | null }): Chat {
    return new Chat({
      id: dbChat.id,
      chatId: dbChat.chatId,
      service: dbChat.service,
      settings: dbChat.settings ? this.toDomainChatSettings(dbChat.settings) : undefined,
    });
  }

  static toDomainChatSettings(dbSettings: DbChatSettings): ChatSettings {
    return new ChatSettings({
      id: dbSettings.id,
      chatId: dbSettings.chatId,
      gameChangeNotification: dbSettings.gameChangeNotification,
      titleChangeNotification: dbSettings.titleChangeNotification,
      gameAndTitleChangeNotification: dbSettings.gameAndTitleChangeNotification,
      offlineNotification: dbSettings.offlineNotification,
      imageInNotification: dbSettings.imageInNotification,
      language: dbSettings.language as SupportedLanguage,
    });
  }

  static toDomainChannel(dbChannel: DbChannel): Channel {
    return new Channel({
      id: dbChannel.id,
      channelId: dbChannel.channelId,
      service: dbChannel.service,
      isLive: dbChannel.isLive,
      title: dbChannel.title ?? undefined,
      category: dbChannel.category ?? undefined,
      updatedAt: dbChannel.updatedAt ? new Date(dbChannel.updatedAt) : undefined,
    });
  }

  static toDomainFollow(dbFollow: DbFollow): Follow {
    return new Follow({
      id: dbFollow.id,
      channelId: dbFollow.channelId,
      chatId: dbFollow.chatId,
    });
  }

  static toDomainStream(dbStream: DbStream): Stream {
    return new Stream({
      id: dbStream.id,
      channelId: dbStream.channelId,
      isLive: dbStream.isLive,
      title: dbStream.title ?? undefined,
      category: dbStream.category ?? undefined,
      titles: dbStream.titles,
      categories: dbStream.categories,
      startedAt: new Date(dbStream.startedAt!),
      updatedAt: dbStream.updatedAt ? new Date(dbStream.updatedAt) : undefined,
      endedAt: dbStream.endedAt ? new Date(dbStream.endedAt) : undefined,
    });
  }
}
