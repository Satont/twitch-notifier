import { command } from '../../src/decorators/command';

it('Should add command to object commands', () => {
  class Test {
    commands: string[]

    @command('testCommand', { description: 'This is test command' })
    test() {}
  }

  expect(new Test().commands).toEqual(
    expect.arrayContaining([
      expect.objectContaining({ name: 'testCommand', description: 'This is test command', fnc: 'test' })
    ])
  )
})