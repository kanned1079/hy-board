<template>
  <div class="max-w-3xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white">{{ t('user_center') }}</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('manage_account_desc') || 'Manage your account credentials and system limits' }}</p>
      </div>
    </header>

    <!-- Loading -->
    <div v-if="loading" class="space-y-4">
      <USkeleton class="h-40 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" />
      <USkeleton class="h-48 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" />
    </div>

    <div v-else class="space-y-5">
      <!-- Profile Card -->
      <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg space-y-4 shadow-sm">
        <h2 class="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2 border-b border-slate-100 dark:border-zinc-800/60 pb-2">
          <UIcon name="i-lucide-user" class="w-5 h-5 text-primary-500" />
          <span>{{ t('account_info') || 'Account Information' }}</span>
        </h2>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
          <!-- Email -->
          <div class="space-y-1">
            <span class="text-slate-400 dark:text-zinc-500 font-medium">{{ t('email') }}</span>
            <p class="font-bold text-slate-800 dark:text-zinc-200">{{ user?.email }}</p>
          </div>
          <!-- Balance -->
          <div class="space-y-1">
            <span class="text-slate-400 dark:text-zinc-500 font-medium">{{ t('balance') }}</span>
            <p class="font-bold text-emerald-650 dark:text-emerald-400 font-mono">${{ (user?.balance || 0).toFixed(2) }}</p>
          </div>
          <!-- Plan expiration -->
          <div class="space-y-1">
            <span class="text-slate-400 dark:text-zinc-500 font-medium">{{ t('plan_expiration') }}</span>
            <p class="font-bold text-slate-800 dark:text-zinc-200 font-mono">{{ formatDate(user?.expired_at) }}</p>
          </div>
        </div>
      </section>

      <!-- Credentials Card -->
      <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg space-y-4 shadow-sm">
        <h2 class="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2 border-b border-slate-100 dark:border-zinc-800/60 pb-2">
          <UIcon name="i-lucide-key" class="w-5 h-5 text-primary-500" />
          <span>{{ t('connection_credentials') || 'Connection Credentials' }}</span>
        </h2>

        <div class="space-y-4">
          <!-- UUID -->
          <div class="space-y-1">
            <span class="text-xs text-slate-400 dark:text-zinc-500 font-medium">V2ray UUID</span>
            <div class="flex gap-2">
              <UInput
                :model-value="user?.v2ray_uuid"
                readonly
                :type="showUuid ? 'text' : 'password'"
                class="flex-1 font-mono text-xs"
                color="primary"
                size="sm"
              />
              <UButton
                color="gray"
                variant="ghost"
                :icon="showUuid ? 'i-lucide-eye-slash' : 'i-lucide-eye'"
                size="xs"
                @click="showUuid = !showUuid"
              />
              <UButton
                color="primary"
                variant="soft"
                icon="i-lucide-copy"
                size="xs"
                @click="copyText(user?.v2ray_uuid, 'UUID')"
              />
            </div>
          </div>

          <!-- Trojan Password -->
          <div class="space-y-1">
            <span class="text-xs text-slate-400 dark:text-zinc-500 font-medium">Trojan Password</span>
            <div class="flex gap-2">
              <UInput
                :model-value="user?.trojan_password"
                readonly
                :type="showTrojan ? 'text' : 'password'"
                class="flex-1 font-mono text-xs"
                color="primary"
                size="sm"
              />
              <UButton
                color="gray"
                variant="ghost"
                :icon="showTrojan ? 'i-lucide-eye-slash' : 'i-lucide-eye'"
                size="xs"
                @click="showTrojan = !showTrojan"
              />
              <UButton
                color="primary"
                variant="soft"
                icon="i-lucide-copy"
                size="xs"
                @click="copyText(user?.trojan_password, 'Trojan Password')"
              />
            </div>
          </div>
        </div>
      </section>

      <!-- System limits -->
      <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg space-y-4 shadow-sm">
        <h2 class="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2 border-b border-slate-100 dark:border-zinc-800/60 pb-2">
          <UIcon name="i-lucide-sliders-horizontal" class="w-5 h-5 text-primary-500" />
          <span>{{ t('limitations') }}</span>
        </h2>

        <div class="grid grid-cols-2 gap-4 text-xs">
          <div class="p-3 bg-slate-50 dark:bg-zinc-950/40 border border-slate-150 dark:border-zinc-800/60 rounded-lg">
            <span class="text-[10px] text-slate-400 dark:text-zinc-500 uppercase tracking-wider font-semibold block">{{ t('speed_limit') }}</span>
            <span class="text-lg font-bold text-slate-800 dark:text-zinc-100 font-mono mt-1 block">
              {{ user?.speed_limit === 0 ? '∞' : user?.speed_limit + ' Mbps' }}
            </span>
          </div>
          <div class="p-3 bg-slate-50 dark:bg-zinc-950/40 border border-slate-150 dark:border-zinc-800/60 rounded-lg">
            <span class="text-[10px] text-slate-400 dark:text-zinc-500 uppercase tracking-wider font-semibold block">{{ t('max_devices') }}</span>
            <span class="text-lg font-bold text-slate-800 dark:text-zinc-100 font-mono mt-1 block">
              {{ user?.device_limit === 0 ? '∞' : user?.device_limit }}
            </span>
          </div>
        </div>
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
const showUuid = ref(false)
const showTrojan = ref(false)

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

useSeoMeta({
  title: () => `${t('user_center')} - HY-Board`
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
          query GetUserProfile {
            userInfo {
              id
              email
              v2ray_uuid
              trojan_password
              expired_at
              speed_limit
              device_limit
              balance
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

const formatDate = (dateStr) => {
  if (!dateStr) return t('lifetime')
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return t('lifetime')
  return date.toLocaleDateString()
}

const copyText = (text, name) => {
  if (!text) return
  navigator.clipboard.writeText(text)
  toast.add({
    id: 'copied_key',
    title: t('copied') || 'Copied',
    description: `${name} has been copied to clipboard.`,
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
