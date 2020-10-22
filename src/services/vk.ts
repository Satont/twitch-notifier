import { ServiceInterface } from './_interface'
import { VK as VKIO, MessageContext  } from 'vk-io'
import { error, info } from '../libs/logger'
import { Chat } from '../entities/Chat'
import { orm } from '../libs/db'

export default new class VK extends ServiceInterface {
  service = 'vk'
  bot: VKIO = null
  commands = [
    { name: 'follow', fnc: this.follow },
  ]

  async init() {
    const token = process.env.VK_GROUP_TOKEN
    if (!token) return false

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
    const repository = orm.em.fork().getRepository(Chat)
    const data = { chatId: ctx.chatId }
    const user = await repository.findOne(data) || repository.assign(new Chat(), data)
    await repository.persistAndFlush(user)
    
    ctx.ChatEntity = user
  }

  async listener(msg: MessageContext) {
    if (!msg.hasText || !msg.text?.startsWith('!')) return
    const commandName = msg.text.substring(1).split(' ')[0]
    const args = msg.text.split(' ').slice(1)
    const arg = msg.text.substring(1).replace(commandName, '')

    const command = this.commands.find(c => c.name === commandName)
    if (!command) return

    command['fnc'].call(VK, msg, args, arg)
    return true
  }

  async follow(msg: MessageContext, args?: string[], arg?: string) {
    return true
  }
}