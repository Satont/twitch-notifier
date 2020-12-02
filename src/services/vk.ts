import { SendMessageOpts, ServiceInterface } from './_interface'
import { VK as VKIO, MessageContext  } from 'vk-io'
import { error, info, warning } from '../libs/logger'
import { Chat, Services } from '../entities/Chat'
import { getConnection } from 'typeorm'
import { command } from '../decorators/command'
import { followCommand } from '../commands/follow'
import { chunk } from 'lodash'

class VK extends ServiceInterface {
  service = Services.VK
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
      this.inited = true
    } catch (e) {
      error(e)
    }
  }

  async ensureUser(ctx: MessageContext) {
    const repository = getConnection().getRepository(Chat)
    const data = { id: String(ctx.chatId), service: Services.VK }
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

  async sendMessage(opts: SendMessageOpts) {
    const targets = Array.isArray(opts.target) ? opts.target : [opts.target]
    const chunks = chunk(targets.map(t => Number(t)), 100)
    const attachment = opts.image ? this.uploadPhoto(opts.image) : undefined
    for (const chunk of chunks) {
      await this.bot.api.messages.send({
        random_id: Math.random() * (1000000000 - 9) + 10,
        user_ids: chunk,
        message: opts.message,
        dont_parse_links: true,
        attachment,
      })
    }
  }

  public async uploadPhoto(source: string) {
    return await this.bot.upload.messagePhoto({
      source: {
        value: source,
      },
    })
  }
}

export default new VK()
