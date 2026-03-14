import type { BotContext } from './types';
import type { Chat } from '../domain/models';
import { InlineKeyboard } from 'grammy';

export async function sendSettingsMenu(ctx: BotContext, chat: Chat) {
  const settings = chat.settings;
  if (!settings) return;

  const createCheckmark = (value: boolean) => value ? '✅' : '❌';

  const keyboard = new InlineKeyboard()
    .text(
      `${createCheckmark(settings.gameChangeNotification)} ${ctx.t('commands.start.game_change_notification_setting.button')}`,
      'toggle_game_change'
    ).row()
    .text(
      `${createCheckmark(settings.offlineNotification)} ${ctx.t('commands.start.offline_notification.button')}`,
      'toggle_offline'
    ).row()
    .text(
      `${createCheckmark(settings.titleChangeNotification)} ${ctx.t('commands.start.title_change_notification_setting.button')}`,
      'toggle_title_change'
    ).row()
    .text(
      `${createCheckmark(settings.gameAndTitleChangeNotification)} ${ctx.t('commands.start.game_and_title_change_notification_setting.button')}`,
      'toggle_game_and_title'
    ).row()
    .text(
      `${createCheckmark(settings.imageInNotification)} ${ctx.t('commands.start.image_in_notification_setting.button')}`,
      'toggle_image'
    ).row()
    .text(
      `🌐 ${ctx.t('commands.start.language.button')}`,
      'language_picker'
    ).row()
    .url('Github', 'https://github.com/Satont/twitch-notifier');

  const description = ctx.t('bot.description');

  if (ctx.callbackQuery) {
    try {
      await ctx.editMessageText(description, { reply_markup: keyboard });
    } catch (error: any) {
      // Игнорируем ошибку "message is not modified"
      if (!error?.message?.includes('message is not modified')) {
        throw error;
      }
    }
  } else {
    await ctx.reply(description, { reply_markup: keyboard });
  }
}

export async function sendLanguagePicker(ctx: BotContext) {
  const keyboard = new InlineKeyboard();

  const locales = ctx.services.i18n.getAvailableLocales();
  for (const locale of locales) {
    const emoji = ctx.services.i18n.t(locale, 'language.emoji');
    const name = ctx.services.i18n.t(locale, 'language.name');
    keyboard.text(`${emoji} ${name}`, `language_picker_set_${locale}`).row();
  }
  keyboard.text('«', 'start_command_menu');

  const text = ctx.t('language.select');

  if (ctx.callbackQuery) {
    await ctx.editMessageText(text, { reply_markup: keyboard });
  } else {
    await ctx.reply(text, { reply_markup: keyboard });
  }
}

export async function buildFollowsKeyboard(ctx: BotContext, chatId: string): Promise<InlineKeyboard> {
  const follows = await ctx.services.followRepo.findByChatId(chatId);
  const keyboard = new InlineKeyboard();

  for (const follow of follows) {
    const channel = await ctx.services.channelRepo.findById(follow.channelId);
    if (!channel) continue;

    const twitchUser = await ctx.services.twitch.getUserById(channel.channelId);
    if (!twitchUser) continue;

    keyboard.text(twitchUser.displayName, `channels_unfollow_${channel.channelId}`).row();
  }

  // Add pagination buttons if needed
  if (ctx.session.followsMenu) {
    const { currentPage, totalPages } = ctx.session.followsMenu;
    if (totalPages > 1) {
      keyboard.text('«', 'channels_unfollow_prev_page');
      keyboard.text('»', 'channels_unfollow_next_page');
    }
  }

  return keyboard;
}

export async function handleToggleSetting(ctx: BotContext, data: string, chat: Chat) {
  const chatId = ctx.chat?.id;
  if (!chatId || !chat.settings) return;

  const updates: any = {};

  switch (data) {
    case 'toggle_game_change':
      updates.gameChangeNotification = !chat.settings.gameChangeNotification;
      chat.settings.gameChangeNotification = updates.gameChangeNotification;
      
      // Если включили game change, а title change тоже включен, то включаем game_and_title
      if (updates.gameChangeNotification && chat.settings.titleChangeNotification) {
        updates.gameAndTitleChangeNotification = true;
        chat.settings.gameAndTitleChangeNotification = true;
      }
      // Если выключили game change, то выключаем game_and_title
      if (!updates.gameChangeNotification) {
        updates.gameAndTitleChangeNotification = false;
        chat.settings.gameAndTitleChangeNotification = false;
      }
      break;
      
    case 'toggle_offline':
      updates.offlineNotification = !chat.settings.offlineNotification;
      chat.settings.offlineNotification = updates.offlineNotification;
      break;
      
    case 'toggle_title_change':
      updates.titleChangeNotification = !chat.settings.titleChangeNotification;
      chat.settings.titleChangeNotification = updates.titleChangeNotification;
      
      // Если включили title change, а game change тоже включен, то включаем game_and_title
      if (updates.titleChangeNotification && chat.settings.gameChangeNotification) {
        updates.gameAndTitleChangeNotification = true;
        chat.settings.gameAndTitleChangeNotification = true;
      }
      // Если выключили title change, то выключаем game_and_title
      if (!updates.titleChangeNotification) {
        updates.gameAndTitleChangeNotification = false;
        chat.settings.gameAndTitleChangeNotification = false;
      }
      break;
      
    case 'toggle_game_and_title':
      updates.gameAndTitleChangeNotification = !chat.settings.gameAndTitleChangeNotification;
      chat.settings.gameAndTitleChangeNotification = updates.gameAndTitleChangeNotification;
      
      // Если включили game_and_title, включаем оба
      if (updates.gameAndTitleChangeNotification) {
        updates.gameChangeNotification = true;
        updates.titleChangeNotification = true;
        chat.settings.gameChangeNotification = true;
        chat.settings.titleChangeNotification = true;
      }
      // Если выключили game_and_title, выключаем оба
      else {
        updates.gameChangeNotification = false;
        updates.titleChangeNotification = false;
        chat.settings.gameChangeNotification = false;
        chat.settings.titleChangeNotification = false;
      }
      break;
      
    case 'toggle_image':
      updates.imageInNotification = !chat.settings.imageInNotification;
      chat.settings.imageInNotification = updates.imageInNotification;
      break;
  }

  if (Object.keys(updates).length > 0) {
    await ctx.services.chatRepo.updateSettings(chat.id, updates);
  }
}

export async function handleUnfollow(ctx: BotContext, chat: Chat, channelIdFromCallback: string) {
  const channel = await ctx.services.channelRepo.findById(channelIdFromCallback);
  if (!channel) {
    await ctx.answerCallbackQuery('Channel not found');
    return;
  }

  const follow = await ctx.services.followRepo.findByChatAndChannel(chat.id, channel.id);
  if (!follow) {
    await ctx.answerCallbackQuery('Already unfollowed');
    return;
  }

  const twitchUser = await ctx.services.twitch.getUserById(channel.channelId);
  const streamerName = twitchUser?.displayName || channel.channelId;

  await ctx.services.followRepo.delete(follow.id);

  // Check if this channel still has followers
  const remainingFollows = await ctx.services.followRepo.findByChannelId(channel.id);

  // If no followers remain, unsubscribe from EventSub
  if (remainingFollows.length === 0) {
    try {
      await ctx.services.eventsub.unsubscribeFromChannel(channel.channelId);
      console.log(`Unsubscribed from EventSub for channel ${channel.channelId}`);
    } catch (error) {
      console.error(`Failed to unsubscribe from EventSub for ${channel.channelId}:`, error);
      // Don't fail the unfollow if EventSub unsubscription fails
    }
  }

  await ctx.answerCallbackQuery(
    ctx.t('commands.unfollow.success', {
      streamer: streamerName,
    })
  );

  // Update keyboard
  const totalFollows = await ctx.services.followRepo.countByChatId(chat.id);

  if (totalFollows === 0) {
    await ctx.editMessageText('You are not following any channels.');
    await ctx.editMessageReplyMarkup({ reply_markup: new InlineKeyboard() });
    return;
  }

  const keyboard = await buildFollowsKeyboard(ctx, chat.id);

  await ctx.editMessageText(
    ctx.t('commands.follows.total', {
      count: totalFollows.toString(),
    }),
    {
      reply_markup: keyboard,
    }
  );
}
