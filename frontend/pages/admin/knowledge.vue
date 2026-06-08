<template>
  <div class="max-w-5xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">{{ t('knowledge_base') }}</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('admin_panel') }}</p>
      </div>
      <UButton
        color="primary"
        icon="i-lucide-plus"
        size="sm"
        @click="openCreateModal"
      >
        {{ t('create_knowledge') }}
      </UButton>
    </header>

    <!-- Knowledge Base List Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-4 shadow-sm">
      <h2 class="text-md font-bold text-slate-900 dark:text-white flex items-center space-x-2">
        <UIcon name="i-lucide-book-open" class="w-4 h-4 text-primary-500 dark:text-primary-400" />
        <span>{{ t('knowledge_list') }}</span>
      </h2>

      <div v-if="loading" class="space-y-3">
        <USkeleton class="h-10 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 3" :key="i" />
      </div>

      <div v-else-if="articles.length === 0" class="text-center py-6 text-slate-500 text-xs">
        {{ t('no_articles') }}
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="border-b border-slate-200 dark:border-zinc-800 text-slate-500 dark:text-zinc-400 font-medium">
              <th class="py-2 px-3 w-16">ID</th>
              <th class="py-2 px-3">{{ t('title') }}</th>
              <th class="py-2 px-3 w-24">{{ t('sort') }}</th>
              <th class="py-2 px-3 w-32">{{ t('date') }}</th>
              <th class="py-2 px-3 w-28 text-center">{{ t('visibility') }}</th>
              <th class="py-2 px-3 w-28 text-right"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in articles" :key="item.id" class="border-b border-slate-200/50 dark:border-zinc-800/40 hover:bg-slate-100/50 dark:hover:bg-zinc-900/20 text-slate-700 dark:text-zinc-200">
              <td class="py-2.5 px-3 font-mono">{{ item.id }}</td>
              <td class="py-2.5 px-3 font-semibold">{{ item.title }}</td>
              <td class="py-2.5 px-3 font-mono">
                <UBadge color="primary" variant="soft" size="xs">{{ item.sort }}</UBadge>
              </td>
              <td class="py-2.5 px-3 text-slate-500 font-mono">{{ formatDate(item.created_at) }}</td>
              <td class="py-2.5 px-3 text-center">
                <USwitch
                  :model-value="item.show"
                  color="primary"
                  size="sm"
                  @update:model-value="(val) => handleToggleShow(item, val)"
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

    <!-- Create / Edit Article Drawer -->
    <USlideover v-model:open="isModalOpen" :title="isEditing ? t('edit_knowledge') : t('create_knowledge')" prevent-close>
      <template #body>
        <form id="article-form" @submit.prevent="saveArticle" class="space-y-4">
          <UFormField :label="t('title')" name="title" class="text-slate-700 dark:text-zinc-300">
            <UInput v-model="form.title" placeholder="How to use the client" color="primary" size="sm" class="w-full" required />
          </UFormField>

          <div class="grid grid-cols-2 gap-3">
            <UFormField :label="t('sort')" name="sort" class="text-slate-700 dark:text-zinc-300">
              <UInput v-model.number="form.sort" type="number" color="primary" size="sm" class="w-full" required />
            </UFormField>
            <div class="flex flex-col justify-end pb-1.5 pl-1">
              <div class="flex items-center space-x-2">
                <span class="text-xs text-slate-500 dark:text-zinc-400">{{ t('visibility') }}:</span>
                <USwitch v-model="form.show" color="primary" size="sm" />
              </div>
            </div>
          </div>

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
                <UTextarea v-model="form.content" placeholder="Write markdown content here..." color="primary" size="sm" rows="16" class="w-full" required />
              </UFormField>
            </div>

            <!-- Preview Tab Content -->
            <div v-show="activeTab === 'preview'" class="p-3 bg-slate-50 dark:bg-zinc-950/60 border border-slate-200 dark:border-zinc-800/60 rounded-lg max-h-[400px] overflow-y-auto">
              <div v-if="form.content" class="text-[11px] text-slate-600 dark:text-zinc-300 leading-relaxed markdown-body" v-html="mdRender(form.content)"></div>
              <p v-else class="text-xs text-slate-400 italic">Nothing to preview</p>
            </div>
          </div>
        </form>
      </template>

      <template #footer="{ close }">
        <UButton color="neutral" variant="ghost" size="sm" @click="close">{{ t('cancel') }}</UButton>
        <UButton type="submit" form="article-form" color="primary" size="sm" :loading="submitLoading">
          {{ isEditing ? t('edit_knowledge') : t('create_knowledge') }}
        </UButton>
      </template>
    </USlideover>

  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})

const { t } = useI18n()
const localePath = useLocalePath()
const mdRender = renderMarkdown

const articles = ref([])
const loading = ref(true)
const isModalOpen = ref(false)
const isEditing = ref(false)
const submitLoading = ref(false)
const activeTab = ref('write')

const form = ref({
  id: null,
  title: '',
  content: '',
  sort: 0,
  show: true
})

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

useSeoMeta({
  title: () => `${t('knowledge_base')} - HY-Board`
})

const fetchArticles = async () => {
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
          query GetAdminKnowledges {
            adminKnowledges {
              id
              title
              content
              show
              sort
              created_at
              updated_at
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    articles.value = response.data?.adminKnowledges || []
  } catch (error) {
    toast.add({
      id: 'session_expired',
      title: t('session_expired'),
      description: error.message || t('login_again_admin'),
      color: 'red'
    })
    router.push(localePath('/login'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchArticles()
})

const openCreateModal = () => {
  isEditing.value = false
  form.value = {
    id: null,
    title: '',
    content: '',
    sort: 0,
    show: true
  }
  activeTab.value = 'write'
  isModalOpen.value = true
}

const openEditModal = (item) => {
  isEditing.value = true
  form.value = {
    id: item.id,
    title: item.title,
    content: item.content,
    sort: item.sort,
    show: item.show
  }
  activeTab.value = 'write'
  isModalOpen.value = true
}

const saveArticle = async () => {
  submitLoading.value = true
  const token = useCookie('auth_token').value
  try {
    let query = ''
    let variables = {}

    if (isEditing.value) {
      query = `
        mutation UpdateKnowledge($id: Int!, $title: String!, $content: String!, $show: Boolean!, $sort: Int!) {
          updateKnowledge(id: $id, title: $title, content: $content, show: $show, sort: $sort) {
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
        show: form.value.show,
        sort: form.value.sort
      }
    } else {
      query = `
        mutation CreateKnowledge($title: String!, $content: String!, $sort: Int!) {
          createKnowledge(title: $title, content: $content, sort: $sort) {
            id
            title
          }
        }
      `
      variables = {
        title: form.value.title,
        content: form.value.content,
        sort: form.value.sort
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
      id: 'article_saved',
      title: isEditing.value ? t('update_success') : t('article_created'),
      description: isEditing.value 
        ? t('update_success') 
        : t('article_created_desc').replace('{title}', form.value.title),
      color: 'success'
    })

    isModalOpen.value = false
    await fetchArticles()
  } catch (error) {
    toast.add({
      id: 'article_save_failed',
      title: t('update_failed'),
      description: error.message || 'Article operation failed',
      color: 'error'
    })
  } finally {
    submitLoading.value = false
  }
}

const handleToggleShow = async (item, showVal) => {
  const token = useCookie('auth_token').value
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation UpdateKnowledge($id: Int!, $title: String!, $content: String!, $show: Boolean!, $sort: Int!) {
            updateKnowledge(id: $id, title: $title, content: $content, show: $show, sort: $sort) {
              id
              show
            }
          }
        `,
        variables: {
          id: item.id,
          title: item.title,
          content: item.content,
          show: showVal,
          sort: item.sort
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

    // Sync locally
    item.show = showVal
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
  if (!confirm(t('confirm_delete_knowledge'))) return

  const token = useCookie('auth_token').value
  try {
    const response = await $fetch(`${config.public.apiBase}/graphql`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        query: `
          mutation DeleteKnowledge($id: Int!) {
            deleteKnowledge(id: $id)
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

    await fetchArticles()
  } catch (error) {
    toast.add({
      id: 'delete_failed',
      title: t('update_failed'),
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
