import { Composer } from 'grammy';
import type { BotContext } from '../types';
import type { Env } from '../../types/env';

export function createBroadcastCommand(env: Env) {
  const broadcast = new Composer<BotContext>();

  const isAdmin = (userId: number): boolean => {
    const admins = env.TELEGRAM_BOT_ADMINS.split(',').map(id => parseInt(id.trim()));
    return admins.includes(userId);
  };

  broadcast.command('broadcast', async (ctx) => {
    const userId = ctx.from?.id;
    if (!userId || !isAdmin(userId)) {
      return;
    }

    const text = ctx.message?.text?.replace('/broadcast', '').trim();
    if (!text) {
      await ctx.reply('Usage: /broadcast <message>');
      return;
    }

    // Get all chats (only positive IDs = private chats/groups)
    const allChats = await ctx.services.chatRepo.findAllByService('telegram');

    let sent = 0;
    let failed = 0;

    for (const chat of allChats) {
      const chatIdNum = parseInt(chat.chatId);
      if (chatIdNum <= 0) continue; // Skip channels/supergroups

      try {
        await ctx.api.sendMessage(chatIdNum, text);
        sent++;
      } catch (error) {
        console.error(`Failed to send to ${chat.chatId}:`, error);
        failed++;
      }
    }

    await ctx.reply(`Broadcast completed!\nSent: ${sent}\nFailed: ${failed}`);
  });

  return broadcast;
}
