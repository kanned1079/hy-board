<template>
  <div class="min-h-screen flex items-center justify-center p-4">
    <!-- Sleek premium loader if checking credentials -->
    <div v-if="checkingAuth" class="w-full max-w-sm backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-8 rounded-lg shadow-xl flex flex-col items-center justify-center space-y-4">
      <UIcon name="i-lucide-refresh-cw" class="w-8 h-8 text-primary-500 animate-spin" />
      <span class="text-xs text-slate-500 dark:text-zinc-400 font-medium tracking-wide">驗證憑據中...</span>
    </div>
    
    <!-- Login container with glassmorphic styling -->
    <div v-else class="w-full max-w-sm backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg shadow-xl relative overflow-hidden transition-all duration-500 hover:border-primary-500/20">
      
      <div class="space-y-1 text-center mb-5">
        <h2 class="text-2xl font-bold text-slate-900 dark:text-white tracking-tight">{{ t('welcome_back') }}</h2>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400">{{ t('login_desc') }}</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <UFormField :label="t('email')" name="email" class="text-slate-700 dark:text-zinc-300 text-xs">
          <UInput
            v-model="email"
            type="email"
            placeholder="name@example.com"
            icon="i-lucide-mail"
            color="primary"
            size="md"
            required
            class="mt-1 w-full"
          />
        </UFormField>

        <UFormField :label="t('password')" name="password" class="text-slate-700 dark:text-zinc-300 text-xs">
          <UInput
            v-model="password"
            type="password"
            placeholder="••••••••"
            icon="i-lucide-lock"
            color="primary"
            size="md"
            required
            class="mt-1 w-full"
          />
        </UFormField>

        <div class="pt-1">
          <UButton
            type="submit"
            color="primary"
            size="md"
            block
            :loading="loading"
            class="font-semibold transition-transform hover:scale-[1.01]"
          >
            {{ t('sign_in') }}
          </UButton>
        </div>
        <p class="text-center text-[11px] text-slate-500 mt-2">
          {{ t('dont_have_account') }} 
          <NuxtLink :to="localePath('/register')" class="text-primary-500 dark:text-primary-400 hover:underline">{{ t('register') }}</NuxtLink>
        </p>
      </form>
    </div>
  </div>
</template>

<script setup>
const { t } = useI18n()
const localePath = useLocalePath()

const email = ref('')
const password = ref('')
const loading = ref(false)
const checkingAuth = ref(false)

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

useSeoMeta({
  title: () => `${t('sign_in')} - HY-Board`
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
    // Clean up invalid cookie and show login form
    const tokenCookie = useCookie('auth_token')
    tokenCookie.value = null
    checkingAuth.value = false
  }
})

const handleLogin = async () => {
  loading.value = true
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      body: {
        query: `
          mutation Login($email: String!, $password: String!) {
            login(email: $email, password: $password) {
              token
              user {
                id
                email
                total_traffic
                used_traffic
                expired_at
                is_admin
              }
            }
          }
        `,
        variables: {
          email: email.value,
          password: password.value
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    const data = response.data.login

    // Store JWT token in a cookie
    const tokenCookie = useCookie('auth_token', { maxAge: 60 * 60 * 24 })
    tokenCookie.value = data.token

    toast.add({
      id: 'login_success',
      title: t('login_success'),
      description: t('login_success_desc'),
      color: 'success',
      timeout: 2000
    })

    // Non-linear visual redirection lag to allow toast to render smoothly
    setTimeout(() => {
      if (data.user && data.user.is_admin) {
        router.push(localePath('/admin/dashboard'))
      } else {
        router.push(localePath('/dashboard'))
      }
    }, 1000)
  } catch (error) {
    toast.add({
      id: 'login_failed',
      title: t('login_failed'),
      description: error.message || t('login_failed_desc'),
      color: 'error',
      timeout: 3000
    })
  } finally {
    loading.value = false
  }
}
</script>
