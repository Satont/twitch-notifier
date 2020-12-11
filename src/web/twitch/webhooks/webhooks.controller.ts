import { Body, Controller, Get, Post, Query, Req, Res } from '@nestjs/common'
import { Request, Response } from 'express'
import { ITwitchStreamChangedPayload } from '../../../typings/twitch'
import TwitchWatcher from '../../../watchers/twitch'

@Controller('twitch/webhooks/callback')
export class WebhooksController {
  cache = new Set()

  constructor() {
    // clear cache every 2 hours
    setInterval(() => this.cache.clear(), 2 * 60 * 60 * 1000)
  }

  @Get()
  async getRequest(@Query() query: any) {
    return query['hub.challenge']
  }

  @Post()
  async postRequest(@Body() body: ITwitchStreamChangedPayload, @Res() res: Response, @Req() req: Request) {
    res.sendStatus(200)

    if (this.cache.has(req.headers['twitch-notification-id'])) {
      return
    }

    this.cache.add(req.headers['twitch-notification-id'])

    // We need to manually set some params if data length is 0. Data length can be 0 because twitch sending empty array if stream goes offline
    if (!body.data.length) {
      const regexp = /(\buser_id=\b)([0-9]+)/gm
      const user_id = regexp.exec(req.headers.link as string)[2]
      body.data[0] = { user_id, type: 'offline' } as any
    }

    TwitchWatcher.processPayload(body.data)
  }
}
