import gl from 'glob'
import { promisify } from 'util'
import { error } from './logger'
import { get, set, template } from 'lodash'
import fs from 'fs'
const glob = promisify(gl)

export class I18n {
  translations = {}
  private lang = 'en'

  constructor(lang?: string, translations?: any) {
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
    return true
  }

  clone(lang: string) {
    return new I18n(lang, this.translations)
  }

  translate(path: string, data?: Record<string, any>) {
    const defaultTranslation = get(this.translations, `en.{path}`, `${this.lang}.errors.translateNotFound`)
    const str = get(this.translations, `${this.lang}.${path}`, defaultTranslation)
    const result = template(str, { interpolate: /{{([\s\S]+?)}}/g })
    return result(data)
  }
}

export const i18n = new I18n()
export default i18n
