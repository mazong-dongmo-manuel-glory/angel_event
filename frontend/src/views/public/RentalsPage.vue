<template>
  <div class="rentals-page">
    <Header />
    <div class="container rentals-container">
      <div class="rentals-header fade-in-up">
        <h1 class="font-script text-gold">{{ t('rentals_public.title') }}</h1>
        <p class="subtitle">{{ t('rentals_public.subtitle') }}</p>
      </div>

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

      <div v-if="cartCount" class="cart-overview fade-in-up">
        <div class="cart-overview-header">
          <div>
            <p class="cart-kicker">{{ t('rentals_public.cart.kicker') }}</p>
            <h2>{{ t('rentals_public.cart.title') }}</h2>
            <p class="cart-subtitle">
              {{ t('rentals_public.cart.count', { count: cartCount }) }}
            </p>
          </div>

          <div class="cart-total">
            <span>{{ t('rentals_public.cart.total') }}</span>
            <strong>{{ formatPrice(cartTotal) }}</strong>
          </div>
        </div>

        <div class="cart-items-preview">
          <div v-for="item in cartItems" :key="item.id" class="cart-preview-item">
            <img
              :src="getImageWithFallback(item.image_url, 'rental', item.category?.slug || item.category_enum)"
              :alt="item.title"
              loading="lazy"
            />

            <div class="cart-preview-details">
              <strong>{{ item.title }}</strong>
              <span>{{ formatPrice(item.price) }}</span>
            </div>

            <button class="remove-item-btn" type="button" @click="removeFromCart(item.id)">
              {{ t('rentals_public.cart.remove') }}
            </button>
          </div>
        </div>

        <div class="cart-overview-actions">
          <Button variant="ghost" size="sm" @click="clearCart">
            {{ t('rentals_public.cart.clear') }}
          </Button>
          <Button variant="secondary" size="md" @click="goToBooking">
            {{ t('rentals_public.cart.checkout') }}
          </Button>
        </div>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>{{ t('rentals_public.loading') }}</p>
      </div>

      <div v-else class="items-grid">
        <div
          v-for="item in filteredItems"
          :key="item.id"
          class="rental-card fade-in-up"
          :class="{ selected: isInCart(item.id) }"
        >
          <div class="card-image">
            <img
              :src="getImageWithFallback(item.image_url, 'rental', item.category?.slug || item.category_enum)"
              :alt="item.title"
              loading="lazy"
            />
            <span v-if="item.featured" class="tag featured">{{ t('rentals_public.featured') }}</span>
            <span v-if="isInCart(item.id)" class="tag selected-tag">{{ t('rentals_public.in_cart') }}</span>
          </div>

          <div class="card-content">
            <div class="card-header">
              <h3>{{ item.title }}</h3>
              <span class="price">{{ formatPrice(item.price) }}</span>
            </div>
            <p v-if="item.description" class="description">{{ item.description }}</p>

            <Button
              class="cart-btn"
              :variant="isInCart(item.id) ? 'secondary' : 'primary'"
              block
              @click="toggleCartItem(item)"
            >
              {{ isInCart(item.id) ? t('rentals_public.remove_btn') : t('rentals_public.add_btn') }}
            </Button>
          </div>
        </div>
      </div>

      <div v-if="!loading && filteredItems.length === 0" class="empty-state">
        <p>{{ t('rentals_public.empty') }}</p>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import api from '../../services/api'
import Header from '../../components/Header.vue'
import Footer from '../../components/Footer.vue'
import Button from '../../components/ui/Button.vue'
import { getImageWithFallback } from '../../config/defaultImages'
import { useRentalCartStore } from '../../stores/rentalCart'

const router = useRouter()
const { t, locale } = useI18n()
const rentalCartStore = useRentalCartStore()
const { cartItems, cartCount, cartTotal } = storeToRefs(rentalCartStore)

const items = ref([])
const loading = ref(true)
const activeCategorySlug = ref('all')
const categories = ref([])

async function fetchCategories() {
  try {
    const res = await api.get('/public/categories')
    const cats = (res.data || []).filter((category) => category.type === 'rental')
    categories.value = [{ name: t('rentals_public.categories.all'), slug: 'all' }, ...cats]
  } catch (err) {
    console.error('Erreur categories', err)
  }
}

const filteredItems = computed(() => {
  if (activeCategorySlug.value === 'all') {
    return items.value
  }

  const selectedCategory = categories.value.find((category) => category.slug === activeCategorySlug.value)
  if (!selectedCategory) {
    return items.value
  }

  return items.value.filter((item) => item.category_id === selectedCategory.id)
})

function formatPrice(price) {
  const activeLocale = locale.value === 'en' ? 'en-CA' : 'fr-CA'
  return new Intl.NumberFormat(activeLocale, { style: 'currency', currency: 'CAD' }).format(price || 0)
}

function isInCart(itemId) {
  return rentalCartStore.hasItem(itemId)
}

function toggleCartItem(item) {
  rentalCartStore.toggleItem(item)
}

function removeFromCart(itemId) {
  rentalCartStore.removeItem(itemId)
}

function clearCart() {
  rentalCartStore.clearCart()
}

function goToBooking() {
  router.push({ name: 'booking' })
}

async function fetchItems() {
  loading.value = true
  try {
    const response = await api.get('/public/rentals')
    items.value = response.data || []
    rentalCartStore.syncCatalog(items.value)
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

.cart-overview {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.96) 0%, rgba(247, 231, 206, 0.92) 100%);
  border: 1px solid rgba(212, 175, 55, 0.25);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
  padding: var(--spacing-2xl);
  margin-bottom: var(--spacing-3xl);
}

.cart-overview-header {
  display: flex;
  justify-content: space-between;
  gap: var(--spacing-xl);
  align-items: flex-start;
  margin-bottom: var(--spacing-xl);
}

.cart-kicker {
  margin: 0 0 var(--spacing-xs);
  text-transform: uppercase;
  letter-spacing: 0.14em;
  font-size: var(--font-size-xs);
  color: var(--color-gold);
}

.cart-overview h2 {
  margin-bottom: var(--spacing-xs);
  font-size: var(--font-size-2xl);
}

.cart-subtitle {
  color: var(--color-gray);
  margin: 0;
}

.cart-total {
  min-width: 180px;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: var(--spacing-xs);
  color: var(--color-gray);
}

.cart-total strong {
  font-size: var(--font-size-2xl);
  color: var(--color-gold);
}

.cart-items-preview {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-xl);
}

.cart-preview-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  background: rgba(255, 255, 255, 0.92);
  border-radius: var(--radius-lg);
  padding: var(--spacing-md);
}

.cart-preview-item img {
  width: 72px;
  height: 72px;
  object-fit: cover;
  border-radius: var(--radius-md);
  flex-shrink: 0;
}

.cart-preview-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.cart-preview-details strong {
  color: var(--color-charcoal);
}

.cart-preview-details span {
  color: var(--color-gold);
  font-weight: var(--font-weight-semibold);
}

.remove-item-btn {
  margin-left: auto;
  border: none;
  background: transparent;
  color: var(--color-error);
  cursor: pointer;
  font-weight: var(--font-weight-semibold);
}

.cart-overview-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-md);
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

.rental-card.selected {
  box-shadow: 0 18px 35px rgba(212, 175, 55, 0.18);
  outline: 2px solid rgba(212, 175, 55, 0.35);
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

.selected-tag {
  left: 1rem;
  right: auto;
  background: var(--color-gold);
  color: white;
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

.cart-btn {
  width: 100%;
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
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-top-color: var(--color-gold);
  border-radius: 50%;
  margin: 0 auto 1rem;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
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

  .cart-overview-header {
    flex-direction: column;
  }

  .cart-total {
    align-items: flex-start;
  }

  .cart-overview-actions {
    flex-direction: column;
  }
}
</style>
