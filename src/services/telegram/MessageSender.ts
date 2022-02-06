import { chatOut } from '../../libs/logger'
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
          .catch(console.error)
      } else {
        TelegramService.bot.api.sendMessage(target, opts.message, {
          disable_web_page_preview: true,
          parse_mode: 'HTML',
        })
          .then(() => log())
          .catch(console.error)
      }
    }
  }
}