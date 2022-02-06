import { GrammyError } from 'grammy'
import { getRepository } from 'typeorm'
import { Chat } from '../../entities/Chat'
import { chatOut, error } from '../../libs/logger'
import { SendMessageOpts } from '../_interface'
import TelegramService from './index'

export class TelegramMessageSender {
  static async sendMessage(opts: SendMessageOpts) {
    const targets = Array.isArray(opts.target) ? opts.target : [opts.target]
    for (const target of targets) {
      const log = () => chatOut(`TG [${target}]: ${opts.message}`.replace(/(\r\n|\n|\r)/gm, ' '))
      if (opts.image) {
        TelegramService.bot.api.sendPhoto(target, opts.image, {
          caption: opts.message,
          parse_mode: 'HTML',
        })
          .then(() => log())
          .catch((e) => this.catchError(e, target))
      } else {
        TelegramService.bot.api.sendMessage(target, opts.message, {
          disable_web_page_preview: true,
          parse_mode: 'HTML',
        })
          .then(() => log())
          .catch((e) => this.catchError(e, target))
      }
    }
  }

  private static async catchError(e: unknown, chatId: string | number) {
    error(e)
    if (e instanceof GrammyError) {
      if (e.error_code === 400 || e.error_code === 403) {
        const chat = await getRepository(Chat).findOne({ id: chatId.toString() })
        if (chat) await chat.remove()
      }
    }
  }
}