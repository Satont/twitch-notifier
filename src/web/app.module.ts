import { Module } from '@nestjs/common'
import { TwitchModule } from './twitch/twitch.module'
import { AppController } from './app.controller'
import { AppService } from './app.service'

@Module({
  imports: [TwitchModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
