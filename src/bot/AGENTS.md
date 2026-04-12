# BOT LAYER

## OVERVIEW
Grammy bot factory, middleware chain, session management, and command registration.

## STRUCTURE
```
bot/
├── index.ts          # createBot() — wires middleware + commands
├── types.ts          # BotContext, BotSession type definitions
├── storage.ts        # DatabaseSessionStorage adapter (D1-backed)
├── helpers.ts        # Shared UI helpers (keyboards, message builders)
└── commands/         # Individual command handlers
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Add command | `commands/<name>.command.ts` | Export from `commands/index.ts`, register in `index.ts` |
| Add session field | `types.ts` `BotSession` | Initial value in `index.ts` `session({initial})` |
| Add service to context | `types.ts` `BotContext.services` + `index.ts` createBot params | |
| Admin-only command | See `createBroadcastCommand` pattern | Pass `env` to command, check `TELEGRAM_BOT_ADMINS` |
| Keyboard/message helpers | `helpers.ts` | 232 lines — check before writing new UI code |

## CONVENTIONS
- Bot created fresh per request — never store state in module scope
- Middleware order in `index.ts`: session → env/services → chat-init/lang-sync → i18n → commands
- `ctx.services.*` — access repos and services inside handlers (never import directly)
- `ctx.t(key)` — i18n translation via middleware-injected helper
- Commands using Grammy `Composer` pattern — export a pre-wired composer, not a function
- Admin commands that need `env` are factory functions: `createXxxCommand(env)`

## ANTI-PATTERNS
- Never import repos or services directly in command files — use `ctx.services.*`
- Never add stateful module-level variables — workers are stateless
- `ConversationFlavor` is in context type but `@grammyjs/conversations` usage is minimal — check before adding conversation flows
