import { Composer } from 'grammy';
import type { BotContext } from '../types';
import { sendSettingsMenu } from '../helpers';

export const startCommand = new Composer<BotContext>();

startCommand.command(['start', 'help', 'info', 'settings'], async (ctx) => {
  const chatId = ctx.chat?.id;
  if (!chatId) return;

  const chat = await ctx.services.chatRepo.findByChatId(chatId, 'telegram');
  if (chat) {
    await sendSettingsMenu(ctx, chat);
  }
});
