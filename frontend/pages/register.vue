<template>
  <div class="min-h-screen flex items-center justify-center p-4">
    <!-- Registration container with glassmorphic styling -->
    <div class="w-full max-w-sm backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg shadow-xl relative overflow-hidden transition-all duration-500 hover:border-primary-500/20">
      
      <div class="space-y-1 text-center mb-5">
        <h2 class="text-2xl font-bold text-slate-900 dark:text-white tracking-tight">{{ t('create_account') }}</h2>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400">{{ t('register_desc') }}</p>
      </div>

      <!-- Registration Disabled Message -->
      <div v-if="stopRegister" class="p-3 mb-4 rounded-lg bg-error-50 dark:bg-error-950/30 border border-error-200 dark:border-error-800 text-error-700 dark:text-error-400 text-xs text-center space-y-1">
        <p class="font-bold">註冊通道已關閉</p>
        <p>目前站點暫停開放新用戶註冊，請聯絡管理員。</p>
      </div>

      <form v-else @submit.prevent="handleRegister" class="space-y-4">
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
            {{ t('register') }}
          </UButton>
        </div>

        <p class="text-center text-[11px] text-slate-500 mt-2">
          {{ t('already_have_account') }} 
          <NuxtLink :to="localePath('/login')" class="text-primary-500 dark:text-primary-400 hover:underline">{{ t('sign_in') }}</NuxtLink>
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
const stopRegister = ref(false)
const siteName = ref('HY-Board')

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

useSeoMeta({
  title: () => `${t('register')} - HY-Board`
})

// Fetch public settings on mount
onMounted(async () => {
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      body: {
        query: `
          query GetPublicSettings {
            publicSettings {
              site_name
              stop_register
            }
          }
        `
      }
    })
    if (response.data && response.data.publicSettings) {
      stopRegister.value = response.data.publicSettings.stop_register
      siteName.value = response.data.publicSettings.site_name
    }
  } catch (err) {
    // Ignore error
  }
})

const handleRegister = async () => {
  if (stopRegister.value) {
    toast.add({
      id: 'register_disabled',
      title: '註冊已暫停',
      description: '目前暫不開放新用戶註冊',
      color: 'error'
    })
    return
  }
  loading.value = true
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      body: {
        query: `
          mutation Register($email: String!, $password: String!) {
            register(email: $email, password: $password)
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

    toast.add({
      id: 'register_success',
      title: t('reg_success'),
      description: t('reg_success_desc'),
      color: 'success',
      timeout: 3000
    })

    setTimeout(() => {
      router.push(localePath('/login'))
    }, 1500)
  } catch (error) {
    toast.add({
      id: 'register_failed',
      title: t('reg_failed'),
      description: error.message || t('reg_failed_desc'),
      color: 'error',
      timeout: 3000
    })
  } finally {
    loading.value = false
  }
}
</script>
