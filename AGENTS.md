# PROJECT KNOWLEDGE BASE

**Generated:** 2026-04-12
**Commit:** d08afd7
**Branch:** main

## OVERVIEW
Telegram bot for Twitch stream notifications, running on Cloudflare Workers (serverless). Stack: TypeScript + Grammy (bot) + Hono (HTTP) + Drizzle ORM + D1 (SQLite) + Twurple (Twitch API) + EventSub webhooks.

## STRUCTURE
```
twitch-notifier/
├── src/
│   ├── index.ts              # Hono app entry — /telegram-webhook, /twitch-webhook
│   ├── bot/                  # Grammy bot setup, commands, middleware
│   ├── db/                   # Repository pattern + Drizzle schema
│   ├── domain/               # Domain models (Chat, Channel, Follow, Stream) + errors
│   ├── services/             # TwitchService, TelegramService, EventSubService, NotificationService, I18nService
│   ├── types/                # Env interface (Cloudflare bindings)
│   ├── utils/                # Utilities
│   └── webhooks/twitch.ts   # Twitch EventSub webhook handler
├── locales/                  # i18n JSON files (en, ru, uk)
├── migrations/               # Drizzle D1 migrations (SQL)
├── wrangler.toml             # CF Workers config: D1 + KV bindings, secrets
├── drizzle.config.ts         # Drizzle Kit config
├── docker-compose.yml        # Local dev with Docker
└── .dev.vars                 # Local secrets (gitignored, use .dev.vars.example)
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Add bot command | `src/bot/commands/` | Export from index.ts, register in `src/bot/index.ts` |
| Add notification type | `src/services/notification.service.ts` | Add handler + call from `src/webhooks/twitch.ts` |
| DB schema change | `src/db/schema.ts` → run `bun run db:create` | SQLite/D1; then add repo interface + drizzle impl |
| Add Twitch event | `src/webhooks/twitch.ts` + `src/services/eventsub.service.ts` | Subscribe in EventSubService, handle in webhook |
| Add translation | `locales/{en,ru,uk}.json` | 3 languages required: en, ru, uk |
| Env vars | `src/types/env.ts` + `wrangler.toml` `[vars]` / secrets | Must add to both |
| Session storage | `src/db/repositories/drizzle/session.d1.repository.ts` | D1-backed, 24h TTL |

## CODE MAP
| Symbol | Type | File | Role |
|--------|------|------|------|
| `Env` | interface | `src/types/env.ts` | All CF bindings + secrets |
| `BotContext` | type | `src/bot/types.ts` | Grammy context with services injected |
| `BotSession` | interface | `src/bot/types.ts` | Session shape persisted in D1 |
| `createBot` | function | `src/bot/index.ts` | Bot factory, middleware chain |
| `NotificationService` | class | `src/services/notification.service.ts` | Core fan-out logic for stream events |
| `DrizzleRepositoryFactory` | class | `src/db/repository.factory.ts` | Creates all repos from single DB connection |
| `handleTwitchWebhook` | function | `src/webhooks/twitch.ts` | Routes EventSub events to NotificationService |
| `Chat/Channel/Follow/Stream` | classes | `src/domain/models.ts` | Domain models, separate from DB schema types |
| `DomainMapper` | class | `src/domain/mapper.ts` | Maps DB schema types → domain model classes (single boundary) |

## CONVENTIONS
- Path alias `~/` maps to `./src/` (tsconfig.json `paths`)
- DB schema types (`Chat`, `Channel`) in `src/db/schema.ts` are Drizzle infer types — distinct from domain classes in `src/domain/models.ts`
- All dates stored as ISO strings in D1 (`text`), not Date objects
- `titles` and `categories` in streams table stored as JSON arrays in text column
- Services receive `Env` in constructor — no global singletons
- Each request to `/telegram-webhook` creates fresh service instances (stateless workers)
- `ctx.services.*` for repo/service access inside bot handlers
- Admin commands guarded by `TELEGRAM_BOT_ADMINS` env var (comma-separated IDs)
- Indentation: **tabs** (size 2), max line length 120 (`.editorconfig`); no ESLint/Prettier enforced

## ANTI-PATTERNS
- No test files exist — no `*.test.ts` or test runner configured
- No ESLint/Prettier config — no enforced formatting
- `IDatabaseConnection.getClient()` returns `any` — intentional flexibility escape hatch, not a bug
- `(sub as any).transport` casts in EventSubService — Twurple type gap, not a pattern to copy
- Cloudflare KV session repository exists (`src/db/repositories/cloudflare-kv/`) but is **not used** — D1 is the active session backend

## UNIQUE STYLES
- **Serverless-agnostic design**: repo interfaces allow swapping D1 ↔ PostgreSQL without changing service layer
- **Factory pattern**: `DrizzleRepositoryFactory` is the only place repos are instantiated
- **No DI container**: dependencies manually wired in `src/index.ts` per request
- **EventSub over polling**: all Twitch state updates via webhooks, no background jobs
- `waitUntil` pattern used in Hono handler to keep worker alive after 200 response to Telegram

## COMMANDS
```bash
bun run dev            # Apply local D1 migrations + start wrangler dev
bun run deploy         # wrangler deploy
bun run deploy:with-migrations  # deploy + apply remote migrations
bun run db:migrate     # Apply migrations to remote D1
bun run db:migrate:local  # Apply migrations to local D1
bun run db:create      # Create new migration file
bun run db:studio      # Drizzle Studio UI
```

## NOTES
- `.dev.vars` holds secrets for local dev — never commit, use `.dev.vars.example` as template
- `wrangler.toml` contains real (non-prod) tokens — treat as dev-only credentials
- `BOT_INFO` env var pre-supplies bot info to Grammy to skip getMe() call on startup
- `compatibility_flags = ["nodejs_compat"]` required for `node:crypto` (randomUUID)
- Sessions stored in D1 (not KV) — `SESSIONS_KV` binding exists but unused by current code
- No build step: Wrangler bundles TypeScript directly
- `Dockerfile` is stale — references Go build artifacts (`go.mod`/`make build`) that don't exist; ignore it
