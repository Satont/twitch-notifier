import { Composer } from 'grammy';
import type { BotContext } from '../types';
import type { Env } from '../../types/env';
import { TwitchService } from '../../services/twitch.service';
import { EventSubService } from '../../services/eventsub.service';

export function createSyncSubscriptionsCommand(env: Env) {
  const sync = new Composer<BotContext>();

  const isAdmin = (userId: number): boolean => {
    const admins = env.TELEGRAM_BOT_ADMINS.split(',').map(id => parseInt(id.trim()));
    return admins.includes(userId);
  };

  sync.command('sync_subscriptions', async (ctx) => {
    const userId = ctx.from?.id;
    if (!userId || !isAdmin(userId)) {
      return;
    }

    await ctx.reply('Checking EventSub subscriptions... This may take a while.');

    try {
      // Get all channels from DB
      const allChannels = await ctx.services.channelRepo.findAll();
      
      const twitchService = new TwitchService(env);
      const eventSubService = new EventSubService(
        twitchService.getApiClient(),
        env,
        env.BASE_URL
      );

      // Get all active subscriptions from Twitch
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

      // Check each channel
      for (const channel of allChannels) {
        if (channel.service !== 'twitch') continue;

        if (subscribedChannelIds.has(channel.channelId)) {
          alreadySubscribed++;
        } else {
          try {
            await eventSubService.subscribeToChannel(channel.channelId);
            created++;
          } catch (error) {
            console.error(`Failed to subscribe to ${channel.channelId}:`, error);
            failed++;
          }
        }
      }

      await ctx.reply(
        `Subscription sync completed!\n` +
        `Total channels: ${allChannels.length}\n` +
        `Already subscribed: ${alreadySubscribed}\n` +
        `Created: ${created}\n` +
        `Failed: ${failed}`
      );
    } catch (error) {
      console.error('Error syncing subscriptions:', error);
      await ctx.reply('Error syncing subscriptions. Check logs.');
    }
  });

  return sync;
}
