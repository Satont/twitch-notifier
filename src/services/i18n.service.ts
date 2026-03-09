import i18next from 'i18next';
import type { MiddlewareFn } from 'grammy';
import enLocale from '../../locales/en.json';
import ruLocale from '../../locales/ru.json';
import ukLocale from '../../locales/uk.json';

export type SupportedLanguage = 'en' | 'ru' | 'uk';

export class I18nService {
  private i18n: typeof i18next;
  private initialized = false;

  constructor() {
    this.i18n = i18next.createInstance();
  }

  /**
   * Initialize i18next instance with locales
   * Must be called before using the service
   */
  async init(): Promise<void> {
    if (this.initialized) return;

    await this.i18n.init({
      lng: 'en',
      fallbackLng: 'en',
      defaultNS: 'translation',
      ns: ['translation'],
      resources: {
        en: { translation: enLocale },
        ru: { translation: ruLocale },
        uk: { translation: ukLocale },
      },
      interpolation: {
        escapeValue: false, // Not needed for Telegram (no XSS risk)
      },
    });

    this.initialized = true;
  }

  /**
   * Get translated string
   * @param locale - Language code
   * @param key - Translation key (dot notation)
   * @param params - Template parameters
   */
  t(locale: SupportedLanguage, key: string, params?: Record<string, any>): string {
    if (!this.initialized) {
      throw new Error('I18nService not initialized. Call init() first.');
    }
    return this.i18n.t(key, { ...params, lng: locale });
  }

  /**
   * Get Grammy middleware that attaches t() function to context
   */
  middleware(): MiddlewareFn<any> {
    return async (ctx, next) => {
      const language = ctx.session?.language || 'en';
      
      // Attach t() function to context that uses session language
      ctx.t = (key: string, params?: Record<string, any>) => {
        return this.t(language, key, params);
      };

      await next();
    };
  }

  /**
   * Get all available locales
   */
  getAvailableLocales(): SupportedLanguage[] {
    return ['en', 'ru', 'uk'];
  }

  /**
   * Check if locale is supported
   */
  isValidLocale(locale: string): locale is SupportedLanguage {
    return ['en', 'ru', 'uk'].includes(locale);
  }
}
