import type { Env } from '../types/env';
import type { DrizzleD1Database } from 'drizzle-orm/d1';
import { TwitchService } from '../services/twitch.service';
import { TelegramService } from '../services/telegram.service';
import { I18nService } from '../services/i18n.service';
import { NotificationService } from '../services/notification.service';
import { EventSubService } from '../services/eventsub.service';
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

// Helper to delay execution
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

// Retry function with exponential backoff
async function retryWithBackoff<T>(
  fn: () => Promise<T>,
  maxRetries: number = 5,
  initialDelay: number = 1000
): Promise<T | null> {
  let lastError: Error | null = null;
  
  for (let attempt = 0; attempt < maxRetries; attempt++) {
    try {
      const result = await fn();
      if (result !== null) {
        return result;
      }
      // Result is null, retry
      if (attempt < maxRetries - 1) {
        const waitTime = initialDelay * Math.pow(2, attempt);
        console.log(`Retry ${attempt + 1}/${maxRetries} after ${waitTime}ms...`);
        await delay(waitTime);
      }
    } catch (error) {
      lastError = error as Error;
      console.error(`Attempt ${attempt + 1} failed:`, error);
      if (attempt < maxRetries - 1) {
        const waitTime = initialDelay * Math.pow(2, attempt);
        await delay(waitTime);
      }
    }
  }
  
  if (lastError) {
    console.error('All retry attempts failed:', lastError);
  }
  return null;
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
      console.error('Missing required Twitch headers:', { messageId: !!messageId, timestamp: !!timestamp, signature: !!signature });
      return new Response('Missing required headers', { status: 400 });
    }

    const body = await request.text();

    // Verify signature
    const hmac = createHmac('sha256', env.TWITCH_EVENTSUB_SECRET);
    hmac.update(messageId + timestamp + body);
    const expectedSignature = 'sha256=' + hmac.digest('hex');

    if (signature !== expectedSignature) {
      console.error('Invalid signature received');
      return new Response('Invalid signature', { status: 403 });
    }

    const payload = JSON.parse(body);

    // Handle verification challenge
    if (messageType === 'webhook_callback_verification') {
      const verification = payload as EventSubVerification;
      console.log('Received verification challenge for subscription:', verification.subscription.id);
      return new Response(verification.challenge, {
        status: 200,
        headers: { 'Content-Type': 'text/plain' },
      });
    }

    // Handle notification
    if (messageType === 'notification') {
      const notification = payload as EventSubNotification;
      console.log('Received Twitch EventSub notification:', {
        type: notification.subscription.type,
        id: notification.subscription.id,
        broadcasterId: notification.event?.broadcaster_user_id,
      });

      // Respond immediately to prevent Twitch from retrying due to timeout
      executionCtx.waitUntil(processNotification(notification, env, db));

      return new Response('OK', { status: 200 });
    }

    // Handle revocation
    if (messageType === 'revocation') {
      const revocation = payload as EventSubNotification;
      console.log('Subscription revoked:', {
        id: revocation.subscription.id,
        type: revocation.subscription.type,
        broadcasterId: revocation.subscription.condition?.broadcaster_user_id,
        status: revocation.subscription.status,
      });

      // Recreate the subscription for this broadcaster
      const broadcasterId = (revocation.subscription.condition as any)?.broadcaster_user_id;
      if (broadcasterId) {
        executionCtx.waitUntil(recreateSubscription(broadcasterId, env));
      }

      return new Response('OK', { status: 200 });
    }

    console.warn('Unknown message type:', messageType);
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
    console.log('Processing Twitch EventSub notification:', {
      type: notification.subscription.type,
      broadcasterId: notification.event?.broadcaster_user_id,
      broadcasterName: notification.event?.broadcaster_user_name,
    });

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
        const broadcasterId = event.broadcaster_user_id;
        const broadcasterName = event.broadcaster_user_name;

        // Check if we already processed this stream (idempotent)
        const existingStream = await streamRepo.findById(event.id);
        if (existingStream) {
          console.log(`Stream ${event.id} already processed, skipping notification`);
          return;
        }

        // Try to get stream data from API with retry
        // Twitch API may not have the stream immediately after webhook is sent
        const stream = await retryWithBackoff(
          () => twitchService.getStreamByUserId(broadcasterId),
          5,
          1000
        );

        if (stream) {
          // API returned data - use it
          console.log(`Got stream data from API for ${broadcasterName}:`, {
            streamId: stream.id,
            title: stream.title,
            category: stream.gameName,
          });
          
          await notificationService.handleStreamOnline({
            channelId: broadcasterId,
            channelName: broadcasterName,
            streamId: stream.id,
            category: stream.gameName,
            title: stream.title,
            thumbnailUrl: stream.thumbnailUrl,
          });
        } else {
          // API still doesn't see the stream, but webhook says it's online
          // Create notification with EventSub data as fallback
          console.warn(`Twitch API returned null for ${broadcasterName} after retries, using EventSub data as fallback`);
          
          // Use EventSub stream ID directly
          await notificationService.handleStreamOnline({
            channelId: broadcasterId,
            channelName: broadcasterName,
            streamId: event.id,
            category: 'Unknown',
            title: 'Stream starting...',
            thumbnailUrl: undefined,
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
        if (!channel) {
          console.log(`Channel ${event.broadcaster_user_id} not found for channel.update`);
          break;
        }

        const stream = await streamRepo.findLatestByChannelId(channel.id);
        if (!stream || !stream.isLive) {
          console.log(`No active stream for channel ${event.broadcaster_user_id}, skipping channel.update`);
          break;
        }

        if (stream.category && event.category_name !== stream.category) {
          console.log(`Category changed for ${event.broadcaster_user_name}: ${stream.category} -> ${event.category_name}`);
          await notificationService.handleCategoryChange({
            channelId: event.broadcaster_user_id,
            channelName: event.broadcaster_user_name,
            oldCategory: stream.category,
            newCategory: event.category_name,
          });
        }

        if (stream.title && event.title !== stream.title) {
          console.log(`Title changed for ${event.broadcaster_user_name}`);
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
    throw error;
  }
}

async function recreateSubscription(broadcasterId: string, env: Env): Promise<void> {
  try {
    console.log(`Recreating EventSub subscription for broadcaster ${broadcasterId}`);
    
    const twitchService = new TwitchService(env);
    const eventSubService = new EventSubService(
      twitchService.getApiClient(),
      env,
      env.BASE_URL
    );

    // Small delay to ensure Twitch has fully processed the revocation
    await new Promise(resolve => setTimeout(resolve, 2000));

    await eventSubService.subscribeToChannel(broadcasterId);
    console.log(`Successfully recreated subscription for ${broadcasterId}`);
  } catch (error) {
    console.error(`Failed to recreate subscription for ${broadcasterId}:`, error);
  }
}
