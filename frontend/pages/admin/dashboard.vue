<template>
  <div class="max-w-5xl mx-auto space-y-5">
      
      <!-- Admin header -->
      <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
        <div>
          <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">{{ t('system_management') }}</h1>
          <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('admin_panel') }}</p>
        </div>
      </header>

      <!-- Users Section -->
      <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-4 shadow-sm">
        <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
          <UIcon name="i-lucide-users" class="w-4 h-4 text-primary-500 dark:text-primary-400" />
          <span>{{ t('registered_users') }} ({{ users.length }})</span>
        </h2>

        <div v-if="loading" class="space-y-3">
          <USkeleton class="h-10 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 3" :key="i" />
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="border-b border-slate-200 dark:border-zinc-800 text-slate-500 dark:text-zinc-400 font-medium">
                <th class="py-2 px-3">ID</th>
                <th class="py-2 px-3">{{ t('email') }}</th>
                <th class="py-2 px-3">{{ t('traffic_usage') }}</th>
                <th class="py-2 px-3">{{ t('status') }}</th>
                <th class="py-2 px-3">{{ t('type') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="usr in users" :key="usr.id" class="border-b border-slate-200/50 dark:border-zinc-800/40 hover:bg-slate-100/50 dark:hover:bg-zinc-900/20 text-slate-700 dark:text-zinc-200">
                <td class="py-2.5 px-3 font-mono">{{ usr.id }}</td>
                <td class="py-2.5 px-3">{{ usr.email }}</td>
                <td class="py-2.5 px-3">
                  <div class="flex items-center space-x-1.5">
                    <span class="font-bold">{{ formatTraffic(usr.used_traffic) }}</span>
                    <span class="text-[10px] text-slate-400 dark:text-zinc-500">/ {{ formatTraffic(usr.total_traffic) }}</span>
                  </div>
                </td>
                <td class="py-2.5 px-3">
                  <UBadge :color="usr.status === 1 ? 'success' : 'error'" variant="subtle" size="xs">
                    {{ usr.status === 1 ? t('active') : t('disabled') }}
                  </UBadge>
                </td>
                <td class="py-2.5 px-3">
                  <UBadge :color="usr.is_admin ? 'primary' : 'neutral'" variant="soft" size="xs">
                    {{ usr.is_admin ? t('admin') : t('client') }}
                  </UBadge>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <!-- Nodes Section -->
      <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-4 shadow-sm">
        <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
          <UIcon name="i-lucide-server" class="w-4 h-4 text-primary-500 dark:text-primary-400" />
          <span>{{ t('active_nodes') }} ({{ nodes.length }})</span>
        </h2>

        <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <USkeleton class="h-24 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 3" :key="i" />
        </div>

        <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="node in nodes" :key="node.id" class="p-4 bg-white dark:bg-zinc-950/60 border border-slate-200 dark:border-zinc-800/60 rounded-lg space-y-3 shadow-sm">
            <div class="flex justify-between items-start">
              <div>
                <p class="font-bold text-slate-900 dark:text-white text-md">{{ node.name }}</p>
                <p class="text-[11px] text-slate-500 font-mono mt-0.5">{{ node.address }}:{{ node.port }}</p>
              </div>
              <UBadge color="primary" variant="subtle" size="xs">{{ node.type }}</UBadge>
            </div>
            
            <div class="text-[11px] text-slate-500 dark:text-zinc-400 flex justify-between">
              <span>{{ t('traffic_rate') }}: <strong class="text-slate-700 dark:text-white">{{ node.traffic_rate }}x</strong></span>
              <span>{{ t('status') }}: <strong class="text-emerald-500 dark:text-emerald-400">{{ t('online') }}</strong></span>
            </div>
          </div>
        </div>
      </section>

      <!-- Add Node Modal -->
      <UModal v-model:open="isNodeModalOpen">
        <div class="p-4 space-y-4 bg-white dark:bg-zinc-900 border border-slate-200 dark:border-zinc-800 rounded-lg">
          <div class="flex justify-between items-center">
            <h3 class="text-md font-bold text-slate-900 dark:text-white">{{ t('create_new_node') }}</h3>
            <UButton
              color="neutral"
              variant="ghost"
              icon="i-lucide-x"
              @click="isNodeModalOpen = false"
            />
          </div>

          <form @submit.prevent="createNode" class="space-y-3">
            <UFormField :label="t('node_name')" name="name" class="text-slate-700 dark:text-zinc-300">
              <UInput v-model="newNode.name" placeholder="SG-01 Vless Premium" color="primary" size="sm" class="w-full" required />
            </UFormField>

            <UFormField :label="t('protocol_type')" name="type" class="text-slate-700 dark:text-zinc-300">
              <USelect
                v-model="newNode.type"
                :items="['V2ray', 'Vless', 'Trojan']"
                color="primary"
                size="sm"
                class="w-full"
              />
            </UFormField>

            <div class="grid grid-cols-3 gap-3">
              <UFormField :label="t('address')" name="address" class="text-slate-700 dark:text-zinc-300 col-span-2">
                <UInput v-model="newNode.address" placeholder="sg1.example.com" color="primary" size="sm" class="w-full" required />
              </UFormField>
              <UFormField :label="t('port')" name="port" class="text-slate-700 dark:text-zinc-300">
                <UInput v-model.number="newNode.port" type="number" placeholder="443" color="primary" size="sm" class="w-full" required />
              </UFormField>
            </div>

            <UFormField :label="t('traffic_rate')" name="traffic_rate" class="text-slate-700 dark:text-zinc-300">
              <UInput v-model.number="newNode.traffic_rate" type="number" step="0.1" placeholder="1.0" color="primary" size="sm" class="w-full" required />
            </UFormField>

            <UFormField :label="t('transport_settings')" name="settings" class="text-slate-700 dark:text-zinc-300">
              <UTextarea v-model="newNode.settings" placeholder='{"transport": "ws", "path": "/v2ray"}' color="primary" size="sm" rows="3" class="w-full" />
            </UFormField>

            <div class="pt-2 flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" size="sm" @click="isNodeModalOpen = false">{{ t('cancel') }}</UButton>
              <UButton type="submit" color="primary" size="sm" :loading="modalLoading">{{ t('create_node') }}</UButton>
            </div>
          </form>
        </div>
      </UModal>

      <!-- Add Announcement Modal -->
      <UModal v-model:open="isAnnouncementModalOpen">
        <div class="p-4 space-y-4 bg-white dark:bg-zinc-900 border border-slate-200 dark:border-zinc-800 rounded-lg">
          <div class="flex justify-between items-center">
            <h3 class="text-md font-bold text-slate-900 dark:text-white">{{ t('announcement_new') }}</h3>
            <UButton
              color="neutral"
              variant="ghost"
              icon="i-lucide-x"
              @click="isAnnouncementModalOpen = false"
            />
          </div>

          <form @submit.prevent="createAnnouncement" class="space-y-3">
            <UFormField :label="t('title')" name="title" class="text-slate-700 dark:text-zinc-300">
              <UInput v-model="newAnnouncement.title" placeholder="Maintenance Notification" color="primary" size="sm" class="w-full" required />
            </UFormField>

            <UFormField :label="t('content')" name="content" class="text-slate-700 dark:text-zinc-300">
              <UTextarea v-model="newAnnouncement.content" placeholder="We will be performing node maintenance..." color="primary" size="sm" rows="4" class="w-full" required />
            </UFormField>

            <div class="pt-2 flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" size="sm" @click="isAnnouncementModalOpen = false">{{ t('cancel') }}</UButton>
              <UButton type="submit" color="primary" size="sm" :loading="announcementLoading">{{ t('create_announcement') }}</UButton>
            </div>
          </form>
        </div>
      </UModal>

  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})

const { t } = useI18n()
const localePath = useLocalePath()

const users = ref([])
const nodes = ref([])
const loading = ref(true)
const modalLoading = ref(false)
const isNodeModalOpen = ref(false)

const isAnnouncementModalOpen = ref(false)
const announcementLoading = ref(false)
const newAnnouncement = ref({
  title: '',
  content: ''
})

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

useSeoMeta({
  title: () => `${t('admin_panel')} - HY-Board`
})

const newNode = ref({
  name: '',
  type: 'Vless',
  address: '',
  port: 443,
  traffic_rate: 1.0,
  settings: '{"transport": "ws", "path": "/vless"}'
})

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
          query GetAdminData {
            adminUsers {
              id
              email
              total_traffic
              used_traffic
              status
              is_admin
            }
            nodes {
              id
              name
              type
              address
              port
              traffic_rate
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    users.value = response.data.adminUsers
    nodes.value = response.data.nodes
  } catch (error) {
    toast.add({
      id: 'session_expired',
      title: t('session_expired'),
      description: error.message || t('login_again_admin'),
      color: 'error'
    })
    logout()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})

const createNode = async () => {
  modalLoading.value = true
  const token = useCookie('auth_token').value
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation CreateNode($name: String!, $type: String!, $address: String!, $port: Int!, $traffic_rate: Float!, $settings: String) {
            createNode(name: $name, type: $type, address: $address, port: $port, traffic_rate: $traffic_rate, settings: $settings) {
              id
              name
            }
          }
        `,
        variables: {
          name: newNode.value.name,
          type: newNode.value.type,
          address: newNode.value.address,
          port: newNode.value.port,
          traffic_rate: newNode.value.traffic_rate,
          settings: newNode.value.settings
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'node_created',
      title: t('node_created'),
      description: t('node_created_desc').replace('{name}', newNode.value.name),
      color: 'success'
    })

    isNodeModalOpen.value = false
    // Reset form
    newNode.value = {
      name: '',
      type: 'Vless',
      address: '',
      port: 443,
      traffic_rate: 1.0,
      settings: '{"transport": "ws", "path": "/vless"}'
    }

    // Refresh nodes
    await fetchData()
  } catch (error) {
    toast.add({
      id: 'node_creation_failed',
      title: t('node_failed'),
      description: error.message || 'Node creation failed',
      color: 'error'
    })
  } finally {
    modalLoading.value = false
  }
}

const createAnnouncement = async () => {
  announcementLoading.value = true
  const token = useCookie('auth_token').value
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation CreateAnnouncement($title: String!, $content: String!) {
            createAnnouncement(title: $title, content: $content) {
              id
              title
            }
          }
        `,
        variables: {
          title: newAnnouncement.value.title,
          content: newAnnouncement.value.content
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'announcement_created',
      title: t('announcement_created'),
      description: t('announcement_created_desc').replace('{title}', newAnnouncement.value.title),
      color: 'success'
    })

    isAnnouncementModalOpen.value = false
    // Reset form
    newAnnouncement.value = {
      title: '',
      content: ''
    }
  } catch (error) {
    toast.add({
      id: 'announcement_creation_failed',
      title: t('announcement_failed'),
      description: error.message || 'Announcement creation failed',
      color: 'error'
    })
  } finally {
    announcementLoading.value = false
  }
}

const formatTraffic = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const logout = () => {
  const tokenCookie = useCookie('auth_token')
  tokenCookie.value = null
  router.push(localePath('/login'))
}
</script>
