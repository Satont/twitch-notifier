import type { D1Database, KVNamespace } from '@cloudflare/workers-types';

export interface Env {
  // D1 Database
  twitch_notifier_db: D1Database;

  // KV Namespace for sessions
  twitch_notifier_kv: KVNamespace;

  // Secrets
  TELEGRAM_TOKEN: string;
  TWITCH_CLIENT_ID: string;
  TWITCH_CLIENT_SECRET: string;
  TELEGRAM_BOT_ADMINS: string; // comma-separated user IDs
  TWITCH_EVENTSUB_SECRET: string;
  BASE_URL: string; // Base URL for webhooks (e.g., https://your-worker.workers.dev)
	BOT_INFO: string;

  // Variables
  APP_ENV: 'development' | 'production';
}
