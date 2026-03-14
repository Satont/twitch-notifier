import type { KVNamespace } from '@cloudflare/workers-types';
import type { ISessionRepository } from '../interfaces/session.repository.interface';

/**
 * Cloudflare KV-based session repository
 * Fast, distributed key-value storage perfect for sessions
 */
export class CloudflareKVSessionRepository implements ISessionRepository {
  constructor(private readonly kv: KVNamespace) {}

  async get(key: string): Promise<string | undefined> {
    const value = await this.kv.get(key);
    return value ?? undefined;
  }

  async set(key: string, value: string, expiresAt?: number): Promise<void> {
    const options: { expirationTtl?: number } = {};

    // Convert expiresAt (unix timestamp) to TTL in seconds
    if (expiresAt) {
      const ttl = Math.floor((expiresAt - Date.now()) / 1000);
      if (ttl > 0) {
        options.expirationTtl = ttl;
      }
    }

    await this.kv.put(key, value, options);
  }

  async delete(key: string): Promise<void> {
    await this.kv.delete(key);
  }

  async cleanup(): Promise<void> {
    // KV automatically cleans up expired keys, no manual cleanup needed
    return;
  }
}
