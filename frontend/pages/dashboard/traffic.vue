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
      <USkeleton class="h-64 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" />
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

      <!-- Traffic Audit Logs -->
      <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg shadow-sm space-y-4">
        <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 border-b border-slate-100 dark:border-zinc-800/60 pb-3">
          <h2 class="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2">
            <UIcon name="i-lucide-receipt-text" class="w-5 h-5 text-primary-500" />
            <span>{{ t('traffic_audit_log') }}</span>
          </h2>
          <div class="w-full sm:w-64">
            <UInput
              v-model="searchQuery"
              icon="i-lucide-search"
              :placeholder="t('search_node')"
              size="sm"
              color="primary"
              variant="outline"
              class="w-full"
            />
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="border-b border-slate-100 dark:border-zinc-800/60 text-slate-400 dark:text-zinc-500 font-medium">
                <th class="py-2.5 px-3">{{ t('node_name') }}</th>
                <th class="py-2.5 px-3">{{ t('rate') }}</th>
                <th class="py-2.5 px-3">{{ t('upload') }}</th>
                <th class="py-2.5 px-3">{{ t('download') }}</th>
                <th class="py-2.5 px-3">{{ t('calculated_total') }}</th>
                <th class="py-2.5 px-3 text-right">{{ t('time') }}</th>
              </tr>
            </thead>
            <tbody v-if="filteredLogs.length > 0">
              <tr
                v-for="log in filteredLogs"
                :key="log.id"
                class="border-b border-slate-50 dark:border-zinc-800/30 hover:bg-slate-50/50 dark:hover:bg-zinc-950/20 transition-colors"
              >
                <td class="py-3 px-3 font-semibold text-slate-800 dark:text-zinc-200">
                  {{ log.node_name }}
                </td>
                <td class="py-3 px-3">
                  <span class="inline-flex items-center px-1.5 py-0.5 rounded-full text-[10px] font-bold bg-primary-100 dark:bg-primary-950/30 text-primary-600 dark:text-primary-400">
                    {{ log.node_rate }}x
                  </span>
                </td>
                <td class="py-3 px-3 font-mono text-slate-500 dark:text-zinc-400">
                  {{ formatTraffic(log.up) }}
                </td>
                <td class="py-3 px-3 font-mono text-slate-500 dark:text-zinc-400">
                  {{ formatTraffic(log.down) }}
                </td>
                <td class="py-3 px-3 font-mono font-bold text-slate-800 dark:text-zinc-200">
                  {{ formatTraffic(log.up + log.down) }}
                </td>
                <td class="py-3 px-3 text-right text-slate-400 dark:text-zinc-500 font-mono">
                  {{ formatDate(log.created_at) }}
                </td>
              </tr>
            </tbody>
            <tbody v-else>
              <tr>
                <td colspan="6" class="py-8 text-center text-slate-400 dark:text-zinc-500">
                  {{ t('no_logs') }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Load More Button -->
        <div v-if="hasMore" class="flex justify-center pt-3 border-t border-slate-100 dark:border-zinc-800/40">
          <UButton
            :loading="loadingMore"
            icon="i-lucide-chevrons-down"
            :label="t('load_more_logs')"
            color="primary"
            variant="ghost"
            size="sm"
            @click="loadMore"
          />
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
const logs = ref([])
const loading = ref(true)
const loadingMore = ref(false)
const searchQuery = ref('')

const limit = 20
const offset = ref(0)
const hasMore = ref(true)

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

useSeoMeta({
  title: () => `${t('traffic_details')} - HY-Board`
})

const fetchTrafficAndLogs = async (isLoadMore = false) => {
  if (isLoadMore) {
    loadingMore.value = true
  } else {
    loading.value = true
  }

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
          query GetUserTrafficAndLogs($limit: Int!, $offset: Int!) {
            userInfo {
              id
              total_traffic
              used_traffic
            }
            trafficLogs(limit: $limit, offset: $offset) {
              id
              node_name
              node_rate
              up
              down
              created_at
            }
          }
        `,
        variables: {
          limit,
          offset: offset.value
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    user.value = response.data.userInfo
    const newLogs = response.data.trafficLogs || []

    if (isLoadMore) {
      logs.value = [...logs.value, ...newLogs]
    } else {
      logs.value = newLogs
    }

    if (newLogs.length < limit) {
      hasMore.value = false
    } else {
      hasMore.value = true
    }
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
    loadingMore.value = false
  }
}

onMounted(async () => {
  await fetchTrafficAndLogs()
})

const loadMore = () => {
  offset.value += limit
  fetchTrafficAndLogs(true)
}

const filteredLogs = computed(() => {
  if (!searchQuery.value) return logs.value
  const query = searchQuery.value.toLowerCase().trim()
  return logs.value.filter(log => log.node_name.toLowerCase().includes(query))
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

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-TW', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
}

const logout = () => {
  const tokenCookie = useCookie('auth_token')
  tokenCookie.value = null
  router.push(localePath('/login'))
}
</script>
