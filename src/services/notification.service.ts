import type { DrizzleD1Database } from 'drizzle-orm/d1';
import type { Env } from '../types/env';
import { TelegramService } from './telegram.service';
import { TwitchService } from './twitch.service';
import { I18nService, type SupportedLanguage } from './i18n.service';
import type {
  IChatRepository,
  IChannelRepository,
  IFollowRepository,
  IStreamRepository,
} from '../db/repositories/interfaces';

export interface StreamOnlineEventData {
  channelId: string;
  channelName: string;
  streamId: string;
  category: string;
  title: string;
  thumbnailUrl?: string;
}

export interface StreamOfflineEventData {
  channelId: string;
  channelName: string;
}

export interface StreamCategoryChangeEventData {
  channelId: string;
  channelName: string;
  oldCategory: string;
  newCategory: string;
}

export interface StreamTitleChangeEventData {
  channelId: string;
  channelName: string;
  oldTitle: string;
  newTitle: string;
}

export class NotificationService {
  constructor(
    private env: Env,
    private db: DrizzleD1Database,
    private telegramService: TelegramService,
    private twitchService: TwitchService,
    private i18nService: I18nService,
    private chatRepo: IChatRepository,
    private channelRepo: IChannelRepository,
    private followRepo: IFollowRepository,
    private streamRepo: IStreamRepository
  ) {}

  async handleStreamOnline(data: StreamOnlineEventData): Promise<void> {
    // Get or create channel
    let channel = await this.channelRepo.findByChannelId(data.channelId, 'twitch');
    if (!channel) {
      channel = await this.channelRepo.create(data.channelId, 'twitch');
    }

    if (!channel.isLive) {
      await this.channelRepo.update(channel.id, { isLive: true });
    }

    // Create stream record
    try {
      await this.streamRepo.create(
        data.streamId,
        channel.id,
        data.category,
        data.title
      );
    } catch (error: any) {
      // If the stream already exists, it's a duplicate webhook - skip notifications
      if (error?.message?.includes('UNIQUE constraint failed') || 
          error?.message?.includes('already exists')) {
        console.log(`Stream ${data.streamId} already exists, skipping duplicate notification`);
        return;
      }
      throw error;
    }

    // Get all followers of this channel
    const follows = await this.followRepo.findByChannelId(channel.id);

    console.log(`Sending online notification for ${data.channelName} to ${follows.length} followers`);

    // Send notifications to all followers
    let sentCount = 0;
    for (const follow of follows) {
      try {
        const chat = await this.chatRepo.findById(follow.chatId);
        if (!chat || !chat.settings) continue;

        await this.telegramService.sendStreamOnlineNotification({
          chatId: parseInt(chat.chatId),
          language: chat.settings.language as SupportedLanguage,
          channelName: data.channelName,
          channelUrl: `https://twitch.tv/${data.channelName}`,
          category: data.category,
          title: data.title,
          thumbnailUrl: data.thumbnailUrl,
          showImage: chat.settings.imageInNotification,
        });
        sentCount++;
      } catch (error) {
        console.error(`Failed to send online notification to chat ${follow.chatId}:`, error);
      }
    }
    console.log(`Sent ${sentCount}/${follows.length} online notifications for ${data.channelName}`);
  }

  async handleStreamOffline(data: StreamOfflineEventData): Promise<void> {
    const channel = await this.channelRepo.findByChannelId(data.channelId, 'twitch');
    if (!channel) {
      console.log(`Channel ${data.channelId} not found for stream.offline`);
      return;
    }

    // Get latest stream
    const stream = await this.streamRepo.findLatestByChannelId(channel.id);
    if (!stream || !stream.isLive) {
      console.log(`No active stream found for ${data.channelName}, skipping offline notification`);
      return;
    }

    // Update stream as offline
    await this.streamRepo.update(stream.id, {
      isLive: false,
      endedAt: new Date().toISOString(),
    });

    if (channel.isLive) {
      await this.channelRepo.update(channel.id, { isLive: false });
    }

    // Get all followers
    const follows = await this.followRepo.findByChannelId(channel.id);

    console.log(`Sending offline notification for ${data.channelName} to ${follows.length} followers`);

    // Send notifications to followers who want offline notifications
    let sentCount = 0;
    for (const follow of follows) {
      try {
        const chat = await this.chatRepo.findById(follow.chatId);
        if (!chat || !chat.settings || !chat.settings.offlineNotification) continue;

        const duration = stream.startedAt
          ? Math.floor((Date.now() - new Date(stream.startedAt).getTime()) / 1000)
          : 0;
        const hours = Math.floor(duration / 3600);
        const minutes = Math.floor((duration % 3600) / 60);
        const seconds = duration % 60;
        const durationStr = `${hours}h ${minutes}m ${seconds}s`;

        await this.telegramService.sendStreamOfflineNotification({
          chatId: parseInt(chat.chatId),
          language: chat.settings.language as SupportedLanguage,
          channelName: data.channelName,
          channelUrl: `https://twitch.tv/${data.channelName}`,
          categories: stream.categories || [],
          duration: durationStr,
        });
        sentCount++;
      } catch (error) {
        console.error(`Failed to send offline notification to chat ${follow.chatId}:`, error);
      }
    }
    console.log(`Sent ${sentCount}/${follows.length} offline notifications for ${data.channelName}`);
  }

  async handleCategoryChange(data: StreamCategoryChangeEventData): Promise<void> {
    const channel = await this.channelRepo.findByChannelId(data.channelId, 'twitch');
    if (!channel) return;

    const stream = await this.streamRepo.findLatestByChannelId(channel.id);
    if (!stream || !stream.isLive) return;

    // Update stream categories
    const categories = [...(stream.categories || []), data.newCategory];
    await this.streamRepo.update(stream.id, {
      category: data.newCategory,
      categories,
    });

    // Get all followers
    const follows = await this.followRepo.findByChannelId(channel.id);

    for (const follow of follows) {
      try {
        const chat = await this.chatRepo.findById(follow.chatId);
        if (!chat || !chat.settings || !chat.settings.gameChangeNotification) continue;

        await this.telegramService.sendCategoryChangeNotification({
          chatId: parseInt(chat.chatId),
          language: chat.settings.language as SupportedLanguage,
          channelName: data.channelName,
          channelUrl: `https://twitch.tv/${data.channelName}`,
          oldCategory: data.oldCategory,
          category: data.newCategory,
        });
      } catch (error) {
        console.error('Failed to send category change notification:', error);
      }
    }
  }

  async handleTitleChange(data: StreamTitleChangeEventData): Promise<void> {
    const channel = await this.channelRepo.findByChannelId(data.channelId, 'twitch');
    if (!channel) return;

    const stream = await this.streamRepo.findLatestByChannelId(channel.id);
    if (!stream || !stream.isLive) return;

    // Update stream titles
    const titles = [...(stream.titles || []), data.newTitle];
    await this.streamRepo.update(stream.id, {
      title: data.newTitle,
      titles,
    });

    // Get all followers
    const follows = await this.followRepo.findByChannelId(channel.id);

    for (const follow of follows) {
      try {
        const chat = await this.chatRepo.findById(follow.chatId);
        if (!chat || !chat.settings || !chat.settings.titleChangeNotification) continue;

        await this.telegramService.sendTitleChangeNotification({
          chatId: parseInt(chat.chatId),
          language: chat.settings.language as SupportedLanguage,
          channelName: data.channelName,
          channelUrl: `https://twitch.tv/${data.channelName}`,
          oldTitle: data.oldTitle,
          title: data.newTitle,
        });
      } catch (error) {
        console.error('Failed to send title change notification:', error);
      }
    }
  }
}
