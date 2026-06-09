<template>
  <div class="max-w-5xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">系統配置</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">管理系統核心設定、安全性、個性化、節點與郵件服務</p>
      </div>
    </header>

    <!-- Tabs Navigation -->
    <nav class="flex space-x-1 p-1 bg-slate-100 dark:bg-zinc-900/60 rounded-xl max-w-full overflow-x-auto border border-slate-200/50 dark:border-zinc-800/50">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        @click="activeTab = tab.id"
        type="button"
        class="flex items-center space-x-1.5 px-3 py-2 text-xs font-semibold rounded-lg transition-all duration-200 whitespace-nowrap cursor-pointer"
        :class="activeTab === tab.id
          ? 'bg-white dark:bg-zinc-800 text-slate-900 dark:text-white shadow-sm'
          : 'text-slate-500 dark:text-zinc-400 hover:text-slate-900 dark:hover:text-white hover:bg-white/50 dark:hover:bg-zinc-800/40'"
      >
        <UIcon :name="tab.icon" class="w-4 h-4 text-primary-500" />
        <span>{{ tab.label }}</span>
      </button>
    </nav>

    <!-- Main Content Form -->
    <form @submit.prevent="saveSettings" class="space-y-5">
      <div v-if="loading" class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-8 rounded-xl space-y-4 shadow-sm flex flex-col items-center justify-center min-h-[300px]">
        <UIcon name="i-lucide-loader-2" class="w-8 h-8 text-primary-500 animate-spin" />
        <p class="text-xs text-slate-500 dark:text-zinc-400">正在加載系統配置...</p>
      </div>

      <div v-else class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-6 rounded-xl space-y-6 shadow-sm">
        <!-- Site Settings Tab -->
        <div v-show="activeTab === 'site'" class="space-y-4">
          <div class="border-b border-slate-200/50 dark:border-zinc-800/50 pb-2">
            <h2 class="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2">
              <UIcon name="i-lucide-globe" class="w-4 h-4 text-primary-500" />
              <span>站點基礎設定</span>
            </h2>
            <p class="text-[11px] text-slate-500 dark:text-zinc-400">配置您網站的名稱、網址以及基本運營參數。</p>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <UFormField label="站點名稱" name="site_name">
              <UInput v-model="settings.site_name" placeholder="例如: HY-Board" color="primary" size="sm" class="w-full" required />
            </UFormField>
            <UFormField label="站點描述" name="site_description">
              <UInput v-model="settings.site_description" placeholder="輸入站點描述或 SEO 標籤" color="primary" size="sm" class="w-full" />
            </UFormField>
            <UFormField label="站點網址 (URL)" name="site_url">
              <UInput v-model="settings.site_url" placeholder="https://yourdomain.com" color="primary" size="sm" class="w-full" required />
            </UFormField>
            <UFormField label="用戶服務條款 (TOS) 網址" name="tos_url">
              <UInput v-model="settings.tos_url" placeholder="https://yourdomain.com/tos" color="primary" size="sm" class="w-full" />
            </UFormField>
            <UFormField label="貨幣單位" name="currency_unit">
              <UInput v-model="settings.currency_unit" placeholder="例如: TWD, CNY, USD" color="primary" size="sm" class="w-full" required />
            </UFormField>
            <UFormField label="貨幣符號" name="currency_symbol">
              <UInput v-model="settings.currency_symbol" placeholder="例如: $, NT$, ¥" color="primary" size="sm" class="w-full" required />
            </UFormField>
            <div class="md:col-span-2 flex items-center justify-between p-3 bg-slate-50 dark:bg-zinc-950/40 rounded-lg border border-slate-200/50 dark:border-zinc-800/50">
              <div>
                <p class="text-xs font-semibold text-slate-900 dark:text-white">停止新用戶註冊</p>
                <p class="text-[10px] text-slate-500 dark:text-zinc-400">開啟後，系統將關閉前台註冊功能，僅限管理員後台新增用戶。</p>
              </div>
              <USwitch v-model="settings.stop_register" color="primary" size="sm" />
            </div>
          </div>
        </div>

        <!-- Security Settings Tab -->
        <div v-show="activeTab === 'security'" class="space-y-4">
          <div class="border-b border-slate-200/50 dark:border-zinc-800/50 pb-2">
            <h2 class="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2">
              <UIcon name="i-lucide-shield-check" class="w-4 h-4 text-primary-500" />
              <span>安全與防護設定</span>
            </h2>
            <p class="text-[11px] text-slate-500 dark:text-zinc-400">配置郵箱驗證規則、防刷註冊以及限制單一 IP 註冊頻率。</p>
          </div>

          <div class="space-y-4">
            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-zinc-950/40 rounded-lg border border-slate-200/50 dark:border-zinc-800/50">
              <div>
                <p class="text-xs font-semibold text-slate-900 dark:text-white">強制郵箱驗證</p>
                <p class="text-[10px] text-slate-500 dark:text-zinc-400">開啟後，用戶註冊必須完成郵箱驗證碼驗證（需要先正確配置 SMTP 郵件伺服器）。</p>
              </div>
              <USwitch v-model="settings.email_verify" color="primary" size="sm" />
            </div>

            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-zinc-950/40 rounded-lg border border-slate-200/50 dark:border-zinc-800/50">
              <div>
                <p class="text-xs font-semibold text-slate-900 dark:text-white">禁止使用 Gmail 多別名 (Gmail Alias)</p>
                <p class="text-[10px] text-slate-500 dark:text-zinc-400">拒絕帶有 '+' 的 Gmail 位址，且在資料庫中將 '.' 點號忽略，防止使用同一個 Gmail 刷取多個帳戶。</p>
              </div>
              <USwitch v-model="settings.ban_gmail_alias" color="primary" size="sm" />
            </div>

            <div class="p-3 bg-slate-50 dark:bg-zinc-950/40 rounded-lg border border-slate-200/50 dark:border-zinc-800/50 space-y-3">
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-xs font-semibold text-slate-900 dark:text-white">IP 註冊頻率限制</p>
                  <p class="text-[10px] text-slate-500 dark:text-zinc-400">限制單個 IP 位址在指定時間內可註冊的帳戶次數，防範惡意機器人註冊。</p>
                </div>
                <USwitch v-model="settings.ip_register_limit" color="primary" size="sm" />
              </div>

              <div v-if="settings.ip_register_limit" class="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
                <UFormField label="限制次數 (默認 5)" name="ip_register_limit_count">
                  <UInput v-model.number="settings.ip_register_limit_count" type="number" color="primary" size="sm" min="1" class="w-full" required />
                </UFormField>
                <UFormField label="懲罰冷卻時間 (分鐘，默認 60)" name="ip_register_limit_penalty">
                  <UInput v-model.number="settings.ip_register_limit_penalty" type="number" color="primary" size="sm" min="1" class="w-full" required />
                </UFormField>
              </div>
            </div>
          </div>
        </div>

        <!-- Personalization Tab -->
        <div v-show="activeTab === 'personalization'" class="space-y-4">
          <div class="border-b border-slate-200/50 dark:border-zinc-800/50 pb-2">
            <h2 class="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2">
              <UIcon name="i-lucide-palette" class="w-4 h-4 text-primary-500" />
              <span>個性化外觀設定</span>
            </h2>
            <p class="text-[11px] text-slate-500 dark:text-zinc-400">自定義您面板的主題色調與背景，提供給用戶極致的美感反饋。</p>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <UFormField label="主題色調" name="theme_color">
              <USelect
                v-model="settings.theme_color"
                :items="['green', 'blue', 'indigo', 'purple', 'rose', 'orange', 'brand']"
                color="primary"
                size="sm"
                class="w-full"
              />
            </UFormField>
            <UFormField label="首頁背景 (網址或 CSS 漸變)" name="home_background">
              <UInput v-model="settings.home_background" placeholder="例如: radial-gradient(circle, #0e1726, #030712) 或圖片網址" color="primary" size="sm" class="w-full" />
            </UFormField>
          </div>
        </div>

        <!-- Node Configuration Tab -->
        <div v-show="activeTab === 'node'" class="space-y-4">
          <div class="border-b border-slate-200/50 dark:border-zinc-800/50 pb-2">
            <h2 class="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2">
              <UIcon name="i-lucide-server" class="w-4 h-4 text-primary-500" />
              <span>節點通訊與對接設定</span>
            </h2>
            <p class="text-[11px] text-slate-500 dark:text-zinc-400">配置對接節點 (XrayR 等後端) 所需的通訊金鑰與同步頻率參數。</p>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="md:col-span-2">
              <UFormField label="UniProxy 通訊密鑰 (Token)" name="uniproxy_token">
                <div class="relative flex items-center">
                  <UInput
                    v-model="settings.uniproxy_token"
                    :type="showToken ? 'text' : 'password'"
                    placeholder="輸入對接通訊密鑰"
                    color="primary"
                    size="sm"
                    class="w-full pr-10"
                    required
                  />
                  <button
                    type="button"
                    class="absolute right-3 text-slate-400 hover:text-slate-600 dark:hover:text-white cursor-pointer"
                    @click="showToken = !showToken"
                  >
                    <UIcon :name="showToken ? 'i-lucide-eye-off' : 'i-lucide-eye'" class="w-4 h-4" />
                  </button>
                </div>
              </UFormField>
            </div>
            <UFormField label="節點拉取輪詢間隔 (秒，默認 60)" name="node_pull_interval">
              <UInput v-model.number="settings.node_pull_interval" type="number" color="primary" size="sm" min="10" class="w-full" required />
            </UFormField>
            <UFormField label="節點推送輪詢間隔 (秒，默認 60)" name="node_push_interval">
              <UInput v-model.number="settings.node_push_interval" type="number" color="primary" size="sm" min="10" class="w-full" required />
            </UFormField>
          </div>
        </div>

        <!-- SMTP Email Tab -->
        <div v-show="activeTab === 'email'" class="space-y-4">
          <div class="border-b border-slate-200/50 dark:border-zinc-800/50 pb-2">
            <h2 class="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2">
              <UIcon name="i-lucide-mail" class="w-4 h-4 text-primary-500" />
              <span>SMTP 郵件伺服器設定</span>
            </h2>
            <p class="text-[11px] text-slate-500 dark:text-zinc-400">配置系統發送郵件通知、驗證碼與密碼重設的 SMTP 連線資訊。</p>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <UFormField label="SMTP 伺服器位址" name="smtp_host">
              <UInput v-model="settings.smtp_host" placeholder="例如: smtp.gmail.com" color="primary" size="sm" class="w-full" />
            </UFormField>
            <UFormField label="SMTP 服務端口 (Port)" name="smtp_port">
              <UInput v-model.number="settings.smtp_port" type="number" color="primary" size="sm" class="w-full" />
            </UFormField>
            <UFormField label="SMTP 加密方式" name="smtp_encryption">
              <USelect
                v-model="settings.smtp_encryption"
                :items="['SSL', 'STARTTLS', 'none']"
                color="primary"
                size="sm"
                class="w-full"
              />
            </UFormField>
            <UFormField label="發件人郵箱位址 (From Address)" name="smtp_from">
              <UInput v-model="settings.smtp_from" placeholder="your-email@domain.com" type="email" color="primary" size="sm" class="w-full" />
            </UFormField>
            <UFormField label="SMTP 驗證帳戶" name="smtp_username">
              <UInput v-model="settings.smtp_username" placeholder="your-email@domain.com" color="primary" size="sm" class="w-full" />
            </UFormField>
            <UFormField label="SMTP 驗證密碼" name="smtp_password">
              <UInput v-model="settings.smtp_password" type="password" placeholder="請輸入密碼或授權碼" color="primary" size="sm" class="w-full" />
            </UFormField>
          </div>

          <!-- Send Test Email section -->
          <div class="mt-6 p-4 rounded-lg bg-slate-50 dark:bg-zinc-950/60 border border-slate-200/50 dark:border-zinc-800/80 space-y-3">
            <h3 class="text-xs font-bold text-slate-900 dark:text-white">發送測試郵件</h3>
            <p class="text-[10px] text-slate-500 dark:text-zinc-400">儲存變更後，您可以輸入一個外部郵箱地址來測試您的 SMTP 連線是否配置正確。</p>
            <div class="flex items-end space-x-2">
              <UFormField label="收件人電子郵箱" name="test_email" class="flex-1">
                <UInput v-model="testEmailAddress" type="email" placeholder="example@domain.com" color="primary" size="sm" class="w-full" />
              </UFormField>
              <UButton
                type="button"
                color="neutral"
                variant="solid"
                size="sm"
                icon="i-lucide-send"
                :loading="testEmailLoading"
                :disabled="!testEmailAddress"
                @click="sendTestEmail"
              >
                發送測試
              </UButton>
            </div>
          </div>
        </div>

        <!-- APP Downloads Tab -->
        <div v-show="activeTab === 'app'" class="space-y-4">
          <div class="border-b border-slate-200/50 dark:border-zinc-800/50 pb-2">
            <h2 class="text-sm font-bold text-slate-900 dark:text-white flex items-center space-x-2">
              <UIcon name="i-lucide-smartphone" class="w-4 h-4 text-primary-500" />
              <span>客戶端下載連結</span>
            </h2>
            <p class="text-[11px] text-slate-500 dark:text-zinc-400">配置用戶前台「下載客戶端」卡片中各個作業系統 APP 的直接下載位址。</p>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <UFormField label="Windows 客戶端網址" name="app_win">
              <UInput v-model="settings.app_win" placeholder="https://..." color="primary" size="sm" class="w-full" />
            </UFormField>
            <UFormField label="macOS 客戶端網址" name="app_macos">
              <UInput v-model="settings.app_macos" placeholder="https://..." color="primary" size="sm" class="w-full" />
            </UFormField>
            <UFormField label="Linux 客戶端網址" name="app_linux">
              <UInput v-model="settings.app_linux" placeholder="https://..." color="primary" size="sm" class="w-full" />
            </UFormField>
            <UFormField label="Android 客戶端網址" name="app_android">
              <UInput v-model="settings.app_android" placeholder="https://..." color="primary" size="sm" class="w-full" />
            </UFormField>
            <UFormField label="iOS 客戶端網址" name="app_ios">
              <UInput v-model="settings.app_ios" placeholder="https://..." color="primary" size="sm" class="w-full" />
            </UFormField>
          </div>
        </div>

        <!-- Save button at the bottom of the card -->
        <div class="pt-4 border-t border-slate-200/50 dark:border-zinc-800/50 flex justify-end">
          <UButton
            type="submit"
            color="primary"
            size="sm"
            icon="i-lucide-save"
            :loading="saveLoading"
          >
            儲存所有設定
          </UButton>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})

const tabs = [
  { id: 'site', label: '站點設定', icon: 'i-lucide-globe' },
  { id: 'security', label: '安全防護', icon: 'i-lucide-shield-check' },
  { id: 'personalization', label: '個性化', icon: 'i-lucide-palette' },
  { id: 'node', label: '節點配置', icon: 'i-lucide-server' },
  { id: 'email', label: '郵件發送', icon: 'i-lucide-mail' },
  { id: 'app', label: 'APP 下載', icon: 'i-lucide-smartphone' }
]

const activeTab = ref('site')
const loading = ref(true)
const saveLoading = ref(false)
const showToken = ref(false)

const settings = ref({
  site_name: '',
  site_description: '',
  site_url: '',
  tos_url: '',
  stop_register: false,
  currency_unit: 'CNY',
  currency_symbol: '¥',
  email_verify: false,
  ban_gmail_alias: false,
  ip_register_limit: false,
  ip_register_limit_count: 5,
  ip_register_limit_penalty: 60,
  theme_color: 'green',
  home_background: '',
  uniproxy_token: '',
  node_pull_interval: 60,
  node_push_interval: 60,
  smtp_host: '',
  smtp_port: 465,
  smtp_encryption: 'SSL',
  smtp_username: '',
  smtp_password: '',
  smtp_from: '',
  app_win: '',
  app_macos: '',
  app_linux: '',
  app_android: '',
  app_ios: ''
})

const testEmailAddress = ref('')
const testEmailLoading = ref(false)

const toast = useToast()
const router = useRouter()
const localePath = useLocalePath()
const config = useRuntimeConfig()

// SEO Meta
useSeoMeta({
  title: '系統配置 - HY-Board'
})

const getHeaders = () => {
  const token = useCookie('auth_token').value
  if (!token) {
    router.push(localePath('/login'))
    return {}
  }
  return {
    Authorization: `Bearer ${token}`
  }
}

// Fetch settings
const fetchSettings = async () => {
  try {
    const headers = getHeaders()
    if (!headers.Authorization) return

    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers,
      body: {
        query: `
          query GetSystemSettings {
            systemSettings {
              site_name
              site_description
              site_url
              tos_url
              stop_register
              currency_unit
              currency_symbol
              email_verify
              ban_gmail_alias
              ip_register_limit
              ip_register_limit_count
              ip_register_limit_penalty
              theme_color
              home_background
              uniproxy_token
              node_pull_interval
              node_push_interval
              smtp_host
              smtp_port
              smtp_encryption
              smtp_username
              smtp_password
              smtp_from
              app_win
              app_macos
              app_linux
              app_android
              app_ios
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    if (response.data && response.data.systemSettings) {
      settings.value = { ...settings.value, ...response.data.systemSettings }
    }
  } catch (error) {
    toast.add({
      id: 'settings_fetch_error',
      title: '讀取系統配置失敗',
      description: error.message,
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}

// Save settings
const saveSettings = async () => {
  saveLoading.value = true
  try {
    const headers = getHeaders()
    if (!headers.Authorization) return

    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers,
      body: {
        query: `
          mutation UpdateSystemSettings(
            $site_name: String
            $site_description: String
            $site_url: String
            $tos_url: String
            $stop_register: Boolean
            $currency_unit: String
            $currency_symbol: String
            $email_verify: Boolean
            $ban_gmail_alias: Boolean
            $ip_register_limit: Boolean
            $ip_register_limit_count: Int
            $ip_register_limit_penalty: Int
            $theme_color: String
            $home_background: String
            $uniproxy_token: String
            $node_pull_interval: Int
            $node_push_interval: Int
            $smtp_host: String
            $smtp_port: Int
            $smtp_encryption: String
            $smtp_username: String
            $smtp_password: String
            $smtp_from: String
            $app_win: String
            $app_macos: String
            $app_linux: String
            $app_android: String
            $app_ios: String
          ) {
            updateSystemSettings(
              site_name: $site_name
              site_description: $site_description
              site_url: $site_url
              tos_url: $tos_url
              stop_register: $stop_register
              currency_unit: $currency_unit
              currency_symbol: $currency_symbol
              email_verify: $email_verify
              ban_gmail_alias: $ban_gmail_alias
              ip_register_limit: $ip_register_limit
              ip_register_limit_count: $ip_register_limit_count
              ip_register_limit_penalty: $ip_register_limit_penalty
              theme_color: $theme_color
              home_background: $home_background
              uniproxy_token: $uniproxy_token
              node_pull_interval: $node_pull_interval
              node_push_interval: $node_push_interval
              smtp_host: $smtp_host
              smtp_port: $smtp_port
              smtp_encryption: $smtp_encryption
              smtp_username: $smtp_username
              smtp_password: $smtp_password
              smtp_from: $smtp_from
              app_win: $app_win
              app_macos: $app_macos
              app_linux: $app_linux
              app_android: $app_android
              app_ios: $app_ios
            ) {
              site_name
            }
          }
        `,
        variables: {
          site_name: settings.value.site_name,
          site_description: settings.value.site_description,
          site_url: settings.value.site_url,
          tos_url: settings.value.tos_url,
          stop_register: settings.value.stop_register,
          currency_unit: settings.value.currency_unit,
          currency_symbol: settings.value.currency_symbol,
          email_verify: settings.value.email_verify,
          ban_gmail_alias: settings.value.ban_gmail_alias,
          ip_register_limit: settings.value.ip_register_limit,
          ip_register_limit_count: settings.value.ip_register_limit_count,
          ip_register_limit_penalty: settings.value.ip_register_limit_penalty,
          theme_color: settings.value.theme_color,
          home_background: settings.value.home_background,
          uniproxy_token: settings.value.uniproxy_token,
          node_pull_interval: settings.value.node_pull_interval,
          node_push_interval: settings.value.node_push_interval,
          smtp_host: settings.value.smtp_host,
          smtp_port: settings.value.smtp_port,
          smtp_encryption: settings.value.smtp_encryption,
          smtp_username: settings.value.smtp_username,
          smtp_password: settings.value.smtp_password,
          smtp_from: settings.value.smtp_from,
          app_win: settings.value.app_win,
          app_macos: settings.value.app_macos,
          app_linux: settings.value.app_linux,
          app_android: settings.value.app_android,
          app_ios: settings.value.app_ios
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'settings_save_success',
      title: '儲存成功',
      description: '系統配置已成功儲存並生效！',
      color: 'success'
    })
  } catch (error) {
    toast.add({
      id: 'settings_save_error',
      title: '儲存失敗',
      description: error.message,
      color: 'error'
    })
  } finally {
    saveLoading.value = false
  }
}

// Send test email
const sendTestEmail = async () => {
  if (!testEmailAddress.value) return
  testEmailLoading.value = true
  try {
    const headers = getHeaders()
    if (!headers.Authorization) return

    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers,
      body: {
        query: `
          mutation SendTestEmail($to: String!) {
            sendTestEmail(to: $to)
          }
        `,
        variables: {
          to: testEmailAddress.value
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    if (response.data && response.data.sendTestEmail) {
      toast.add({
        id: 'test_email_success',
        title: '發送成功',
        description: `測試郵件已成功發送至 ${testEmailAddress.value}，請檢查收件匣！`,
        color: 'success'
      })
    } else {
      throw new Error('發送失敗，請確認 SMTP 設定')
    }
  } catch (error) {
    toast.add({
      id: 'test_email_error',
      title: '發送測試郵件失敗',
      description: error.message,
      color: 'error'
    })
  } finally {
    testEmailLoading.value = false
  }
}

onMounted(() => {
  fetchSettings()
})
</script>
