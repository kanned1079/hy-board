<template>
  <div class="max-w-5xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">訂閱管理</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">設計、上架並編輯不同流量、限制與效期的套餐計劃</p>
      </div>
    </header>

    <!-- Content Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-6 rounded-lg space-y-4 shadow-sm">
      <div class="flex justify-between items-center">
        <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
          <UIcon name="i-lucide-credit-card" class="w-5 h-5 text-primary-500" />
          <span>套餐訂閱計劃</span>
        </h2>
        <UButton color="primary" icon="i-lucide-plus" size="sm" disabled>新建套餐</UButton>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 pt-2">
        <div class="p-4 bg-slate-50 dark:bg-zinc-950/60 border border-slate-200 dark:border-zinc-800/60 rounded-lg space-y-2">
          <h3 class="font-bold text-slate-800 dark:text-zinc-200">基礎套餐 (Demo)</h3>
          <p class="text-[11px] text-slate-500">100 GB 流量 / 月</p>
          <p class="text-xs font-semibold text-primary-500">$5.00 / 月</p>
        </div>
        <div class="p-4 bg-slate-50 dark:bg-zinc-950/60 border border-slate-200 dark:border-zinc-800/60 rounded-lg space-y-2">
          <h3 class="font-bold text-slate-800 dark:text-zinc-200">專業套餐 (Demo)</h3>
          <p class="text-[11px] text-slate-500">500 GB 流量 / 月</p>
          <p class="text-xs font-semibold text-primary-500">$15.00 / 月</p>
        </div>
        <div class="p-4 bg-slate-50 dark:bg-zinc-950/60 border border-slate-200 dark:border-zinc-800/60 rounded-lg space-y-2 border-dashed flex items-center justify-center">
          <span class="text-xs text-slate-400">自訂更多套餐...</span>
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
