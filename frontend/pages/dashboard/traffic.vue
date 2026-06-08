<template>
  <div class="max-w-3xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white">{{ t('traffic_details') }}</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('traffic_details_desc') || 'Monitor your data usage statistics in real-time' }}</p>
      </div>
    </header>

    <!-- Skeletons -->
    <div v-if="loading" class="space-y-4">
      <USkeleton class="h-32 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" />
      <USkeleton class="h-48 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" />
    </div>

    <div v-else class="space-y-5">
      <!-- Traffic Stats Card -->
      <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg shadow-sm space-y-5">
        <h2 class="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2 border-b border-slate-100 dark:border-zinc-800/60 pb-2">
          <UIcon name="i-lucide-bar-chart-3" class="w-5 h-5 text-primary-500" />
          <span>{{ t('traffic_usage') }}</span>
        </h2>

        <div class="space-y-4">
          <!-- Progress Bar -->
          <div class="space-y-1.5">
            <div class="flex justify-between text-xs font-semibold text-slate-700 dark:text-zinc-300">
              <span>{{ t('used') || 'Used' }} ({{ trafficPercent }}%)</span>
              <span>{{ t('remaining') || 'Remaining' }} ({{ 100 - trafficPercent }}%)</span>
            </div>
            <UProgress :value="trafficPercent" color="primary" size="sm" />
          </div>

          <!-- Numbers Grid -->
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs pt-2">
            <div class="p-4 bg-slate-50 dark:bg-zinc-950/40 border border-slate-150 dark:border-zinc-800/65 rounded-lg space-y-1">
              <span class="text-slate-450 dark:text-zinc-500 font-medium uppercase tracking-wider text-[9px]">{{ t('total_traffic') }}</span>
              <p class="text-md font-bold text-slate-800 dark:text-zinc-200 font-mono">{{ formatTraffic(user?.total_traffic) }}</p>
            </div>
            <div class="p-4 bg-slate-50 dark:bg-zinc-950/40 border border-slate-150 dark:border-zinc-800/65 rounded-lg space-y-1">
              <span class="text-slate-450 dark:text-zinc-500 font-medium uppercase tracking-wider text-[9px]">{{ t('used_traffic') }}</span>
              <p class="text-md font-bold text-primary-600 dark:text-primary-400 font-mono">{{ formatTraffic(user?.used_traffic) }}</p>
            </div>
            <div class="p-4 bg-slate-50 dark:bg-zinc-950/40 border border-slate-150 dark:border-zinc-800/65 rounded-lg space-y-1">
              <span class="text-slate-450 dark:text-zinc-500 font-medium uppercase tracking-wider text-[9px]">{{ t('remaining_traffic') || 'Remaining Traffic' }}</span>
              <p class="text-md font-bold text-slate-800 dark:text-zinc-200 font-mono">{{ formatTraffic(remainingTraffic) }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Usage Tips / Notice -->
      <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg shadow-sm space-y-3">
        <h3 class="text-xs font-bold text-slate-900 dark:text-white flex items-center space-x-1.5">
          <UIcon name="i-lucide-info" class="w-4.5 h-4.5 text-primary-500" />
          <span>{{ t('traffic_tips') || 'Usage Details & Notice' }}</span>
        </h3>
        <ul class="text-[11px] text-slate-500 dark:text-zinc-400 space-y-2 list-disc pl-4 leading-relaxed">
          <li>{{ t('traffic_tip_1') || 'Traffic usage calculations include both upload and download data streams multiplied by the respective node multiplier rate.' }}</li>
          <li>{{ t('traffic_tip_2') || 'Your traffic quota resets automatically upon billing plan renewal or manual subscription reload.' }}</li>
          <li>{{ t('traffic_tip_3') || 'Please monitor your background tasks. Once your usage exceeds the 100% threshold, proxy connections will be automatically suspended.' }}</li>
        </ul>
      </section>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'user'
})

const { t } = useI18n()
const localePath = useLocalePath()

const user = ref(null)
const loading = ref(true)

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

useSeoMeta({
  title: () => `${t('traffic_details')} - HY-Board`
})

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
          query GetUserTraffic {
            userInfo {
              id
              total_traffic
              used_traffic
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    user.value = response.data.userInfo
  } catch (error) {
    toast.add({
      id: 'session_expired',
      title: t('session_expired'),
      description: error.message || t('login_again'),
      color: 'red'
    })
    logout()
  } finally {
    loading.value = false
  }
})

const remainingTraffic = computed(() => {
  if (!user.value) return 0
  const remaining = user.value.total_traffic - user.value.used_traffic
  return remaining > 0 ? remaining : 0
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

const logout = () => {
  const tokenCookie = useCookie('auth_token')
  tokenCookie.value = null
  router.push(localePath('/login'))
}
</script>
