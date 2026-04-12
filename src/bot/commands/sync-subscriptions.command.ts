import { Composer } from "grammy";
import type { BotContext } from "../types";
import type { Env } from "../../types/env";
import { TwitchService } from "../../services/twitch.service";
import { EventSubService } from "../../services/eventsub.service";

const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));
const syncLockKey = 'sync_subscriptions_lock';

function formatProgress(
  processed: number,
  total: number,
  created: number,
  alreadySubscribed: number,
  failed: number,
  done: boolean = false,
): string {
  if (done) {
    return [
      'Subscription sync completed!',
      `Total channels: ${total}`,
      `Channels already complete: ${alreadySubscribed}`,
      `Created subscriptions: ${created}`,
      `Failed: ${failed}`,
    ].join('\n');
  }

  return [
    'Syncing subscriptions...',
    `Progress: ${processed}/${total}`,
    `Created subscriptions: ${created}`,
    `Channels already complete: ${alreadySubscribed}`,
    `Failed: ${failed}`,
  ].join('\n');
}

async function* syncSubscriptionProgress(
  ctx: BotContext,
  env: Env,
): AsyncGenerator<string> {
  const allChannels = await ctx.services.channelRepo.findAll();
  const twitchChannels = allChannels.filter(c => c.service === "twitch");

  const twitchService = new TwitchService(env);
  const eventSubService = new EventSubService(twitchService.getApiClient(), env, env.BASE_URL);

  const activeSubscriptions = await eventSubService.getActiveSubscriptions();
  const subscribedChannelIds = new Set<string>();

  for (const sub of activeSubscriptions) {
    const broadcastId = (sub.condition as any)?.broadcaster_user_id;
    if (broadcastId) {
      subscribedChannelIds.add(broadcastId);
    }
  }

  let created = 0;
  let alreadySubscribed = 0;
  let failed = 0;
  let processed = 0;

  const total = twitchChannels.length;
  const progressInterval = Math.max(1, Math.floor(total / 50));

  yield formatProgress(processed, total, created, alreadySubscribed, failed);

  for (const channel of twitchChannels) {
    processed++;

    if (subscribedChannelIds.has(channel.channelId)) {
      alreadySubscribed++;
    } else {
      const createdCount = await retrySubscribe(eventSubService, channel.channelId);
      if (createdCount >= 0) {
        created += createdCount;

        if (createdCount === 0) {
          alreadySubscribed++;
        }
      } else {
        failed++;
      }

      await delay(100);
    }

    if (processed % progressInterval === 0 || processed === total) {
      yield formatProgress(processed, total, created, alreadySubscribed, failed);
    }
  }

  yield formatProgress(processed, total, created, alreadySubscribed, failed, true);
}

async function retrySubscribe(
  eventSubService: EventSubService,
  channelId: string,
  maxRetries: number = 10
): Promise<number> {
  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      return await eventSubService.ensureSubscriptions(channelId);
    } catch (error: any) {
      const isRateLimit = error?.status === 429 || 
                         error?.message?.includes('rate limit') ||
                         error?.message?.includes('Too Many Requests');
      
      const isAlreadyExists = error?.status === 409 ||
                             error?.message?.includes('subscription already exists') ||
                             error?.message?.includes('Conflict');
      
      if (isAlreadyExists) {
        return 0;
      }
      
      if (isRateLimit && attempt < maxRetries) {
        const waitTime = Math.min(1000 * Math.pow(2, attempt - 1), 30000);
        console.log(`Rate limit for ${channelId}, retry ${attempt}/${maxRetries} after ${waitTime}ms...`);
        await delay(waitTime);
        continue;
      }
      
      if (attempt === maxRetries) {
        console.error(`Failed to subscribe to ${channelId} after ${maxRetries} attempts:`, error);
        return -1;
      }
      
      throw error;
    }
  }
  return -1;
}

export function createSyncSubscriptionsCommand(env: Env) {
  const sync = new Composer<BotContext>();

  const isAdmin = (userId: number): boolean => {
    const admins = env.TELEGRAM_BOT_ADMINS.split(",").map((id) => parseInt(id.trim()));
    return admins.includes(userId);
  };

  sync.command("sync_subscriptions", async (ctx) => {
    const userId = ctx.from?.id;
    if (!userId || !isAdmin(userId)) {
      return;
    }

    const activeLock = await ctx.env.twitch_notifier_kv.get(syncLockKey);
    if (activeLock) {
      await ctx.reply("Sync already in progress. Please wait.");
      return;
    }

    await ctx.env.twitch_notifier_kv.put(syncLockKey, "1", { expirationTtl: 60 * 15 });

    try {
      if (ctx.chat?.type === 'private') {
        await ctx.replyWithStream(syncSubscriptionProgress(ctx, env));
      } else {
        for await (const message of syncSubscriptionProgress(ctx, env)) {
          console.log(message);
        }
      }
    } catch (error) {
      console.error("Error syncing subscriptions:", error);
      await ctx.reply("Error syncing subscriptions. Check logs.");
    } finally {
      await ctx.env.twitch_notifier_kv.delete(syncLockKey);
    }
  });

  return sync;
}
