export interface SendMessageOpts {
  target: string | string[]
  message: string
  image?: string
}

export const services: ServiceInterface[] = []

export class ServiceInterface {
  service!: string
  commands: Array<{ name: string, fnc: string }>

  constructor() {
    services.push(this)
  }

  protected init(): void {
    throw new Error('Method init not implemented')
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  async sendMessage(opts: SendMessageOpts): Promise<any> {
    throw new Error('Method sendMessage not implemented')
  }
}
