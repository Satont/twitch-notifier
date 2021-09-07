import createDbConnection from './helpers/createDbConnection'

beforeAll(async () => {
  await createDbConnection()
})

import { Telegram } from '../src/services/telegram'
import { createTelegramClient, createUpdate } from './helpers/telegram'
import { Context } from 'telegraf'

describe('telegram', function() {
  let instance: Telegram
  
  beforeEach(() => {
    instance = null
    instance = createTelegramClient()
  })

  it('class command should exists and not equals 0 length', async () => {
    expect(instance.commands.length).not.toBeFalsy()
  })

  it('chat entity should be created', async () => {
    const update = createUpdate({ text: '/start'})
    
    instance.bot.use((ctx: Context) => {
      expect(ctx.ChatEntity).not.toBeFalsy()
      expect(ctx.ChatEntity.settings.language).toEqual('en')
    })

    await instance.bot.handleUpdate({ message: update } as any)
  })
})
