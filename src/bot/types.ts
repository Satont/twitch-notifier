import type { Context, SessionFlavor } from 'grammy';
import type { ConversationFlavor } from '@grammyjs/conversations';
import type { StreamFlavor } from '@grammyjs/stream';
import type { SupportedLanguage } from '../services/i18n.service';
import type { Env } from '../types/env';
import type { DrizzleD1Database } from 'drizzle-orm/d1';
import type { I18nService, TwitchService, EventSubService } from '../services';
import type { IChatRepository, IChannelRepository, IFollowRepository } from '../db/repositories/interfaces';

export interface BotSession {
  chatId?: number;
  language: SupportedLanguage;
  scene?: string;
  followsMenu?: {
    currentPage: number;
    totalPages: number;
  };
}

export type BotContext = Context &
  SessionFlavor<BotSession> &
  StreamFlavor<Context> &
  ConversationFlavor<Context> & {
    t: (key: string, params?: Record<string, any>) => string;
    env: Env;
    db: DrizzleD1Database;
    services: {
      i18n: I18nService;
      twitch: TwitchService;
      eventsub: EventSubService;
      chatRepo: IChatRepository;
      channelRepo: IChannelRepository;
      followRepo: IFollowRepository;
    };
  };
