import type { DrizzleD1Database } from 'drizzle-orm/d1';

/**
 * Database connection abstraction
 * This allows us to support different database implementations (D1, PostgreSQL, etc.)
 */
export interface IDatabaseConnection {
  getClient(): any; // Returns the underlying database client (DrizzleD1Database, etc.)
}

/**
 * Cloudflare D1 database connection
 */
export class CloudflareD1Connection implements IDatabaseConnection {
  constructor(private client: DrizzleD1Database) {}

  getClient(): DrizzleD1Database {
    return this.client;
  }
}
