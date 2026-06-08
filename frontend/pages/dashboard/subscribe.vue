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
    <div class="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl mx-auto pt-4">
      <div v-if="loading" class="col-span-2 text-center py-10">
        <USkeleton class="h-48 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" />
      </div>

      <div v-else-if="plans.length === 0" class="col-span-2 text-center py-10 text-slate-400 text-xs">
        暫無可用的訂閱套餐上架。
      </div>

      <div 
        v-else 
        v-for="plan in plans" 
        :key="plan.id" 
        class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border p-6 rounded-lg flex flex-col justify-between shadow-sm relative overflow-hidden group hover:border-primary-500/40 transition-all duration-300"
        :class="plan.group_id === 2 ? 'border-primary-500/40 bg-primary-500/5 dark:bg-primary-500/10' : 'border-slate-200 dark:border-zinc-800/80'"
      >
        <!-- Popular tag for Group 2 or expensive ones -->
        <div v-if="plan.group_id === 2" class="absolute top-0 right-0 bg-primary-500 text-white text-[9px] font-bold uppercase tracking-wider py-1 px-3 rounded-bl-lg">
          熱門推薦
        </div>

        <div class="space-y-4">
          <div class="flex justify-between items-start">
            <div class="pr-6">
              <h3 class="text-md font-bold text-slate-900 dark:text-white">{{ plan.name }}</h3>
              <p class="text-[10px] text-slate-500 mt-0.5">{{ plan.description }}</p>
            </div>
            <UBadge :color="plan.group_id === 2 ? 'indigo' : 'sky'" variant="subtle" size="xs">
              Group {{ plan.group_id }}
            </UBadge>
          </div>
          
          <div class="flex items-baseline space-x-1.5">
            <span class="text-3xl font-extrabold text-slate-900 dark:text-white">${{ plan.price.toFixed(2) }}</span>
            <span class="text-[10px] text-slate-500">/ {{ plan.expiry_days }}天</span>
          </div>

          <ul class="space-y-3.5 text-[11px] text-slate-500 dark:text-zinc-400 border-t border-slate-100 dark:border-zinc-800/60 pt-4">
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-emerald-500 shrink-0" />
              <span><strong>{{ plan.traffic }} GB</strong> 每月流量配額</span>
            </li>
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-emerald-500 shrink-0" />
              <span>速率限制: <strong>{{ plan.speed_limit > 0 ? plan.speed_limit + ' Mbps' : '無限制' }}</strong></span>
            </li>
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-emerald-500 shrink-0" />
              <span>最大設備限制: <strong>{{ plan.device_limit > 0 ? plan.device_limit + ' 台' : '無限制' }}</strong></span>
            </li>
            <li class="flex items-center space-x-2">
              <UIcon name="i-lucide-check" class="w-4 h-4 text-emerald-500 shrink-0" />
              <span>訂閱天數: <strong>{{ plan.expiry_days }} 天</strong></span>
            </li>
          </ul>
        </div>

        <div class="pt-6 mt-4">
          <UButton 
            :color="plan.group_id === 2 ? 'primary' : 'sky'" 
            block 
            variant="soft" 
            size="sm" 
            :loading="purchaseLoadingId === plan.id"
            @click="handlePurchase(plan)"
          >
            {{ t('order_now') || '立即訂購' }}
          </UButton>
        </div>
      </div>
    </div>

    <!-- Purchase Confirmation Modal -->
    <UModal v-model:open="isConfirmOpen">
      <template #content>
        <div class="p-6 space-y-4 bg-white dark:bg-zinc-900 border border-slate-200 dark:border-zinc-800 rounded-lg shadow-lg">
          <div class="flex items-center space-x-3 text-amber-500">
            <UIcon name="i-lucide-alert-triangle" class="w-6 h-6 shrink-0" />
            <h3 class="text-md font-bold text-slate-900 dark:text-white">確認訂購訂閱套餐</h3>
          </div>
          
          <div class="text-xs text-slate-600 dark:text-zinc-300 space-y-2 leading-relaxed">
            <p>您確定要訂購 <strong>{{ selectedPlanToBuy?.name }}</strong> 嗎？</p>
            <div class="p-3 bg-slate-50 dark:bg-zinc-950/40 rounded-lg border border-slate-200/50 dark:border-zinc-800/50 space-y-1">
              <p>• 扣除餘額: <strong class="text-primary-500">${{ selectedPlanToBuy?.price.toFixed(2) }}</strong></p>
              <p>• 流量配額: <strong>{{ selectedPlanToBuy?.traffic }} GB</strong></p>
              <p>• 訂閱效期: <strong>{{ selectedPlanToBuy?.expiry_days }} 天</strong></p>
            </div>
            <p class="text-[10px] text-slate-400 dark:text-zinc-500">訂購成功後，您的流量限制與速率將立即生效，原相同權限組訂閱期將自動順延。</p>
          </div>

          <div class="pt-3 flex justify-end gap-2.5">
            <UButton color="neutral" variant="ghost" size="sm" @click="isConfirmOpen = false">取消</UButton>
            <UButton 
              color="primary" 
              size="sm" 
              :loading="purchaseLoadingId === selectedPlanToBuy?.id"
              @click="confirmPurchase"
            >
              確認付款
            </UButton>
          </div>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'user'
})

const { t } = useI18n()
const toast = useToast()

const plans = ref([])
const loading = ref(true)
const purchaseLoadingId = ref(null)
const isConfirmOpen = ref(false)
const selectedPlanToBuy = ref(null)

const router = useRouter()
const localePath = useLocalePath()
const config = useRuntimeConfig()

const fetchPlans = async () => {
  const token = useCookie('auth_token').value
  if (!token) {
    router.push(localePath('/login'))
    return
  }

  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          query GetAvailablePlans {
            plans {
              id
              name
              description
              price
              traffic
              speed_limit
              device_limit
              expiry_days
              group_id
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    plans.value = response.data.plans || []
  } catch (error) {
    toast.add({
      id: 'fetch_plans_failed',
      title: '加載失敗',
      description: error.message || '無法獲取套餐數據，請重新登入',
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchPlans()
})

const handlePurchase = (plan) => {
  selectedPlanToBuy.value = plan
  isConfirmOpen.value = true
}

const confirmPurchase = async () => {
  if (!selectedPlanToBuy.value) return
  
  const plan = selectedPlanToBuy.value
  purchaseLoadingId.value = plan.id
  const token = useCookie('auth_token').value

  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation PurchasePlan($plan_id: Int!) {
            purchasePlan(plan_id: $plan_id) {
              id
              email
              balance
            }
          }
        `,
        variables: {
          plan_id: plan.id
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'purchase_success',
      title: '訂購成功',
      description: `您已成功訂購 "${plan.name}"！`,
      color: 'success'
    })

    isConfirmOpen.value = false
    router.push(localePath('/dashboard'))
  } catch (error) {
    toast.add({
      id: 'purchase_failed',
      title: '訂購失敗',
      description: error.message || '無法完成訂購，請確認您的餘額是否充足',
      color: 'error'
    })
  } finally {
    purchaseLoadingId.value = null
  }
}
</script>
