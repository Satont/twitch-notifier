import i18n from '../src/libs/i18n'

describe('i18n', function() {
  const i18: typeof i18n = i18n

  beforeAll(async () => {
    await i18n.init()
  })

  it('by default should return "English"', async () => {
    expect(i18.translate('language.name')).toEqual('English')
  })

  it('clone should switch language to russian', async () => {
    const newI18n = i18.clone('ru')

    expect(newI18n.lang).toEqual('ru')
  })

  it('after clone should return "Русский"', async () => {
    const newI18n = i18.clone('ru')

    expect(newI18n.translate('language.name')).toEqual('Русский')
  })

})
