import { Body, Controller, Get, Post, Query, Res } from '@nestjs/common'
import { Response } from 'express'
import { ITwitchStreamChangedPayload } from '../../../typings/twitch'
import TwitchWatcher from '../../../watchers/twitch'

@Controller('twitch/webhooks/callback')
export class WebhooksController {
  @Get()
  async getRequest(@Query() query: any) {
    return query['hub.challenge']
  }

  @Post()
  async postRequest(@Body() body: ITwitchStreamChangedPayload, @Res() res: Response) {
    TwitchWatcher.processPayload(body.data)
    res.sendStatus(200)
  }
}
