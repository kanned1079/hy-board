<template>
  <div class="max-w-5xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">權限組管理</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">劃分不同用戶權限等級，對應不同的訂閱與訪問權限</p>
      </div>
    </header>

    <!-- Content Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-6 rounded-lg space-y-4 shadow-sm">
      <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
        <UIcon name="i-lucide-shield" class="w-5 h-5 text-primary-500" />
        <span>訪問控制權限組</span>
      </h2>
      <p class="text-xs text-slate-500">當前預設權限組：管理員 (Admin)、普通用戶 (Client)。</p>

      <div class="overflow-x-auto pt-2">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="border-b border-slate-200 dark:border-zinc-800 text-slate-500 font-medium">
              <th class="py-2 px-3">權限組名稱</th>
              <th class="py-2 px-3">等級</th>
              <th class="py-2 px-3">權限範圍描述</th>
              <th class="py-2 px-3">狀態</th>
            </tr>
          </thead>
          <tbody>
            <tr class="border-b border-slate-200/50 dark:border-zinc-800/40">
              <td class="py-2.5 px-3 font-semibold text-primary-500">Admin</td>
              <td class="py-2.5 px-3 font-mono">99</td>
              <td class="py-2.5 px-3 text-slate-500">擁有最高管理權限，管理節點、系統、財務與用戶資訊。</td>
              <td class="py-2.5 px-3"><UBadge color="success" size="xs">系統預設</UBadge></td>
            </tr>
            <tr>
              <td class="py-2.5 px-3 font-semibold">Client</td>
              <td class="py-2.5 px-3 font-mono">1</td>
              <td class="py-2.5 px-3 text-slate-500">普通用戶等級，僅能使用訂閱鏈接與閱讀知識庫。</td>
              <td class="py-2.5 px-3"><UBadge color="success" size="xs">系統預設</UBadge></td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})
</script>
