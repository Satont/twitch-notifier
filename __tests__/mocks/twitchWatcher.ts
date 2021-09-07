jest.mock('../../src/watchers/twitch', () => ({
  TwitchWatcher: {
    addChannelToWatch: jest.fn().mockImplementation(() => true)
  },
}))