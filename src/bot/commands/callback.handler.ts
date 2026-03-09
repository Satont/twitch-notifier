import { Composer } from 'grammy';
import type { BotContext } from '../types';
import type { SupportedLanguage } from '../../services/i18n.service';
import {
  sendSettingsMenu,
  sendLanguagePicker,
  handleToggleSetting,
  handleUnfollow,
  buildFollowsKeyboard
} from '../helpers';

export const callbackQueryHandler = new Composer<BotContext>();

callbackQueryHandler.on('callback_query:data', async (ctx) => {
  const data = ctx.callbackQuery.data;
  const chatId = ctx.chat?.id;
  if (!chatId) return;

  let chat = await ctx.services.chatRepo.findByChatId(chatId, 'telegram');
  if (!chat || !chat.settings) return;

  // Handle toggle settings
  if (data.startsWith('toggle_')) {
    await handleToggleSetting(ctx, data, chat);
    // Перезагрузить чат из БД чтобы получить актуальные настройки
    chat = await ctx.services.chatRepo.findByChatId(chatId, 'telegram');
    if (!chat) return;
    await sendSettingsMenu(ctx, chat);
  }

  // Handle language picker
  else if (data === 'language_picker') {
    await sendLanguagePicker(ctx);
  }

  // Handle language selection
  else if (data.startsWith('language_picker_set_')) {
    const lang = data.replace('language_picker_set_', '') as SupportedLanguage;
    if (ctx.services.i18n.isValidLocale(lang)) {
      await ctx.services.chatRepo.updateSettings(chat.settings.id, { language: lang });
      ctx.session.language = lang;
      
      // Обновляем ctx.t() для использования нового языка
      ctx.t = (key: string, params?: Record<string, any>) => {
        return ctx.services.i18n.t(lang, key, params);
      };
      
      // Перезагрузить чат из БД
      chat = await ctx.services.chatRepo.findByChatId(chatId, 'telegram');
      if (!chat) return;
      
      await ctx.answerCallbackQuery(
        ctx.services.i18n.t(lang, 'language.changed')
      );
      
      // Вернуться в главное меню с новым языком
      await sendSettingsMenu(ctx, chat);
    }
  }

  // Handle back to main menu
  else if (data === 'start_command_menu') {
    await sendSettingsMenu(ctx, chat);
  }

  // Handle unfollow
  else if (data.startsWith('channels_unfollow_')) {
    const channelId = data.replace('channels_unfollow_', '');
    await handleUnfollow(ctx, chat, channelId);
  }

  // Handle pagination
  else if (data === 'channels_unfollow_prev_page') {
    if (ctx.session.followsMenu && ctx.session.followsMenu.currentPage > 1) {
      ctx.session.followsMenu.currentPage--;
    }
    const keyboard = await buildFollowsKeyboard(ctx, chat.id);
    await ctx.editMessageReplyMarkup({ reply_markup: keyboard });
  }
  else if (data === 'channels_unfollow_next_page') {
    if (ctx.session.followsMenu && ctx.session.followsMenu.currentPage < ctx.session.followsMenu.totalPages) {
      ctx.session.followsMenu.currentPage++;
    }
    const keyboard = await buildFollowsKeyboard(ctx, chat.id);
    await ctx.editMessageReplyMarkup({ reply_markup: keyboard });
  }

  await ctx.answerCallbackQuery();
});
