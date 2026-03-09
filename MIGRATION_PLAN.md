# План миграции на Cloudflare Workers

## Обзор проекта

**Текущий стек:**
- Go 1.19+
- PostgreSQL + Ent ORM
- Polling Telegram Bot
- Standalone приложение с собственным воркером проверки стримов

**Целевой стек:**
- TypeScript/JavaScript
- Cloudflare Workers + D1 (SQLite)
- Grammy Bot (Telegram)
- Drizzle ORM с паттерном репозиториев
- pnpm как пакетный менеджер

## Существующие фичи для сохранения

### 1. База данных (5 таблиц)

#### Chat
- `id` (UUID)
- `chat_id` (string) - ID чата в Telegram
- `service` (enum: "telegram")
- Уникальный индекс: `chat_id + service`

#### ChatSettings
- `id` (UUID)
- `chat_id` (UUID FK -> Chat)
- `game_change_notification` (boolean, default: true)
- `title_change_notification` (boolean, default: false)
- `game_and_title_change_notification` (boolean, default: false)
- `offline_notification` (boolean, default: true)
- `image_in_notification` (boolean, default: true)
- `chat_language` (enum: "ru", "en", "uk", default: "en")

#### Channel
- `id` (UUID)
- `channel_id` (string) - ID канала на Twitch
- `service` (enum: "twitch")
- `is_live` (boolean, default: false)
- `title` (string, nullable)
- `category` (string, nullable)
- `updated_at` (timestamp)
- Уникальный индекс: `channel_id + service`

#### Follow
- `id` (UUID)
- `channel_id` (UUID FK -> Channel)
- `chat_id` (UUID FK -> Chat)
- Уникальный индекс: `channel_id + chat_id`

#### Stream
- `id` (string, unique) - ID стрима от Twitch
- `channel_id` (UUID FK -> Channel)
- `titles` (string[], default: [])
- `categories` (string[], default: [])
- `started_at` (timestamp)
- `updated_at` (timestamp)
- `ended_at` (timestamp, nullable)

### 2. Telegram команды

#### Пользовательские команды:
- `/start` (aliases: /help, /info, /settings) - главное меню с настройками
- `/follow <username|url>` - подписка на Twitch канал
- `/follows` (alias: /unfollow) - список подписок с возможностью отписки (пагинация)
- `/live` - список онлайн стримов из подписок

#### Административные команды:
- `/broadcast <message>` - массовая рассылка
- `/change_channel_id <old_id> <new_id>` - изменение ID канала

#### Интерактивные элементы:
- Callback buttons для настроек (в /start)
- Callback buttons для отписки (в /follows)
- Кнопки пагинации для списка подписок
- Выбор языка через inline-кнопки

### 3. Система уведомлений

Проверка стримов каждую минуту (в dev - 10 секунд) с уведомлениями:

#### При запуске стрима:
- Сообщение с названием, категорией, стримером
- Превью (thumbnail) если включено
- Кнопка отписки

#### При завершении стрима:
- Сообщение о завершении
- Список категорий за стрим
- Длительность стрима
- Кнопка отписки

#### Изменения во время стрима:
- Смена категории (если включено)
- Смена названия (если включено)
- Смена категории И названия одновременно (если включено)

### 4. Интернационализация (i18n)

Поддержка языков:
- Русский (ru)
- Английский (en)
- Украинский (uk)

Переводы хранятся в директории `locales/`

### 5. Twitch API интеграция

- Получение информации о каналах
- Получение информации о стримах
- Батчинг запросов (chunked requests)
- OAuth авторизация с автоматическим обновлением токена

### 6. Конфигурация

Переменные окружения:
- `TWITCH_CLIENTID` - ID приложения Twitch
- `TWITCH_CLIENTSECRET` - Secret приложения Twitch
- `TELEGRAM_TOKEN` - токен Telegram бота
- `TELEGRAM_BOT_ADMINS` - список ID администраторов (через запятую)
- `DATABASE_URL` - URL базы данных (старый PostgreSQL)
- `SENTRY_DSN` - (опционально) для мониторинга ошибок

## Архитектура Cloudflare Workers решения

### Workers

#### 1. **bot-worker** (основной)
- Обработка Telegram Webhook
- Обработка всех команд
- Grammy bot + conversations для multi-step команд
- Использует KV для хранения session данных

#### 2. **streams-checker-worker** (Cron Worker)
- Запускается каждую минуту (Cron Trigger)
- Проверяет статус стримов через Twitch API
- Отправляет уведомления через Telegram API
- Использует D1 для чтения/записи данных

### Cloudflare сервисы

- **D1** - SQLite база данных для всех таблиц
- **KV** (опционально) - для session storage Grammy
- **Cron Triggers** - для периодической проверки стримов
- **Workers Analytics** (опционально) - для мониторинга

## Структура проекта

```
twitch-notifier/
├── src/
│   ├── bot/                      # Telegram bot
│   │   ├── index.ts             # Entry point для bot-worker
│   │   ├── bot.ts               # Grammy bot инициализация
│   │   ├── commands/            # Команды
│   │   │   ├── start.ts
│   │   │   ├── follow.ts
│   │   │   ├── follows.ts
│   │   │   ├── live.ts
│   │   │   ├── broadcast.ts
│   │   │   └── change-channel-id.ts
│   │   ├── keyboards/           # Inline клавиатуры
│   │   │   ├── settings.ts
│   │   │   ├── language.ts
│   │   │   └── follows.ts
│   │   └── middlewares/         # Миддлвары
│   │       ├── logger.ts
│   │       ├── admin.ts
│   │       └── chat.ts
│   ├── checker/                 # Streams checker
│   │   ├── index.ts            # Entry point для checker-worker
│   │   └── checker.ts          # Логика проверки стримов
│   ├── db/                      # База данных
│   │   ├── schema.ts           # Drizzle схемы
│   │   ├── migrations/         # SQL миграции
│   │   └── repositories/       # Паттерн репозиториев
│   │       ├── chat.repository.ts
│   │       ├── channel.repository.ts
│   │       ├── follow.repository.ts
│   │       └── stream.repository.ts
│   ├── services/                # Сервисы
│   │   ├── twitch.service.ts   # Twitch API клиент
│   │   ├── telegram.service.ts # Telegram message sender
│   │   └── i18n.service.ts     # Интернационализация
│   ├── types/                   # TypeScript типы
│   │   ├── env.ts
│   │   └── index.ts
│   └── utils/                   # Утилиты
│       ├── thumbnail.ts
│       └── helpers.ts
├── locales/                     # Переводы (скопировать из Go проекта)
│   ├── en.json
│   ├── ru.json
│   └── uk.json
├── migrations/                  # Скрипты миграции
│   └── migrate-users.ts        # Скрипт миграции из PostgreSQL в D1
├── drizzle.config.ts           # Конфигурация Drizzle
├── wrangler.toml               # Конфигурация Cloudflare Workers
├── package.json
├── tsconfig.json
└── README.md
```

## План реализации

### Этап 1: Настройка проекта ✓
1. Инициализация проекта с pnpm
2. Установка зависимостей:
   - `wrangler` - Cloudflare CLI
   - `grammy` + `@grammyjs/conversations` - Telegram bot
   - `drizzle-orm` + `drizzle-kit` - ORM
   - Другие зависимости
3. Создание структуры директорий
4. Настройка TypeScript
5. Настройка wrangler.toml для обоих workers

### Этап 2: База данных ✓
1. Создание Drizzle схем на основе Ent схем
2. Имплементация паттерна репозиториев
3. Создание D1 базы через Wrangler
4. Генерация и применение миграций

### Этап 3: Сервисы ✓
1. Twitch API клиент с OAuth
2. Telegram message sender
3. i18n сервис (адаптация с Go проекта)
4. Thumbnail builder

### Этап 4: Telegram Bot ✓
1. Настройка Grammy bot
2. Имплементация всех команд
3. Создание inline клавиатур
4. Настройка миддлваров (логгирование, admin check, chat persistence)
5. Настройка conversations для multi-step команд (/follow)

### Этап 5: Streams Checker Worker ✓
1. Имплементация логики проверки стримов
2. Отправка уведомлений
3. Настройка Cron Trigger

### Этап 6: Деплой и тестирование ✓
1. Деплой bot-worker
2. Деплой streams-checker-worker
3. Настройка Telegram Webhook
4. Тестирование всех команд
5. Тестирование уведомлений

### Этап 7: Миграция данных ✓
1. Создание скрипта миграции из PostgreSQL в D1
2. Тестовая миграция
3. Продакшн миграция

## Скрипт миграции данных

Создать отдельный скрипт `migrations/migrate-users.ts` который:
1. Подключается к старой PostgreSQL базе
2. Читает все данные из таблиц (Chat, ChatSettings, Channel, Follow, Stream)
3. Трансформирует данные если нужно
4. Записывает в D1 через Wrangler API или D1 HTTP API

Особенности миграции:
- UUID в PostgreSQL -> сохраняются как есть (D1 поддерживает текстовые UUID)
- Массивы (titles, categories) -> JSON в SQLite
- Timestamps -> ISO 8601 строки в SQLite
- Enum values -> остаются как есть

## Отличия от Go версии

### Архитектурные:
- **Polling -> Webhook**: Cloudflare Workers работает по модели request/response, используем Telegram Webhook вместо Long Polling
- **Отдельный воркер для проверки стримов**: вместо горутины - отдельный Worker с Cron Trigger
- **Serverless**: нет постоянно запущенного процесса, оплата за выполнение
- **SQLite вместо PostgreSQL**: D1 - управляемая SQLite база

### Технические:
- **Grammy вместо go-tg**: официальная TypeScript библиотека для Telegram
- **Drizzle вместо Ent**: type-safe ORM для TypeScript
- **Паттерн репозиториев**: изоляция логики БД для простой замены ORM в будущем

## Следующие шаги

После прочтения этого плана:

1. Убедитесь что у вас есть:
   - Аккаунт Cloudflare (Workers Paid plan для Cron Triggers)
   - Доступ к текущей PostgreSQL базе для миграции

2. Дайте подтверждение для начала реализации:
   ```
   Да, начинаем! Начни с Этапа 1.
   ```

3. Я буду реализовывать каждый этап последовательно, показывая прогресс

## Важные замечания

- Cloudflare Workers Free tier имеет лимиты (100k запросов/день)
- Cron Triggers требуют Workers Paid plan ($5/месяц)
- D1 пока в beta, но стабильна для продакшн использования
- Webhook требует HTTPS домен (можно использовать workers.dev)
- Grammy conversations требуют хранилище для session (используем KV или D1)

## Вопросы для уточнения

1. Есть ли у вас уже Cloudflare аккаунт?
2. Сколько пользователей у текущего бота? (для оценки нагрузки)
3. Нужна ли интеграция с Sentry для мониторинга?
4. Хотите ли сохранить историю стримов (таблица Stream) или только активные?

---

**Готовы начать миграцию?** Дайте команду и я начну с Этапа 1!
