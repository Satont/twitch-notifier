import { getConnection, In } from 'typeorm'
import { Chat, Services } from '../entities/Chat'

export interface SendMessageOpts {
  target: string | string[] | number | number[]
  message: string
  image?: string
}

export const services: ServiceInterface[] = []

export class ServiceInterface {
  inited = false
  service!: Services

  constructor({ service }: { service: Services }) {
    services.push(this)
    this.service = service
  }

  protected init(): void {
    throw new Error('Method init not implemented')
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  async sendMessage(opts: SendMessageOpts): Promise<any> {
    throw new Error('Method sendMessage not implemented')
  }

  async makeAnnounce(opts: SendMessageOpts) {
    if (!this.inited) return
    const targets = Array.isArray(opts.target) ? opts.target : [opts.target]
    const repository = getConnection().getRepository(Chat)
    const chats = (await repository.find({ service: this.service, id: In(targets.map(t => t.toString())) })).map(c => c.id)
    this.sendMessage({ target: chats, ...opts })
  }
}
