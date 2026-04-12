# BOT COMMANDS

## OVERVIEW
Eight command/handler files — user commands, admin commands, and callback query handler.

## COMMANDS
| File | Command | Type | Notes |
|------|---------|------|-------|
| `start.command.ts` | `/start`, `/help`, `/info`, `/settings` | User | Entry point, settings menu |
| `follow.command.ts` | `/follow <username>` | User | Subscribe to Twitch channel; 122 lines |
| `follows.command.ts` | `/follows`, `/unfollow` | User | Paginated list + unfollow |
| `live.command.ts` | `/live` | User | Show currently live channels; 102 lines |
| `broadcast.command.ts` | `/broadcast <msg>` | Admin | Mass message all chats |
| `change-channel-id.command.ts` | `/change_channel_id <old> <new>` | Admin | Migrate Twitch channel ID |
| `callback.handler.ts` | `callback_query` | — | Handles inline keyboard callbacks |
| `index.ts` | — | — | Re-exports all commands |

## CONVENTIONS
- Regular commands: export named `Composer` instance (e.g. `export const followCommand = new Composer<BotContext>(...)`)
- Admin commands: export factory `createXxxCommand(env: Env)` — need `env` for `TELEGRAM_BOT_ADMINS` check
- All commands registered in `../index.ts` via `bot.use()`
- Access repos/services via `ctx.services.*` only
- Translations via `ctx.t('key')` — keys must exist in all 3 locales

## ANTI-PATTERNS
- Never access `process.env` — use `ctx.env` (CF Workers binding)
- Never import DB/service classes directly — always `ctx.services.*`
