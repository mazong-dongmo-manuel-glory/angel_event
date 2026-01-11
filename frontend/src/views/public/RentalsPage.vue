<template>
  <div class="rentals-page">
    <div class="hero-section">
      <div class="hero-content">
        <h1>Nos Articles en Location</h1>
        <p>Sublimez votre événement avec notre sélection d'articles</p>
      </div>
    </div>

    <div class="container">
      <!-- Filters -->
      <div class="filters">
        <button 
          v-for="cat in categories" 
          :key="cat.value" 
          @click="activeCategory = cat.value"
          :class="{ active: activeCategory === cat.value }"
        >
          {{ cat.label }}
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Chargement des articles...</p>
      </div>

      <!-- Items Grid -->
      <div v-else class="items-grid">
        <div v-for="item in filteredItems" :key="item.id" class="rental-card">
          <div class="card-image">
            <img :src="item.image_url" :alt="item.title" loading="lazy" />
            <span v-if="item.featured" class="tag featured">Coup de ❤️</span>
          </div>
          <div class="card-content">
            <div class="card-header">
              <h3>{{ item.title }}</h3>
              <span class="price">{{ formatPrice(item.price) }}</span>
            </div>
            <p v-if="item.description" class="description">{{ item.description }}</p>
            <button class="contact-btn" @click="contactForItem(item)">
              Réserver cet article
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="!loading && filteredItems.length === 0" class="empty-state">
        <p>Aucun article trouvé dans cette catégorie.</p>
      </div>
    </div>
  </div>
</template>



<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../../services/api'
import { useRouter } from 'vue-router'

const router = useRouter()
const items = ref([])
const loading = ref(true)
const activeCategory = ref('all')

const categories = [
  { label: 'Tout voir', value: 'all' },
  { label: 'Centres de table', value: 'centerpiece' },
  { label: 'Backdrops / Panneaux', value: 'backdrop' },
  { label: 'Fleurs', value: 'flower' },
  { label: 'Autres articles', value: 'other' }
]

const filteredItems = computed(() => {
  if (activeCategory.value === 'all') return items.value
  return items.value.filter(item => item.category === activeCategory.value)
})

function formatPrice(price) {
  return new Intl.NumberFormat('fr-CA', { style: 'currency', currency: 'CAD' }).format(price)
}

function contactForItem(item) {
  router.push({
    name: 'Booking',
    query: { rental_item_id: item.id }
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
  fetchItems()
})
</script>

<style scoped>
.rentals-page {
  min-height: 100vh;
  background-color: var(--color-background);
}

.hero-section {
  background-color: var(--color-black);
  color: white;
  padding: 6rem 2rem 4rem;
  text-align: center;
  margin-bottom: 3rem;
}

.hero-content h1 {
  font-family: var(--font-heading);
  font-size: 3.5rem;
  margin-bottom: 1rem;
  color: var(--color-gold);
}

.hero-content p {
  font-size: 1.2rem;
  color: rgba(255, 255, 255, 0.8);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem 4rem;
}

.filters {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 3rem;
}

.filters button {
  padding: 0.8rem 1.5rem;
  border: 1px solid #ddd;
  background: white;
  border-radius: 50px;
  cursor: pointer;
  font-family: var(--font-body);
  transition: all 0.3s ease;
}

.filters button.active,
.filters button:hover {
  background: var(--color-gold);
  color: white;
  border-color: var(--color-gold);
}

.items-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
}

.rental-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
  transition: transform 0.3s ease;
}

.rental-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 20px rgba(0,0,0,0.1);
}

.card-image {
  position: relative;
  aspect-ratio: 4/3;
}

.card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.tag {
  position: absolute;
  top: 1rem;
  right: 1rem;
  padding: 0.4rem 0.8rem;
  border-radius: 50px;
  font-size: 0.8rem;
  font-weight: bold;
}

.tag.featured {
  background: white;
  color: #e74c3c;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.card-content {
  padding: 1.5rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.5rem;
}

.card-header h3 {
  font-family: var(--font-heading);
  font-size: 1.2rem;
  margin: 0;
  color: var(--color-charcoal);
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
  margin-bottom: 1.5rem;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.contact-btn {
  width: 100%;
  padding: 0.8rem;
  background: var(--color-black);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.3s ease;
  font-weight: 500;
}

.contact-btn:hover {
  background: #333;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 4rem 0;
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

@media (max-width: 768px) {
  .hero-content h1 {
    font-size: 2.5rem;
  }
}
</style>
