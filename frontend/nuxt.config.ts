// https://nuxt.com/docs/api/configuration/nuxt-config

export default defineNuxtConfig({
  compatibilityDate: '2024-04-03',
  devtools: { enabled: true },
  modules: ['@nuxt/ui', '@nuxtjs/i18n'],

  css: [
    '~/assets/css/main.css',
    'flag-icons/css/flag-icons.min.css'
  ],

  // i18n Internationalization Config
  i18n: {
    bundle: {
      optimizeTranslationDirective: false
    },
    locales: [
      { code: 'en-US', iso: 'en-US', name: 'English', file: 'en-US.json' },
      { code: 'zh-TW', iso: 'zh-TW', name: '繁體中文', file: 'zh-TW.json' }
    ],
    defaultLocale: 'zh-TW',
    langDir: 'locales/',
    lazy: true,
    strategy: 'prefix',
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'i18n_redirected',
      redirectOn: 'root',
      alwaysRedirect: true,
      fallbackLocale: 'zh-TW'
    }
  },

  // App-level page transition configurations
  app: {
    pageTransition: {
      name: 'page',
      mode: 'out-in'
    },
    layoutTransition: {
      name: 'layout',
      mode: 'out-in'
    }
  },

  // Runtime config using environment variables
  runtimeConfig: {
    public: {
      apiBase: '/api/v1' // Points to local proxy in dev, or custom domain in prod
    }
  },

  // Route rules for forwarding API requests to the Go backend, preserving the prefix
  routeRules: {
    '/api/v1/**': {
      proxy: `${process.env.BACKEND_API_BASE || 'http://localhost:8080'}/api/v1/**`
    }
  }
})
