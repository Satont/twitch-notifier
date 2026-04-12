import { ApiClient } from '@twurple/api';
import type { Env } from '../types/env';

export class EventSubService {
  private apiClient: ApiClient;
  private webhookUrl: string;
  private secret: string;

  private readonly requiredSubscriptionTypes = [
    'stream.online',
    'stream.offline',
    'channel.update',
  ] as const;

  constructor(apiClient: ApiClient, env: Env, baseUrl: string) {
    this.apiClient = apiClient;
    this.webhookUrl = `${baseUrl}/twitch-webhook`;
    this.secret = env.TWITCH_EVENTSUB_SECRET;
  }

  private getSubscriptionCallback(sub: any): string | undefined {
    return sub.transport?.callback || sub._transport?.callback;
  }

  private getSubscriptionBroadcasterId(sub: any): string | undefined {
    return sub.condition?.broadcaster_user_id;
  }

  private isUsableSubscriptionStatus(status: string): boolean {
    return status === 'enabled' || status === 'webhook_callback_verification_pending';
  }

  private isConflictError(error: any): boolean {
    return error?.status === 409 ||
      error?.message?.includes('subscription already exists') ||
      error?.message?.includes('Conflict');
  }

  private async getAllWebhookSubscriptions() {
    const subscriptions = await this.apiClient.eventSub.getSubscriptionsPaginated().getAll();

    return subscriptions.filter((sub: any) => this.getSubscriptionCallback(sub) === this.webhookUrl);
  }

  private getExistingTypesForBroadcaster(subscriptions: any[], broadcasterId: string): Set<string> {
    return new Set(
      subscriptions
        .filter((sub: any) => {
          return this.getSubscriptionBroadcasterId(sub) === broadcasterId && this.isUsableSubscriptionStatus(sub.status);
        })
        .map((sub: any) => sub.type)
    );
  }

  private async subscribeToMissingType(type: string, broadcasterId: string): Promise<boolean> {
    try {
      if (type === 'stream.online') {
        await this.apiClient.eventSub.subscribeToStreamOnlineEvents(broadcasterId, {
          method: 'webhook',
          callback: this.webhookUrl,
          secret: this.secret,
        });
        return true;
      }

      if (type === 'stream.offline') {
        await this.apiClient.eventSub.subscribeToStreamOfflineEvents(broadcasterId, {
          method: 'webhook',
          callback: this.webhookUrl,
          secret: this.secret,
        });
        return true;
      }

      if (type === 'channel.update') {
        await this.apiClient.eventSub.subscribeToChannelUpdateEvents(broadcasterId, {
          method: 'webhook',
          callback: this.webhookUrl,
          secret: this.secret,
        });
        return true;
      }

      return false;
    } catch (error) {
      if (this.isConflictError(error)) {
        return false;
      }

      throw error;
    }
  }

  async ensureSubscriptions(broadcasterId: string, existingTypes?: Iterable<string>): Promise<number> {
    const presentTypes = new Set(existingTypes ?? []);
    let createdCount = 0;

    for (const type of this.requiredSubscriptionTypes) {
      if (presentTypes.has(type)) {
        continue;
      }

      const created = await this.subscribeToMissingType(type, broadcasterId);
      if (created) {
        createdCount++;
        presentTypes.add(type);
      }
    }

    return createdCount;
  }

  /**
   * Subscribe to all events for a broadcaster (stream.online, stream.offline, channel.update)
   */
  async subscribeToChannel(broadcasterId: string): Promise<void> {
    try {
      const subscriptions = await this.getAllWebhookSubscriptions();
      const existingTypes = this.getExistingTypesForBroadcaster(subscriptions, broadcasterId);
      await this.ensureSubscriptions(broadcasterId, existingTypes);
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
      const subscriptions = await this.getAllWebhookSubscriptions();
      
      const broadcasterSubs = subscriptions.filter(
        (sub: any) => this.getSubscriptionBroadcasterId(sub) === broadcasterId
      );

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
      const subscriptions = await this.getAllWebhookSubscriptions();
      const existingTypes = this.getExistingTypesForBroadcaster(subscriptions, broadcasterId);

      return this.requiredSubscriptionTypes.every((type) => existingTypes.has(type));
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
      const subscriptions = await this.getAllWebhookSubscriptions();
      return subscriptions.filter((sub: any) => this.isUsableSubscriptionStatus(sub.status));
    } catch (error) {
      console.error('Failed to get active subscriptions:', error);
      return [];
    }
  }
}
