import { Composer } from 'grammy';
import type { BotContext } from '../types';
import { buildFollowsKeyboard } from '../helpers';

export const followsCommand = new Composer<BotContext>();

followsCommand.command(['follows', 'unfollow'], async (ctx) => {
  const chatId = ctx.chat?.id;
  if (!chatId) return;

  const chat = await ctx.services.chatRepo.findByChatId(chatId, 'telegram');
  if (!chat) return;

  ctx.session.followsMenu = {
    currentPage: 1,
    totalPages: 1,
  };

  const totalFollows = await ctx.services.followRepo.countByChatId(chat.id);

  if (totalFollows === 0) {
    await ctx.reply('You are not following any channels.');
    return;
  }

  const keyboard = await buildFollowsKeyboard(ctx, chat.id);

  await ctx.reply(
    ctx.t(
      'commands.follows.total',
      { count: totalFollows.toString() }
    ),
    {
      reply_markup: keyboard,
    }
  );
});
