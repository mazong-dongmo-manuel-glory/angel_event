<template>
  <div class="rentals-page">
    <Header />
    <div class="container rentals-container">
      <div class="rentals-header fade-in-up">
        <h1 class="font-script text-gold">{{ t('rentals_public.title') }}</h1>
        <p class="subtitle">{{ t('rentals_public.subtitle') }}</p>
      </div>

      <!-- Filters -->
      <div class="filters fade-in-up">
        <button 
          v-for="cat in categories" 
          :key="cat.slug" 
          @click="activeCategorySlug = cat.slug"
          :class="{ active: activeCategorySlug === cat.slug }"
        >
          {{ cat.slug === 'all' ? t('rentals_public.categories.all') : cat.name }}
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>{{ t('rentals_public.loading') }}</p>
      </div>

      <!-- Items Grid -->
      <div v-else class="items-grid">
        <div v-for="item in filteredItems" :key="item.id" class="rental-card fade-in-up">
          <div class="card-image">
            <img :src="getImageWithFallback(item.image_url, 'rental', item.category_enum)" :alt="item.title" loading="lazy" />
            <span v-if="item.featured" class="tag featured">{{ t('rentals_public.featured') }}</span>
          </div>
          <div class="card-content">
            <div class="card-header">
              <h3>{{ item.title }}</h3>
              <span class="price">{{ formatPrice(item.price) }}</span>
            </div>
            <p v-if="item.description" class="description">{{ item.description }}</p>
            <button class="contact-btn" @click="contactForItem(item)">
              {{ t('rentals_public.book_btn') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="!loading && filteredItems.length === 0" class="empty-state">
        <p>{{ t('rentals_public.empty') }}</p>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '../../services/api'
import { useRouter } from 'vue-router'
import Header from '../../components/Header.vue'
import Footer from '../../components/Footer.vue'
import { getImageWithFallback } from '../../config/defaultImages'

const router = useRouter()
const { t } = useI18n()
const items = ref([])
const loading = ref(true)
const activeCategorySlug = ref('all')
const categories = ref([])

async function fetchCategories() {
  try {
    const res = await api.get('/public/categories')
    const cats = res.data.filter(c => c.type === 'rental')
    categories.value = [
      { name: 'Tout voir', slug: 'all' },
      ...cats
    ]
  } catch (err) {
    console.error('Erreur categories', err)
  }
}

const filteredItems = computed(() => {
  if (activeCategorySlug.value === 'all') return items.value
  // Find category ID from slug
  const cat = categories.value.find(c => c.slug === activeCategorySlug.value)
  if (!cat) return items.value
  
  return items.value.filter(item => item.category_id === cat.id)
})

function formatPrice(price) {
  return new Intl.NumberFormat('fr-CA', { style: 'currency', currency: 'CAD' }).format(price)
}

function contactForItem(item) {
  router.push({
    name: 'contact',
    query: { 
      rental_item_name: item.title,
      rental_item_id: item.id
    }
  })
}

async function fetchItems() {
  loading.value = true
  try {
    const response = await api.get('/public/rentals')
    items.value = response.data
  } catch (error) {
    console.error('Failed to fetch rentals:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCategories()
  fetchItems()
})
</script>

<style scoped>
.rentals-page {
  min-height: 100vh;
  background-color: var(--color-ivory);
}

.rentals-container {
  padding: calc(80px + var(--spacing-5xl)) 0 var(--spacing-5xl);
}

.rentals-header {
  text-align: center;
  margin-bottom: var(--spacing-4xl);
}

.rentals-header h1 {
  font-size: var(--font-size-5xl);
  margin-bottom: var(--spacing-md);
}

.subtitle {
  font-size: var(--font-size-xl);
  color: var(--color-gray);
}

.filters {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-4xl);
}

.filters button {
  padding: var(--spacing-sm) var(--spacing-xl);
  border: 1px solid var(--color-border);
  background: white;
  border-radius: var(--radius-full);
  cursor: pointer;
  font-family: var(--font-sans);
  font-weight: var(--font-weight-medium);
  transition: all var(--transition-base);
  color: var(--color-gray);
}

.filters button.active,
.filters button:hover {
  background: var(--color-gold);
  color: white;
  border-color: var(--color-gold);
  transform: translateY(-2px);
}

.items-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--spacing-xl);
}

.rental-card {
  background: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-md);
  transition: transform var(--transition-base), box-shadow var(--transition-base);
}

.rental-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-xl);
}

.card-image {
  position: relative;
  aspect-ratio: 4/3;
  overflow: hidden;
}

.card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-slow);
}

.rental-card:hover .card-image img {
  transform: scale(1.05);
}

.tag {
  position: absolute;
  top: 1rem;
  right: 1rem;
  padding: 0.4rem 0.8rem;
  border-radius: var(--radius-full);
  font-size: 0.8rem;
  font-weight: bold;
}

.tag.featured {
  background: white;
  color: #e74c3c;
  box-shadow: var(--shadow-sm);
}

.card-content {
  padding: var(--spacing-lg);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-xs);
}

.card-header h3 {
  font-family: var(--font-heading);
  font-size: 1.2rem;
  margin: 0;
  color: var(--color-charcoal);
  font-weight: var(--font-weight-semibold);
}

.price {
  font-weight: bold;
  color: var(--color-gold);
  font-size: 1.1rem;
}

.description {
  color: var(--color-gray);
  font-size: 0.9rem;
  line-height: 1.5;
  margin-bottom: var(--spacing-lg);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.contact-btn {
  width: 100%;
  padding: var(--spacing-md);
  background: var(--color-black);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background var(--transition-base);
  font-weight: 500;
  font-family: var(--font-sans);
}

.contact-btn:hover {
  background: var(--color-charcoal);
}

.loading-state,
.empty-state {
  text-align: center;
  padding: var(--spacing-4xl) 0;
  color: var(--color-gray);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(0,0,0,0.1);
  border-top-color: var(--color-gold);
  border-radius: 50%;
  margin: 0 auto 1rem;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.fade-in-up {
  animation: fadeInUp 0.5s ease-out forwards;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 768px) {
  .rentals-header h1 {
    font-size: 2.5rem;
  }
}
</style>

