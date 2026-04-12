import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '../services/api'
import { normalizeImageUrl } from '../config/defaultImages'

export function useSiteContent(sectionFilter) {
  const { locale, fallbackLocale } = useI18n()
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
      const languages = [locale.value]
      const configuredFallback = typeof fallbackLocale.value === 'string'
        ? fallbackLocale.value
        : Array.isArray(fallbackLocale.value)
          ? fallbackLocale.value[0]
          : null

      const alternateLanguage = configuredFallback && configuredFallback !== locale.value
        ? configuredFallback
        : locale.value === 'fr' ? 'en' : 'fr'

      if (!languages.includes(alternateLanguage)) {
        languages.push(alternateLanguage)
      }

      const requests = languages.map((language) => {
        const params = { language }
        if (normalizedSections.value.length > 0) {
          params.section = normalizedSections.value.join(',')
        }

        return api.get('/public/content', { params })
      })

      const responses = await Promise.all(requests)
      const mergedByKey = new Map()

      for (let i = responses.length - 1; i >= 0; i -= 1) {
        const responseItems = responses[i].data || []
        responseItems.forEach((item) => {
          mergedByKey.set(item.key, item)
        })
      }

      items.value = Array.from(mergedByKey.values())
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
    return normalizeImageUrl(getContent(key, fallback))
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
