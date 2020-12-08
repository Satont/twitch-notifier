import { expect } from 'chai'
import { getConnection } from 'typeorm'
import createDbConnection from './helpers/createDbConnection'
import { TestTelegramClient } from './helpers/telegramСlient'

describe('telegram', function() {
  let instance: TestTelegramClient

  before(async () => {
    await createDbConnection()
    instance = new TestTelegramClient()
    await instance.init()
  })

  it('class command should exists and not equals 0 length', async () => {
    expect(instance.client.commands).to.be.an('array')
    expect(instance.client.commands).to.be.not.empty
  })

  it('/start should reply marukup inline keyboard', async () => {
    await instance.sendCommand('/start')

    const recieved = instance.received[0]
    expect(recieved).to.be.exist
    expect(recieved).to.have.nested.property('data.reply_markup.inline_keyboard')
    expect(recieved.data.reply_markup.inline_keyboard).to.be.not.empty
  })

  after(() => {
    getConnection().close()
    instance.client.bot.stop()
  })

})
