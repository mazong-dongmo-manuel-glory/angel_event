import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '../services/api'

export function useSiteContent(sectionFilter) {
  const { locale } = useI18n()
  const items = ref([])
  const loading = ref(false)

  const normalizedSections = computed(() => {
    if (!sectionFilter) {
      return []
    }

    return Array.isArray(sectionFilter) ? sectionFilter : [sectionFilter]
  })

  const contentMap = computed(() => {
    const entries = items.value.map((item) => [item.key, item.value])
    return Object.fromEntries(entries)
  })

  async function fetchContent() {
    loading.value = true

    try {
      const params = { language: locale.value }
      if (normalizedSections.value.length > 0) {
        params.section = normalizedSections.value.join(',')
      }

      const response = await api.get('/public/content', { params })
      items.value = response.data || []
    } catch (error) {
      console.error('Failed to fetch site content:', error)
      items.value = []
    } finally {
      loading.value = false
    }
  }

  function getContent(key, fallback = '') {
    const value = contentMap.value[key]
    return value && value.trim() !== '' ? value : fallback
  }

  function getImage(key, fallback = '') {
    return getContent(key, fallback)
  }

  watch(locale, fetchContent)
  onMounted(fetchContent)

  return {
    items,
    loading,
    contentMap,
    getContent,
    getImage,
    refreshContent: fetchContent,
  }
}
