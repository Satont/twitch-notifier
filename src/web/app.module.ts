import { Module } from '@nestjs/common'
import { TwitchModule } from './twitch/twitch.module'
import { AppController } from './app.controller'
import { AppService } from './app.service'
import { LoggerModule } from 'nestjs-pino'

@Module({
  imports: [
    TwitchModule, 
    LoggerModule.forRoot({
      pinoHttp: { prettifier: true, prettyPrint: true }
    })
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
