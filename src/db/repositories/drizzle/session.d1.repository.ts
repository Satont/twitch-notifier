import type { DrizzleD1Database } from 'drizzle-orm/d1';
import { sql } from 'drizzle-orm';
import type { ISessionRepository } from '../interfaces/session.repository.interface';

export class D1SessionRepository implements ISessionRepository {
  constructor(private readonly db: DrizzleD1Database) {}

  async get(key: string): Promise<string | undefined> {
    const now = Date.now();
    const result = await this.db.run(
      sql`SELECT value FROM sessions WHERE key = ${key} AND (expires_at IS NULL OR expires_at > ${now})`
    );
    const row = (result.results as any[])[0];
    return row?.value ?? undefined;
  }

  async set(key: string, value: string, expiresAt?: number): Promise<void> {
    await this.db.run(
      sql`INSERT INTO sessions (key, value, expires_at) VALUES (${key}, ${value}, ${expiresAt ?? null})
          ON CONFLICT(key) DO UPDATE SET value = excluded.value, expires_at = excluded.expires_at`
    );
  }

  async delete(key: string): Promise<void> {
    await this.db.run(sql`DELETE FROM sessions WHERE key = ${key}`);
  }

  async cleanup(): Promise<void> {
    const now = Date.now();
    await this.db.run(sql`DELETE FROM sessions WHERE expires_at IS NOT NULL AND expires_at <= ${now}`);
  }
}
