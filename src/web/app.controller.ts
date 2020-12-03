import { Controller, Get, Render } from '@nestjs/common'
import { AppService } from './app.service'

@Controller('')
export class AppController {
  constructor(
    private readonly service: AppService
  ){}

  @Get()
  @Render('pages/home')
  async root() {
    const counts = await this.service.counts()
    return counts
  }

  @Get('/top')
  @Render('pages/top')
  async top() {
    const list = await this.service.top(10)
    return { list }
  }
}
