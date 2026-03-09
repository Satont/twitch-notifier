import { Composer } from 'grammy';
import type { BotContext } from '../types';
import type { Env } from '../../types/env';

export function createChangeChannelIdCommand(env: Env) {
  const changeChannelId = new Composer<BotContext>();

  const isAdmin = (userId: number): boolean => {
    const admins = env.TELEGRAM_BOT_ADMINS.split(',').map(id => parseInt(id.trim()));
    return admins.includes(userId);
  };

  changeChannelId.command('change_channel_id', async (ctx) => {
    const userId = ctx.from?.id;
    if (!userId || !isAdmin(userId)) {
      return;
    }

    const text = ctx.message?.text?.replace('/change_channel_id', '').trim();

    if (!text) {
      await ctx.reply('Usage: /change_channel_id <old_id> <new_id>');
      return;
    }

    const parts = text.split(' ');

    if (parts.length !== 2) {
      await ctx.reply('Usage: /change_channel_id <old_id> <new_id>');
      return;
    }

    const [oldId, newId] = parts;

    try {
      await ctx.services.channelRepo.updateChannelId(oldId, newId, 'twitch');
      await ctx.reply('Channel ID updated successfully!');
    } catch (error) {
      console.error('Error updating channel ID:', error);
      await ctx.reply('Error updating channel ID.');
    }
  });

  return changeChannelId;
}
