import { expect } from 'chai'
import { getConnection } from 'typeorm'
import i18n from '../src/libs/i18n'

describe('i18n', function() {
  const i18: typeof i18n = i18n

  before(async () => {
    await i18n.init()
  })

  it('by default should return "English"', async () => {
    expect(i18.translate('language.name')).to.eq('English')
  })

  it('clone should switch language to russian', async () => {
    const newI18n = i18.clone('ru')

    expect(newI18n.lang).to.eq('ru')
  })

  it('after clone should return "Русский"', async () => {
    const newI18n = i18.clone('ru')

    expect(newI18n.translate('language.name')).to.eq('Русский')
  })

})
