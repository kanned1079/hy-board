<template>
  <div class="max-w-5xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">系統配置</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">管理系統核心設定、安全性與 API 配置</p>
      </div>
    </header>

    <!-- Content Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-6 rounded-lg space-y-6 shadow-sm">
      <div class="space-y-4">
        <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
          <UIcon name="i-lucide-settings" class="w-5 h-5 text-primary-500" />
          <span>核心參數配置</span>
        </h2>
        <p class="text-xs text-slate-500 dark:text-zinc-400">系統核心參數以及代理伺服器通訊密鑰配置。</p>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
          <UFormField label="站點名稱" name="site_name">
            <UInput value="HY-Board Subscription Panel" disabled color="primary" size="sm" class="w-full" />
          </UFormField>
          <UFormField label="站點網址" name="site_url">
            <UInput value="https://example.com" disabled color="primary" size="sm" class="w-full" />
          </UFormField>
          <UFormField label="JWT 金鑰" name="jwt_secret">
            <UInput value="••••••••••••••••••••••••••••••••" type="password" disabled color="primary" size="sm" class="w-full" />
          </UFormField>
          <UFormField label="UniProxy 通訊密鑰" name="uniproxy_token">
            <UInput value="••••••••••••••••••••••••" type="password" disabled color="primary" size="sm" class="w-full" />
          </UFormField>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})
</script>
