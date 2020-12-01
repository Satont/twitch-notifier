export interface SendMessageOpts {
  target: number | number[]
  message: string
  image?: string
}

export class ServiceInterface {
  service!: string
  commands: Array<{ name: string, fnc: string }>

  protected init(): void {
    throw new Error('Method init not implemented')
  }

  async follow(args: any): Promise<any> {
    throw new Error('Method follow not implemented')
  }

  async unfollow(args: any): Promise<any> {
    throw new Error('Method unfollow not implemented')
  }

  async commandsList(args: any): Promise<any> {
    throw new Error('Method commandsList not implemented')
  }

  async live(args: any): Promise<any> {
    throw new Error('Method live not implemented')
  }

  async follows(args: any): Promise<any> {
    throw new Error('Method follows not implemented')
  }
}
