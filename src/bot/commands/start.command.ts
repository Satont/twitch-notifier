import { Composer } from 'grammy';
import type { BotContext } from '../types';
import type { SupportedLanguage } from '../../services/i18n.service';
import { sendSettingsMenu } from '../helpers';

export const startCommand = new Composer<BotContext>();

startCommand.command(['start', 'help', 'info', 'settings'], async (ctx) => {
  const chatId = ctx.chat?.id;
  if (!chatId) return;

  // Get or create chat in database
  let chat = await ctx.services.chatRepo.findByChatId(chatId, 'telegram');
  if (!chat) {
    await ctx.services.chatRepo.create(chatId.toString(), 'telegram');
    // Fetch the chat again to get it with settings
    chat = await ctx.services.chatRepo.findByChatId(chatId, 'telegram');
  }

  // Update session language
  if (chat?.settings) {
    ctx.session.language = chat.settings.language as SupportedLanguage;
  }

  if (chat) {
    await sendSettingsMenu(ctx, chat);
  }
});
