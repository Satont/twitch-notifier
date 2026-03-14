import { Bot, InputFile } from 'grammy';
import type { Env } from '../types/env';
import type { I18nService, SupportedLanguage } from './i18n.service';
import { ThumbnailBuilder } from '../utils/thumbnail';

export interface StreamOnlineNotification {
  chatId: number;
  language: SupportedLanguage;
  channelName: string;
  channelUrl: string;
  category: string;
  title: string;
  thumbnailUrl?: string;
  showImage: boolean;
}

export interface StreamOfflineNotification {
  chatId: number;
  language: SupportedLanguage;
  channelName: string;
  channelUrl: string;
  categories: string[];
  duration: string;
}

export interface CategoryChangeNotification {
  chatId: number;
  language: SupportedLanguage;
  channelName: string;
  channelUrl: string;
  oldCategory: string;
  category: string;
}

export interface TitleChangeNotification {
  chatId: number;
  language: SupportedLanguage;
  channelName: string;
  channelUrl: string;
  oldTitle: string;
  title: string;
}

export interface TitleAndCategoryChangeNotification {
  chatId: number;
  language: SupportedLanguage;
  channelName: string;
  channelUrl: string;
  oldTitle: string;
  title: string;
  oldCategory: string;
  category: string;
}

export class TelegramService {
  private bot: Bot;
  private i18n: I18nService;
  private thumbnailBuilder: ThumbnailBuilder;

  constructor(env: Env, i18n: I18nService) {
    this.bot = new Bot(env.TELEGRAM_TOKEN, { client: { timeoutSeconds: 60 } });
    this.i18n = i18n;
    this.thumbnailBuilder = new ThumbnailBuilder();
  }

  async sendStreamOnlineNotification(notification: StreamOnlineNotification): Promise<void> {
    const channelLink = `<a href="${notification.channelUrl}">${notification.channelName}</a>`;

    const text = this.i18n.t(notification.language, 'notifications.streams.nowOnline', {
      channelLink,
      category: notification.category,
      title: notification.title,
    });

    if (notification.showImage && notification.thumbnailUrl) {
      try {
        const thumbnailUrl = this.thumbnailBuilder.build(notification.thumbnailUrl);
        await this.bot.api.sendPhoto(notification.chatId, new InputFile(new URL(thumbnailUrl)), {
          caption: text,
          parse_mode: 'HTML',
        });
        return;
      } catch (error) {
        // Fallback to text message if image fails
        console.error('Failed to send photo:', error);
      }
    }

    await this.bot.api.sendMessage(notification.chatId, text, {
      parse_mode: 'HTML',
      link_preview_options: { is_disabled: false },
    });
  }

  async sendStreamOfflineNotification(notification: StreamOfflineNotification): Promise<void> {
    const channelLink = `<a href="${notification.channelUrl}">${notification.channelName}</a>`;
    const categories = notification.categories.join(', ');

    const text = this.i18n.t(notification.language, 'notifications.streams.nowOffline', {
      channelLink,
      categories,
      duration: notification.duration,
    });

    await this.bot.api.sendMessage(notification.chatId, text, {
      parse_mode: 'HTML',
      link_preview_options: { is_disabled: true },
    });
  }

  async sendCategoryChangeNotification(notification: CategoryChangeNotification): Promise<void> {
    const channelLink = `<a href="${notification.channelUrl}">${notification.channelName}</a>`;

    const text = this.i18n.t(notification.language, 'notifications.streams.newCategory', {
      channelLink,
      oldCategory: notification.oldCategory,
      category: notification.category,
    });

    await this.bot.api.sendMessage(notification.chatId, text, {
      parse_mode: 'HTML',
      link_preview_options: { is_disabled: true },
    });
  }

  async sendTitleChangeNotification(notification: TitleChangeNotification): Promise<void> {
    const channelLink = `<a href="${notification.channelUrl}">${notification.channelName}</a>`;

    const text = this.i18n.t(notification.language, 'notifications.streams.titleChanged', {
      channelLink,
      oldTitle: notification.oldTitle,
      title: notification.title,
    });

    await this.bot.api.sendMessage(notification.chatId, text, {
      parse_mode: 'HTML',
      link_preview_options: { is_disabled: true },
    });
  }

  async sendTitleAndCategoryChangeNotification(
    notification: TitleAndCategoryChangeNotification
  ): Promise<void> {
    const channelLink = `<a href="${notification.channelUrl}">${notification.channelName}</a>`;

    const text = this.i18n.t(notification.language, 'notifications.streams.titleAndCategoryChanged', {
      channelLink,
      oldTitle: notification.oldTitle,
      title: notification.title,
      oldCategory: notification.oldCategory,
      category: notification.category,
    });

    await this.bot.api.sendMessage(notification.chatId, text, {
      parse_mode: 'HTML',
      link_preview_options: { is_disabled: true },
    });
  }

  getBot(): Bot {
    return this.bot;
  }
}
