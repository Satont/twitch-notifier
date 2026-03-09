import { Bot } from 'grammy';
import { session } from 'grammy'
import type { Env } from '../types/env';
import type { BotSession, BotContext } from './types';
import { I18nService } from '../services/i18n.service';
import { TwitchService } from '../services/twitch.service';
import { EventSubService } from '../services/eventsub.service';
import { DatabaseSessionStorage } from './storage';
import type { IChatRepository, IChannelRepository, IFollowRepository, ISessionRepository } from '../db/repositories/interfaces';
import {
  startCommand,
  followCommand,
  followsCommand,
  liveCommand,
  createBroadcastCommand,
  createChangeChannelIdCommand,
  callbackQueryHandler
} from './commands';

export function createBot(
  env: Env,
  services: {
    i18n: I18nService;
    twitch: TwitchService;
    eventsub: EventSubService;
    chatRepo: IChatRepository;
    channelRepo: IChannelRepository;
    followRepo: IFollowRepository;
    sessionRepo: ISessionRepository;
  }
): Bot<BotContext> {
  const bot = new Bot<BotContext>(env.TELEGRAM_TOKEN);

  // Use database session storage
  const sessionStorage = new DatabaseSessionStorage<BotSession>(
    services.sessionRepo,
    86400 // 24 hours TTL
  );

	bot.use(session({
		initial: (): BotSession => ({
			language: 'en',
			followsMenu: {
				currentPage: 1,
				totalPages: 1,
			},
		}),
		storage: sessionStorage,
	}))

  // Attach environment and services to context
  bot.use(async (ctx, next) => {
    ctx.env = env;
    ctx.services = services;
    await next();
  });

  // Use i18n middleware
  bot.use(services.i18n.middleware());

  // Register commands
  bot.use(startCommand);
  bot.use(followCommand);
  bot.use(followsCommand);
  bot.use(liveCommand);
  bot.use(createBroadcastCommand(env));
  bot.use(createChangeChannelIdCommand(env));
  bot.use(callbackQueryHandler);

  return bot;
}
