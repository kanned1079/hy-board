<template>
  <div class="min-h-screen flex flex-col items-center justify-center p-4">
    <!-- Sleek premium loader if checking credentials -->
    <div v-if="checkingAuth" class="flex flex-col items-center space-y-4">
      <UIcon name="i-lucide-refresh-cw" class="w-8 h-8 text-primary-500 animate-spin" />
      <span class="text-xs text-slate-500 dark:text-zinc-400 font-medium tracking-wide">驗證憑據中...</span>
    </div>
    
    <!-- Welcome landing content -->
    <div v-else class="text-center max-w-lg mx-auto space-y-6">
      <div class="inline-flex items-center space-x-2 bg-primary-500/10 border border-primary-500/20 px-3 py-1 rounded-full text-primary-500 dark:text-primary-400 text-xs font-medium tracking-wide">
        <span>v1.0.0 Stable</span>
      </div>
      
      <h1 class="text-4xl sm:text-6xl font-extrabold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-slate-900 via-slate-800 to-primary-600 dark:from-white dark:via-slate-100 dark:to-primary-400">
        HY-Board
      </h1>
      
      <p class="text-slate-600 dark:text-zinc-400 text-sm sm:text-base max-w-md mx-auto leading-relaxed">
        {{ t('welcome_desc') }}
      </p>

      <div class="pt-4 flex flex-col sm:flex-row justify-center items-center gap-4">
        <UButton
          :to="localePath('/login')"
          size="lg"
          color="primary"
          variant="solid"
          class="font-semibold transition-all duration-300 transform hover:scale-[1.02]"
        >
          {{ t('access_portal') }}
        </UButton>
      </div>
    </div>
  </div>
</template>

<script setup>
const { t } = useI18n()
const localePath = useLocalePath()
const router = useRouter()
const config = useRuntimeConfig()

const checkingAuth = ref(false)

// Page metadata config
useSeoMeta({
  title: () => `${t('welcome')} - HY-Board`,
  description: () => t('welcome_desc')
})

onMounted(async () => {
  const token = useCookie('auth_token').value
  if (!token) return

  checkingAuth.value = true
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          query CheckAuth {
            userInfo {
              id
              is_admin
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    const user = response.data?.userInfo
    if (user) {
      if (user.is_admin) {
        router.push(localePath('/admin/dashboard'))
      } else {
        router.push(localePath('/dashboard'))
      }
    } else {
      // Clean up invalid cookie
      const tokenCookie = useCookie('auth_token')
      tokenCookie.value = null
      checkingAuth.value = false
    }
  } catch (error) {
    // Clean up invalid cookie and show landing
    const tokenCookie = useCookie('auth_token')
    tokenCookie.value = null
    checkingAuth.value = false
  }
})
</script>
