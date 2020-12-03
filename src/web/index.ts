import { ValidationPipe } from '@nestjs/common'
import { NestFactory } from '@nestjs/core'
import { NestExpressApplication } from '@nestjs/platform-express'
import { resolve } from 'path'
import { AppModule } from './app.module'
import hbs from 'hbs'

async function bootstrap() {
  const app = await NestFactory.create<NestExpressApplication>(AppModule, {
    logger: ['error', 'warn', 'log'],
  })
  app.useGlobalPipes(new ValidationPipe())
  app.useStaticAssets(resolve(process.cwd(), 'public'))

  app.set('view engine', 'hbs')
  app.set('views', resolve(process.cwd(), 'views'))
  app.set('view options', {
    layout: 'layouts/index',
    templates: resolve(process.cwd(), 'views'),
  })
  hbs.registerPartials(resolve(process.cwd(), 'views', 'partials'))

  await app.listen(3000, '0.0.0.0')
}

bootstrap()
