import { ValidationPipe } from '@nestjs/common'
import { NestFactory } from '@nestjs/core'
import { NestExpressApplication } from '@nestjs/platform-express'
import { resolve } from 'path'
import hbs from 'hbs'
import { Logger } from 'nestjs-pino'

let app: NestExpressApplication
// eslint-disable-next-line prefer-const
export let listened = false

export async function bootstrap() {
  app = await NestFactory.create<NestExpressApplication>((await import('./app.module')).AppModule, { 
    logger: false,
  })

  app.useLogger(app.get(Logger))

  app.useGlobalPipes(new ValidationPipe())
  app.useStaticAssets(resolve(process.cwd(), 'public'))
  app.set('view engine', 'hbs')
  app.set('views', resolve(process.cwd(), 'views'))
  app.set('view options', {
    layout: 'layouts/index',
    templates: resolve(process.cwd(), 'views'),
  })
  hbs.registerPartials(resolve(process.cwd(), 'views', 'partials'))
}

export const getAppLication = () => app
