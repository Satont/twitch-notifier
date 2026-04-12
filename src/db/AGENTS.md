# DATABASE LAYER

## OVERVIEW
Repository pattern with interface/implementation split. Active backend: Drizzle ORM + Cloudflare D1 (SQLite). KV implementation exists but unused.

## STRUCTURE
```
db/
├── schema.ts                  # Drizzle table definitions + inferred TS types
├── connection.ts              # IDatabaseConnection interface + CloudflareD1Connection impl
├── repository.factory.ts      # DrizzleRepositoryFactory — single instantiation point
├── index.ts                   # Re-exports
└── repositories/
    ├── interfaces/            # IChatRepository, IChannelRepository, IFollowRepository, IStreamRepository, ISessionRepository
    ├── drizzle/               # Active implementations: Chat, Channel, Follow, Stream, Session (D1)
    └── cloudflare-kv/         # UNUSED — KV session repo, kept for reference
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Add table | `schema.ts` → `bun run db:create` | SQLite/D1; add Drizzle relations too |
| Add repo method | `interfaces/<entity>.repository.interface.ts` + `drizzle/<entity>.drizzle.repository.ts` | Interface first, then impl |
| Add new repo | Add interface → add drizzle impl → add factory method in `repository.factory.ts` | |
| Session logic | `drizzle/session.d1.repository.ts` | 24h TTL, D1-backed |

## CONVENTIONS
- Schema types (`Chat`, `Channel` etc.) in `schema.ts` = Drizzle `$inferSelect` — NOT the same as domain classes in `src/domain/models.ts`
- All dates as ISO strings (`text`) in D1 — never `Date` objects in schema
- `titles`/`categories` on streams = JSON arrays stored in `text` column (`mode: 'json'`)
- `IDatabaseConnection.getClient()` returns `any` — intentional, allows D1 ↔ PostgreSQL swap
- `DrizzleRepositoryFactory` is the ONLY place repo classes are instantiated
- Streams table: primary key is Twitch stream ID (string), not UUID

## ANTI-PATTERNS
- Never instantiate repo classes directly outside the factory
- Don't use `cloudflare-kv/` session impl — D1 is active
- Don't store Date objects in schema — use ISO string and convert in domain layer
