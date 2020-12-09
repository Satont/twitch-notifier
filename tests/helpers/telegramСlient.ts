import httpMocks from 'node-mocks-http'
import { EventEmitter } from 'events'
import createTelegramNock from './createTelegramNock'

const getRandomId = () => Math.floor(Math.random() * Math.floor(Math.pow(10, 10)))

export class TestTelegramClient {
  received: any[]
  clientId: any
  updateId: number
  client: import('../../src/services/telegram').Telegram

  async init() {
    const ee = new EventEmitter() 
    await createTelegramNock(ee)

    const { Telegram } = await import('../../src/services/telegram')
    this.client = new Telegram('123456')
    await this.client.init()

    this.received = []
    this.clientId = getRandomId()
    this.updateId = 0
    ee.on('sendMessage', data => {
      if (data.chat_id != this.clientId) return
      this.received.push({

        data,
      })
    })
  }

  async send(params: {
    message?: { text: string, message_id: number, chat: any, date: number, from?: any },
  }) {
    const response = httpMocks.createResponse()
    const update = {
      update_id: ++this.updateId,
      ...params,
    }

    if (update.message) {
      update.message.date = Date.now()
      update.message.from = { id: 42, is_bot: false, username: 'Test', first_name: 'User' }
    }

    if (update.message && !update.message.message_id) {
      update.message.message_id = this.updateId
    }

    if (update.message && !update.message.chat) {
      update.message.chat = {
        id: this.clientId,
      }
    }

    await this.client.bot.handleUpdate(update, response)

    this.received.push({
      data: response._getJSONData(),
    })
  }

  sendCommand(command: string) {
    const params = {
      message: {
        text: command,
        entities: [
          {
            type: 'bot_command',
            offset: 0,
            length: command.length,
          }
        ],
      },
    }
    return this.send(params as any)
  }

  sendText(text: string) {
    const params = {
      message: {
        text,
      },
    }
    return this.send(params as any)
  }

  sendSticker(sticker: { file_id: string, file_unique_id: string }) {
    const params = {
      message: {
        sticker,
      },
    }
    return this.send(params as any)
  }
}
