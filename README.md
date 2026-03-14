# Twitch Notifier Bot

Telegram бот для уведомлений о стримах Twitch с использованием Cloudflare Workers, D1 и KV.

## Архитектура

### Serverless-Agnostic Design
Проект построен с учетом возможности запуска в разных окружениях:
- **Cloudflare Workers** (основная платформа)
- **Docker** (для локальной разработки)
- Другие serverless платформы (AWS Lambda, Vercel, etc.)

### Repository Pattern
```
src/db/
├── connection.ts                    # IDatabaseConnection интерфейс
├── repository.factory.ts            # Factory для создания репозиториев
├── repositories/
│   ├── interfaces/                  # Интерфейсы репозиториев
│   ├── drizzle/                     # Реализации для Drizzle ORM (D1, PostgreSQL)
│   └── cloudflare-kv/              # Реализации для Cloudflare KV
```

### Технологический стек
- **Runtime**: Cloudflare Workers (Node.js compatible)
- **Database**: Cloudflare D1 (SQLite)
- **Cache/Sessions**: Cloudflare KV
- **ORM**: Drizzle ORM
- **Bot Framework**: Grammy
- **HTTP Framework**: Hono
- **Twitch API**: Twurple

## Команды бота

### Пользовательские команды:
- `/start`, `/help`, `/info`, `/settings` - Меню настроек
- `/follow <username>` - Подписаться на канал Twitch
- `/follows`, `/unfollow` - Управление подписками
- `/live` - Показать онлайн стримы

### Админские команды:
- `/broadcast <message>` - Рассылка всем пользователям
- `/change_channel_id <old> <new>` - Обновить Twitch ID канала

## Установка и деплой

### 1. Установка зависимостей
```bash
bun install
```

### 2. Создание Cloudflare D1 базы данных
```bash
wrangler d1 create twitch-notifier-db
```

Скопируйте `database_id` из вывода команды и вставьте в `wrangler.toml`:
```toml
[[d1_databases]]
binding = "DB"
database_name = "twitch-notifier-db"
database_id = "YOUR_DATABASE_ID_HERE"
```

### 3. Создание Cloudflare KV namespace для сессий
```bash
wrangler kv:namespace create SESSIONS_KV
```

Скопируйте `id` из вывода команды и вставьте в `wrangler.toml`:
```toml
[[kv_namespaces]]
binding = "SESSIONS_KV"
id = "YOUR_KV_ID_HERE"
```

### 4. Применение миграций
```bash
wrangler d1 execute twitch-notifier-db --file=./drizzle/0000_init.sql
```

### 5. Настройка переменных окружения

**Через Cloudflare Dashboard** или с помощью `wrangler secret put`:

```bash
wrangler secret put TELEGRAM_TOKEN
wrangler secret put BASE_URL  # URL вашего воркера, например: https://twitch-notifier.yourname.workers.dev
```

Остальные переменные можно задать в `wrangler.toml`:
```toml
[vars]
TWITCH_CLIENT_ID = "your_client_id"
TWITCH_CLIENT_SECRET = "your_client_secret"
TELEGRAM_BOT_ADMINS = "123456789,987654321"  # Telegram user IDs через запятую
TWITCH_EVENTSUB_SECRET = "your_eventsub_secret"
```

### 6. Деплой
```bash
bun run deploy
```

### 7. Настройка Telegram webhook
После деплоя настройте webhook для бота:
```bash
curl -X POST "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://your-worker.workers.dev/telegram-webhook"}'
```

### 8. Настройка Twitch EventSub
Webhook для Twitch EventSub настроится автоматически при подписке на каналы через команду `/follow`.

URL для EventSub: `https://your-worker.workers.dev/twitch-webhook`

## Разработка

### Локальный запуск
```bash
bun run dev
```

### Генерация миграций
```bash
bun drizzle-kit generate
```

### Применение миграций локально
```bash
bun drizzle-kit migrate
```

### Проверка типов
```bash
bun run typecheck
```

## Структура проекта

```
src/
├── bot/
│   ├── commands/           # Команды через Composer
│   ├── helpers.ts          # Вспомогательные функции
│   ├── storage.ts          # Storage adapter для Grammy
│   └── types.ts            # Типы контекста
├── db/
│   ├── connection.ts       # Абстракция подключения к БД
│   ├── schema.ts           # Drizzle схема
│   ├── repository.factory.ts
│   └── repositories/
│       ├── interfaces/     # Интерфейсы репозиториев
│       ├── drizzle/        # Реализации для D1
│       └── cloudflare-kv/  # Реализации для KV
├── domain/
│   ├── models.ts           # Доменные модели
│   └── mapper.ts           # Маппер DB → Domain
├── services/               # Сервисы (Twitch, Telegram, etc.)
├── webhooks/              # Обработчики webhook'ов
└── index.ts               # Hono приложение
```

## Особенности реализации

### Персистентные сессии через Cloudflare KV
Сессии Grammy хранятся в Cloudflare KV с автоматическим TTL. Это решает проблему сброса сессий в serverless окружении. KV обеспечивает:
- Низкую латентность (читается с ближайшего edge)
- Автоматическое истечение ключей
- Глобальное распределение

### EventSub вместо polling
Используются Twitch EventSub webhooks для получения событий в реальном времени:
- `stream.online` - стример начал трансляцию
- `stream.offline` - стример закончил трансляцию
- `channel.update` - изменились название или категория

### Domain-Driven Design
Разделение между DB schema и domain models для чистой архитектуры.

### Factory Pattern
Единая точка создания репозиториев для простой замены реализаций.

## Лицензия

MIT
