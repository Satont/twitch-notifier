import { expect } from 'chai'
import { getConnection } from 'typeorm'
import createDbConnection from './helpers/createDbConnection'
import { TestTelegramClient } from './helpers/telegram-client'

describe('telegram', function() {
  let client: TestTelegramClient

  before(async () => {
    await createDbConnection()
    client = new TestTelegramClient()
    await client.init()
  })

  it('class command should exists and not equals 0 length', async () => {
    expect(client.client.commands).to.be.an('array')
    expect(client.client.commands).to.be.not.empty
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
    client.client.bot.stop()
  })

})
