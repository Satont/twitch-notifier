import { Composer } from 'grammy';
import type { BotContext } from '../types';

export const followCommand = new Composer<BotContext>();

followCommand.command('follow', async (ctx) => {
  const text = ctx.message?.text?.replace('/follow', '').trim();

  if (!text) {
    await ctx.reply(
      ctx.t('commands.follow.enter')
    );
    ctx.session.scene = 'follow';
    return;
  }

  await handleFollow(ctx, text);
});

// Handle follow scene
followCommand.on('message:text', async (ctx, next) => {
  if (ctx.session.scene === 'follow') {
    await handleFollow(ctx, ctx.message.text);
    ctx.session.scene = undefined;
    return;
  }
  await next();
});

async function handleFollow(ctx: BotContext, text: string) {
  const chatId = ctx.chat?.id;
  if (!chatId) return;

  const chat = await ctx.services.chatRepo.findByChatId(chatId, 'telegram');
  if (!chat) return;

  // Extract Twitch username from text or URL
  const twitchLinkRegex = /(?:https?:\/\/)?(?:www\.)?twitch\.tv\/(\w+)/g;
  const matches = Array.from(text.matchAll(twitchLinkRegex));

  const usernames = matches.length > 0
    ? matches.map(m => m[1])
    : [text.trim()];

  const results: string[] = [];

  for (const username of usernames) {
    // Validate username
    if (!/^[a-zA-Z0-9_]{3,25}$/.test(username)) {
      results.push(
        ctx.t(
          'commands.follow.errors.badUsername',
          { streamer: username }
        )
      );
      continue;
    }

    try {
      // Get Twitch user
      const twitchUser = await ctx.services.twitch.getUserByLogin(username);

      if (!twitchUser) {
        results.push(
          ctx.t(
            'commands.follow.errors.streamerNotFound',
            { streamer: username }
          )
        );
        continue;
      }

      // Get or create channel
      let channel = await ctx.services.channelRepo.findByChannelId(twitchUser.id, 'twitch');
      if (!channel) {
        channel = await ctx.services.channelRepo.create(twitchUser.id, 'twitch');
      }

      // Create follow
      try {
        await ctx.services.followRepo.create(chat.id, channel.id);

        // Subscribe to EventSub events for this channel
        // Check if we already have subscriptions for this channel
        const hasSubscriptions = await ctx.services.eventsub.hasActiveSubscriptions(twitchUser.id);
        if (!hasSubscriptions) {
          try {
            await ctx.services.eventsub.subscribeToChannel(twitchUser.id);
            console.log(`Subscribed to EventSub for channel ${twitchUser.id}`);
          } catch (eventSubError) {
            console.error(`Failed to subscribe to EventSub for ${twitchUser.id}:`, eventSubError);
            // Don't fail the follow if EventSub subscription fails
          }
        }

        results.push(
          ctx.t(
            'commands.follow.success',
            { streamer: username }
          )
        );
      } catch (error: any) {
        // Check if it's FollowAlreadyExistsError by checking error name or message
        if (error.constructor.name === 'FollowAlreadyExistsError' || error.message === 'Follow already exists') {
          results.push(
            ctx.t(
              'commands.follow.errors.alreadyFollowed',
              { streamer: username }
            )
          );
        } else {
          throw error;
        }
      }
    } catch (error) {
      console.error('Error following user:', error);
      results.push(`${username} - internal error`);
    }
  }

  await ctx.reply(results.join('\n'));
}
