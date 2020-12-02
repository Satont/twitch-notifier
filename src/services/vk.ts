import { SendMessageOpts, ServiceInterface } from './_interface'
import { VK as VKIO, MessageContext, Keyboard  } from 'vk-io'
import { error, info, warning } from '../libs/logger'
import { Chat, Services } from '../entities/Chat'
import { getConnection } from 'typeorm'
import { command } from '../decorators/command'
import { followCommand } from '../commands/follow'
import { chunk } from 'lodash'
import { Languages } from '../entities/ChatSettings'
import { HearManager } from '@vk-io/hear'
import { cpuUsage } from 'process'

class VK extends ServiceInterface {
  service = Services.VK
  bot: VKIO = null
  hearManager: HearManager<MessageContext> = null

  async init() {
    const token = process.env.VK_GROUP_TOKEN
    if (!token) {
      warning('VK: group token not setuped, library will not works.')
      return
    }

    try {
      this.bot = new VKIO({ token })
      this.hearManager = new HearManager<MessageContext>()

      this.bot.updates.on('message_new', this.hearManager.middleware)
      this.bot.updates.on('message', async (ctx, next) => {
        await this.ensureUser(ctx)
        await this.listener(ctx)
        next()
      })
      this.hearManager.hear('/qwe', async (context, next) => {
        next()
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
    const data = { chatId: String(ctx.chatId || ctx.senderId), service: Services.VK }
    const chat = await repository.findOne(data) || repository.create({ ...data, settings: { language: Languages.RUSSIAN } })
    chat.save()

    ctx.ChatEntity = chat
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

  @command('start', { description: 'Start command' })
  async startCommand(ctx: MessageContext, args?: string[], arg?: string) {
    ctx.reply({
      message: 'menu',
      keyboard: Keyboard.builder()
        .textButton({
          label: 'The help',
          payload: {
            command: 'help',
          },
        })
        .row()
        .textButton({
          label: 'The current date',
          payload: {
            command: 'time',
          },
        })
        .row()
        .textButton({
          label: 'Cat photo',
          payload: {
            command: 'cat',
          },
          color: Keyboard.PRIMARY_COLOR,
        })
        .textButton({
          label: 'Cat purring',
          payload: {
            command: 'purr',
          },
          color: Keyboard.PRIMARY_COLOR,
        }),
    })
    return true
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
