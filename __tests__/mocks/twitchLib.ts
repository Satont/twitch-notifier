import { createHelixStream } from '../helpers/createHelixStream'
import { createHelixUser } from '../helpers/createHelixUser'

jest.mock('../../src/libs/twitch', () => ({
  Twitch: {
    getUser: jest.fn().mockImplementation((opts) => Promise.resolve(createHelixUser(opts))),
    getUsers: jest.fn().mockImplementation(({ ids, names }: { ids?: string[], names?: string[] }) => {
      const type = ids ? 'ids' : 'names'
      const users = [...(ids || names)].map(arg => {
        if (type === 'ids') {
          return createHelixUser({ id: arg })
        } else {
          return createHelixUser({ name: arg })
        }
      })

      return Promise.resolve(users)
    }),
    getStreams: jest.fn().mockImplementation((usersIds: string[]) => {
      const streams = usersIds.map(id => createHelixStream({ id }))

      return Promise.resolve(streams)
    })
  }
}))