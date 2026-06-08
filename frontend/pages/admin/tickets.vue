<template>
  <div class="max-w-5xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">工單管理</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">管理客戶支援請求，點擊進入會話開啟即時聊天視窗</p>
      </div>
    </header>

    <!-- Filter and Search Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3">
      <!-- Status Filters -->
      <div class="flex gap-1.5">
        <UButton
          v-for="status in ['all', 'open', 'replied', 'closed']"
          :key="status"
          size="xs"
          :color="filterStatus === status ? 'primary' : 'neutral'"
          :variant="filterStatus === status ? 'solid' : 'ghost'"
          class="capitalize"
          @click="filterStatus = status"
        >
          {{ status === 'all' ? '全部' : status === 'open' ? '待處理' : status === 'replied' ? '已回覆' : '已關閉' }}
        </UButton>
      </div>
    </div>

    <!-- Tickets List Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-4 shadow-sm">
      <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
        <UIcon name="i-lucide-message-square" class="w-4 h-4 text-primary-500 dark:text-primary-400" />
        <span>支援工單列表</span>
      </h2>

      <div v-if="loading" class="space-y-3">
        <USkeleton class="h-10 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 3" :key="i" />
      </div>

      <div v-else-if="filteredTickets.length === 0" class="text-center py-8 text-slate-500 text-xs">
        暫無工單記錄。
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="border-b border-slate-200 dark:border-zinc-800 text-slate-500 dark:text-zinc-400 font-medium">
              <th class="py-2.5 px-3 w-16">ID</th>
              <th class="py-2.5 px-3">主題</th>
              <th class="py-2.5 px-3 w-48">提案者</th>
              <th class="py-2.5 px-3 w-32">更新時間</th>
              <th class="py-2.5 px-3 w-24 text-center">狀態</th>
              <th class="py-2.5 px-3 w-28 text-right"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in filteredTickets" :key="item.id" class="border-b border-slate-200/50 dark:border-zinc-800/40 hover:bg-slate-100/50 dark:hover:bg-zinc-900/20 text-slate-700 dark:text-zinc-200">
              <td class="py-3 px-3 font-mono">{{ item.id }}</td>
              <td class="py-3 px-3 font-semibold">{{ item.title }}</td>
              <td class="py-3 px-3 text-slate-500">{{ item.user?.email || 'Guest User' }}</td>
              <td class="py-3 px-3 text-slate-550 dark:text-zinc-400 font-mono">{{ formatDate(item.updated_at || item.created_at) }}</td>
              <td class="py-3 px-3 text-center">
                <UBadge :color="getStatusColor(item.status)" variant="soft" size="xs">
                  {{ getStatusText(item.status) }}
                </UBadge>
              </td>
              <td class="py-3 px-3 text-right">
                <UButton
                  color="primary"
                  variant="ghost"
                  icon="i-lucide-external-link"
                  size="xs"
                  @click="openChatWindow(item.id)"
                >
                  進入會話
                </UButton>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})

const { t } = useI18n()
const localePath = useLocalePath()
const token = useCookie('auth_token').value
const config = useRuntimeConfig()

const tickets = ref([])
const loading = ref(true)
const filterStatus = ref('all')

const filteredTickets = computed(() => {
  if (filterStatus.value === 'all') return tickets.value
  return tickets.value.filter(t => t.status === filterStatus.value)
})

const fetchTickets = async () => {
  try {
    const response = await $fetch(`${config.public.apiBase}/tickets`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    tickets.value = response || []
  } catch (error) {
    console.error('Failed to fetch tickets:', error)
  } finally {
    loading.value = false
  }
}

const openChatWindow = (ticketId) => {
  const width = 450
  const height = 620
  const left = (window.screen.width - width) / 2
  const top = (window.screen.height - height) / 2
  
  const chatWindow = window.open(
    localePath(`/admin/tickets-chat?id=${ticketId}`),
    `chat_${ticketId}`,
    `width=${width},height=${height},left=${left},top=${top},status=no,menubar=no,toolbar=no,location=no`
  )

  // Polling check to refresh lists on popup close if status changes
  const timer = setInterval(() => {
    if (chatWindow && chatWindow.closed) {
      clearInterval(timer)
      fetchTickets() // Refresh ticket status
    }
  }, 1000)
}

onMounted(() => {
  fetchTickets()
})

const getStatusColor = (status) => {
  if (status === 'open') return 'emerald'
  if (status === 'replied') return 'primary'
  return 'neutral'
}

const getStatusText = (status) => {
  if (status === 'open') return '待處理'
  if (status === 'replied') return '已回覆'
  return '已關閉'
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleString()
}
</script>
