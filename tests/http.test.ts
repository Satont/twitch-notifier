import { NestExpressApplication } from '@nestjs/platform-express'
import { expect } from 'chai'
import { getConnection } from 'typeorm'
import createDbConnection from './helpers/createDbConnection'

describe('http server should be runned at port 3000', function() {
  let web: typeof import('../src/web')
  let app: NestExpressApplication

  before(async () => {
    await createDbConnection()

    web = await import('../src/web')
  })

  it('bootstrap web', async () => {
    app = await web.bootstrap()
    expect(app).to.not.eq(undefined)
  })

  after(() => {
    getConnection().close()
    app.close()
  })
})
