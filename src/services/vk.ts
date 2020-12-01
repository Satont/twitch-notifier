import { ServiceInterface } from './_interface'
import { VK as VKIO, MessageContext  } from 'vk-io'
import { error, info, warning } from '../libs/logger'
import { Chat } from '../entities/Chat'
import { getConnection } from 'typeorm'
import { command } from '../decorators/command'
import { followCommand } from '../commands/follow'

class VK extends ServiceInterface {
  service = 'vk'
  bot: VKIO = null

  async init() {
    const token = process.env.VK_GROUP_TOKEN
    if (!token) {
      warning('VK: group token not setuped, telegram library will not works.')
      return
    }

    try {
      this.bot = new VKIO({ token })

      this.bot.updates.on('message_new', async (msg) => {
        await this.ensureUser(msg)
        await this.listener(msg)
      })
      await this.bot.updates.start()
      info('VK Service initialized.')
    } catch (e) {
      error(e)
    }
  }

  async ensureUser(ctx: MessageContext) {
    const repository = getConnection().getRepository(Chat)
    const data = { id: String(ctx.chatId) }
    const user = await repository.findOne(data) || await repository.create(data).save()

    ctx.ChatEntity = user
  }

  async listener(msg: MessageContext) {
    if (!msg.hasText || !msg.text?.startsWith('!')) return
    const commandName = msg.text.substring(1).split(' ')[0]
    const args = msg.text.split(' ').slice(1)
    const arg = msg.text.substring(1).replace(commandName, '')

    const command = this.commands.find(c => c.name === commandName)
    if (!command) return

    await this[command.fnc](msg, args, arg)
    return true
  }

  @command('follow')
  async follow(msg: MessageContext, args?: string[], arg?: string) {
    msg.reply(await followCommand({ chat: msg.ChatEntity, channelName: arg }))
  }
}

export default new VK()
