<template>
  <div class="max-w-4xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white">{{ t('node_status') }}</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('proxy_nodes') }}</p>
      </div>
    </header>

    <!-- Skeletons -->
    <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <USkeleton class="h-24 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 4" :key="i" />
    </div>

    <!-- Empty State -->
    <div v-else-if="nodes.length === 0" class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-8 rounded-lg text-center text-slate-500 text-xs">
      {{ t('no_nodes') || 'No active nodes available.' }}
    </div>

    <!-- Nodes list -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div 
        v-for="node in nodes" 
        :key="node.id" 
        class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-3 shadow-sm hover:border-primary-500/40 transition-all duration-300"
      >
        <div class="flex justify-between items-start">
          <div>
            <h3 class="font-bold text-slate-900 dark:text-white text-sm flex items-center gap-1.5">
              <span class="inline-block w-2 h-2 rounded-full bg-primary-500 animate-pulse"></span>
              {{ node.name }}
            </h3>
            <p class="text-[10px] text-slate-500 dark:text-zinc-400 font-mono mt-0.5">{{ node.address }}:{{ node.port }}</p>
          </div>
          <UBadge color="primary" variant="subtle" size="xs">{{ node.type }}</UBadge>
        </div>

        <div class="grid grid-cols-2 gap-2 text-[10px] border-t border-slate-100 dark:border-zinc-800/50 pt-2.5">
          <div class="text-slate-500 dark:text-zinc-400">
            {{ t('traffic_rate') }}: <span class="font-semibold text-slate-800 dark:text-zinc-200 font-mono">{{ node.traffic_rate }}x</span>
          </div>
          <div class="text-slate-500 dark:text-zinc-400 text-right">
            {{ t('status') }}: <span class="font-semibold text-primary-500 dark:text-primary-400">{{ t('online') }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'user'
})

const { t } = useI18n()
const localePath = useLocalePath()

const nodes = ref([])
const loading = ref(true)

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

useSeoMeta({
  title: () => `${t('node_status')} - HY-Board`
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
          query GetNodesData {
            nodes {
              id
              name
              type
              address
              port
              traffic_rate
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    nodes.value = response.data.nodes || []
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

const logout = () => {
  const tokenCookie = useCookie('auth_token')
  tokenCookie.value = null
  router.push(localePath('/login'))
}
</script>
