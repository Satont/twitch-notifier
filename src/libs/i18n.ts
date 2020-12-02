import i18next, { i18n as i18 } from 'i18next'
import Backend from 'i18next-fs-backend'
import { join } from 'path'
import { lstatSync, readdirSync } from 'fs'

export let i18n: i18
i18next
  .use(Backend)
  .init({
    lng: 'english',
    fallbackLng: 'english',
    debug: false,
    preload: readdirSync(join(__dirname, '../../locales')).filter((fileName) => {
      const joinedPath = join(join(__dirname, '../../locales'), fileName)
      const isDirectory = lstatSync(joinedPath).isDirectory()
      return isDirectory
    }),
    backend: {
      loadPath: join(__dirname, '../../locales/{{lng}}/{{ns}}.json'),
    },
  }).then(() => i18n = i18next)
