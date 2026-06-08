<template>
  <div class="h-screen flex flex-col lg:flex-row bg-slate-50 dark:bg-zinc-950 text-slate-900 dark:text-zinc-100 relative transition-colors duration-300 overflow-hidden">
    
    <!-- Mobile USlideover Sidebar Drawer -->
    <USlideover v-model:open="isDrawerOpen" side="left" class="lg:hidden">
      <div class="p-5 h-full flex flex-col bg-white dark:bg-zinc-900 border-r border-slate-200 dark:border-zinc-800 overflow-y-auto">
        <div class="flex justify-between items-center pb-4 border-b border-slate-150 dark:border-zinc-800 mb-6">
          <span class="text-lg font-bold bg-gradient-to-r from-primary-500 to-primary-400 dark:from-primary-400 dark:to-primary-300 bg-clip-text text-transparent">HY-Board</span>
          <UButton
            color="gray"
            variant="ghost"
            icon="i-lucide-x"
            @click="isDrawerOpen = false"
          />
        </div>

        <nav class="flex-1 space-y-6">
          <!-- Standalone items -->
          <div class="space-y-1">
            <NuxtLink
              :to="localePath('/dashboard')"
              class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
              :class="route.path === localePath('/dashboard') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
            >
              <UIcon name="i-lucide-layout-dashboard" class="w-5 h-5" />
              <span>{{ t('dashboard') }}</span>
            </NuxtLink>
            <NuxtLink
              :to="localePath('/dashboard/knowledge')"
              class="flex items-center space-x-3 px-4 py-2.5 rounded-lg transition-all duration-300 text-sm font-medium"
              :class="route.path === localePath('/dashboard/knowledge') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
            >
              <UIcon name="i-lucide-book-open" class="w-5 h-5" />
              <span>{{ t('knowledge_base') }}</span>
            </NuxtLink>
          </div>

          <!-- Group: Subscription -->
          <div class="space-y-1.5">
            <div class="px-4 text-[10px] font-bold uppercase tracking-wider text-slate-400 dark:text-zinc-500">
              {{ t('subscription') }}
            </div>
            <NuxtLink
              :to="localePath('/dashboard/subscribe')"
              class="flex items-center space-x-3 px-4 py-2 rounded-lg transition-all duration-300 text-xs font-medium"
              :class="route.path === localePath('/dashboard/subscribe') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
            >
              <UIcon name="i-lucide-shopping-bag" class="w-4.5 h-4.5" />
              <span>{{ t('purchase_subscription') }}</span>
            </NuxtLink>
            <NuxtLink
              :to="localePath('/dashboard/nodes')"
              class="flex items-center space-x-3 px-4 py-2 rounded-lg transition-all duration-300 text-xs font-medium"
              :class="route.path === localePath('/dashboard/nodes') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
            >
              <UIcon name="i-lucide-cpu" class="w-4.5 h-4.5" />
              <span>{{ t('node_status') }}</span>
            </NuxtLink>
          </div>

          <!-- Group: Billing -->
          <div class="space-y-1.5">
            <div class="px-4 text-[10px] font-bold uppercase tracking-wider text-slate-400 dark:text-zinc-500">
              {{ t('billing') }}
            </div>
            <NuxtLink
              :to="localePath('/dashboard/billing')"
              class="flex items-center space-x-3 px-4 py-2 rounded-lg transition-all duration-300 text-xs font-medium"
              :class="route.path === localePath('/dashboard/billing') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
            >
              <UIcon name="i-lucide-receipt" class="w-4.5 h-4.5" />
              <span>{{ t('my_billing') }}</span>
            </NuxtLink>
            <NuxtLink
              :to="localePath('/dashboard/invite')"
              class="flex items-center space-x-3 px-4 py-2 rounded-lg transition-all duration-300 text-xs font-medium"
              :class="route.path === localePath('/dashboard/invite') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
            >
              <UIcon name="i-lucide-gift" class="w-4.5 h-4.5" />
              <span>{{ t('my_invite') }}</span>
            </NuxtLink>
          </div>

          <!-- Group: Account -->
          <div class="space-y-1.5">
            <div class="px-4 text-[10px] font-bold uppercase tracking-wider text-slate-400 dark:text-zinc-500">
              {{ t('account') }}
            </div>
            <NuxtLink
              :to="localePath('/dashboard/profile')"
              class="flex items-center space-x-3 px-4 py-2 rounded-lg transition-all duration-300 text-xs font-medium"
              :class="route.path === localePath('/dashboard/profile') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
            >
              <UIcon name="i-lucide-user" class="w-4.5 h-4.5" />
              <span>{{ t('user_center') }}</span>
            </NuxtLink>
            <NuxtLink
              :to="localePath('/dashboard/tickets')"
              class="flex items-center space-x-3 px-4 py-2 rounded-lg transition-all duration-300 text-xs font-medium"
              :class="route.path === localePath('/dashboard/tickets') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
            >
              <UIcon name="i-lucide-message-square" class="w-4.5 h-4.5" />
              <span>{{ t('my_tickets') }}</span>
            </NuxtLink>
            <NuxtLink
              :to="localePath('/dashboard/traffic')"
              class="flex items-center space-x-3 px-4 py-2 rounded-lg transition-all duration-300 text-xs font-medium"
              :class="route.path === localePath('/dashboard/traffic') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
            >
              <UIcon name="i-lucide-bar-chart-3" class="w-4.5 h-4.5" />
              <span>{{ t('traffic_details') }}</span>
            </NuxtLink>
          </div>
        </nav>

        <div class="border-t border-slate-200 dark:border-zinc-800/80 pt-4 mt-auto">
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
      </div>
    </USlideover>

    <!-- Desktop Left Sidebar -->
    <aside class="hidden lg:flex flex-col w-64 border-r border-slate-200 dark:border-zinc-800/80 backdrop-blur-md bg-white/50 dark:bg-zinc-900/30 p-6 space-y-6 h-screen sticky top-0">
      <div class="flex justify-between items-center gap-1">
        <span class="text-xl font-bold bg-gradient-to-r from-primary-500 to-primary-400 dark:from-primary-400 dark:to-primary-300 bg-clip-text text-transparent mr-2">HY-Board</span>
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
      
      <nav class="flex-1 space-y-5 overflow-y-auto pr-1 select-none custom-scrollbar">
        <!-- Standalone items -->
        <div class="space-y-1">
          <NuxtLink
            :to="localePath('/dashboard')"
            class="flex items-center space-x-3 px-4 py-2 rounded-lg transition-all duration-300 text-xs font-semibold"
            :class="route.path === localePath('/dashboard') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-layout-dashboard" class="w-4.5 h-4.5" />
            <span>{{ t('dashboard') }}</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/dashboard/knowledge')"
            class="flex items-center space-x-3 px-4 py-2 rounded-lg transition-all duration-300 text-xs font-semibold"
            :class="route.path === localePath('/dashboard/knowledge') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-200/50 dark:hover:bg-zinc-900/50'"
          >
            <UIcon name="i-lucide-book-open" class="w-4.5 h-4.5" />
            <span>{{ t('knowledge_base') }}</span>
          </NuxtLink>
        </div>

        <!-- Group: Subscription -->
        <div class="space-y-1">
          <div class="px-4 text-[9px] font-bold uppercase tracking-wider text-slate-400 dark:text-zinc-500">
            {{ t('subscription') }}
          </div>
          <NuxtLink
            :to="localePath('/dashboard/subscribe')"
            class="flex items-center space-x-3 px-4 py-1.5 rounded-lg transition-all duration-300 text-[11px] font-medium"
            :class="route.path === localePath('/dashboard/subscribe') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-905 dark:hover:text-white hover:bg-slate-200/40 dark:hover:bg-zinc-900/40'"
          >
            <UIcon name="i-lucide-shopping-bag" class="w-4 h-4 text-slate-400 dark:text-zinc-500" />
            <span>{{ t('purchase_subscription') }}</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/dashboard/nodes')"
            class="flex items-center space-x-3 px-4 py-1.5 rounded-lg transition-all duration-300 text-[11px] font-medium"
            :class="route.path === localePath('/dashboard/nodes') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-905 dark:hover:text-white hover:bg-slate-200/40 dark:hover:bg-zinc-900/40'"
          >
            <UIcon name="i-lucide-cpu" class="w-4 h-4 text-slate-400 dark:text-zinc-500" />
            <span>{{ t('node_status') }}</span>
          </NuxtLink>
        </div>

        <!-- Group: Billing -->
        <div class="space-y-1">
          <div class="px-4 text-[9px] font-bold uppercase tracking-wider text-slate-400 dark:text-zinc-500">
            {{ t('billing') }}
          </div>
          <NuxtLink
            :to="localePath('/dashboard/billing')"
            class="flex items-center space-x-3 px-4 py-1.5 rounded-lg transition-all duration-300 text-[11px] font-medium"
            :class="route.path === localePath('/dashboard/billing') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-905 dark:hover:text-white hover:bg-slate-200/40 dark:hover:bg-zinc-900/40'"
          >
            <UIcon name="i-lucide-receipt" class="w-4 h-4 text-slate-400 dark:text-zinc-500" />
            <span>{{ t('my_billing') }}</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/dashboard/invite')"
            class="flex items-center space-x-3 px-4 py-1.5 rounded-lg transition-all duration-300 text-[11px] font-medium"
            :class="route.path === localePath('/dashboard/invite') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-905 dark:hover:text-white hover:bg-slate-200/40 dark:hover:bg-zinc-900/40'"
          >
            <UIcon name="i-lucide-gift" class="w-4 h-4 text-slate-400 dark:text-zinc-500" />
            <span>{{ t('my_invite') }}</span>
          </NuxtLink>
        </div>

        <!-- Group: Account -->
        <div class="space-y-1">
          <div class="px-4 text-[9px] font-bold uppercase tracking-wider text-slate-400 dark:text-zinc-500">
            {{ t('account') }}
          </div>
          <NuxtLink
            :to="localePath('/dashboard/profile')"
            class="flex items-center space-x-3 px-4 py-1.5 rounded-lg transition-all duration-300 text-[11px] font-medium"
            :class="route.path === localePath('/dashboard/profile') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-905 dark:hover:text-white hover:bg-slate-200/40 dark:hover:bg-zinc-900/40'"
          >
            <UIcon name="i-lucide-user" class="w-4 h-4 text-slate-400 dark:text-zinc-500" />
            <span>{{ t('user_center') }}</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/dashboard/tickets')"
            class="flex items-center space-x-3 px-4 py-1.5 rounded-lg transition-all duration-300 text-[11px] font-medium"
            :class="route.path === localePath('/dashboard/tickets') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-905 dark:hover:text-white hover:bg-slate-200/40 dark:hover:bg-zinc-900/40'"
          >
            <UIcon name="i-lucide-message-square" class="w-4 h-4 text-slate-400 dark:text-zinc-500" />
            <span>{{ t('my_tickets') }}</span>
          </NuxtLink>
          <NuxtLink
            :to="localePath('/dashboard/traffic')"
            class="flex items-center space-x-3 px-4 py-1.5 rounded-lg transition-all duration-300 text-[11px] font-medium"
            :class="route.path === localePath('/dashboard/traffic') ? 'bg-primary-500/10 dark:bg-primary-500/15 text-primary-600 dark:text-primary-400 border border-primary-500/20 dark:border-primary-500/25' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-905 dark:hover:text-white hover:bg-slate-200/40 dark:hover:bg-zinc-900/40'"
          >
            <UIcon name="i-lucide-bar-chart-3" class="w-4 h-4 text-slate-400 dark:text-zinc-500" />
            <span>{{ t('traffic_details') }}</span>
          </NuxtLink>
        </div>
      </nav>

      <div class="border-t border-slate-200 dark:border-zinc-800/80 pt-4">
        <UButton
          color="red"
          variant="ghost"
          block
          icon="i-lucide-log-out"
          class="justify-start px-4 text-xs"
          @click="logout"
        >
          {{ t('logout') }}
        </UButton>
      </div>
    </aside>

    <!-- Mobile Top Navigation Header -->
    <header class="lg:hidden flex justify-between items-center px-6 py-4 border-b border-slate-200 dark:border-zinc-800/80 backdrop-blur-md bg-white/80 dark:bg-zinc-900/20 sticky top-0 z-40">
      <div class="flex items-center space-x-2">
        <!-- Drawer Trigger -->
        <UButton
          color="gray"
          variant="ghost"
          icon="i-lucide-menu"
          size="sm"
          class="p-1"
          @click="isDrawerOpen = true"
        />
        <span class="text-md font-bold bg-gradient-to-r from-primary-500 to-primary-400 dark:from-primary-400 dark:to-primary-300 bg-clip-text text-transparent">HY-Board</span>
      </div>
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

    <!-- Mobile Bottom Navigation Bar (Stick to bottom on mobile, quick actions) -->
    <nav class="lg:hidden fixed bottom-0 left-0 right-0 z-40 border-t border-slate-200 dark:border-zinc-800/80 backdrop-blur-lg bg-white/80 dark:bg-zinc-950/80 px-4 py-2 flex justify-around items-center">
      <NuxtLink
        :to="localePath('/dashboard')"
        class="flex flex-col items-center space-y-1 transition-all duration-300 text-[9px] font-medium w-16"
        :class="route.path === localePath('/dashboard') ? 'text-primary-650 dark:text-primary-400' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white'"
      >
        <UIcon name="i-lucide-layout-dashboard" class="w-5 h-5" />
        <span>{{ t('dashboard') }}</span>
      </NuxtLink>
      <NuxtLink
        :to="localePath('/dashboard/knowledge')"
        class="flex flex-col items-center space-y-1 transition-all duration-300 text-[9px] font-medium w-16"
        :class="route.path === localePath('/dashboard/knowledge') ? 'text-primary-650 dark:text-primary-400' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white'"
      >
        <UIcon name="i-lucide-book-open" class="w-5 h-5" />
        <span>{{ t('knowledge_base') }}</span>
      </NuxtLink>
      <NuxtLink
        :to="localePath('/dashboard/subscribe')"
        class="flex flex-col items-center space-y-1 transition-all duration-300 text-[9px] font-medium w-16"
        :class="route.path === localePath('/dashboard/subscribe') ? 'text-primary-650 dark:text-primary-400' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white'"
      >
        <UIcon name="i-lucide-shopping-bag" class="w-5 h-5" />
        <span>{{ t('purchase_subscription') }}</span>
      </NuxtLink>
      <NuxtLink
        :to="localePath('/dashboard/profile')"
        class="flex flex-col items-center space-y-1 transition-all duration-300 text-[9px] font-medium w-16"
        :class="route.path === localePath('/dashboard/profile') ? 'text-primary-650 dark:text-primary-400' : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white'"
      >
        <UIcon name="i-lucide-user" class="w-5 h-5" />
        <span>{{ t('user_center') }}</span>
      </NuxtLink>
    </nav>
  </div>
</template>

<script setup>
const route = useRoute()
const router = useRouter()
const { t, locale, setLocale } = useI18n()
const localePath = useLocalePath()

const isDrawerOpen = ref(false)

// Close drawer automatically on route transition
watch(() => route.path, () => {
  isDrawerOpen.value = false
})

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

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 2px;
}
.dark .custom-scrollbar::-webkit-scrollbar-thumb {
  background: #334155;
}
</style>
