import nock from 'nock'
import { EventEmitter } from 'typeorm/platform/PlatformTools'

export default async(ee: EventEmitter) => {
  const nocked = nock('https://api.telegram.org/bot123456', { allowUnmocked: true })
    .post('/sendMessage')
    .reply((uri, requestBody) => {
      ee.emit('sendMessage', {
        method: 'sendMessage',
        ...requestBody as any,
      })
      return [
        200,
        { ok: true },
      ]
    })
    .post('/getMe')
    .reply(200, { ok: true, result: { id: 42, is_bot: true, username: 'bot', first_name: 'Bot' } })
    .post('/setMyCommands')
    .reply(200, { ok: true })
    .post('/deleteWebhook')
    .reply(200, { ok: true })
    .post('/getUpdates')
    .query({ offset: 0, limit: 100, timeout: 30 })
    .reply(200, { ok: true, result: [] })

  return nocked
}
