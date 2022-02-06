import { NextFunction } from 'grammy'
import { getRepository } from 'typeorm'
import { Chat, Services } from '../../../entities/Chat'
import { ChatSettings } from '../../../entities/ChatSettings'
import { Context } from '../types'

export const middleware = async (ctx: Context, next: NextFunction) => {
  if (!ctx.from?.id) return
  const repository = getRepository(Chat)

  const data = { chatId: String(ctx.from?.id), service: Services.TELEGRAM }
  const chat = await repository.findOne(data, { relations: ['follows', 'follows.channel'] })
    ?? repository.create({ ...data, settings: new ChatSettings() })
  await chat.save()

  ctx.session.entity = chat
  return await next()
}