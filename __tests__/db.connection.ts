import createDbConnection from './helpers/createDbConnection'

describe('Db connection should be created', () => {
  test('Create sqlite connection', async () => {
    const connection = await createDbConnection()

    expect(connection.isConnected).toBe(true)
  })
})