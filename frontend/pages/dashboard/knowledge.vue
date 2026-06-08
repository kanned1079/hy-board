<template>
  <div class="max-w-5xl mx-auto h-[calc(100vh-8rem)] lg:h-[calc(100vh-4rem)] flex flex-col space-y-4">
    <!-- Header & Search -->
    <header class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3 border-b border-slate-200 dark:border-zinc-800/80 pb-3 shrink-0">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white">{{ t('knowledge_base') }}</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('knowledge_list') }}</p>
      </div>
      <!-- Search Input -->
      <div class="w-full sm:w-64">
        <UInput
          v-model="searchQuery"
          icon="i-lucide-search"
          size="sm"
          color="primary"
          :placeholder="t('search_placeholder') || 'Search articles...'"
          class="w-full"
        />
      </div>
    </header>

    <!-- Main Content Area -->
    <div class="flex-1 min-h-0 flex flex-col lg:flex-row gap-4">
      <!-- Left: Article List -->
      <div 
        class="w-full lg:w-80 flex flex-col min-h-0 border border-slate-200 dark:border-zinc-800/80 rounded-lg bg-white/70 dark:bg-zinc-900/40 backdrop-blur-md overflow-hidden shrink-0"
        :class="{ 'hidden lg:flex': selectedArticle && isMobile }"
      >
        <div class="p-3 border-b border-slate-200 dark:border-zinc-800/80 bg-slate-50/50 dark:bg-zinc-900/50 shrink-0">
          <span class="text-xs font-semibold text-slate-500 dark:text-zinc-400 uppercase tracking-wider">{{ t('knowledge_list') }}</span>
        </div>

        <!-- Skeletons -->
        <div v-if="loading" class="flex-1 overflow-y-auto p-3 space-y-3">
          <USkeleton class="h-16 w-full rounded-lg bg-slate-200 dark:bg-zinc-900/40" v-for="i in 4" :key="i" />
        </div>

        <!-- Empty State -->
        <div v-else-if="filteredArticles.length === 0" class="flex-1 flex flex-col items-center justify-center p-6 text-slate-500 dark:text-zinc-400 text-xs">
          <UIcon name="i-lucide-file-text" class="w-8 h-8 opacity-40 mb-2" />
          <span>{{ t('no_articles') }}</span>
        </div>

        <!-- Article items -->
        <div v-else class="flex-1 overflow-y-auto divide-y divide-slate-100 dark:divide-slate-800/60 custom-scrollbar">
          <button
            v-for="article in filteredArticles"
            :key="article.id"
            @click="selectArticle(article)"
            class="w-full p-4 text-left transition-all duration-200 flex flex-col space-y-1.5"
            :class="selectedArticle?.id === article.id 
              ? 'bg-primary-500/10 dark:bg-primary-500/15 border-l-4 border-primary-500 pl-3' 
              : 'hover:bg-slate-100/50 dark:hover:bg-zinc-900/30 pl-4 border-l-4 border-transparent'"
          >
            <h3 class="text-xs font-bold text-slate-800 dark:text-zinc-100 line-clamp-1">
              {{ article.title }}
            </h3>
            <div class="flex items-center justify-between text-[9px] text-slate-400 dark:text-zinc-500 font-mono">
              <span>{{ formatDate(article.created_at) }}</span>
              <span class="bg-slate-200/50 dark:bg-zinc-800/50 px-1.5 py-0.5 rounded-sm">#{{ article.sort }}</span>
            </div>
          </button>
        </div>
      </div>

      <!-- Right: Article Reader -->
      <div 
        class="flex-1 min-h-0 flex flex-col border border-slate-200 dark:border-zinc-800/80 rounded-lg bg-white/70 dark:bg-zinc-900/40 backdrop-blur-md overflow-hidden"
        :class="{ 'flex': selectedArticle || !isMobile, 'hidden lg:flex': !selectedArticle }"
      >
        <!-- Reader Header -->
        <div class="p-3 border-b border-slate-200 dark:border-zinc-800/80 bg-slate-50/50 dark:bg-zinc-900/50 flex items-center justify-between shrink-0">
          <div class="flex items-center space-x-2">
            <UButton
              v-if="isMobile && selectedArticle"
              icon="i-lucide-arrow-left"
              color="gray"
              variant="ghost"
              size="xs"
              @click="selectedArticle = null"
            />
            <span class="text-xs font-semibold text-slate-500 dark:text-zinc-400 uppercase tracking-wider">{{ t('reader_mode') || 'Reading Pane' }}</span>
          </div>
          <span v-if="selectedArticle" class="text-[9px] font-mono text-slate-400 dark:text-zinc-500">
            {{ formatDate(selectedArticle.created_at) }}
          </span>
        </div>

        <!-- Reader Content -->
        <div class="flex-1 overflow-y-auto p-5 custom-scrollbar">
          <div v-if="selectedArticle" class="space-y-4">
            <h2 class="text-lg sm:text-xl font-extrabold text-slate-900 dark:text-white border-b border-slate-100 dark:border-zinc-850 pb-2">
              {{ selectedArticle.title }}
            </h2>
            <article 
              class="text-xs text-slate-700 dark:text-zinc-300 leading-relaxed markdown-body" 
              v-html="mdRender(selectedArticle.content)"
            ></article>
          </div>
          <div v-else class="h-full flex flex-col items-center justify-center text-center p-6 text-slate-400 dark:text-zinc-500">
            <UIcon name="i-lucide-book-open" class="w-12 h-12 opacity-30 mb-2" />
            <p class="text-xs">{{ t('select_article_to_read') || 'Select an article from the list to start reading' }}</p>
          </div>
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
const localePath = useLocalePath()
const mdRender = renderMarkdown

const articles = ref([])
const loading = ref(true)
const selectedArticle = ref(null)
const searchQuery = ref('')

const toast = useToast()
const router = useRouter()
const config = useRuntimeConfig()

useSeoMeta({
  title: () => `${t('knowledge_base')} - HY-Board`
})

// Mobile detection
const isMobile = ref(false)
const checkMobile = () => {
  if (typeof window !== 'undefined') {
    isMobile.value = window.innerWidth < 1024
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchArticles()
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('resize', checkMobile)
  }
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
          query GetKnowledges {
            knowledges {
              id
              title
              content
              sort
              created_at
            }
          }
        `
      }
    })

    if (response.errors && response.errors.length > 0) {
      throw new Error(response.errors[0].message)
    }

    articles.value = response.data?.knowledges || []
    
    // Auto-select first article on desktop
    if (!isMobile.value && articles.value.length > 0) {
      selectedArticle.value = articles.value[0]
    }
  } catch (error) {
    toast.add({
      id: 'knowledge_fetch_failed',
      title: t('session_expired'),
      description: error.message || 'Failed to fetch articles',
      color: 'red'
    })
    router.push(localePath('/login'))
  } finally {
    loading.value = false
  }
}

const selectArticle = (article) => {
  selectedArticle.value = article
}

const filteredArticles = computed(() => {
  if (!searchQuery.value) return articles.value
  const query = searchQuery.value.toLowerCase()
  return articles.value.filter(a => 
    a.title.toLowerCase().includes(query) || 
    a.content.toLowerCase().includes(query)
  )
})

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return '-'
  return date.toLocaleDateString()
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
