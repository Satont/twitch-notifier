import { Commands } from './_commands'

export interface SendMessageOpts {
  target: number | number[]
  message: string
  image?: string 
}

export class ServiceInterface extends Commands { 
  service!: string

  protected init(): void {
    throw new Error('Method init not implemented')
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  protected async sendMessage(opts: SendMessageOpts): Promise<boolean> {
    throw new Error('Method sendMessage not implemented')
  }
}