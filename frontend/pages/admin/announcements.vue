<template>
  <div class="max-w-5xl mx-auto space-y-5">
      <!-- Header -->
      <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
        <div>
          <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">{{ t('announcements_manage') }}</h1>
          <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('admin_panel') }}</p>
        </div>
        <UButton
          color="primary"
          icon="i-lucide-megaphone"
          size="sm"
          @click="openCreateModal"
        >
          {{ t('announcement_new') }}
        </UButton>
      </header>

      <!-- Announcement List Card -->
      <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-4 shadow-sm">
        <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
          <UIcon name="i-lucide-megaphone" class="w-4 h-4 text-primary-500 dark:text-primary-400" />
          <span>{{ t('announcement_list') }}</span>
        </h2>

        <div v-if="loading" class="space-y-3">
          <USkeleton class="h-10 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 3" :key="i" />
        </div>

        <div v-else-if="announcements.length === 0" class="text-center py-6 text-slate-500 text-xs">
          {{ t('no_announcements') }}
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="border-b border-slate-200 dark:border-zinc-800 text-slate-500 dark:text-zinc-400 font-medium">
                <th class="py-2 px-3 w-16">ID</th>
                <th class="py-2 px-3">Title</th>
                <th class="py-2 px-3 w-32">Date</th>
                <th class="py-2 px-3 w-28 text-center">Visibility</th>
                <th class="py-2 px-3 w-32 text-right"></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in announcements" :key="item.id" class="border-b border-slate-200/50 dark:border-zinc-800/40 hover:bg-slate-100/50 dark:hover:bg-zinc-900/20 text-slate-700 dark:text-zinc-200">
                <td class="py-2.5 px-3 font-mono">{{ item.id }}</td>
                <td class="py-2.5 px-3 font-semibold">{{ item.title }}</td>
                <td class="py-2.5 px-3 text-slate-500 font-mono">{{ formatDate(item.created_at) }}</td>
                <td class="py-2.5 px-3 text-center">
                  <USwitch
                    :model-value="item.show"
                    color="primary"
                    size="sm"
                    @update:model-value="(val) => handleToggleShow(item.id, val)"
                  />
                </td>
                <td class="py-2.5 px-3 text-right space-x-1">
                  <UButton
                    color="primary"
                    variant="ghost"
                    icon="i-lucide-edit"
                    size="xs"
                    @click="openEditModal(item)"
                  />
                  <UButton
                    color="error"
                    variant="ghost"
                    icon="i-lucide-trash-2"
                    size="xs"
                    @click="handleDelete(item.id)"
                  />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
 
      <!-- Publish/Edit Announcement Modal -->
      <UModal v-model:open="isAnnouncementModalOpen">
        <template #content>
          <div class="p-4 space-y-4 bg-white dark:bg-zinc-900 border border-slate-200 dark:border-zinc-800 rounded-lg">
            <div class="flex justify-between items-center">
              <h3 class="text-md font-bold text-slate-900 dark:text-white">
                {{ isEditing ? t('edit_announcement') : t('announcement_new') }}
              </h3>
              <UButton
                color="neutral"
                variant="ghost"
                icon="i-lucide-x"
                @click="isAnnouncementModalOpen = false"
              />
            </div>
   
            <form @submit.prevent="saveAnnouncement" class="space-y-3">
              <UFormField :label="t('title')" name="title" class="text-slate-700 dark:text-zinc-300">
                <UInput v-model="form.title" placeholder="Maintenance Notification" color="primary" size="sm" class="w-full" required />
              </UFormField>
   
              <!-- Editor / Preview Tabs -->
              <div class="space-y-2">
                <div class="flex border-b border-slate-200 dark:border-zinc-800 mb-2">
                  <button 
                    type="button" 
                    class="px-3 py-1.5 text-xs font-semibold transition-colors"
                    :class="activeTab === 'write' ? 'border-b-2 border-primary-500 text-primary-500 dark:text-primary-400' : 'text-slate-500 hover:text-slate-700'"
                    @click="activeTab = 'write'"
                  >
                    {{ t('write') }} (Markdown)
                  </button>
                  <button 
                    type="button" 
                    class="px-3 py-1.5 text-xs font-semibold transition-colors"
                    :class="activeTab === 'preview' ? 'border-b-2 border-primary-500 text-primary-500 dark:text-primary-400' : 'text-slate-500 hover:text-slate-700'"
                    @click="activeTab = 'preview'"
                  >
                    {{ t('preview') }}
                  </button>
                </div>
   
                <!-- Write Tab Content -->
                <div v-show="activeTab === 'write'">
                  <UFormField :label="t('content')" name="content" class="text-slate-700 dark:text-zinc-300">
                    <UTextarea v-model="form.content" placeholder="Write markdown content here..." color="primary" size="sm" rows="6" class="w-full" required />
                  </UFormField>
                </div>
   
                <!-- Preview Tab Content -->
                <div v-show="activeTab === 'preview'" class="p-3 bg-slate-50 dark:bg-zinc-950/60 border border-slate-200 dark:border-zinc-800/60 rounded-lg max-h-[220px] overflow-y-auto">
                  <div v-if="form.content" class="text-[11px] text-slate-600 dark:text-zinc-300 leading-relaxed markdown-body" v-html="mdRender(form.content)"></div>
                  <p v-else class="text-xs text-slate-400 italic">Nothing to preview</p>
                </div>
              </div>

              <!-- Visibility status for edit mode -->
              <div v-if="isEditing" class="flex items-center space-x-2 py-1">
                <span class="text-xs text-slate-500 dark:text-zinc-400">{{ t('visibility') }}:</span>
                <USwitch v-model="form.show" color="primary" size="sm" />
              </div>
   
              <div class="pt-2 flex justify-end gap-2">
                <UButton color="neutral" variant="ghost" size="sm" @click="isAnnouncementModalOpen = false">{{ t('cancel') }}</UButton>
                <UButton type="submit" color="primary" size="sm" :loading="announcementLoading">
                  {{ isEditing ? t('edit_announcement') : t('create_announcement') }}
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
const mdRender = renderMarkdown

const announcements = ref([])
const loading = ref(true)
const isAnnouncementModalOpen = ref(false)
const announcementLoading = ref(false)
const activeTab = ref('write')
const isEditing = ref(false)

const form = ref({
  id: null,
  title: '',
  content: '',
  show: true
})

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

useSeoMeta({
  title: () => `${t('announcements_manage')} - HY-Board`
})

const fetchAnnouncements = async () => {
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
          query GetAdminAnnouncements {
            adminAnnouncements {
              id
              title
              content
              show
              created_at
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    announcements.value = response.data?.adminAnnouncements || []
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
  fetchAnnouncements()
})

const openCreateModal = () => {
  isEditing.value = false
  form.value = {
    id: null,
    title: '',
    content: '',
    show: true
  }
  activeTab.value = 'write'
  isAnnouncementModalOpen.value = true
}

const openEditModal = (item) => {
  isEditing.value = true
  form.value = {
    id: item.id,
    title: item.title,
    content: item.content,
    show: item.show
  }
  activeTab.value = 'write'
  isAnnouncementModalOpen.value = true
}

const saveAnnouncement = async () => {
  announcementLoading.value = true
  const token = useCookie('auth_token').value
  try {
    let query = ''
    let variables = {}

    if (isEditing.value) {
      query = `
        mutation UpdateAnnouncement($id: Int!, $title: String!, $content: String!, $show: Boolean!) {
          updateAnnouncement(id: $id, title: $title, content: $content, show: $show) {
            id
            title
            show
          }
        }
      `
      variables = {
        id: form.value.id,
        title: form.value.title,
        content: form.value.content,
        show: form.value.show
      }
    } else {
      query = `
        mutation CreateAnnouncement($title: String!, $content: String!) {
          createAnnouncement(title: $title, content: $content) {
            id
            title
          }
        }
      `
      variables = {
        title: form.value.title,
        content: form.value.content
      }
    }

    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query,
        variables
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'announcement_saved',
      title: isEditing.value ? t('announcement_updated') : t('announcement_created'),
      description: isEditing.value 
        ? t('announcement_updated_desc').replace('{title}', form.value.title)
        : t('announcement_created_desc').replace('{title}', form.value.title),
      color: 'success'
    })

    isAnnouncementModalOpen.value = false
    await fetchAnnouncements()
  } catch (error) {
    toast.add({
      id: 'announcement_save_failed',
      title: t('announcement_failed'),
      description: error.message || 'Announcement operation failed',
      color: 'error'
    })
  } finally {
    announcementLoading.value = false
  }
}

const handleToggleShow = async (id, showVal) => {
  const token = useCookie('auth_token').value
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation ToggleAnnouncement($id: Int!, $show: Boolean!) {
            toggleAnnouncement(id: $id, show: $show) {
              id
              show
            }
          }
        `,
        variables: {
          id: id,
          show: showVal
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'toggle_success',
      title: t('toggle_success'),
      color: 'success',
      timeout: 1500
    })

    // Local sync state
    const item = announcements.value.find(a => a.id === id)
    if (item) item.show = showVal
  } catch (error) {
    toast.add({
      id: 'toggle_failed',
      title: t('toggle_failed'),
      description: error.message,
      color: 'error'
    })
  }
}

const handleDelete = async (id) => {
  if (!confirm(t('confirm_delete'))) return

  const token = useCookie('auth_token').value
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation DeleteAnnouncement($id: Int!) {
            deleteAnnouncement(id: $id)
          }
        `,
        variables: {
          id: id
        }
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    toast.add({
      id: 'delete_success',
      title: t('delete_success'),
      color: 'success',
      timeout: 1500
    })

    await fetchAnnouncements()
  } catch (error) {
    toast.add({
      id: 'delete_failed',
      title: t('delete_failed'),
      description: error.message,
      color: 'error'
    })
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return '-'
  return date.toLocaleString()
}
</script>
