<template>
  <div class="max-w-5xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white">{{ t('purchase_subscription') }}</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('select_plan_desc') || 'Select the perfect data package tailored to your speed and traffic needs' }}</p>
      </div>
    </header>

    <!-- Pricing Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 pt-4">
      <!-- Plan 1: Starter -->
      <div class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg flex flex-col justify-between shadow-sm relative overflow-hidden group hover:border-primary-500/40 transition-all duration-300">
        <div class="space-y-4">
          <div>
            <h3 class="text-md font-bold text-slate-900 dark:text-white">Starter Plan</h3>
            <p class="text-[10px] text-slate-500 mt-0.5">Best for casual browsing & text messaging</p>
          </div>
          
          <div class="flex items-baseline space-x-1.5">
            <span class="text-2xl font-extrabold text-slate-900 dark:text-white">$3.99</span>
            <span class="text-[10px] text-slate-500">/ month</span>
          </div>

          <ul class="space-y-2 text-[11px] text-slate-500 dark:text-zinc-400 border-t border-slate-100 dark:border-zinc-800/60 pt-4">
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-primary-500 shrink-0" />
              <span><strong>50 GB</strong> Traffic Quota</span>
            </li>
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-primary-500 shrink-0" />
              <span><strong>30 Mbps</strong> Speed Limit</span>
            </li>
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-primary-500 shrink-0" />
              <span><strong>2</strong> Max Active Devices</span>
            </li>
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-primary-500 shrink-0" />
              <span>Access to Standard Nodes</span>
            </li>
          </ul>
        </div>

        <div class="pt-6">
          <UButton color="primary" block variant="soft" size="sm" @click="handlePurchase('Starter')">
            {{ t('order_now') || 'Order Now' }}
          </UButton>
        </div>
      </div>

      <!-- Plan 2: Pro -->
      <div class="backdrop-blur-md bg-primary-500/5 dark:bg-primary-500/10 border-2 border-primary-500/80 p-5 rounded-lg flex flex-col justify-between shadow-md relative overflow-hidden group transition-all duration-300">
        <!-- Badge -->
        <div class="absolute top-0 right-0 bg-primary-500 text-white text-[9px] font-bold uppercase tracking-wider py-1 px-3 rounded-bl-lg">
          Popular
        </div>

        <div class="space-y-4">
          <div>
            <h3 class="text-md font-bold text-slate-900 dark:text-white">Pro Plan</h3>
            <p class="text-[10px] text-slate-500 mt-0.5">Perfect for high-definition streaming & gaming</p>
          </div>
          
          <div class="flex items-baseline space-x-1.5">
            <span class="text-2xl font-extrabold text-slate-900 dark:text-white">$7.99</span>
            <span class="text-[10px] text-slate-500">/ month</span>
          </div>

          <ul class="space-y-2 text-[11px] text-slate-500 dark:text-zinc-400 border-t border-slate-100 dark:border-zinc-800/60 pt-4">
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-primary-500 shrink-0" />
              <span><strong>200 GB</strong> Traffic Quota</span>
            </li>
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-primary-500 shrink-0" />
              <span><strong>100 Mbps</strong> Speed Limit</span>
            </li>
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-primary-500 shrink-0" />
              <span><strong>5</strong> Max Active Devices</span>
            </li>
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-primary-500 shrink-0" />
              <span>Access to VIP & Premium Nodes</span>
            </li>
          </ul>
        </div>

        <div class="pt-6">
          <UButton color="primary" block size="sm" @click="handlePurchase('Pro')">
            {{ t('order_now') || 'Order Now' }}
          </UButton>
        </div>
      </div>

      <!-- Plan 3: Ultra -->
      <div class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg flex flex-col justify-between shadow-sm relative overflow-hidden group hover:border-primary-500/40 transition-all duration-300">
        <div class="space-y-4">
          <div>
            <h3 class="text-md font-bold text-slate-900 dark:text-white">Ultra Plan</h3>
            <p class="text-[10px] text-slate-500 mt-0.5">Ultimate experience with infinite power</p>
          </div>
          
          <div class="flex items-baseline space-x-1.5">
            <span class="text-2xl font-extrabold text-slate-900 dark:text-white">$14.99</span>
            <span class="text-[10px] text-slate-500">/ month</span>
          </div>

          <ul class="space-y-2 text-[11px] text-slate-500 dark:text-zinc-400 border-t border-slate-100 dark:border-zinc-800/60 pt-4">
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-primary-500 shrink-0" />
              <span><strong>1 TB (1024GB)</strong> Traffic Quota</span>
            </li>
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-primary-500 shrink-0" />
              <span><strong>Unlimited</strong> Speed Limit</span>
            </li>
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-primary-500 shrink-0" />
              <span><strong>Unlimited</strong> Active Devices</span>
            </li>
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-primary-500 shrink-0" />
              <span>Access to all worldwide Nodes</span>
            </li>
          </ul>
        </div>

        <div class="pt-6">
          <UButton color="primary" block variant="soft" size="sm" @click="handlePurchase('Ultra')">
            {{ t('order_now') || 'Order Now' }}
          </UButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'user'
})

const { t } = useI18n()
const toast = useToast()

useSeoMeta({
  title: () => `${t('purchase_subscription')} - HY-Board`
})

const handlePurchase = (planName) => {
  toast.add({
    id: `purchase_${planName.toLowerCase()}`,
    title: 'Purchase Request',
    description: `Billing connection is currently in sandbox mode. Selected plan: ${planName}`,
    color: 'primary',
    timeout: 3000
  })
}
</script>
