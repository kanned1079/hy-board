<template>
  <div class="max-w-5xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">訂單管理</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">查看付款流水、手動完成訂單與查詢訂單詳情</p>
      </div>
    </header>

    <!-- Content Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-6 rounded-lg space-y-4 shadow-sm">
      <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
        <UIcon name="i-lucide-receipt" class="w-5 h-5 text-primary-500" />
        <span>訂單歷史記錄</span>
      </h2>

      <div class="text-center py-8 text-slate-500 text-xs">
        暫無銷售或儲值訂單記錄。
      </div>
    </section>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})
</script>
