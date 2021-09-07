const watcherObject = {
  addChannelToWatch: jest.fn().mockImplementation(() => true),
}

jest.mock('../../src/watchers/twitch', () => ({
  ...watcherObject,
  TwitchWatcher: watcherObject,
}))