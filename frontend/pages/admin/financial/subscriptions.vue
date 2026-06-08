<template>
  <div class="max-w-5xl mx-auto space-y-5">
    
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">訂閱管理</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">設計、上架並編輯不同流量、限制與效期的套餐計劃</p>
      </div>
      <UButton
        color="primary"
        icon="i-lucide-plus"
        size="sm"
        @click="openAddModal"
      >
        添加套餐
      </UButton>
    </header>

    <!-- Plans Grid -->
    <section class="space-y-4">
      <div v-if="loading" class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <USkeleton class="h-56 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 3" :key="i" />
      </div>

      <div v-else-if="plans.length === 0" class="text-center py-12 backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 rounded-lg text-slate-400 text-xs">
        暫無已配置的套餐計劃，請點擊右上角「添加套餐」創建。
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <div 
          v-for="plan in plans" 
          :key="plan.id" 
          class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border p-5 rounded-lg flex flex-col justify-between shadow-sm relative overflow-hidden group hover:border-primary-500/40 transition-all duration-300"
          :class="plan.show ? 'border-slate-200 dark:border-zinc-800/80' : 'border-slate-200 dark:border-zinc-800/80 opacity-60'"
        >
          <!-- Hidden Badge -->
          <div v-if="!plan.show" class="absolute top-0 right-0 bg-rose-500 text-white text-[9px] font-bold uppercase tracking-wider py-0.5 px-2 rounded-bl">
            已下架
          </div>

          <div class="space-y-3.5">
            <div class="flex justify-between items-start">
              <div class="pr-6">
                <h3 class="text-sm font-bold text-slate-900 dark:text-white line-clamp-1">{{ plan.name }}</h3>
                <p class="text-[10px] text-slate-500 mt-0.5 line-clamp-2">{{ plan.description || '暫無描述' }}</p>
              </div>
              <UBadge color="sky" variant="subtle" size="xs" class="shrink-0">
                {{ getGroupName(plan.group_id) }}
              </UBadge>
            </div>
            
            <div class="flex items-baseline space-x-1 border-b border-slate-100 dark:border-zinc-800/50 pb-3">
              <span class="text-2xl font-extrabold text-slate-900 dark:text-white">${{ plan.price.toFixed(2) }}</span>
              <span class="text-[10px] text-slate-500">/ {{ plan.expiry_days }}天</span>
            </div>

            <ul class="space-y-2.5 text-[11px] text-slate-600 dark:text-zinc-400 pt-1">
              <li class="flex items-center space-x-2">
                <UIcon name="i-lucide-database" class="w-3.5 h-3.5 text-primary-500 shrink-0" />
                <span>流量配額: <strong>{{ plan.traffic }} GB</strong></span>
              </li>
              <li class="flex items-center space-x-2">
                <UIcon name="i-lucide-gauge" class="w-3.5 h-3.5 text-primary-500 shrink-0" />
                <span>速率限制: <strong>{{ plan.speed_limit > 0 ? plan.speed_limit + ' Mbps' : '無限制' }}</strong></span>
              </li>
              <li class="flex items-center space-x-2">
                <UIcon name="i-lucide-laptop" class="w-3.5 h-3.5 text-primary-500 shrink-0" />
                <span>設備限制: <strong>{{ plan.device_limit > 0 ? plan.device_limit + ' 台' : '無限制' }}</strong></span>
              </li>
              <li class="flex items-center space-x-2">
                <UIcon name="i-lucide-calendar" class="w-3.5 h-3.5 text-primary-500 shrink-0" />
                <span>有效期限: <strong>{{ plan.expiry_days }} 天</strong></span>
              </li>
            </ul>
          </div>

          <div class="pt-5 mt-4 border-t border-slate-100 dark:border-zinc-800/40 flex justify-end gap-2">
            <UButton
              color="neutral"
              variant="ghost"
              icon="i-lucide-edit"
              size="xs"
              @click="openEditModal(plan)"
            >
              編輯
            </UButton>
            <UButton
              color="red"
              variant="ghost"
              icon="i-lucide-trash-2"
              size="xs"
              :loading="deleteLoadingId === plan.id"
              @click="deletePlan(plan.id)"
            >
              刪除
            </UButton>
          </div>
        </div>
      </div>
    </section>

    <!-- Add/Edit Plan Modal -->
    <UModal v-model:open="isPlanModalOpen">
      <template #content>
        <div class="p-4 space-y-4 bg-white dark:bg-zinc-900 border border-slate-200 dark:border-zinc-800 rounded-lg">
          <div class="flex justify-between items-center">
            <h3 class="text-md font-bold text-slate-900 dark:text-white">
              {{ isEditMode ? '編輯套餐計劃' : '創建新套餐計劃' }}
            </h3>
            <UButton
              color="neutral"
              variant="ghost"
              icon="i-lucide-x"
              @click="isPlanModalOpen = false"
            />
          </div>

          <form @submit.prevent="handleSubmit" class="space-y-3">
            <UFormField label="套餐名稱" name="name" class="text-slate-700 dark:text-zinc-300">
              <UInput v-model="newPlan.name" placeholder="S1 亞太標準版" color="primary" size="sm" class="w-full" required />
            </UFormField>

            <UFormField label="套餐描述" name="description" class="text-slate-700 dark:text-zinc-300">
              <UInput v-model="newPlan.description" placeholder="適合日常網頁瀏覽、社交媒體與普通通訊" color="primary" size="sm" class="w-full" />
            </UFormField>

            <div class="grid grid-cols-2 gap-3">
              <UFormField label="價格 (USD)" name="price" class="text-slate-700 dark:text-zinc-300">
                <UInput v-model.number="newPlan.price" type="number" step="0.01" min="0" placeholder="5.99" color="primary" size="sm" class="w-full" required />
              </UFormField>
              <UFormField label="流量 (GB)" name="traffic" class="text-slate-700 dark:text-zinc-300">
                <UInput v-model.number="newPlan.traffic" type="number" min="1" placeholder="150" color="primary" size="sm" class="w-full" required />
              </UFormField>
            </div>

            <div class="grid grid-cols-3 gap-3">
              <UFormField label="速率限制 (Mbps)" name="speed_limit" class="text-slate-700 dark:text-zinc-300">
                <UInput v-model.number="newPlan.speed_limit" type="number" min="0" placeholder="0 = 無限制" color="primary" size="sm" class="w-full" required />
              </UFormField>
              <UFormField label="最大設備數" name="device_limit" class="text-slate-700 dark:text-zinc-300">
                <UInput v-model.number="newPlan.device_limit" type="number" min="0" placeholder="0 = 無限制" color="primary" size="sm" class="w-full" required />
              </UFormField>
              <UFormField label="有效天數" name="expiry_days" class="text-slate-700 dark:text-zinc-300">
                <UInput v-model.number="newPlan.expiry_days" type="number" min="1" placeholder="30" color="primary" size="sm" class="w-full" required />
              </UFormField>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <UFormField label="對應權限組" name="group_id" class="text-slate-700 dark:text-zinc-300">
                <select v-model.number="newPlan.group_id" class="w-full bg-white dark:bg-zinc-950 border border-slate-200 dark:border-zinc-800 rounded-md py-1.5 px-2.5 text-xs focus:outline-none focus:ring-1 focus:ring-primary-500 text-slate-700 dark:text-zinc-300">
                  <option v-for="g in groups" :key="g.id" :value="g.id">
                    {{ g.name }} (Group {{ g.id }})
                  </option>
                </select>
              </UFormField>
              <UFormField label="是否上架可見" name="show" class="text-slate-700 dark:text-zinc-300 flex items-center justify-between pt-5">
                <USwitch v-model="newPlan.show" color="primary" />
              </UFormField>
            </div>

            <div class="pt-2 flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" size="sm" @click="isPlanModalOpen = false">取消</UButton>
              <UButton type="submit" color="primary" size="sm" :loading="modalLoading">
                {{ isEditMode ? '保存修改' : '創建套餐' }}
              </UButton>
            </div>
          </form>
        </div>
      </template>
    </UModal>

  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})

const plans = ref([])
const groups = ref([])
const loading = ref(true)
const modalLoading = ref(false)
const deleteLoadingId = ref(null)
const isPlanModalOpen = ref(false)
const isEditMode = ref(false)
const editingPlanId = ref(null)

const toast = useToast()
const router = useRouter()
const localePath = useLocalePath()
const config = useRuntimeConfig()

useSeoMeta({
  title: '套餐計劃與訂閱管理 - HY Admin'
})

const getGroupName = (groupId) => {
  const g = groups.value.find(item => item.id === groupId)
  return g ? g.name : `Group ${groupId}`
}

const newPlan = ref({
  name: '',
  description: '',
  price: 5.99,
  traffic: 150,
  speed_limit: 0,
  device_limit: 0,
  expiry_days: 30,
  group_id: 1,
  show: true
})

const openAddModal = () => {
  isEditMode.value = false
  editingPlanId.value = null
  newPlan.value = {
    name: '',
    description: '',
    price: 5.99,
    traffic: 150,
    speed_limit: 0,
    device_limit: 0,
    expiry_days: 30,
    group_id: groups.value[0]?.id || 1,
    show: true
  }
  isPlanModalOpen.value = true
}

const openEditModal = (plan) => {
  isEditMode.value = true
  editingPlanId.value = plan.id
  newPlan.value = {
    name: plan.name,
    description: plan.description || '',
    price: plan.price,
    traffic: plan.traffic,
    speed_limit: plan.speed_limit,
    device_limit: plan.device_limit,
    expiry_days: plan.expiry_days,
    group_id: plan.group_id,
    show: plan.show
  }
  isPlanModalOpen.value = true
}

const fetchData = async () => {
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
          query GetPlansData {
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
              show
              created_at
            }
            groups {
              id
              name
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    plans.value = response.data.plans || []
    groups.value = response.data.groups || []
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
  fetchData()
})

const handleSubmit = async () => {
  if (isEditMode.value) {
    await savePlanChanges()
  } else {
    await createPlan()
  }
}

const createPlan = async () => {
  modalLoading.value = true
  const token = useCookie('auth_token').value

  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation CreatePlan($name: String!, $description: String, $price: Float!, $traffic: Int!, $speed_limit: Int, $device_limit: Int, $expiry_days: Int, $group_id: Int!, $show: Boolean) {
            createPlan(name: $name, description: $description, price: $price, traffic: $traffic, speed_limit: $speed_limit, device_limit: $device_limit, expiry_days: $expiry_days, group_id: $group_id, show: $show) {
              id
              name
            }
          }
        `,
        variables: {
          name: newPlan.value.name,
          description: newPlan.value.description,
          price: newPlan.value.price,
          traffic: newPlan.value.traffic,
          speed_limit: newPlan.value.speed_limit,
          device_limit: newPlan.value.device_limit,
          expiry_days: newPlan.value.expiry_days,
          group_id: newPlan.value.group_id,
          show: newPlan.value.show
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'plan_created',
      title: '套餐已上架',
      description: `成功添加套餐：${newPlan.value.name}`,
      color: 'success'
    })

    isPlanModalOpen.value = false
    await fetchData()
  } catch (error) {
    toast.add({
      id: 'create_plan_failed',
      title: '創建失敗',
      description: error.message || '創建套餐失敗',
      color: 'error'
    })
  } finally {
    modalLoading.value = false
  }
}

const savePlanChanges = async () => {
  modalLoading.value = true
  const token = useCookie('auth_token').value

  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation UpdatePlan($id: Int!, $name: String, $description: String, $price: Float, $traffic: Int, $speed_limit: Int, $device_limit: Int, $expiry_days: Int, $group_id: Int, $show: Boolean) {
            updatePlan(id: $id, name: $name, description: $description, price: $price, traffic: $traffic, speed_limit: $speed_limit, device_limit: $device_limit, expiry_days: $expiry_days, group_id: $group_id, show: $show) {
              id
              name
            }
          }
        `,
        variables: {
          id: editingPlanId.value,
          name: newPlan.value.name,
          description: newPlan.value.description,
          price: newPlan.value.price,
          traffic: newPlan.value.traffic,
          speed_limit: newPlan.value.speed_limit,
          device_limit: newPlan.value.device_limit,
          expiry_days: newPlan.value.expiry_days,
          group_id: newPlan.value.group_id,
          show: newPlan.value.show
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'plan_updated',
      title: '套餐已更新',
      description: `成功更新套餐資訊：${newPlan.value.name}`,
      color: 'success'
    })

    isPlanModalOpen.value = false
    await fetchData()
  } catch (error) {
    toast.add({
      id: 'update_plan_failed',
      title: '更新失敗',
      description: error.message || '更新套餐失敗',
      color: 'error'
    })
  } finally {
    modalLoading.value = false
  }
}

const deletePlan = async (id) => {
  if (!confirm('確認刪除該套餐計劃？這將使得該套餐不再可被新用戶購買！')) return
  
  deleteLoadingId.value = id
  const token = useCookie('auth_token').value
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation DeletePlan($id: Int!) {
            deletePlan(id: $id)
          }
        `,
        variables: { id }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'plan_deleted',
      title: '套餐已刪除',
      description: '套餐計劃已成功從系統中移除。',
      color: 'success'
    })

    await fetchData()
  } catch (error) {
    toast.add({
      id: 'delete_plan_failed',
      title: '刪除失敗',
      description: error.message || '刪除套餐失敗',
      color: 'error'
    })
  } finally {
    deleteLoadingId.value = null
  }
}
</script>
