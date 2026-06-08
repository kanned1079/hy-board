<template>
  <div class="max-w-5xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">路由管理</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">配置客戶端分流規則、直連與阻擋規則（Audit 審計功能）</p>
      </div>
    </header>

    <!-- Content Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-6 rounded-lg space-y-4 shadow-sm">
      <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
        <UIcon name="i-lucide-git-fork" class="w-5 h-5 text-primary-500" />
        <span>分流路由設定</span>
      </h2>
      <p class="text-xs text-slate-500">審計與自訂分流路由規則暫未開啟，目前為默認直連/代理模式。</p>
    </section>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})
</script>
