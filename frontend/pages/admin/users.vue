<template>
  <div class="max-w-6xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">{{ t('user_management') }}</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('admin_panel') }}</p>
      </div>
    </header>

    <!-- Search & Stats -->
    <div class="flex flex-col sm:flex-row justify-between gap-3 items-center">
      <div class="w-full sm:w-72">
        <UInput
          v-model="searchQuery"
          icon="i-lucide-search"
          :placeholder="t('search_placeholder') || 'Search users...'"
          size="sm"
          class="w-full"
        />
      </div>
      <div class="text-[11px] text-slate-550 dark:text-zinc-400">
        Total Users: <span class="font-bold text-primary-500">{{ filteredUsers.length }}</span>
      </div>
    </div>

    <!-- Users Table Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-4 shadow-sm">
      <div v-if="loading" class="space-y-3">
        <USkeleton class="h-10 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 5" :key="i" />
      </div>

      <div v-else-if="filteredUsers.length === 0" class="text-center py-8 text-slate-500 text-xs">
        No users found.
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="border-b border-slate-200 dark:border-zinc-800 text-slate-500 dark:text-zinc-400 font-medium">
              <th class="py-2 px-3 w-12">ID</th>
              <th class="py-2 px-3">Email</th>
              <th class="py-2 px-3">{{ t('balance') }}</th>
              <th class="py-2 px-3">{{ t('traffic_usage') }}</th>
              <th class="py-2 px-3">{{ t('limitations') }}</th>
              <th class="py-2 px-3">{{ t('plan_expiration') }}</th>
              <th class="py-2 px-3 w-20 text-center">{{ t('status') }}</th>
              <th class="py-2 px-3 w-28 text-right"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="usr in filteredUsers" :key="usr.id" class="border-b border-slate-200/50 dark:border-zinc-800/40 hover:bg-slate-100/50 dark:hover:bg-zinc-900/20 text-slate-700 dark:text-zinc-200">
              <td class="py-3 px-3 font-mono">{{ usr.id }}</td>
              <td class="py-3 px-3">
                <div class="font-semibold flex items-center space-x-1.5">
                  <span>{{ usr.email }}</span>
                  <UBadge v-if="usr.is_admin" color="primary" variant="subtle" size="xs" class="ml-1">Admin</UBadge>
                </div>
              </td>
              <td class="py-3 px-3 font-mono font-bold text-emerald-600 dark:text-emerald-450">${{ usr.balance.toFixed(2) }}</td>
              <td class="py-3 px-3">
                <div class="flex flex-col space-y-1">
                  <div class="flex items-center space-x-1.5">
                    <span class="font-bold">{{ formatTraffic(usr.used_traffic) }}</span>
                    <span class="text-[10px] text-slate-400 dark:text-zinc-500">/ {{ usr.total_traffic ? formatTraffic(usr.total_traffic) : t('unlimited') }}</span>
                  </div>
                  <!-- Progress Bar -->
                  <div v-if="usr.total_traffic" class="w-24 bg-slate-200 dark:bg-zinc-800 h-1 rounded overflow-hidden">
                    <div class="bg-primary-500 h-full" :style="{ width: Math.min(100, (usr.used_traffic / usr.total_traffic) * 100) + '%' }"></div>
                  </div>
                </div>
              </td>
              <td class="py-3 px-3">
                <div class="flex flex-col space-y-0.5 text-[10px] text-slate-500 dark:text-zinc-400">
                  <div>Speed: <span class="font-mono font-semibold text-slate-700 dark:text-zinc-300">{{ usr.speed_limit ? `${usr.speed_limit} Mbps` : t('unlimited') }}</span></div>
                  <div>Devices: <span class="font-mono font-semibold text-slate-700 dark:text-zinc-300">{{ usr.device_limit ? `${usr.device_limit}` : t('unlimited') }}</span></div>
                </div>
              </td>
              <td class="py-3 px-3 font-mono text-slate-500">{{ usr.expired_at ? formatDate(usr.expired_at) : t('lifetime') }}</td>
              <td class="py-3 px-3 text-center">
                <USwitch
                  :model-value="usr.status === 1"
                  color="primary"
                  size="sm"
                  @update:model-value="(val) => handleToggleStatus(usr, val)"
                />
              </td>
              <td class="py-3 px-3 text-right">
                <div class="flex justify-end space-x-1.5">
                  <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-lucide-edit"
                    size="xs"
                    @click="openEditModal(usr)"
                  />
                  <UButton
                    color="error"
                    variant="ghost"
                    icon="i-lucide-trash-2"
                    size="xs"
                    :disabled="usr.is_admin"
                    @click="handleDelete(usr.id)"
                  />
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <UModal v-model:open="isEditModalOpen">
      <template #content>
        <div class="p-5 space-y-4 bg-white dark:bg-zinc-900 border border-slate-200 dark:border-zinc-800 rounded-lg max-h-[90vh] overflow-y-auto custom-scrollbar">
          <div class="flex justify-between items-center pb-2 border-b border-slate-100 dark:border-zinc-800">
            <h3 class="text-md font-bold text-slate-900 dark:text-white">{{ t('edit_user') }}</h3>
            <UButton
              color="neutral"
              variant="ghost"
              icon="i-lucide-x"
              @click="isEditModalOpen = false"
            />
          </div>

          <form @submit.prevent="updateUser" class="space-y-4">
            <UFormField :label="t('email')" name="email">
              <UInput v-model="editingUser.email" type="email" size="sm" class="w-full" required />
            </UFormField>

            <UFormField label="New Password (leave blank if unchanged)" name="password">
              <UInput v-model="editingUser.password" type="password" size="sm" class="w-full" placeholder="••••••••" />
            </UFormField>

            <div class="grid grid-cols-2 gap-4">
              <UFormField :label="t('balance') + ' ($)'" name="balance">
                <UInput v-model.number="editingUser.balance" type="number" step="0.01" size="sm" class="w-full" required />
              </UFormField>

              <UFormField :label="t('speed_limit') + ' (Mbps)'" name="speed_limit">
                <UInput v-model.number="editingUser.speed_limit" type="number" size="sm" class="w-full" required />
              </UFormField>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <UFormField :label="t('max_devices')" name="device_limit">
                <UInput v-model.number="editingUser.device_limit" type="number" size="sm" class="w-full" required />
              </UFormField>

              <UFormField label="Total Traffic (GB)" name="total_traffic">
                <UInput v-model.number="editingUser.total_traffic_gb" type="number" step="0.1" size="sm" class="w-full" required />
              </UFormField>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <UFormField label="Used Traffic (GB)" name="used_traffic">
                <UInput v-model.number="editingUser.used_traffic_gb" type="number" step="0.1" size="sm" class="w-full" required />
              </UFormField>

              <UFormField :label="t('plan_expiration')" name="expired_at">
                <UInput v-model="editingUser.expired_at_str" type="date" size="sm" class="w-full" />
              </UFormField>
            </div>

            <div class="flex space-x-6 items-center pt-2">
              <div class="flex items-center space-x-2">
                <span class="text-xs text-slate-500 dark:text-zinc-400">Is Admin</span>
                <USwitch v-model="editingUser.is_admin" color="primary" size="sm" />
              </div>
              <div class="flex items-center space-x-2">
                <span class="text-xs text-slate-500 dark:text-zinc-400">Status Active</span>
                <USwitch v-model="editingUser.status_active" color="primary" size="sm" />
              </div>
            </div>

            <div class="pt-3 flex justify-end gap-2 border-t border-slate-100 dark:border-zinc-800">
              <UButton color="neutral" variant="ghost" size="sm" @click="isEditModalOpen = false">{{ t('cancel') }}</UButton>
              <UButton type="submit" color="primary" size="sm" :loading="updateLoading">Save Changes</UButton>
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

const { t } = useI18n()
const localePath = useLocalePath()
const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

const users = ref([])
const loading = ref(true)
const searchQuery = ref('')
const isEditModalOpen = ref(false)
const updateLoading = ref(false)

const editingUser = ref({
  id: 0,
  email: '',
  password: '',
  balance: 0.0,
  speed_limit: 0,
  device_limit: 0,
  total_traffic_gb: 0,
  used_traffic_gb: 0,
  expired_at_str: '',
  is_admin: false,
  status_active: true
})

useSeoMeta({
  title: () => `${t('user_management')} - HY Admin`
})

const formatTraffic = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateStr) => {
  if (!dateStr || dateStr.startsWith('0001-01-01')) return ''
  const d = new Date(dateStr)
  return d.toISOString().split('T')[0]
}

const filteredUsers = computed(() => {
  if (!searchQuery.value.trim()) return users.value
  const q = searchQuery.value.toLowerCase()
  return users.value.filter(u => u.email.toLowerCase().includes(q))
})

const fetchUsers = async () => {
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
          query GetAdminUsers {
            adminUsers {
              id
              email
              speed_limit
              device_limit
              total_traffic
              used_traffic
              expired_at
              status
              is_admin
              balance
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    users.value = response.data?.adminUsers || []
  } catch (error) {
    toast.add({
      id: 'session_expired',
      title: t('session_expired'),
      description: error.message || t('login_again_admin'),
      color: 'error'
    })
    router.push(localePath('/login'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchUsers()
})

const openEditModal = (user) => {
  editingUser.value = {
    id: user.id,
    email: user.email,
    password: '',
    balance: user.balance,
    speed_limit: user.speed_limit,
    device_limit: user.device_limit,
    total_traffic_gb: parseFloat((user.total_traffic / (1024 * 1024 * 1024)).toFixed(2)),
    used_traffic_gb: parseFloat((user.used_traffic / (1024 * 1024 * 1024)).toFixed(2)),
    expired_at_str: formatDate(user.expired_at),
    is_admin: user.is_admin,
    status_active: user.status === 1
  }
  isEditModalOpen.value = true
}

const handleToggleStatus = async (user, newStatus) => {
  const token = useCookie('auth_token').value
  const statusInt = newStatus ? 1 : 0
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation ToggleUserStatus($id: Int!, $status: Int!) {
            updateUser(id: $id, status: $status) {
              id
              status
            }
          }
        `,
        variables: {
          id: user.id,
          status: statusInt
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    user.status = statusInt
    toast.add({
      title: 'Status Updated',
      description: `User ${user.email} status changed.`,
      color: 'success'
    })
  } catch (error) {
    toast.add({
      title: 'Failed to toggle status',
      description: error.message,
      color: 'error'
    })
  }
}

const updateUser = async () => {
  updateLoading.value = true
  const token = useCookie('auth_token').value
  try {
    const totalTrafficBytes = editingUser.value.total_traffic_gb * 1024 * 1024 * 1024
    const usedTrafficBytes = editingUser.value.used_traffic_gb * 1024 * 1024 * 1024
    
    let expiredAtIso = ''
    if (editingUser.value.expired_at_str) {
      expiredAtIso = new Date(editingUser.value.expired_at_str).toISOString()
    }

    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation UpdateUser(
            $id: Int!, 
            $email: String!, 
            $password: String, 
            $balance: Float!, 
            $speed_limit: Int!, 
            $device_limit: Int!, 
            $total_traffic: Float!, 
            $used_traffic: Float!, 
            $expired_at: String!, 
            $status: Int!, 
            $is_admin: Boolean!
          ) {
            updateUser(
              id: $id, 
              email: $email, 
              password: $password, 
              balance: $balance, 
              speed_limit: $speed_limit, 
              device_limit: $device_limit, 
              total_traffic: $total_traffic, 
              used_traffic: $used_traffic, 
              expired_at: $expired_at, 
              status: $status, 
              is_admin: $is_admin
            ) {
              id
            }
          }
        `,
        variables: {
          id: editingUser.value.id,
          email: editingUser.value.email,
          password: editingUser.value.password || null,
          balance: editingUser.value.balance,
          speed_limit: editingUser.value.speed_limit,
          device_limit: editingUser.value.device_limit,
          total_traffic: totalTrafficBytes,
          used_traffic: usedTrafficBytes,
          expired_at: expiredAtIso,
          status: editingUser.value.status_active ? 1 : 0,
          is_admin: editingUser.value.is_admin
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      title: t('user_updated') || 'User updated',
      description: 'Changes have been saved successfully.',
      color: 'success'
    })

    isEditModalOpen.value = false
    await fetchUsers()
  } catch (error) {
    toast.add({
      title: t('user_update_failed') || 'Failed to update user',
      description: error.message,
      color: 'error'
    })
  } finally {
    updateLoading.value = false
  }
}

const handleDelete = async (id) => {
  if (!confirm(t('confirm_delete_user') || 'Are you sure you want to delete this user?')) return

  const token = useCookie('auth_token').value
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation DeleteUser($id: Int!) {
            deleteUser(id: $id)
          }
        `,
        variables: { id }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      title: t('user_deleted') || 'User deleted',
      color: 'success'
    })

    await fetchUsers()
  } catch (error) {
    toast.add({
      title: t('user_delete_failed') || 'Failed to delete user',
      description: error.message,
      color: 'error'
    })
  }
}
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 2px;
}
.dark .custom-scrollbar::-webkit-scrollbar-thumb {
  background: #334155;
}
</style>
