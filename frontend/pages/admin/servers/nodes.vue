<template>
  <div class="max-w-5xl mx-auto space-y-5">
    
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">節點管理</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">新增、修改、刪除並監控節點狀態與協議設定</p>
      </div>
      <UButton
        color="primary"
        icon="i-lucide-plus"
        size="sm"
        @click="openAddModal"
      >
        添加節點
      </UButton>
    </header>

    <!-- Nodes List Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg space-y-4 shadow-sm">
      <div v-if="loading" class="space-y-3">
        <USkeleton class="h-10 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 3" :key="i" />
      </div>

      <div v-else-if="nodes.length === 0" class="text-center py-10 text-slate-400 text-xs">
        暫無已配置的節點，請點擊右上角「添加節點」創建。
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="border-b border-slate-200 dark:border-zinc-800 text-slate-500 dark:text-zinc-400 font-medium">
              <th class="py-2.5 px-3">狀態</th>
              <th class="py-2.5 px-3">ID</th>
              <th class="py-2.5 px-3">節點名稱</th>
              <th class="py-2.5 px-3">通訊協議</th>
              <th class="py-2.5 px-3">伺服器地址</th>
              <th class="py-2.5 px-3">端口</th>
              <th class="py-2.5 px-3">流量倍率</th>
              <th class="py-2.5 px-3">要求訂閱等級</th>
              <th class="py-2.5 px-3 text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="node in nodes" :key="node.id" class="border-b border-slate-200/50 dark:border-zinc-800/40 hover:bg-slate-100/50 dark:hover:bg-zinc-900/20 text-slate-700 dark:text-zinc-200">
              <td class="py-3 px-3">
                <div class="flex items-center justify-start pl-1">
                  <span
                    class="w-2.5 h-2.5 rounded-full inline-block"
                    :class="{
                      'bg-rose-500': node.status === 'offline',
                      'bg-amber-500': node.status === 'idle',
                      'bg-emerald-500 animate-pulse': node.status === 'active'
                    }"
                    :title="node.status === 'active' ? '活躍 (有人連結)' : (node.status === 'idle' ? '閒置 (正常連線)' : '離線 (無法連線)')"
                  ></span>
                </div>
              </td>
              <td class="py-3 px-3 font-mono font-bold">{{ node.id }}</td>
              <td class="py-3 px-3 font-medium">{{ node.name }}</td>
              <td class="py-3 px-3">
                <UBadge color="primary" variant="subtle" size="xs">{{ node.type }}</UBadge>
              </td>
              <td class="py-3 px-3 font-mono text-slate-500 dark:text-zinc-400">{{ node.address }}</td>
              <td class="py-3 px-3 font-mono">{{ node.port }}</td>
              <td class="py-3 px-3 font-bold text-slate-600 dark:text-zinc-300">{{ node.traffic_rate }}x</td>
              <td class="py-3 px-3">
                <div class="flex flex-wrap gap-1">
                  <span v-if="!node.group_ids || node.group_ids.trim() === ''" class="text-slate-400">-</span>
                  <UBadge
                    v-else
                    v-for="gId in (node.group_ids || '').split(',').map(s => s.trim()).filter(Boolean)"
                    :key="gId"
                    color="sky"
                    variant="soft"
                    size="xs"
                  >
                    {{ getGroupName(parseInt(gId)) }}
                  </UBadge>
                </div>
              </td>
              <td class="py-3 px-3 text-right flex justify-end gap-1.5">
                <UButton
                  color="neutral"
                  variant="ghost"
                  icon="i-lucide-edit"
                  size="xs"
                  @click="openEditModal(node)"
                />
                <UButton
                  color="red"
                  variant="ghost"
                  icon="i-lucide-trash-2"
                  size="xs"
                  :loading="deleteLoadingId === node.id"
                  @click="deleteNode(node.id)"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- Add/Edit Node Modal -->
    <UModal v-model:open="isNodeModalOpen">
      <template #content>
        <div class="p-4 space-y-4 bg-white dark:bg-zinc-900 border border-slate-200 dark:border-zinc-800 rounded-lg">
          <div class="flex justify-between items-center">
            <h3 class="text-md font-bold text-slate-900 dark:text-white">
              {{ isEditMode ? '編輯節點資訊' : '創建新節點' }}
            </h3>
            <UButton
              color="neutral"
              variant="ghost"
              icon="i-lucide-x"
              @click="isNodeModalOpen = false"
            />
          </div>

          <form @submit.prevent="handleSubmit" class="space-y-3">
            <UFormField label="節點名稱" name="name" class="text-slate-700 dark:text-zinc-300">
              <UInput v-model="newNode.name" placeholder="SG-01 Vless Premium" color="primary" size="sm" class="w-full" required />
            </UFormField>

            <UFormField label="協議類型" name="type" class="text-slate-700 dark:text-zinc-300">
              <USelect
                v-model="newNode.type"
                :items="['V2ray', 'Vless', 'Trojan', 'Shadowsocks']"
                color="primary"
                size="sm"
                class="w-full"
              />
            </UFormField>

            <div class="grid grid-cols-3 gap-3">
              <UFormField label="域名/IP" name="address" class="text-slate-700 dark:text-zinc-300 col-span-2">
                <UInput v-model="newNode.address" placeholder="sg1.example.com" color="primary" size="sm" class="w-full" required />
              </UFormField>
              <UFormField label="端口" name="port" class="text-slate-700 dark:text-zinc-300">
                <UInput v-model.number="newNode.port" type="number" placeholder="443" color="primary" size="sm" class="w-full" required />
              </UFormField>
            </div>

            <UFormField label="流量倍率" name="traffic_rate" class="text-slate-700 dark:text-zinc-300">
              <UInput v-model.number="newNode.traffic_rate" type="number" step="0.1" placeholder="1.0" color="primary" size="sm" class="w-full" required />
            </UFormField>

            <!-- Shadowsocks settings -->
            <div v-if="newNode.type === 'Shadowsocks'" class="space-y-3 p-3 bg-slate-50 dark:bg-zinc-950/40 rounded-lg border border-slate-200/50 dark:border-zinc-800/50">
              <p class="text-[11px] font-bold text-slate-500 dark:text-zinc-400">Shadowsocks 加密配置</p>
              <UFormField label="加密方法 (Cipher)" name="ssMethod" class="text-slate-700 dark:text-zinc-300">
                <USelect
                  v-model="newNode.ssMethod"
                  :items="['aes-256-gcm', 'aes-128-gcm', 'chacha20-ietf-poly1305']"
                  color="primary"
                  size="sm"
                  class="w-full"
                />
              </UFormField>
            </div>

            <!-- Vmess / Vless / Trojan settings -->
            <div v-else class="space-y-3 p-3 bg-slate-50 dark:bg-zinc-950/40 rounded-lg border border-slate-200/50 dark:border-zinc-800/50">
              <p class="text-[11px] font-bold text-slate-500 dark:text-zinc-400">傳輸與安全配置</p>
              
              <div class="grid grid-cols-2 gap-3">
                <UFormField label="傳輸協議 (Network)" name="transport" class="text-slate-700 dark:text-zinc-300">
                  <USelect
                    v-model="newNode.transport"
                    :items="['tcp', 'ws', 'grpc']"
                    color="primary"
                    size="sm"
                    class="w-full"
                  />
                </UFormField>
                
                <UFormField label="啟用 TLS 加密" name="enableTls" class="text-slate-700 dark:text-zinc-300 flex items-center justify-between pt-5">
                  <USwitch v-model="newNode.enableTls" color="primary" />
                </UFormField>
              </div>

              <!-- Path configuration for ws / grpc -->
              <UFormField
                v-if="newNode.transport === 'ws' || newNode.transport === 'grpc'"
                :label="newNode.transport === 'ws' ? 'WebSocket 路徑 (Path)' : 'gRPC 服務名稱 (ServiceName)'"
                name="path"
                class="text-slate-700 dark:text-zinc-300"
              >
                <UInput v-model="newNode.path" placeholder="/vless" color="primary" size="sm" class="w-full" />
              </UFormField>

              <!-- TLS SNI config -->
              <UFormField
                v-if="newNode.enableTls"
                label="TLS 伺服器名稱 (SNI)"
                name="host"
                class="text-slate-700 dark:text-zinc-300"
              >
                <UInput v-model="newNode.host" placeholder="example.com" color="primary" size="sm" class="w-full" />
              </UFormField>
            </div>

            <UFormField label="要求權限組 (可多選)" name="group_ids" class="text-slate-700 dark:text-zinc-300">
              <div class="grid grid-cols-2 gap-2 mt-1 p-2.5 border border-slate-200 dark:border-zinc-800 rounded-md bg-slate-50 dark:bg-zinc-950/40">
                <label v-for="g in groups" :key="g.id" class="flex items-center space-x-2 text-xs text-slate-700 dark:text-zinc-300 cursor-pointer select-none">
                  <input
                    type="checkbox"
                    :value="g.id"
                    v-model="newNode.group_ids"
                    class="rounded border-slate-300 dark:border-zinc-700 text-primary-600 focus:ring-primary-500"
                  />
                  <span>{{ g.name }} (Group {{ g.id }})</span>
                </label>
              </div>
            </UFormField>

            <div class="pt-2 flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" size="sm" @click="isNodeModalOpen = false">取消</UButton>
              <UButton type="submit" color="primary" size="sm" :loading="modalLoading">
                {{ isEditMode ? '保存修改' : '創建節點' }}
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

const { t } = useI18n()
const localePath = useLocalePath()

const nodes = ref([])
const groups = ref([])
const loading = ref(true)
const modalLoading = ref(false)
const deleteLoadingId = ref(null)
const isNodeModalOpen = ref(false)
const isEditMode = ref(false)
const editingNodeId = ref(null)

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

useSeoMeta({
  title: '節點管理 - HY-Board'
})

const getGroupName = (groupId) => {
  const g = groups.value.find(item => item.id === groupId)
  return g ? g.name : `Group ${groupId}`
}

const newNode = ref({
  name: '',
  type: 'Vless',
  address: '',
  port: 443,
  traffic_rate: 1.0,
  settings: '',
  group_id: 1,
  group_ids: [1],
  transport: 'tcp',
  path: '/vless',
  ssMethod: 'aes-256-gcm',
  enableTls: false,
  host: ''
})

const openAddModal = () => {
  isEditMode.value = false
  editingNodeId.value = null
  newNode.value = {
    name: '',
    type: 'Vless',
    address: '',
    port: 443,
    traffic_rate: 1.0,
    settings: '',
    group_id: 1,
    group_ids: [1],
    transport: 'tcp',
    path: '/vless',
    ssMethod: 'aes-256-gcm',
    enableTls: false,
    host: ''
  }
  isNodeModalOpen.value = true
}

const openEditModal = (node) => {
  isEditMode.value = true
  editingNodeId.value = node.id
  
  // Set default settings
  let transport = 'tcp'
  let path = '/vless'
  let ssMethod = 'aes-256-gcm'
  let enableTls = false
  let host = ''
  
  if (node.settings) {
    try {
      const s = JSON.parse(node.settings)
      if (node.type === 'Shadowsocks') {
        if (s.method) ssMethod = s.method
      } else {
        if (s.network) transport = s.network
        if (s.network === 'ws' && s.wsSettings) {
          path = s.wsSettings.path || ''
        } else if (s.network === 'grpc' && s.grpcSettings) {
          path = s.grpcSettings.serviceName || ''
        }
        if (s.security === 'tls') {
          enableTls = true
          if (s.tlsSettings) {
            host = s.tlsSettings.serverName || ''
          }
        }
      }
    } catch (e) {
      console.error('Failed to parse node settings for editing:', e)
    }
  }
  
  let groupIds = [1]
  if (node.group_ids) {
    groupIds = node.group_ids.split(',').map(s => s.trim()).filter(Boolean).map(Number)
  } else if (node.group_id) {
    groupIds = [node.group_id]
  }

  newNode.value = {
    name: node.name,
    type: node.type,
    address: node.address,
    port: node.port,
    traffic_rate: node.traffic_rate,
    settings: node.settings || '',
    group_id: node.group_id || 1,
    group_ids: groupIds,
    transport,
    path,
    ssMethod,
    enableTls,
    host
  }
  isNodeModalOpen.value = true
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
          query GetNodesData {
            nodes {
              id
              name
              type
              address
              port
              traffic_rate
              group_id
              group_ids
              settings
              online
              status
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

    nodes.value = response.data.nodes || []
    groups.value = response.data.groups || []
  } catch (error) {
    toast.add({
      id: 'fetch_failed',
      title: '加載失敗',
      description: error.message || '無法獲取節點數據，請重新登入',
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}

let refreshInterval = null

onMounted(() => {
  fetchData()
  refreshInterval = setInterval(() => {
    fetchData()
  }, 15000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})

const handleSubmit = async () => {
  if (isEditMode.value) {
    await saveNodeChanges()
  } else {
    await createNode()
  }
}

const createNode = async () => {
  modalLoading.value = true
  const token = useCookie('auth_token').value

  const settingsObj = {}
  if (newNode.value.type === 'Shadowsocks') {
    settingsObj.method = newNode.value.ssMethod
  } else {
    settingsObj.network = newNode.value.transport
    if (newNode.value.transport === 'ws') {
      settingsObj.wsSettings = {
        path: newNode.value.path || '/'
      }
    } else if (newNode.value.transport === 'grpc') {
      settingsObj.grpcSettings = {
        serviceName: newNode.value.path || 'GrpcService'
      }
    }

    if (newNode.value.enableTls) {
      settingsObj.security = 'tls'
      settingsObj.tlsSettings = {
        serverName: newNode.value.host || '',
        allowInsecure: false
      }
    } else {
      settingsObj.security = 'none'
    }
  }

  const settingsJson = JSON.stringify(settingsObj)

  if (!newNode.value.group_ids || newNode.value.group_ids.length === 0) {
    toast.add({
      id: 'validation_failed',
      title: '驗證失敗',
      description: '請至少選擇一個要求權限組',
      color: 'error'
    })
    modalLoading.value = false
    return
  }

  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation CreateNode($name: String!, $type: String!, $address: String!, $port: Int!, $traffic_rate: Float!, $settings: String, $group_id: Int, $group_ids: String) {
            createNode(name: $name, type: $type, address: $address, port: $port, traffic_rate: $traffic_rate, settings: $settings, group_id: $group_id, group_ids: $group_ids) {
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
          settings: settingsJson,
          group_id: newNode.value.group_ids[0] || 1,
          group_ids: newNode.value.group_ids.join(',')
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'node_created',
      title: '節點已創建',
      description: `成功添加節點：${newNode.value.name}`,
      color: 'success'
    })

    isNodeModalOpen.value = false
    await fetchData()
  } catch (error) {
    toast.add({
      id: 'create_failed',
      title: '創建失敗',
      description: error.message || '創建節點失敗',
      color: 'error'
    })
  } finally {
    modalLoading.value = false
  }
}

const saveNodeChanges = async () => {
  modalLoading.value = true
  const token = useCookie('auth_token').value

  const settingsObj = {}
  if (newNode.value.type === 'Shadowsocks') {
    settingsObj.method = newNode.value.ssMethod
  } else {
    settingsObj.network = newNode.value.transport
    if (newNode.value.transport === 'ws') {
      settingsObj.wsSettings = {
        path: newNode.value.path || '/'
      }
    } else if (newNode.value.transport === 'grpc') {
      settingsObj.grpcSettings = {
        serviceName: newNode.value.path || 'GrpcService'
      }
    }

    if (newNode.value.enableTls) {
      settingsObj.security = 'tls'
      settingsObj.tlsSettings = {
        serverName: newNode.value.host || '',
        allowInsecure: false
      }
    } else {
      settingsObj.security = 'none'
    }
  }

  const settingsJson = JSON.stringify(settingsObj)

  if (!newNode.value.group_ids || newNode.value.group_ids.length === 0) {
    toast.add({
      id: 'validation_failed',
      title: '驗證失敗',
      description: '請至少選擇一個要求權限組',
      color: 'error'
    })
    modalLoading.value = false
    return
  }

  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation UpdateNode($id: Int!, $name: String!, $type: String!, $address: String!, $port: Int!, $traffic_rate: Float!, $settings: String, $group_id: Int, $group_ids: String) {
            updateNode(id: $id, name: $name, type: $type, address: $address, port: $port, traffic_rate: $traffic_rate, settings: $settings, group_id: $group_id, group_ids: $group_ids) {
              id
              name
            }
          }
        `,
        variables: {
          id: editingNodeId.value,
          name: newNode.value.name,
          type: newNode.value.type,
          address: newNode.value.address,
          port: newNode.value.port,
          traffic_rate: newNode.value.traffic_rate,
          settings: settingsJson,
          group_id: newNode.value.group_ids[0] || 1,
          group_ids: newNode.value.group_ids.join(',')
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'node_updated',
      title: '節點已更新',
      description: `成功修改節點資訊：${newNode.value.name}`,
      color: 'success'
    })

    isNodeModalOpen.value = false
    await fetchData()
  } catch (error) {
    toast.add({
      id: 'update_failed',
      title: '更新失敗',
      description: error.message || '更新節點失敗',
      color: 'error'
    })
  } finally {
    modalLoading.value = false
  }
}

const deleteNode = async (id) => {
  if (!confirm('確認刪除該節點？這將導致使用此節點的客戶端斷開連線！')) return
  
  deleteLoadingId.value = id
  const token = useCookie('auth_token').value
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation DeleteNode($id: Int!) {
            deleteNode(id: $id)
          }
        `,
        variables: { id }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'node_deleted',
      title: '節點已刪除',
      description: '節點已成功從面板移除。',
      color: 'success'
    })

    await fetchData()
  } catch (error) {
    toast.add({
      id: 'delete_failed',
      title: '刪除失敗',
      description: error.message || '刪除節點失敗',
      color: 'error'
    })
  } finally {
    deleteLoadingId.value = null
  }
}
</script>
