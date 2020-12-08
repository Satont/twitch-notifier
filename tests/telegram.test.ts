import { expect } from 'chai'
import { getConnection } from 'typeorm'
import createDbConnection from './helpers/createDbConnection'
import { TestTelegramClient } from './helpers/telegram-client'

describe('telegram', function() {
  let client: TestTelegramClient

  before(async () => {
    await createDbConnection()
  })

  beforeEach(async () => {
    client = new TestTelegramClient()
    await client.init()
  })

  it('/start should reply marukup inline keyboard', async () => {
    await client.sendCommand('/start')

    const recieved = client.received[0]
    expect(recieved).to.be.exist
    expect(recieved).to.have.nested.property('data.reply_markup.inline_keyboard')
    expect(recieved.data.reply_markup.inline_keyboard).to.be.not.empty
  })

  after(() => {
    getConnection().close()
    client.bot.stop()
  })
})
