import { Module } from '@nestjs/common'
import { WebhooksController } from './webhooks/webhooks.controller'

@Module({
  controllers: [WebhooksController],
})
export class TwitchModule {}
