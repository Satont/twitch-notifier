import { Module } from '@nestjs/common'
import { TwitchModule } from './twitch/twitch.module'

@Module({
  imports: [TwitchModule],
  controllers: [],
  providers: [],
})
export class AppModule {}
