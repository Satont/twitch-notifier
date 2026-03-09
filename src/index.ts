import { Hono } from 'hono';
import { webhookCallback } from 'grammy';
import { drizzle } from 'drizzle-orm/d1';
import type { Env } from './types';
import { createBot } from './bot';
import { I18nService } from './services/i18n.service';
import { TwitchService } from './services/twitch.service';
import { TelegramService } from './services/telegram.service';
import { EventSubService } from './services/eventsub.service';
import { CloudflareD1Connection } from './db/connection';
import { DrizzleRepositoryFactory } from './db/repository.factory';
import { CloudflareKVSessionRepository } from './db/repositories/cloudflare-kv';
import { handleTwitchWebhook } from './webhooks/twitch';

const app = new Hono<{ Bindings: Env }>();

// Health check
app.get('/', (c) => {
  return c.json({ status: 'ok', service: 'twitch-notifier' });
});

// Telegram webhook endpoint
app.post('/telegram-webhook', async (c) => {
  const env = c.env;

  // Create database connection (serverless-agnostic)
  const dbClient = drizzle(env.DB);
  const dbConnection = new CloudflareD1Connection(dbClient);

  // Create repository factory
  const repositoryFactory = new DrizzleRepositoryFactory(dbConnection);

  // Create repositories
  const chatRepo = repositoryFactory.createChatRepository();
  const channelRepo = repositoryFactory.createChannelRepository();
  const followRepo = repositoryFactory.createFollowRepository();
  const streamRepo = repositoryFactory.createStreamRepository();

  // Create session repository using Cloudflare KV
  const sessionRepo = new CloudflareKVSessionRepository(env.twitch_notifier_kv);

  // Initialize services
  const i18nService = new I18nService();
  await i18nService.init(); // Initialize i18next
  const twitchService = new TwitchService(env);
  const telegramService = new TelegramService(env, i18nService);
  const eventSubService = new EventSubService(
    twitchService.getApiClient(),
    env,
    env.BASE_URL
  );

  // Create bot instance
  const bot = createBot(env, {
    i18n: i18nService,
    twitch: twitchService,
    eventsub: eventSubService,
    chatRepo,
    channelRepo,
    followRepo,
    sessionRepo,
  });

  // Handle webhook
  const handler = webhookCallback(bot, 'hono');
  return handler(c);
});

// Twitch EventSub webhook endpoint
app.post('/twitch-webhook', async (c) => {
  const env = c.env;
  const db = drizzle(env.DB);

  return await handleTwitchWebhook(c.req.raw, env, db);
});

export default app;
