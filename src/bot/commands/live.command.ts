import { Composer } from 'grammy';
import type { BotContext } from '../types';

export const liveCommand = new Composer<BotContext>();

liveCommand.command('live', async (ctx) => {
  const chatId = ctx.chat?.id;
  if (!chatId) return;

  const chat = await ctx.services.chatRepo.findByChatId(chatId, 'telegram');
  if (!chat) return;

  const follows = await ctx.services.followRepo.findByChatId(chat.id);

  if (follows.length === 0) {
    await ctx.reply('You are not following any channels.');
    return;
  }

  // Get all followed channel IDs
  const channelIds: string[] = [];
  for (const follow of follows) {
    const channel = await ctx.services.channelRepo.findById(follow.channelId);
    if (channel) {
      channelIds.push(channel.channelId);
    }
  }

  if (channelIds.length === 0) {
    await ctx.reply('No channels found.');
    return;
  }

  // Get live streams
  const liveChannels: Array<{
    name: string;
    login: string;
    startedAt: Date;
    title: string;
    category: string;
    viewers: number;
  }> = [];

  for (const channelId of channelIds) {
    const stream = await ctx.services.twitch.getStreamByUserId(channelId);
    if (stream) {
      const user = await ctx.services.twitch.getUserById(channelId);
      if (user) {
        liveChannels.push({
          name: user.displayName,
          login: user.name,
          startedAt: stream.startDate,
          title: stream.title,
          category: stream.gameName,
          viewers: stream.viewers,
        });
      }
    }
  }

  if (liveChannels.length === 0) {
    await ctx.reply('No one is online.');
    return;
  }

  // Build message
  const messages: string[] = [];
  for (const channel of liveChannels) {
    const channelMessage: string[] = [];

    channelMessage.push(
      `🟢 <a href="https://twitch.tv/${channel.login}">${channel.name}</a> - ${channel.viewers} 👁️️`
    );

    if (channel.category) {
      channelMessage.push(`🎮 ${channel.category}`);
    }

    if (channel.title) {
      channelMessage.push(`📝 ${channel.title}`);
    }

    // Calculate uptime
    const uptime = Date.now() - channel.startedAt.getTime();
    const hours = Math.floor(uptime / 3600000);
    const minutes = Math.floor((uptime % 3600000) / 60000);
    const seconds = Math.floor((uptime % 60000) / 1000);

    let uptimeStr = '⌛ ';
    if (hours > 0) uptimeStr += `${hours}h `;
    if (minutes > 0) uptimeStr += `${minutes}m `;
    if (seconds > 0) uptimeStr += `${seconds}s `;

    channelMessage.push(uptimeStr);
    messages.push(channelMessage.join('\n'));
  }

  await ctx.reply(messages.join('\n\n'), {
    parse_mode: 'HTML',
    link_preview_options: { is_disabled: true },
  });
});
