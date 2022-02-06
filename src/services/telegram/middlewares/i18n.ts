import { NextFunction } from 'grammy'
import { i18n } from '../../../libs/i18n'
import { Context } from '../types'

export const middleware = async (ctx: Context, next: NextFunction) => {
  ctx.session.i18n = i18n.clone(ctx.session.entity.settings.language)

  await next()
}