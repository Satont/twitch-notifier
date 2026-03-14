import type { Env } from '../types/env';
import type { DrizzleD1Database } from 'drizzle-orm/d1';
import { TwitchService } from '../services/twitch.service';
import { TelegramService } from '../services/telegram.service';
import { I18nService } from '../services/i18n.service';
import { NotificationService } from '../services/notification.service';
import { CloudflareD1Connection } from '../db/connection';
import { DrizzleRepositoryFactory } from '../db/repository.factory';
import { createHmac } from 'node:crypto';

interface EventSubNotification {
  subscription: {
    id: string;
    type: string;
    version: string;
    status: string;
    cost: number;
    condition: Record<string, unknown>;
    transport: {
      method: string;
      callback: string;
    };
    created_at: string;
  };
  event: Record<string, any>;
}

interface EventSubVerification {
  challenge: string;
  subscription: {
    id: string;
    type: string;
    version: string;
    status: string;
    cost: number;
    condition: Record<string, unknown>;
    transport: {
      method: string;
      callback: string;
    };
    created_at: string;
  };
}

export async function handleTwitchWebhook(
  request: Request,
  env: Env,
  db: DrizzleD1Database,
  executionCtx: ExecutionContext
): Promise<Response> {
  try {
    // Verify the signature
    const messageId = request.headers.get('Twitch-Eventsub-Message-Id');
    const timestamp = request.headers.get('Twitch-Eventsub-Message-Timestamp');
    const signature = request.headers.get('Twitch-Eventsub-Message-Signature');
    const messageType = request.headers.get('Twitch-Eventsub-Message-Type');

    if (!messageId || !timestamp || !signature) {
      return new Response('Missing required headers', { status: 400 });
    }

    const body = await request.text();

    // Verify signature
    const hmac = createHmac('sha256', env.TWITCH_EVENTSUB_SECRET);
    hmac.update(messageId + timestamp + body);
    const expectedSignature = 'sha256=' + hmac.digest('hex');

    if (signature !== expectedSignature) {
      return new Response('Invalid signature', { status: 403 });
    }

    const payload = JSON.parse(body);

    // Handle verification challenge
    if (messageType === 'webhook_callback_verification') {
      const verification = payload as EventSubVerification;
      return new Response(verification.challenge, {
        status: 200,
        headers: { 'Content-Type': 'text/plain' },
      });
    }

    // Handle notification
    if (messageType === 'notification') {
      const notification = payload as EventSubNotification;

      // Respond immediately to prevent Twitch from retrying due to timeout
      executionCtx.waitUntil(processNotification(notification, env, db));

      return new Response('OK', { status: 200 });
    }

    // Handle revocation
    if (messageType === 'revocation') {
      console.log('Subscription revoked:', payload);
      return new Response('OK', { status: 200 });
    }

    return new Response('Unknown message type', { status: 400 });
  } catch (error) {
    console.error('Error handling Twitch webhook:', error);
    return new Response('Internal Server Error', { status: 500 });
  }
}

async function processNotification(
  notification: EventSubNotification,
  env: Env,
  db: DrizzleD1Database
): Promise<void> {
  try {
    console.log('Received Twitch EventSub notification:', notification.subscription.type, notification);

    const i18nService = new I18nService();
    await i18nService.init();

    const twitchService = new TwitchService(env);
    const telegramService = new TelegramService(env, i18nService);

    const dbConnection = new CloudflareD1Connection(db);
    const repositoryFactory = new DrizzleRepositoryFactory(dbConnection);

    const chatRepo = repositoryFactory.createChatRepository();
    const channelRepo = repositoryFactory.createChannelRepository();
    const followRepo = repositoryFactory.createFollowRepository();
    const streamRepo = repositoryFactory.createStreamRepository();

    const notificationService = new NotificationService(
      env,
      db,
      telegramService,
      twitchService,
      i18nService,
      chatRepo,
      channelRepo,
      followRepo,
      streamRepo
    );

    switch (notification.subscription.type) {
      case 'stream.online': {
        const event = notification.event;
        const stream = await twitchService.getStreamByUserId(event.broadcaster_user_id);
        if (stream) {
          await notificationService.handleStreamOnline({
            channelId: event.broadcaster_user_id,
            channelName: event.broadcaster_user_name,
            streamId: stream.id,
            category: stream.gameName,
            title: stream.title,
            thumbnailUrl: stream.thumbnailUrl,
          });
        }
        break;
      }

      case 'stream.offline': {
        const event = notification.event;
        await notificationService.handleStreamOffline({
          channelId: event.broadcaster_user_id,
          channelName: event.broadcaster_user_name,
        });
        break;
      }

      case 'channel.update': {
        const event = notification.event;
        const channel = await channelRepo.findByChannelId(event.broadcaster_user_id, 'twitch');
        if (!channel) break;

        const stream = await streamRepo.findLatestByChannelId(channel.id);
        if (!stream || !stream.isLive) break;

        if (stream.category && event.category_name !== stream.category) {
          await notificationService.handleCategoryChange({
            channelId: event.broadcaster_user_id,
            channelName: event.broadcaster_user_name,
            oldCategory: stream.category,
            newCategory: event.category_name,
          });
        }

        if (stream.title && event.title !== stream.title) {
          await notificationService.handleTitleChange({
            channelId: event.broadcaster_user_id,
            channelName: event.broadcaster_user_name,
            oldTitle: stream.title,
            newTitle: event.title,
          });
        }
        break;
      }
    }
  } catch (error) {
    console.error('Error processing Twitch notification:', error);
  }
}
