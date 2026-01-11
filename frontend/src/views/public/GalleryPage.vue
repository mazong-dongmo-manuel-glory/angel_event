<template>
  <div class="gallery-page">
    <Header />
    <div class="container gallery-container">
      <div class="gallery-header fade-in-up">
        <h1 class="font-script text-gold">Notre Galerie</h1>
        <p class="subtitle">Découvrez nos plus belles créations</p>
      </div>

      <!-- Category Sections -->
      <div class="categories-wrapper">
        <div 
          v-for="category in categories" 
          :key="category.value"
          class="category-section fade-in-up"
        >
          <div class="category-header">
            <h2 class="category-title">{{ category.label }}</h2>
            <button 
              v-if="categoryImages[category.value]?.length > 0"
              class="next-btn"
              @click="loadNextImages(category.value)"
              :disabled="loadingCategories[category.value]"
            >
              <span v-if="!loadingCategories[category.value]">Suivant →</span>
              <span v-else class="spinner-small"></span>
            </button>
          </div>

          <!-- Loading State -->
          <div v-if="loadingCategories[category.value] && !categoryImages[category.value]" class="loading-state">
            <div class="spinner"></div>
          </div>

          <!-- Gallery Grid -->
          <div v-else-if="categoryImages[category.value]?.length > 0" class="gallery-grid">
            <div
              v-for="(image, index) in categoryImages[category.value]"
              :key="image.id"
              class="gallery-item fade-in-up"
              @click="openLightbox(category.value, index)"
            >
              <img :src="image.image_url" :alt="category.label" />
            </div>
          </div>

          <!-- Empty State -->
          <div v-else class="empty-state">
            <p>Aucune image disponible pour cette catégorie.</p>
          </div>
        </div>
      </div>

      <!-- Lightbox -->
      <Teleport to="body">
        <Transition name="lightbox">
          <div v-if="lightboxOpen" class="lightbox" @click="closeLightbox">
            <button class="lightbox-close" @click="closeLightbox">×</button>
            <button class="lightbox-prev" @click.stop="prevImage">‹</button>
            <button class="lightbox-next" @click.stop="nextImage">›</button>
            <div class="lightbox-content" @click.stop>
              <img :src="currentImage.image_url" :alt="currentCategory" />
            </div>
          </div>
        </Transition>
      </Teleport>
    </div>
    <Footer />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Header from '../../components/Header.vue'
import Footer from '../../components/Footer.vue'
import api from '../../services/api'

const categoryImages = ref({})
const loadingCategories = ref({})
const lightboxOpen = ref(false)
const currentCategory = ref('')
const currentImageIndex = ref(0)
const excludedIds = ref({}) // Track excluded IDs per category

const categories = [
  { label: 'Baby Shower', value: 'baby_shower' },
  { label: 'Baptême', value: 'bapteme' },
  { label: 'Anniversaire', value: 'birthday' },
  { label: 'Félicitations', value: 'congrats' },
  { label: 'Love Room', value: 'loveroom' },
  { label: 'Demande en Mariage', value: 'marryme' },
  { label: 'Mariage', value: 'wedding' }
]

const currentImage = computed(() => {
  const images = categoryImages.value[currentCategory.value] || []
  return images[currentImageIndex.value] || {}
})

async function loadCategoryImages(category, isNext = false) {
  loadingCategories.value[category] = true
  
  try {
    // Get excluded IDs for this category
    const excluded = excludedIds.value[category] || []
    const excludeParam = excluded.join(',')
    
    const response = await api.get('/public/gallery/random', {
      params: {
        category,
        limit: 3,
        exclude: excludeParam
      }
    })
    
    const newImages = response.data
    
    if (isNext) {
      // Add new image IDs to excluded list
      newImages.forEach(img => {
        if (!excluded.includes(img.id)) {
          excluded.push(img.id)
        }
      })
      excludedIds.value[category] = excluded
    }
    
    categoryImages.value[category] = newImages
  } catch (error) {
    console.error(`Failed to fetch images for ${category}:`, error)
    categoryImages.value[category] = []
  } finally {
    loadingCategories.value[category] = false
  }
}

async function loadNextImages(category) {
  await loadCategoryImages(category, true)
}

function openLightbox(category, index) {
  currentCategory.value = category
  currentImageIndex.value = index
  lightboxOpen.value = true
  document.body.style.overflow = 'hidden'
}

function closeLightbox() {
  lightboxOpen.value = false
  document.body.style.overflow = ''
}

function nextImage() {
  const images = categoryImages.value[currentCategory.value] || []
  currentImageIndex.value = (currentImageIndex.value + 1) % images.length
}

function prevImage() {
  const images = categoryImages.value[currentCategory.value] || []
  currentImageIndex.value = (currentImageIndex.value - 1 + images.length) % images.length
}

onMounted(() => {
  // Load initial images for all categories
  categories.forEach(cat => {
    loadCategoryImages(cat.value, false)
  })
  
  // Keyboard navigation
  document.addEventListener('keydown', (e) => {
    if (!lightboxOpen.value) return
    if (e.key === 'Escape') closeLightbox()
    if (e.key === 'ArrowRight') nextImage()
    if (e.key === 'ArrowLeft') prevImage()
  })
})
</script>

<style scoped>
.gallery-page {
  min-height: 100vh;
  background: var(--color-ivory);
}

.gallery-container {
  padding: calc(80px + var(--spacing-5xl)) 0 var(--spacing-5xl);
}

.gallery-header {
  text-align: center;
  margin-bottom: var(--spacing-4xl);
}

.gallery-header h1 {
  font-size: var(--font-size-5xl);
  margin-bottom: var(--spacing-md);
}

.subtitle {
  font-size: var(--font-size-xl);
  color: var(--color-gray);
}

/* Categories */
.categories-wrapper {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-5xl);
}

.category-section {
  background: white;
  border-radius: var(--radius-xl);
  padding: var(--spacing-2xl);
  box-shadow: var(--shadow-md);
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-xl);
  padding-bottom: var(--spacing-lg);
  border-bottom: 2px solid var(--color-gold);
}

.category-title {
  font-family: var(--font-script);
  font-size: var(--font-size-3xl);
  color: var(--color-gold);
  margin: 0;
}

.next-btn {
  padding: var(--spacing-sm) var(--spacing-xl);
  background: var(--color-gold);
  color: white;
  border: none;
  border-radius: var(--radius-full);
  font-family: var(--font-sans);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: all var(--transition-base);
  min-width: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.next-btn:hover:not(:disabled) {
  background: var(--color-gold-dark, #c4a962);
  transform: translateX(5px);
}

.next-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.spinner-small {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

/* Gallery Grid */
.gallery-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--spacing-xl);
}

.gallery-item {
  position: relative;
  border-radius: var(--radius-lg);
  overflow: hidden;
  cursor: pointer;
  aspect-ratio: 4/3;
  box-shadow: var(--shadow-md);
  transition: transform var(--transition-base);
}

.gallery-item:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-xl);
}

.gallery-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-slow);
}

.gallery-item:hover img {
  transform: scale(1.05);
}

/* Loading & Empty States */
.loading-state,
.empty-state {
  text-align: center;
  padding: var(--spacing-4xl) 0;
  color: var(--color-gray);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-top-color: var(--color-gold);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto var(--spacing-md);
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Lightbox */
.lightbox {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.95);
  z-index: var(--z-modal);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-xl);
}

.lightbox-close,
.lightbox-prev,
.lightbox-next {
  position: absolute;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.3);
  color: white;
  font-size: var(--font-size-3xl);
  cursor: pointer;
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-base);
  z-index: 1;
  font-weight: 300;
  line-height: 1;
}

.lightbox-close:hover,
.lightbox-prev:hover,
.lightbox-next:hover {
  background: rgba(255, 255, 255, 0.25);
  border-color: rgba(255, 255, 255, 0.5);
  transform: scale(1.1);
}

.lightbox-close {
  top: var(--spacing-2xl);
  right: var(--spacing-2xl);
}

.lightbox-prev {
  left: var(--spacing-2xl);
  top: 50%;
  transform: translateY(-50%);
}

.lightbox-prev:hover {
  transform: translateY(-50%) scale(1.1);
}

.lightbox-next {
  right: var(--spacing-2xl);
  top: 50%;
  transform: translateY(-50%);
}

.lightbox-next:hover {
  transform: translateY(-50%) scale(1.1);
}

.lightbox-content {
  max-width: 90vw;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.lightbox-content img {
  max-width: 100%;
  max-height: 90vh;
  object-fit: contain;
  border-radius: var(--radius-lg);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

/* Transitions */
.lightbox-enter-active,
.lightbox-leave-active {
  transition: opacity var(--transition-base);
}

.lightbox-enter-from,
.lightbox-leave-to {
  opacity: 0;
}

.fade-in-up {
  animation: fadeInUp 0.6s ease-out;
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
  .gallery-grid {
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: var(--spacing-md);
  }
  
  .category-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-md);
  }
  
  .next-btn {
    width: 100%;
  }
  
  .lightbox-prev,
  .lightbox-next {
    display: none;
  }
}
</style>
