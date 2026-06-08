<template>
  <div class="max-w-5xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">主題配置</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">配置首頁風格、自訂主色調與 Glassmorphism 毛玻璃強度</p>
      </div>
    </header>

    <!-- Content Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-6 rounded-lg space-y-6 shadow-sm">
      <div class="space-y-4">
        <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
          <UIcon name="i-lucide-palette" class="w-5 h-5 text-primary-500" />
          <span>美學介面設定</span>
        </h2>
        <p class="text-xs text-slate-500 dark:text-zinc-400">配置首頁展示的主題、漸層色以及背景透明度效果。</p>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
          <UFormField label="面板主色調" name="primary_color">
            <USelect :items="['Emerald (默認)', 'Purple', 'Indigo', 'Blue', 'Rose']" value="Emerald (默認)" disabled color="primary" size="sm" class="w-full" />
          </UFormField>
          <UFormField label="首頁毛玻璃強度" name="blur_intensity">
            <USelect :items="['Medium (默認)', 'Heavy', 'Light', 'None']" value="Medium (默認)" disabled color="primary" size="sm" class="w-full" />
          </UFormField>
          <UFormField label="默認深色/淺色模式" name="default_mode">
            <USelect :items="['跟隨系統', '深色模式', '淺色模式']" value="跟隨系統" disabled color="primary" size="sm" class="w-full" />
          </UFormField>
          <UFormField label="背景動態光暈" name="dynamic_background">
            <USwitch :model-value="true" disabled color="primary" size="sm" class="pt-2" />
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
