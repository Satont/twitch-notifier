import { getConnection, In } from 'typeorm'
import { Chat, Services } from '../entities/Chat'

export interface SendMessageOpts {
  target: string | string[]
  message: string
  image?: string
}

export const services: ServiceInterface[] = []

export class ServiceInterface {
  service!: Services
  commands: Array<{ name: string, fnc: string, description?: string }>

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

  async makeAnnounce(opts: SendMessageOpts) {
    const targets = Array.isArray(opts.target) ? opts.target : [opts.target]
    const repository = getConnection().getRepository(Chat)
    const chats = (await repository.find({ service: this.service, id: In(targets) })).map(c => c.id)
    this.sendMessage({ target: chats, ...opts })
  }
}
