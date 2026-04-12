import { Hono } from "hono";
import { drizzle } from "drizzle-orm/d1";
import type { Env } from "./types";
import { createBot } from "./bot";
import { TelegramService, I18nService, TwitchService, EventSubService } from "~/services";
import { CloudflareD1Connection } from "./db/connection";
import { DrizzleRepositoryFactory } from "./db/repository.factory";
import { D1SessionRepository } from "./db/repositories/drizzle/session.d1.repository";
import { handleTwitchWebhook } from "./webhooks/twitch";

const app = new Hono<{ Bindings: Env }>();

// Health check
app.get("/", (c) => {
	return c.json({ status: "ok", service: "twitch-notifier" });
});

// Telegram webhook endpoint
app.post("/telegram-webhook", async (c) => {
	const env = c.env;

	// Create database connection (serverless-agnostic)
	const dbClient = drizzle(env.twitch_notifier_db);
	const dbConnection = new CloudflareD1Connection(dbClient);

	// Create repository factory
	const repositoryFactory = new DrizzleRepositoryFactory(dbConnection);

	// Create repositories
	const chatRepo = repositoryFactory.createChatRepository();
	const channelRepo = repositoryFactory.createChannelRepository();
	const followRepo = repositoryFactory.createFollowRepository();
	const streamRepo = repositoryFactory.createStreamRepository();

	// Create session repository using D1
	const sessionRepo = new D1SessionRepository(dbClient);

	// Initialize services
	const i18nService = new I18nService();
	await i18nService.init(); // Initialize i18next
	const twitchService = new TwitchService(env);
	const telegramService = new TelegramService(env, i18nService);
	const eventSubService = new EventSubService(twitchService.getApiClient(), env, env.BASE_URL);

	// Create bot instance
	const bot = createBot(env, {
		i18n: i18nService,
		twitch: twitchService,
		eventsub: eventSubService,
		chatRepo,
		channelRepo,
		followRepo,
		sessionRepo,
	});

	const update = await c.req.json();
	const updateId = typeof update?.update_id === "number" ? update.update_id.toString() : undefined;

	if (updateId) {
		const updateKey = `telegram_update:${updateId}`;
		const seen = await env.twitch_notifier_kv.get(updateKey);
		if (seen) {
			return new Response("OK", { status: 200 });
		}

		await env.twitch_notifier_kv.put(updateKey, "1", { expirationTtl: 60 * 60 });
	}

	c.executionCtx.waitUntil(
		bot.handleUpdate(update).catch((error) => {
			console.error("Bot handler error:", error);
		}),
	);

	return new Response("OK", { status: 200 });
});

// Twitch EventSub webhook endpoint
app.post("/twitch-webhook", async (c) => {
	const env = c.env;
	const db = drizzle(env.twitch_notifier_db);

	return await handleTwitchWebhook(c.req.raw, env, db, c.executionCtx);
});

console.log("App initialized");

export default app;
