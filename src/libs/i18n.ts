/* import i18next, { i18n as i18 } from 'i18next'
import Backend from 'i18next-fs-backend'
import { join, resolve } from 'path'
import { lstatSync, readdirSync } from 'fs'

const defaultLanguage = 'english'

export let i18n: i18
i18next
  .use(Backend)
  .init({
    lng: defaultLanguage,
    fallbackLng: defaultLanguage,
    debug: true,
    preload: readdirSync(join(__dirname, '../../locales')).filter((fileName) => {
      const joinedPath = join(join(__dirname, '../../locales'), fileName)
      const isDirectory = lstatSync(joinedPath).isDirectory()
      return isDirectory
    }),
    backend: {
      loadPath: join(__dirname, '../../locales/{{lng}}/{{ns}}.json'),
    },
  }).then(() => i18n = i18next)
 */

import gl from 'glob'
import { promisify } from 'util'
import { error } from './logger'
import { get, set, template } from 'lodash'
import fs from 'fs'
import { Languages } from '../entities/ChatSettings'
const glob = promisify(gl)

export class I18n {
  private translations = {}
  private lang = Languages.ENGLISH

  constructor(lang?: Languages, translations?: any) {
    if (lang) this.lang = lang
    if (translations) this.translations = translations
  }

  private async init() {
    const files = await glob('./locales/**')
    for (const f of files) {
      if (!f.endsWith('.json')) {
        continue
      }
      const withoutLocales = f.replace('./locales/', '').replace('.json', '')
      try {
        set(this.translations, withoutLocales.split('/').join('.'), JSON.parse(fs.readFileSync(f, 'utf8')))
      } catch (e) {
        error('Incorrect JSON file: ' + f)
        error(e.stack)
      }
    }
  }

  clone(lang: Languages) {
    return new I18n(lang, this.translations)
  }

  translate(path: string, data?: Record<string, any>) {
    const str = get(this.translations, `${this.lang}.${path}`, `${this.lang}.errors.translateNotFound`)
    const result = template(str, { interpolate: /{{([\s\S]+?)}}/g })
    return result(data)
  }
}

export const i18n = new I18n()
export default i18n
