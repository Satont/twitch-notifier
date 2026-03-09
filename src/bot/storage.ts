import type { StorageAdapter } from 'grammy';
import type { ISessionRepository } from '../db/repositories/interfaces';

/**
 * Storage adapter for Grammy sessions using database persistence
 * Works with any ISessionRepository implementation (D1, PostgreSQL, etc.)
 */
export class DatabaseSessionStorage<T> implements StorageAdapter<T> {
  constructor(
    private sessionRepo: ISessionRepository,
    private ttl?: number // Time to live in seconds
  ) {}

  async read(key: string): Promise<T | undefined> {
    const value = await this.sessionRepo.get(key);
    if (!value) return undefined;

    try {
      return JSON.parse(value) as T;
    } catch (error) {
      console.error('Failed to parse session data:', error);
      return undefined;
    }
  }

  async write(key: string, value: T): Promise<void> {
    const expiresAt = this.ttl ? Date.now() + this.ttl * 1000 : undefined;
    await this.sessionRepo.set(key, JSON.stringify(value), expiresAt);
  }

  async delete(key: string): Promise<void> {
    await this.sessionRepo.delete(key);
  }

  async has(key: string): Promise<boolean> {
    const value = await this.sessionRepo.get(key);
    return value !== undefined;
  }

  /**
   * Clean up expired sessions
   * Should be called periodically (e.g., via cron job)
   */
  async cleanup(): Promise<void> {
    await this.sessionRepo.cleanup();
  }
}
