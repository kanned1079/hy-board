<template>
  <div class="max-w-6xl mx-auto space-y-6">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white flex items-center space-x-2">
          <UIcon name="i-lucide-shield" class="w-6 h-6 text-primary-500 shrink-0" />
          <span>權限組管理 (Subscription Groups)</span>
        </h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">自訂多個權限組作為訂閱套餐等級，設置節點訪問範圍與客戶端限速</p>
      </div>
      <UButton
        color="primary"
        icon="i-lucide-plus"
        size="sm"
        @click="openAddModal"
      >
        新增權限組
      </UButton>
    </header>

    <!-- Stats Overview -->
    <div v-if="!loading" class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg flex items-center space-x-3 shadow-sm">
        <div class="p-2 bg-primary-500/10 rounded-md text-primary-500">
          <UIcon name="i-lucide-folder-key" class="w-5 h-5" />
        </div>
        <div>
          <p class="text-[10px] text-slate-500 uppercase font-semibold">權限組總數</p>
          <p class="text-lg font-black text-slate-900 dark:text-white">{{ groups.length }} <span class="text-xs font-normal text-slate-400">個</span></p>
        </div>
      </div>
      <div class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg flex items-center space-x-3 shadow-sm">
        <div class="p-2 bg-sky-500/10 rounded-md text-sky-500">
          <UIcon name="i-lucide-users" class="w-5 h-5" />
        </div>
        <div>
          <p class="text-[10px] text-slate-500 uppercase font-semibold">總訂閱人數</p>
          <p class="text-lg font-black text-slate-900 dark:text-white">{{ users.length }} <span class="text-xs font-normal text-slate-400">人</span></p>
        </div>
      </div>
      <div class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg flex items-center space-x-3 shadow-sm">
        <div class="p-2 bg-emerald-500/10 rounded-md text-emerald-500">
          <UIcon name="i-lucide-server" class="w-5 h-5" />
        </div>
        <div>
          <p class="text-[10px] text-slate-500 uppercase font-semibold">佈署節點數</p>
          <p class="text-lg font-black text-slate-900 dark:text-white">{{ nodes.length }} <span class="text-xs font-normal text-slate-400">個</span></p>
        </div>
      </div>
    </div>

    <!-- Skeleton Loading -->
    <div v-if="loading" class="space-y-3">
      <USkeleton class="h-12 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 4" :key="i" />
    </div>

    <!-- Content Table -->
    <section v-else class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg shadow-sm">
      <div v-if="groups.length === 0" class="text-center py-8 text-slate-500 text-xs">
        尚未建立任何權限組，請點擊上方按鈕新增。
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="border-b border-slate-200 dark:border-zinc-800 text-slate-500 dark:text-zinc-400 font-medium">
              <th class="py-2.5 px-3 w-16">ID / Level</th>
              <th class="py-2.5 px-3 w-48">權限組名稱</th>
              <th class="py-2.5 px-3">描述與訪問說明</th>
              <th class="py-2.5 px-3 w-28 text-center">關聯節點數</th>
              <th class="py-2.5 px-3 w-28 text-center">訂閱用戶數</th>
              <th class="py-2.5 px-3 w-24 text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="group in groups" :key="group.id" class="border-b border-slate-200/50 dark:border-zinc-800/40 hover:bg-slate-100/50 dark:hover:bg-zinc-900/20 text-slate-700 dark:text-zinc-200">
              <td class="py-3.5 px-3 font-mono font-bold">{{ group.id }}</td>
              <td class="py-3.5 px-3">
                <span class="font-bold text-slate-900 dark:text-white text-sm">{{ group.name }}</span>
                <UBadge v-if="group.id === 99" color="success" variant="subtle" size="xs" class="ml-1.5">系統預設</UBadge>
              </td>
              <td class="py-3.5 px-3 text-slate-550 dark:text-zinc-400 leading-relaxed">{{ group.description || '無描述資訊。' }}</td>
              <td class="py-3.5 px-3 text-center font-mono">
                <UBadge color="neutral" variant="soft" size="xs">
                  {{ getNodeCount(group.id) }} 個節點
                </UBadge>
              </td>
              <td class="py-3.5 px-3 text-center font-mono font-semibold">
                <span class="text-primary-600 dark:text-primary-400">{{ getUserCount(group.id) }}</span> 人
              </td>
              <td class="py-3.5 px-3 text-right">
                <div class="flex justify-end space-x-1.5">
                  <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-lucide-edit"
                    size="xs"
                    @click="openEditModal(group)"
                  />
                  <UButton
                    color="error"
                    variant="ghost"
                    icon="i-lucide-trash-2"
                    size="xs"
                    :disabled="group.id === 99"
                    @click="handleDelete(group.id)"
                  />
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- Create/Edit Modal -->
    <UModal v-model:open="isModalOpen">
      <template #content>
        <div class="p-5 space-y-4 bg-white dark:bg-zinc-900 border border-slate-200 dark:border-zinc-800 rounded-lg">
          <div class="flex justify-between items-center pb-2 border-b border-slate-100 dark:border-zinc-800">
            <h3 class="text-md font-bold text-slate-900 dark:text-white">
              {{ isEditMode ? '修改權限組' : '新增權限組' }}
            </h3>
            <UButton
              color="neutral"
              variant="ghost"
              icon="i-lucide-x"
              @click="isModalOpen = false"
            />
          </div>

          <form @submit.prevent="handleSubmit" class="space-y-4">
            <UFormField label="權限組名稱" name="name" class="text-slate-700 dark:text-zinc-350">
              <UInput v-model="formModel.name" placeholder="例如: S3 歐洲專線版" color="primary" size="sm" class="w-full" required />
            </UFormField>

            <UFormField label="描述與備忘" name="description" class="text-slate-700 dark:text-zinc-350">
              <UTextarea v-model="formModel.description" placeholder="例如: 提供歐洲跨國精品線路，速率上限 500Mbps" color="primary" size="sm" rows="3" class="w-full" />
            </UFormField>

            <div class="pt-3 flex justify-end gap-2 border-t border-slate-100 dark:border-zinc-800">
              <UButton color="neutral" variant="ghost" size="sm" @click="isModalOpen = false">取消</UButton>
              <UButton type="submit" color="primary" size="sm" :loading="submitLoading">儲存</UButton>
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

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()
const localePath = useLocalePath()

const groups = ref([])
const users = ref([])
const nodes = ref([])
const loading = ref(true)

const isModalOpen = ref(false)
const isEditMode = ref(false)
const submitLoading = ref(false)

const formModel = ref({
  id: null,
  name: '',
  description: ''
})

useSeoMeta({
  title: '權限組與訂閱等級管理 - HY Admin'
})

// Statistics helpers
const getUserCount = (groupId) => users.value.filter(u => u.group_id === groupId).length
const getNodeCount = (groupId) => nodes.value.filter(n => {
  if (!n.group_ids) return n.group_id === groupId
  return n.group_ids.split(',').map(s => s.trim()).filter(Boolean).map(Number).includes(groupId)
}).length

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
          query GetGroupsManagementData {
            groups {
              id
              name
              description
              created_at
            }
            adminUsers {
              id
              group_id
            }
            nodes {
              id
              group_id
              group_ids
              show
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    groups.value = response.data?.groups || []
    users.value = response.data?.adminUsers || []
    nodes.value = response.data?.nodes || []
  } catch (error) {
    toast.add({
      id: 'groups_fetch_error',
      title: '讀取權限組資料失敗',
      description: error.message || '請檢查網路連線或重新登入',
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})

const openAddModal = () => {
  isEditMode.value = false
  formModel.value = {
    id: null,
    name: '',
    description: ''
  }
  isModalOpen.value = true
}

const openEditModal = (group) => {
  isEditMode.value = true
  formModel.value = {
    id: group.id,
    name: group.name,
    description: group.description
  }
  isModalOpen.value = true
}

const handleSubmit = async () => {
  submitLoading.value = true
  const token = useCookie('auth_token').value
  
  const query = isEditMode.value ? `
    mutation UpdateGroup($id: Int!, $name: String!, $description: String) {
      updateGroup(id: $id, name: $name, description: $description) {
        id
        name
      }
    }
  ` : `
    mutation CreateGroup($name: String!, $description: String) {
      createGroup(name: $name, description: $description) {
        id
        name
      }
    }
  `

  const variables = isEditMode.value ? {
    id: formModel.value.id,
    name: formModel.value.name,
    description: formModel.value.description
  } : {
    name: formModel.value.name,
    description: formModel.value.description
  }

  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: { query, variables }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      title: isEditMode.value ? '權限組修改成功' : '權限組新增成功',
      description: `群組「${formModel.value.name}」已儲存。`,
      color: 'success'
    })

    isModalOpen.value = false
    await fetchData()
  } catch (error) {
    toast.add({
      title: '儲存權限組失敗',
      description: error.message,
      color: 'error'
    })
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = async (id) => {
  if (!confirm('您確定要刪除此權限組嗎？刪除後，綁定此群組的用戶與節點可能會失去訪問連線。')) return

  const token = useCookie('auth_token').value
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation DeleteGroup($id: Int!) {
            deleteGroup(id: $id)
          }
        `,
        variables: { id }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      title: '權限組刪除成功',
      color: 'success'
    })

    await fetchData()
  } catch (error) {
    toast.add({
      title: '刪除權限組失敗',
      description: error.message,
      color: 'error'
    })
  }
}
</script>
