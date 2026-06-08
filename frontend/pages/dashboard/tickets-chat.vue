<template>
  <div class="h-screen flex flex-col bg-slate-50 dark:bg-zinc-950 text-slate-900 dark:text-zinc-100 font-sans">
    <!-- Header -->
    <header class="p-4 border-b border-slate-200 dark:border-zinc-800/80 backdrop-blur-md bg-white/50 dark:bg-zinc-900/30 flex justify-between items-center shrink-0">
      <div v-if="ticket">
        <h1 class="text-sm font-extrabold text-slate-900 dark:text-white flex items-center gap-1.5">
          <span class="truncate max-w-[200px]">{{ ticket.title }}</span>
          <UBadge :color="getStatusColor(ticket.status)" variant="soft" size="xs">
            {{ getStatusText(ticket.status) }}
          </UBadge>
        </h1>
        <p class="text-[10px] text-slate-500 mt-0.5">工單編號: #{{ ticket.id }}</p>
      </div>
      <div v-else class="space-y-1">
        <USkeleton class="h-4 w-28 bg-slate-200 dark:bg-zinc-800" />
        <USkeleton class="h-3 w-36 bg-slate-200 dark:bg-zinc-800" />
      </div>

      <div class="flex items-center space-x-2">
        <span class="flex h-2 w-2 relative">
          <span :class="isConnected ? 'bg-emerald-400' : 'bg-rose-400'" class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75"></span>
          <span :class="isConnected ? 'bg-emerald-500' : 'bg-rose-500'" class="relative inline-flex rounded-full h-2 w-2"></span>
        </span>
      </div>
    </header>

    <!-- Chat Message History -->
    <div ref="chatContainer" class="flex-1 overflow-y-auto p-4 space-y-4">
      <div v-if="loading" class="space-y-4">
        <div class="flex flex-col items-start max-w-[70%]">
          <USkeleton class="h-3 w-16 mb-1 bg-slate-200 dark:bg-zinc-800" />
          <USkeleton class="h-12 w-48 rounded-2xl bg-slate-200 dark:bg-zinc-800" />
        </div>
        <div class="flex flex-col items-end max-w-[70%] ml-auto">
          <USkeleton class="h-3 w-16 mb-1 bg-slate-200 dark:bg-zinc-800" />
          <USkeleton class="h-12 w-48 rounded-2xl bg-slate-200 dark:bg-zinc-800" />
        </div>
      </div>

      <div v-else-if="ticket && ticket.messages">
        <div 
          v-for="msg in ticket.messages" 
          :key="msg.id" 
          class="flex flex-col mb-4"
          :class="msg.is_admin ? 'items-start' : 'items-end'"
        >
          <!-- Meta Info -->
          <div class="flex items-center space-x-1 mb-1 text-[9px] text-slate-400">
            <span class="font-semibold">{{ msg.is_admin ? '客服人員' : '我' }}</span>
            <span>•</span>
            <span class="font-mono">{{ formatTime(msg.created_at) }}</span>
          </div>

          <!-- Message bubble -->
          <div 
            class="max-w-[85%] rounded-2xl px-3.5 py-2 text-xs leading-relaxed shadow-sm break-all"
            :class="[
              msg.is_admin 
                ? 'bg-white dark:bg-zinc-900 border border-slate-200 dark:border-zinc-800/80 text-slate-800 dark:text-zinc-200 rounded-tl-none'
                : 'bg-primary-500 text-white rounded-tr-none'
            ]"
          >
            {{ msg.message }}
          </div>
        </div>
      </div>
    </div>

    <!-- Footer Input Bar -->
    <div class="p-3 border-t border-slate-200 dark:border-zinc-800/80 bg-white/50 dark:bg-zinc-900/30 backdrop-blur-md shrink-0">
      <div v-if="ticket && ticket.status === 'closed'" class="text-center text-xs text-slate-400 py-1.5">
        此工單已關閉，無法再發送訊息
      </div>
      <form v-else @submit.prevent="sendReply" class="flex gap-2">
        <UInput
          v-model="replyMessage"
          placeholder="輸入您的問題..."
          color="primary"
          size="sm"
          class="flex-1"
          required
        />
        <UButton
          type="submit"
          color="primary"
          size="sm"
          icon="i-lucide-send"
        >
          發送
        </UButton>
      </form>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: false
})

const route = useRoute()
const { t } = useI18n()
const token = useCookie('auth_token').value
const config = useRuntimeConfig()

const ticket = ref(null)
const loading = ref(true)
const isConnected = ref(false)
const replyMessage = ref('')
const chatContainer = ref(null)

const ticketId = parseInt(route.query.id)

let socket = null

const fetchTicketDetails = async () => {
  try {
    const response = await $fetch(`${config.public.apiBase}/tickets/${ticketId}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    ticket.value = response
  } catch (error) {
    console.error('Failed to fetch ticket:', error)
  } finally {
    loading.value = false
    scrollToBottom()
  }
}

const connectWS = () => {
  if (import.meta.server) return

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  let wsHost = window.location.host
  if (wsHost.includes('localhost:3001') || wsHost.includes('127.0.0.1:3001')) {
    wsHost = 'localhost:8080'
  }
  const wsUrl = `${protocol}//${wsHost}/api/v1/ws/tickets?token=${token}`

  socket = new WebSocket(wsUrl)

  socket.onopen = () => {
    isConnected.value = true
    console.log('User WS chat connection opened')
  }

  socket.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data)
      if ((msg.event === 'ticket_updated' || msg.event === 'ticket_created') && msg.data.id === ticketId) {
        ticket.value = msg.data
        scrollToBottom()
      }
    } catch (e) {
      console.error('WS parse error:', e)
    }
  }

  socket.onclose = () => {
    isConnected.value = false
    console.log('User WS connection closed, reconnecting...')
    setTimeout(connectWS, 3000)
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
  })
}

const sendReply = () => {
  if (!replyMessage.value.trim() || !socket) return
  const payload = {
    action: 'reply_ticket',
    ticket_id: ticketId,
    message: replyMessage.value
  }
  socket.send(JSON.stringify(payload))
  replyMessage.value = ''
}

onMounted(async () => {
  await fetchTicketDetails()
  connectWS()
})

onUnmounted(() => {
  if (socket) {
    socket.close()
  }
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

const formatTime = (dateStr) => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return ''
  return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
}
</script>

<style scoped>
::-webkit-scrollbar {
  width: 5px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: rgba(148, 163, 184, 0.2);
  border-radius: 9999px;
}
.dark ::-webkit-scrollbar-thumb {
  background: rgba(63, 63, 70, 0.4);
}
</style>
