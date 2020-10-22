/* eslint-disable @typescript-eslint/no-unused-vars */
export class Commands {
  commands = [
    { name: 'follow', fnc: this.follow },
    { name: 'commands', fnc: this.commandsList },
    { name: 'follows', fnc: this.follows },
    { name: 'unfollow', fnc: this.unfollow },
  ]

  protected async follow(args: any): Promise<any> {
    throw new Error('Method follow not implemented')
  }

  protected async unfollow(args: any): Promise<any> {
    throw new Error('Method unfollow not implemented')
  }

  protected async commandsList(args: any): Promise<any> {
    throw new Error('Method commandsList not implemented')
  }

  protected async live(args: any): Promise<any> {
    throw new Error('Method live not implemented')
  }

  protected async follows(args: any): Promise<any> {
    throw new Error('Method follows not implemented')
  }
}