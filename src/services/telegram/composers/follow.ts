import { Composer } from 'grammy'
import { followCommand } from '../../../commands/follow'
import { TelegramMessageSender } from '../MessageSender'
import { Context } from '../types'
import { StatelessQuestion } from '@grammyjs/stateless-question'
import { Chat } from '../../../entities/Chat'
import { I18n } from '../../../libs/i18n'
import { ReplyToMessageContext } from '@grammyjs/stateless-question/dist/source/identifier'


export const composer = new Composer<Context>()

const followQuestion = new StatelessQuestion('followQuestion', (ctx: ReplyToMessageContext<Context>) => {
  followHelper(ctx.session.entity, ctx.message.text, ctx.session.i18n)
})

composer.use(followQuestion.middleware())

composer.command('follow', (ctx) => {
  const channelName = ctx.match
  if (!channelName) return followQuestion.replyWithMarkdown(ctx, ctx.session.i18n.translate('scenes.follow.enter'))

  followHelper(ctx.session.entity, channelName, ctx.session.i18n)
})

async function followHelper(chat: Chat, channelName: string, i18n: I18n) {
  const { message } = await followCommand({ chat, channelName, i18n })
  
  await TelegramMessageSender.sendMessage({
    target: chat.chatId,
    message,
  })
}