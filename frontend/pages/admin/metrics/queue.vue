<template>
  <div class="max-w-5xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">隊列監控</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">即時監控系統異步任務隊列、流量結算與 Redis 連線指標</p>
      </div>
    </header>

    <!-- Content Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-6 rounded-lg space-y-4 shadow-sm">
      <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
        <UIcon name="i-lucide-activity" class="w-5 h-5 text-primary-500" />
        <span>任務隊列狀態</span>
      </h2>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2 text-xs">
        <div class="p-4 bg-slate-50 dark:bg-zinc-950/60 border border-slate-200 dark:border-zinc-800/60 rounded-lg flex justify-between items-center">
          <div>
            <p class="text-slate-500 font-semibold">異步隊列模式</p>
            <p class="text-[10px] text-slate-400">當前流量上報結算機制</p>
          </div>
          <span class="font-bold text-slate-700 dark:text-white font-mono">SQLite (直接寫入)</span>
        </div>
        <div class="p-4 bg-slate-50 dark:bg-zinc-950/60 border border-slate-200 dark:border-zinc-800/60 rounded-lg flex justify-between items-center">
          <div>
            <p class="text-slate-500 font-semibold">待處理異步任務</p>
            <p class="text-[10px] text-slate-400">當前異步等待任務數</p>
          </div>
          <span class="font-bold text-emerald-500 font-mono">0 個</span>
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
