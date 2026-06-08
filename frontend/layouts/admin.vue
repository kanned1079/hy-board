<template>
  <div class="h-screen flex flex-col lg:flex-row bg-slate-50 dark:bg-zinc-950 text-slate-900 dark:text-zinc-100 relative transition-colors duration-300 overflow-hidden">
    <!-- Desktop Left Sidebar (Admin Panel) -->
    <aside class="hidden lg:flex flex-col w-64 border-r border-slate-200 dark:border-zinc-800/80 backdrop-blur-md bg-white/50 dark:bg-zinc-900/30 p-6 space-y-8 h-screen sticky top-0">
      <div class="flex justify-between items-center gap-1">
        <span class="text-xl font-bold bg-gradient-to-r from-primary-500 to-primary-300 dark:from-primary-400 dark:to-primary-200 bg-clip-text text-transparent mr-2">HY Admin</span>
        <div class="flex items-center gap-1">
          <!-- Language Switcher -->
          <UDropdownMenu :items="langItems" :content="{ align: 'end' }" :ui="{ item: 'items-center' }">
            <UButton color="neutral" variant="ghost" icon="i-lucide-languages" class="flex items-center p-1.5" />

            <template #item-leading="{ item }">
              <span v-if="item.flag" :class="item.flag" class="w-4 h-3 bg-cover bg-center rounded-sm shrink-0 self-center"></span>
              <UIcon v-else-if="item.icon" :name="item.icon" class="w-4 h-4 shrink-0 self-center" />
            </template>
          </UDropdownMenu>

          <!-- Theme Toggle Button -->
          <UButton
            :icon="isDark ? 'i-lucide-moon' : 'i-lucide-sun'"
            color="neutral"
            variant="ghost"
            @click="isDark = !isDark"
          />
        </div>
      </div>
      
      <nav class="flex-1 overflow-y-auto pr-2 space-y-4 max-h-[calc(100vh-180px)] scrollbar-thin scrollbar-thumb-slate-200 dark:scrollbar-thumb-zinc-800">
        <!-- Dashboard -->
        <NuxtLink
          :to="localePath('/admin/dashboard')"
          class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
          :class="route.path === localePath('/admin/dashboard') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
        >
          <UIcon name="i-lucide-layout-dashboard" class="w-5 h-5" />
          <span>{{ t('dashboard') }}</span>
        </NuxtLink>

        <!-- Group: 設置 (Settings) -->
        <div class="space-y-1">
          <p class="px-4 text-[10px] font-bold uppercase tracking-wider text-slate-400 dark:text-zinc-500">設置</p>
          <NuxtLink
            :to="localePath('/admin/settings/system')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/settings/system') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-settings" class="w-4 h-4" />
            <span>系統配置</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/admin/settings/theme')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/settings/theme') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-palette" class="w-4 h-4" />
            <span>主題配置</span>
          </NuxtLink>
        </div>

        <!-- Group: 服務器 (Servers) -->
        <div class="space-y-1">
          <p class="px-4 text-[10px] font-bold uppercase tracking-wider text-slate-400 dark:text-zinc-500">服務器</p>
          <NuxtLink
            :to="localePath('/admin/servers/nodes')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/servers/nodes') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-server" class="w-4 h-4" />
            <span>節點管理</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/admin/servers/groups')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/servers/groups') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-shield" class="w-4 h-4" />
            <span>權限組管理</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/admin/servers/routes')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/servers/routes') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-git-fork" class="w-4 h-4" />
            <span>路由管理</span>
          </NuxtLink>
        </div>

        <!-- Group: 財務 (Financial) -->
        <div class="space-y-1">
          <p class="px-4 text-[10px] font-bold uppercase tracking-wider text-slate-400 dark:text-zinc-500">財務</p>
          <NuxtLink
            :to="localePath('/admin/financial/subscriptions')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/financial/subscriptions') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-credit-card" class="w-4 h-4" />
            <span>訂閱管理</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/admin/financial/orders')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/financial/orders') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-receipt" class="w-4 h-4" />
            <span>訂單管理</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/admin/financial/coupons')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/financial/coupons') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-ticket" class="w-4 h-4" />
            <span>優惠券管理</span>
          </NuxtLink>
        </div>

        <!-- Group: 用戶 (Users) -->
        <div class="space-y-1">
          <p class="px-4 text-[10px] font-bold uppercase tracking-wider text-slate-400 dark:text-zinc-500">用戶</p>
          <NuxtLink
            :to="localePath('/admin/users')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/users') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-users" class="w-4 h-4" />
            <span>用戶管理</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/admin/announcements')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/announcements') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-megaphone" class="w-4 h-4" />
            <span>公告管理</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/admin/tickets')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/tickets') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-message-square" class="w-4 h-4" />
            <span>工單管理</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/admin/knowledge')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/knowledge') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-book-open" class="w-4 h-4" />
            <span>知識庫管理</span>
          </NuxtLink>
        </div>

        <!-- Group: 指標 (Metrics) -->
        <div class="space-y-1">
          <p class="px-4 text-[10px] font-bold uppercase tracking-wider text-slate-400 dark:text-zinc-500">指標</p>
          <NuxtLink
            :to="localePath('/admin/metrics/queue')"
            class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
            :class="route.path === localePath('/admin/metrics/queue') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-activity" class="w-4 h-4" />
            <span>隊列監控</span>
          </NuxtLink>
        </div>
      </nav>

      <div class="border-t border-slate-200 dark:border-zinc-800/80 pt-4">
        <UButton
          color="red"
          variant="ghost"
          block
          icon="i-lucide-log-out"
          class="justify-start px-4"
          @click="logout"
        >
          {{ t('logout') }}
        </UButton>
      </div>
    </aside>

    <header class="lg:hidden flex justify-between items-center px-6 py-4 border-b border-slate-200 dark:border-zinc-800/80 backdrop-blur-md bg-white/80 dark:bg-zinc-900/20 sticky top-0 z-40">
      <span class="text-lg font-bold bg-gradient-to-r from-primary-500 to-primary-300 dark:from-primary-400 dark:to-primary-200 bg-clip-text text-transparent">HY Admin</span>
      <div class="flex items-center space-x-2">
        <!-- Language Switcher -->
        <UDropdownMenu :items="langItems" :content="{ align: 'end' }" :ui="{ item: 'items-center' }">
          <UButton color="neutral" variant="ghost" icon="i-lucide-languages" class="flex items-center p-1.5" />

          <template #item-leading="{ item }">
            <span v-if="item.flag" :class="item.flag" class="w-4 h-3 bg-cover bg-center rounded-sm shrink-0 self-center"></span>
            <UIcon v-else-if="item.icon" :name="item.icon" class="w-4 h-4 shrink-0 self-center" />
          </template>
        </UDropdownMenu>

        <!-- Theme Toggle Button -->
        <UButton
          :icon="isDark ? 'i-lucide-moon' : 'i-lucide-sun'"
          color="neutral"
          variant="ghost"
          @click="isDark = !isDark"
        />
        <UButton
          color="error"
          variant="ghost"
          icon="i-lucide-log-out"
          @click="logout"
        />
      </div>
    </header>

    <!-- Main Content Area -->
    <main class="flex-1 min-h-0 relative overflow-y-auto">
      <div class="p-4 sm:p-6 pb-20 sm:pb-6">
        <slot />
      </div>
    </main>

    <nav class="lg:hidden fixed bottom-0 left-0 right-0 z-40 border-t border-slate-200 dark:border-zinc-800/80 backdrop-blur-lg bg-white/80 dark:bg-zinc-950/80 px-4 py-3 flex justify-around items-center">
      <NuxtLink
        :to="localePath('/admin/dashboard')"
        class="flex flex-col items-center space-y-1 transition-all duration-300 text-xs"
        :class="route.path === localePath('/admin/dashboard') ? 'text-primary-600 dark:text-primary-400' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white'"
      >
        <UIcon name="i-lucide-layout-dashboard" class="w-5 h-5" />
        <span class="text-[10px]">{{ t('dashboard') }}</span>
      </NuxtLink>
      <NuxtLink
        :to="localePath('/admin/users')"
        class="flex flex-col items-center space-y-1 transition-all duration-300 text-xs"
        :class="route.path === localePath('/admin/users') ? 'text-primary-600 dark:text-primary-400' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white'"
      >
        <UIcon name="i-lucide-users" class="w-5 h-5" />
        <span class="text-[10px]">{{ t('user_management') }}</span>
      </NuxtLink>
      <NuxtLink
        :to="localePath('/admin/tickets')"
        class="flex flex-col items-center space-y-1 transition-all duration-300 text-xs"
        :class="route.path === localePath('/admin/tickets') ? 'text-primary-600 dark:text-primary-400' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white'"
      >
        <UIcon name="i-lucide-message-square" class="w-5 h-5" />
        <span class="text-[10px]">工單管理</span>
      </NuxtLink>
      <NuxtLink
        :to="localePath('/admin/announcements')"
        class="flex flex-col items-center space-y-1 transition-all duration-300 text-xs"
        :class="route.path === localePath('/admin/announcements') ? 'text-primary-600 dark:text-primary-400' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white'"
      >
        <UIcon name="i-lucide-megaphone" class="w-5 h-5" />
        <span class="text-[10px]">{{ t('announcements_manage') }}</span>
      </NuxtLink>
      <NuxtLink
        :to="localePath('/admin/knowledge')"
        class="flex flex-col items-center space-y-1 transition-all duration-300 text-xs"
        :class="route.path === localePath('/admin/knowledge') ? 'text-primary-600 dark:text-primary-400' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white'"
      >
        <UIcon name="i-lucide-book-open" class="w-5 h-5" />
        <span class="text-[10px]">{{ t('knowledge_base') }}</span>
      </NuxtLink>
    </nav>
  </div>
</template>

<script setup>
const route = useRoute()
const router = useRouter()
const { t, locale, setLocale } = useI18n()
const localePath = useLocalePath()

const colorMode = useColorMode()
const isDark = computed({
  get () {
    return colorMode.value === 'dark'
  },
  set () {
    colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark'
  }
})

const selectLanguage = (code) => {
  if (code === 'auto') {
    localStorage.setItem('user_lang_pref', 'auto')
    const browserLang = navigator.language || 'zh-TW'
    if (browserLang.toLowerCase().includes('zh')) {
      setLocale('zh-TW')
    } else {
      setLocale('en-US')
    }
  } else {
    localStorage.setItem('user_lang_pref', code)
    setLocale(code)
  }
}

const langItems = [
  [
    {
      label: '自動',
      icon: 'i-lucide-cpu',
      onSelect: () => selectLanguage('auto'),
      click: () => selectLanguage('auto')
    },
    {
      label: '英語 (美國)',
      flag: 'fi fi-us',
      onSelect: () => selectLanguage('en-US'),
      click: () => selectLanguage('en-US')
    },
    {
      label: '正體中文 (台灣)',
      flag: 'fi fi-tw',
      onSelect: () => selectLanguage('zh-TW'),
      click: () => selectLanguage('zh-TW')
    }
  ]
]

onMounted(() => {
  const pref = localStorage.getItem('user_lang_pref')
  if (pref === 'auto') {
    const browserLang = navigator.language || 'zh-TW'
    if (browserLang.toLowerCase().includes('zh')) {
      setLocale('zh-TW')
    } else {
      setLocale('en-US')
    }
  } else if (pref) {
    setLocale(pref)
  }
})

const logout = () => {
  const tokenCookie = useCookie('auth_token')
  tokenCookie.value = null
  router.push(localePath('/login'))
}
</script>
