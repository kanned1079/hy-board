<template>
  <div class="max-w-5xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">節點管理</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">新增、修改、刪除並監控節點狀態與協議設定</p>
      </div>
    </header>

    <!-- Content Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-6 rounded-lg space-y-4 shadow-sm">
      <div class="flex justify-between items-center">
        <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
          <UIcon name="i-lucide-server" class="w-5 h-5 text-primary-500" />
          <span>節點列表</span>
        </h2>
        <UButton color="primary" icon="i-lucide-plus" size="sm" disabled>添加節點</UButton>
      </div>

      <div class="text-center py-8 text-slate-500 text-xs">
        請返回首頁查看與操作當前活動節點，完整節點新增與編輯功能正在準備中。
      </div>
    </section>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})
</script>
