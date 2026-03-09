import { ApiClient } from '@twurple/api';
import type { Env } from '../types/env';

export class EventSubService {
  private apiClient: ApiClient;
  private webhookUrl: string;
  private secret: string;

  constructor(apiClient: ApiClient, env: Env, baseUrl: string) {
    this.apiClient = apiClient;
    this.webhookUrl = `${baseUrl}/twitch-webhook`;
    this.secret = env.TWITCH_EVENTSUB_SECRET;
  }

  /**
   * Subscribe to all events for a broadcaster (stream.online, stream.offline, channel.update)
   */
  async subscribeToChannel(broadcasterId: string): Promise<void> {
    try {
      // Subscribe to stream online events
      await this.apiClient.eventSub.subscribeToStreamOnlineEvents(
        broadcasterId,
        {
          method: 'webhook',
          callback: this.webhookUrl,
          secret: this.secret,
        }
      );

      // Subscribe to stream offline events
      await this.apiClient.eventSub.subscribeToStreamOfflineEvents(
        broadcasterId,
        {
          method: 'webhook',
          callback: this.webhookUrl,
          secret: this.secret,
        }
      );

      // Subscribe to channel update events (title/category changes)
      await this.apiClient.eventSub.subscribeToChannelUpdateEvents(
        broadcasterId,
        {
          method: 'webhook',
          callback: this.webhookUrl,
          secret: this.secret,
        }
      );
    } catch (error) {
      console.error(`Failed to subscribe to events for broadcaster ${broadcasterId}:`, error);
      throw error;
    }
  }

  /**
   * Unsubscribe from all events for a broadcaster
   */
  async unsubscribeFromChannel(broadcasterId: string): Promise<void> {
    try {
      // Get all subscriptions
      const subscriptions = await this.apiClient.eventSub.getSubscriptions();
      
      // Filter subscriptions for this broadcaster and our webhook URL
      const broadcasterSubs = subscriptions.data.filter(
        (sub) => {
          const transportMethod = (sub as any).transport?.callback || (sub as any)._transport?.callback;
          const broadcastId = (sub.condition as any).broadcaster_user_id;
          return transportMethod === this.webhookUrl && broadcastId === broadcasterId;
        }
      );

      // Delete each subscription
      for (const sub of broadcasterSubs) {
        await this.apiClient.eventSub.deleteSubscription(sub.id);
      }
    } catch (error) {
      console.error(`Failed to unsubscribe from events for broadcaster ${broadcasterId}:`, error);
      throw error;
    }
  }

  /**
   * Check if we already have active subscriptions for a broadcaster
   */
  async hasActiveSubscriptions(broadcasterId: string): Promise<boolean> {
    try {
      const subscriptions = await this.apiClient.eventSub.getSubscriptions();
      
      return subscriptions.data.some(
        (sub) => {
          const transportMethod = (sub as any).transport?.callback || (sub as any)._transport?.callback;
          const broadcastId = (sub.condition as any).broadcaster_user_id;
          return transportMethod === this.webhookUrl && broadcastId === broadcasterId && sub.status === 'enabled';
        }
      );
    } catch (error) {
      console.error(`Failed to check subscriptions for broadcaster ${broadcasterId}:`, error);
      return false;
    }
  }

  /**
   * Delete a specific subscription by ID
   */
  async deleteSubscription(subscriptionId: string): Promise<void> {
    try {
      await this.apiClient.eventSub.deleteSubscription(subscriptionId);
    } catch (error) {
      console.error(`Failed to delete subscription ${subscriptionId}:`, error);
      throw error;
    }
  }

  /**
   * Get all active subscriptions for our webhook
   */
  async getActiveSubscriptions() {
    try {
      const subscriptions = await this.apiClient.eventSub.getSubscriptions();
      return subscriptions.data.filter((sub) => {
        const transportMethod = (sub as any).transport?.callback || (sub as any)._transport?.callback;
        return transportMethod === this.webhookUrl;
      });
    } catch (error) {
      console.error('Failed to get active subscriptions:', error);
      return [];
    }
  }
}
