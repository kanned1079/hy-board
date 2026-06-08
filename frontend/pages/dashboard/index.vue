<template>
  <div class="max-w-4xl mx-auto space-y-5">
      <!-- Header with user profile -->
      <header class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2 border-b border-slate-200 dark:border-zinc-800/80 pb-3">
        <div>
          <h1 class="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white">{{ t('welcome_client') }}</h1>
          <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('status') }}: <span class="text-primary-500 dark:text-primary-400 font-semibold">{{ t('connected') }}</span></p>
        </div>
      </header>

      <!-- Announcements Banner Module -->
      <div v-if="announcements && announcements.length > 0" class="space-y-2.5">
        <div 
          v-for="item in announcements" 
          :key="item.id" 
          class="backdrop-blur-md bg-primary-500/5 dark:bg-primary-500/10 border border-primary-500/20 dark:border-primary-500/30 p-3.5 rounded-lg shadow-sm flex items-start gap-3 relative overflow-hidden transition-all duration-300 hover:border-primary-500/40"
        >
          <UIcon name="i-lucide-megaphone" class="w-5 h-5 text-primary-500 dark:text-primary-400 shrink-0 mt-0.5" />
          <div class="space-y-1 flex-1">
            <div class="flex justify-between items-baseline gap-2">
              <h4 class="text-xs font-bold text-slate-900 dark:text-white">{{ item.title }}</h4>
              <span class="text-[9px] text-slate-400 dark:text-zinc-500 font-mono">{{ formatDate(item.created_at) }}</span>
            </div>
            <div class="text-[11px] text-slate-600 dark:text-zinc-300 leading-relaxed markdown-body" v-html="renderMarkdown(item.content)"></div>
          </div>
        </div>
      </div>

    <!-- Stats Cards Layout -->
    <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <USkeleton class="h-24 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 4" :key="i" />
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Balance Card -->
      <div class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-3 shadow-sm">
        <h3 class="text-[10px] font-semibold text-slate-500 dark:text-zinc-400 uppercase tracking-wider">{{ t('balance') }}</h3>
        <div class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">
          ${{ (user?.balance || 0).toFixed(2) }}
        </div>
        <p class="text-[10px] text-slate-500">Available funds</p>
      </div>

      <!-- Traffic Card -->
      <div class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-3 shadow-sm">
        <h3 class="text-[10px] font-semibold text-slate-500 dark:text-zinc-400 uppercase tracking-wider">{{ t('traffic_usage') }}</h3>
        <div class="flex items-baseline space-x-1.5">
          <span class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">{{ formatTraffic(user?.used_traffic) }}</span>
          <span class="text-[10px] text-slate-500">/ {{ formatTraffic(user?.total_traffic) }}</span>
        </div>
        <!-- Progress bar using Nuxt UI -->
        <UProgress :value="trafficPercent" color="primary" size="xs" />
      </div>

      <!-- Expiration Card -->
      <div class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-3 shadow-sm">
        <h3 class="text-[10px] font-semibold text-slate-500 dark:text-zinc-400 uppercase tracking-wider">{{ t('plan_expiration') }}</h3>
        <div class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">
          {{ formatDate(user?.expired_at) }}
        </div>
        <p class="text-[10px] text-slate-500">{{ t('auto_renew_desc') }}</p>
      </div>

      <!-- Node Limits Card -->
      <div class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-3 shadow-sm">
        <h3 class="text-[10px] font-semibold text-slate-500 dark:text-zinc-400 uppercase tracking-wider">{{ t('limitations') }}</h3>
        <div class="grid grid-cols-2 gap-2">
          <div>
            <span class="block text-xl font-extrabold text-slate-900 dark:text-white">{{ user?.speed_limit === 0 ? '∞' : user?.speed_limit + 'M' }}</span>
            <span class="text-[9px] text-slate-500 uppercase tracking-wider">{{ t('speed_limit') }}</span>
          </div>
          <div>
            <span class="block text-xl font-extrabold text-slate-900 dark:text-white">{{ user?.device_limit === 0 ? '∞' : user?.device_limit }}</span>
            <span class="text-[9px] text-slate-500 uppercase tracking-wider">{{ t('max_devices') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Subscription Section -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 sm:p-5 rounded-lg space-y-4 shadow-sm">
      <div>
        <h2 class="text-lg font-bold text-slate-900 dark:text-white">{{ t('import_config') }}</h2>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('import_desc') }}</p>
      </div>

      <div class="flex flex-col sm:flex-row gap-3">
        <UInput
          :model-value="subUrl"
          readonly
          class="flex-1"
          color="primary"
          size="md"
          icon="i-lucide-link"
        />
        <UButton
          color="primary"
          size="md"
          icon="i-lucide-copy"
          @click="copySubscription"
        >
          {{ t('copy_link') }}
        </UButton>
      </div>

      <div class="flex flex-wrap gap-2">
        <UButton color="primary" variant="soft" size="sm" @click="isOpen = true">
          {{ t('show_nodes') }}
        </UButton>
      </div>
    </section>

    <!-- Custom Animated Modal detailing nodes -->
    <UModal v-model:open="isOpen">
      <div class="p-4 space-y-4 bg-white dark:bg-zinc-900 border border-slate-200 dark:border-zinc-800 rounded-lg">
        <div class="flex justify-between items-center">
          <h3 class="text-md font-bold text-slate-900 dark:text-white">{{ t('proxy_nodes') }}</h3>
          <UButton
            color="neutral"
            variant="ghost"
            icon="i-lucide-x"
            @click="isOpen = false"
          />
        </div>

        <div class="space-y-3 max-h-[350px] overflow-y-auto pr-2">
          <div v-for="node in nodes" :key="node.id" class="p-3 bg-slate-50 dark:bg-zinc-950/60 border border-slate-200/60 dark:border-zinc-800/60 rounded-lg flex justify-between items-center">
            <div>
              <p class="text-sm font-bold text-slate-900 dark:text-white">{{ node.name }}</p>
              <p class="text-[11px] text-slate-500">{{ node.address }}:{{ node.port }}</p>
            </div>
            <UBadge color="primary" variant="subtle" size="xs">{{ node.type }}</UBadge>
          </div>
        </div>
      </div>
    </UModal>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'user'
})

const { t } = useI18n()
const localePath = useLocalePath()

const user = ref(null)
const nodes = ref([])
const announcements = ref([])
const loading = ref(true)
const isOpen = ref(false)

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

useSeoMeta({
  title: () => `${t('dashboard')} - HY-Board`
})

// Fetch User and Node info
onMounted(async () => {
  const token = useCookie('auth_token').value
  if (!token) {
    router.push(localePath('/login'))
    return
  }

  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          query GetDashboardData {
            userInfo {
              id
              email
              v2ray_uuid
              total_traffic
              used_traffic
              expired_at
              speed_limit
              device_limit
              balance
            }
            nodes {
              id
              name
              type
              address
              port
              traffic_rate
            }
            announcements {
              id
              title
              content
              created_at
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    user.value = response.data.userInfo
    nodes.value = response.data.nodes
    announcements.value = response.data.announcements
  } catch (error) {
    toast.add({
      id: 'session_expired',
      title: t('session_expired'),
      description: error.message || t('login_again'),
      color: 'error'
    })
    logout()
  } finally {
    loading.value = false
  }
})

const subUrl = computed(() => {
  if (!user.value) return ''
  // Returns sub link based on user uuid token
  return `${window.location.origin}${config.public.apiBase}/client/subscribe?token=${user.value.v2ray_uuid}`
})

const trafficPercent = computed(() => {
  if (!user.value || user.value.total_traffic === 0) return 0
  return Math.min(100, Math.round((user.value.used_traffic / user.value.total_traffic) * 100))
})

const formatTraffic = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateStr) => {
  if (!dateStr) return t('lifetime')
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return t('lifetime')
  return date.toLocaleDateString()
}

const copySubscription = () => {
  navigator.clipboard.writeText(subUrl.value)
  toast.add({
    id: 'link_copied',
    title: t('copied'),
    description: t('copied_desc'),
    color: 'primary',
    timeout: 2000
  })
}

const logout = () => {
  const tokenCookie = useCookie('auth_token')
  tokenCookie.value = null
  router.push(localePath('/login'))
}
</script>
