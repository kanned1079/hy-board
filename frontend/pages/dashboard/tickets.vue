<template>
  <div class="max-w-4xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white">{{ t('my_tickets') }}</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">提交並追蹤技術支援，點擊進入會話開啟即時聊天視窗</p>
      </div>
      <UButton
        color="primary"
        icon="i-lucide-plus"
        size="sm"
        @click="isModalOpen = true"
      >
        建立新工單
      </UButton>
    </header>

    <!-- Tickets List Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-4 shadow-sm">
      <h2 class="text-xs font-bold text-slate-900 dark:text-white flex items-center space-x-2">
        <UIcon name="i-lucide-message-square" class="w-4 h-4 text-primary-500 dark:text-primary-400" />
        <span>我的工單列表</span>
      </h2>

      <div v-if="loading" class="space-y-3">
        <USkeleton class="h-10 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 3" :key="i" />
      </div>

      <div v-else-if="tickets.length === 0" class="text-center py-8 text-slate-500 text-xs">
        您目前暫無任何工單記錄。
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="border-b border-slate-200 dark:border-zinc-800 text-slate-500 dark:text-zinc-400 font-medium">
              <th class="py-2.5 px-3 w-24">工單 ID</th>
              <th class="py-2.5 px-3">主題</th>
              <th class="py-2.5 px-3 w-40">更新時間</th>
              <th class="py-2.5 px-3 w-28 text-center">狀態</th>
              <th class="py-2.5 px-3 w-28 text-right"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="ticket in tickets" :key="ticket.id" class="border-b border-slate-200/50 dark:border-zinc-800/40 hover:bg-slate-100/50 dark:hover:bg-zinc-900/20 text-slate-700 dark:text-zinc-200">
              <td class="py-3 px-3 font-mono font-bold text-slate-500 dark:text-zinc-400">#{{ ticket.id }}</td>
              <td class="py-3 px-3 font-semibold">{{ ticket.title }}</td>
              <td class="py-3 px-3 text-slate-550 dark:text-zinc-400 font-mono">{{ formatDate(ticket.updated_at || ticket.created_at) }}</td>
              <td class="py-3 px-3 text-center">
                <UBadge :color="getStatusColor(ticket.status)" variant="soft" size="xs">
                  {{ getStatusText(ticket.status) }}
                </UBadge>
              </td>
              <td class="py-3 px-3 text-right">
                <UButton
                  color="primary"
                  variant="ghost"
                  icon="i-lucide-external-link"
                  size="xs"
                  @click="openChatWindow(ticket.id)"
                >
                  進入會話
                </UButton>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- Create Ticket Modal (wrapped in template #content to avoid template leaking) -->
    <UModal v-model:open="isModalOpen">
      <template #content>
        <div class="p-4 space-y-4 bg-white dark:bg-zinc-900 border border-slate-200 dark:border-zinc-800 rounded-lg">
          <div class="flex justify-between items-center">
            <h3 class="text-md font-bold text-slate-900 dark:text-white">建立技術支援工單</h3>
            <UButton
              color="neutral"
              variant="ghost"
              icon="i-lucide-x"
              @click="isModalOpen = false"
            />
          </div>

          <form @submit.prevent="submitTicket" class="space-y-3">
            <UFormField label="主題" name="title" class="text-slate-700 dark:text-zinc-300">
              <UInput v-model="newTicket.title" placeholder="例如：新加坡節點連線逾時" color="primary" size="sm" class="w-full" required />
            </UFormField>
            
            <UFormField label="詳細描述" name="message" class="text-slate-700 dark:text-zinc-300">
              <UTextarea v-model="newTicket.message" placeholder="請詳細描述您的問題或遇到的狀況..." color="primary" size="sm" rows="5" class="w-full" required />
            </UFormField>

            <div class="pt-2 flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" size="sm" @click="isModalOpen = false">取消</UButton>
              <UButton type="submit" color="primary" size="sm" :loading="submitLoading">提交工單</UButton>
            </div>
          </form>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'user'
})

const { t } = useI18n()
const toast = useToast()
const localePath = useLocalePath()
const token = useCookie('auth_token').value
const config = useRuntimeConfig()

const tickets = ref([])
const loading = ref(true)
const isModalOpen = ref(false)
const submitLoading = ref(false)

const newTicket = ref({
  title: '',
  message: ''
})

useSeoMeta({
  title: () => `${t('my_tickets')} - HY-Board`
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

const submitTicket = async () => {
  if (!newTicket.value.title.trim() || !newTicket.value.message.trim()) return
  submitLoading.value = true

  try {
    const response = await $fetch(`${config.public.apiBase}/tickets`, {
      method: 'POST',
      headers: { 
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        title: newTicket.value.title,
        message: newTicket.value.message
      })
    })

    toast.add({
      id: 'ticket_created',
      title: '工單已成功提交',
      description: '技術客服將在收到通知後的第一時間回覆您，您可以點擊工單開啟會話。',
      color: 'success',
      timeout: 3000
    })

    isModalOpen.value = false
    newTicket.value = { title: '', message: '' }
    
    // Refresh list and open chat immediately
    await fetchTickets()
    if (response && response.id) {
      openChatWindow(response.id)
    }
  } catch (error) {
    toast.add({
      id: 'ticket_creation_failed',
      title: '工單建立失敗',
      description: error.message || '請稍後再試',
      color: 'error'
    })
  } finally {
    submitLoading.value = false
  }
}

const openChatWindow = (ticketId) => {
  const width = 450
  const height = 620
  const left = (window.screen.width - width) / 2
  const top = (window.screen.height - height) / 2
  
  const chatWindow = window.open(
    localePath(`/dashboard/tickets-chat?id=${ticketId}`),
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
  if (status === 'open') return '處理中'
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
