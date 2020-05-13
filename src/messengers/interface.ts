import { info } from "../helpers/logs"

export type SendMessageOpts = {
  target: number | number[],
  message: string,
  image?: string 
}

export type Commands = ['follow', 'unfollow', 'live', 'follows']

export class IService { 
  service?: string

  protected init(): void {
    throw new Error('Method init not implemented')
  }

  public async sendMessage(opts: SendMessageOpts): Promise<boolean> {
    throw new Error('Method sendMessage not implemented')
  }
  protected async loadCommands(): Promise<void> {
    info(`Commands for ${this.service} messenger loaded`)
  }
  protected async loadMiddlewares?(): Promise<void> {
    info(`Middlewares for ${this.service} messenger loaded`)
  }
  protected async registerScenes?(): Promise<void> {
    info(`Scenes for ${this.service} messenger loaded`)
  }

}