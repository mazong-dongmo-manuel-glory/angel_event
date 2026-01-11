<template>
  <div class="gallery-manager">
    <div class="manager-header">
      <h1>Gestion de la Galerie</h1>
      <div class="header-actions">
        <button class="scan-btn" @click="scanStorageFolder" :disabled="scanning">
          <span v-if="!scanning">📁 Scanner le dossier Storage</span>
          <span v-else>⏳ Scan en cours...</span>
        </button>
        <button class="upload-btn" @click="openUploadModal">
          📤 Uploader une image
        </button>
      </div>
    </div>

    <!-- Category Statistics -->
    <div class="stats-grid">
      <div v-for="stat in categoryStats" :key="stat.category" class="stat-card">
        <h3>{{ getCategoryLabel(stat.category) }}</h3>
        <p class="count">{{ stat.count }}</p>
        <p class="label">images</p>
      </div>
    </div>

    <!-- Filters and Search -->
    <div class="controls">
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Rechercher par titre ou description..."
          class="search-input"
        />
      </div>
      <select v-model="filterCategory" class="category-filter">
        <option value="">Toutes les catégories</option>
        <option v-for="cat in categories" :key="cat.value" :value="cat.value">
          {{ cat.label }}
        </option>
      </select>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement des images...</p>
    </div>

    <!-- Images Grid -->
    <div v-else class="images-grid">
      <div 
        v-for="image in filteredImages" 
        :key="image.id" 
        class="image-card"
      >
        <div class="image-preview">
          <img :src="image.image_url" :alt="image.title" />
          <div class="image-overlay">
            <button class="edit-btn" @click="editImage(image)">✏️ Éditer</button>
            <button class="delete-btn" @click="confirmDelete(image)">🗑️ Supprimer</button>
          </div>
        </div>
        <div class="image-info">
          <h4>{{ image.title }}</h4>
          <p class="category-badge">{{ getCategoryLabel(image.category) }}</p>
          <p v-if="image.description" class="description">{{ image.description }}</p>
          <p class="meta">
            <span v-if="image.is_from_storage" class="badge">📁 Storage</span>
            <span v-if="image.featured" class="badge featured">⭐ Vedette</span>
          </p>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="!loading && filteredImages.length === 0" class="empty-state">
      <p>Aucune image trouvée.</p>
    </div>

    <!-- Upload Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="uploadingImage" class="modal-overlay" @click="closeUploadModal">
          <div class="modal-content" @click.stop>
            <div class="modal-header">
              <h2>Nouvelle Image</h2>
              <button class="close-btn" @click="closeUploadModal">×</button>
            </div>
            <div class="modal-body">
              <form @submit.prevent="uploadImage" class="upload-form">
                <div class="form-group">
                  <label>Image *</label>
                  <input type="file" ref="fileInput" @change="handleFileChange" accept="image/*" class="form-input" required />
                </div>
                <div class="form-group">
                  <label>Titre</label>
                  <input v-model="uploadForm.title" type="text" class="form-input" placeholder="Laisser vide pour utiliser le nom du fichier" />
                </div>
                <div class="form-group">
                  <label>Description</label>
                  <textarea v-model="uploadForm.description" class="form-textarea" rows="3"></textarea>
                </div>
                <div class="form-group">
                  <label>Catégorie *</label>
                  <select v-model="uploadForm.category" class="form-select" required>
                    <option value="" disabled>Choisir une catégorie</option>
                    <option v-for="cat in categories" :key="cat.value" :value="cat.value">
                      {{ cat.label }}
                    </option>
                  </select>
                </div>
                <div class="form-group">
                  <label class="checkbox-label">
                    <input v-model="uploadForm.featured" type="checkbox" />
                    <span>Image vedette</span>
                  </label>
                </div>
                <div class="modal-footer">
                  <button type="button" class="btn-secondary" @click="closeUploadModal">Annuler</button>
                  <button type="submit" class="btn-primary" :disabled="uploading">
                    {{ uploading ? 'Upload en cours...' : 'Uploader' }}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Edit Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="editingImage" class="modal-overlay" @click="closeEditModal">
          <div class="modal-content" @click.stop>
            <div class="modal-header">
              <h2>Éditer l'image</h2>
              <button class="close-btn" @click="closeEditModal">×</button>
            </div>
            <div class="modal-body">
              <div class="form-group">
                <label>Titre</label>
                <input v-model="editForm.title" type="text" class="form-input" />
              </div>
              <div class="form-group">
                <label>Description</label>
                <textarea v-model="editForm.description" class="form-textarea" rows="3"></textarea>
              </div>
              <div class="form-group">
                <label>Catégorie</label>
                <select v-model="editForm.category" class="form-select">
                  <option v-for="cat in categories" :key="cat.value" :value="cat.value">
                    {{ cat.label }}
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label class="checkbox-label">
                  <input v-model="editForm.featured" type="checkbox" />
                  <span>Image vedette</span>
                </label>
              </div>
            </div>
            <div class="modal-footer">
              <button class="btn-secondary" @click="closeEditModal">Annuler</button>
              <button class="btn-primary" @click="saveImage" :disabled="saving">
                {{ saving ? 'Enregistrement...' : 'Enregistrer' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="deletingImage" class="modal-overlay" @click="deletingImage = null">
          <div class="modal-content small" @click.stop>
            <div class="modal-header">
              <h2>Confirmer la suppression</h2>
              <button class="close-btn" @click="deletingImage = null">×</button>
            </div>
            <div class="modal-body">
              <p>Êtes-vous sûr de vouloir supprimer cette image ?</p>
              <p class="warning">Cette action est irréversible.</p>
            </div>
            <div class="modal-footer">
              <button class="btn-secondary" @click="deletingImage = null">Annuler</button>
              <button class="btn-danger" @click="deleteImage" :disabled="deleting">
                {{ deleting ? 'Suppression...' : 'Supprimer' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../../services/api'

const images = ref([])
const categoryStats = ref([])
const loading = ref(true)
const scanning = ref(false)
const uploading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const searchQuery = ref('')
const filterCategory = ref('')
const uploadingImage = ref(false)
const editingImage = ref(null)
const deletingImage = ref(null)
const selectedFile = ref(null)

const uploadForm = ref({
  title: '',
  description: '',
  category: '',
  featured: false
})

const editForm = ref({
  title: '',
  description: '',
  category: '',
  featured: false
})

const categories = [
  { label: 'Baby Shower', value: 'baby_shower' },
  { label: 'Baptême', value: 'bapteme' },
  { label: 'Anniversaire', value: 'birthday' },
  { label: 'Félicitations', value: 'congrats' },
  { label: 'Love Room', value: 'loveroom' },
  { label: 'Demande en Mariage', value: 'marryme' },
  { label: 'Mariage', value: 'wedding' }
]

const filteredImages = computed(() => {
  let filtered = images.value

  // Filter by category
  if (filterCategory.value) {
    filtered = filtered.filter(img => img.category === filterCategory.value)
  }

  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(img => 
      img.title?.toLowerCase().includes(query) ||
      img.description?.toLowerCase().includes(query)
    )
  }

  return filtered
})

function getCategoryLabel(categoryValue) {
  const cat = categories.find(c => c.value === categoryValue)
  return cat ? cat.label : categoryValue
}

async function fetchImages() {
  loading.value = true
  try {
    const response = await api.get('/public/gallery')
    images.value = response.data
  } catch (error) {
    console.error('Failed to fetch images:', error)
    alert('Erreur lors du chargement des images')
  } finally {
    loading.value = false
  }
}

async function fetchCategoryStats() {
  try {
    const response = await api.get('/admin/gallery/categories')
    categoryStats.value = response.data
  } catch (error) {
    console.error('Failed to fetch category stats:', error)
  }
}

async function scanStorageFolder() {
  if (!confirm('Scanner le dossier storage et importer les nouvelles images ?')) {
    return
  }

  scanning.value = true
  try {
    const response = await api.post('/admin/gallery/scan')
    alert(`Scan terminé!\nImportées: ${response.data.imported}\nIgnorées: ${response.data.skipped}`)
    await fetchImages()
    await fetchCategoryStats()
  } catch (error) {
    console.error('Failed to scan storage:', error)
    alert('Erreur lors du scan du dossier storage')
  } finally {
    scanning.value = false
  }
}

function openUploadModal() {
  uploadingImage.value = true
  selectedFile.value = null
  uploadForm.value = {
    title: '',
    description: '',
    category: '',
    featured: false
  }
}

function closeUploadModal() {
  uploadingImage.value = false
  selectedFile.value = null
}

function handleFileChange(event) {
  selectedFile.value = event.target.files[0]
}

async function uploadImage() {
  if (!selectedFile.value) {
    alert('Veuillez sélectionner une image')
    return
  }
  if (!uploadForm.value.category) {
    alert('Veuillez sélectionner une catégorie')
    return
  }

  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('image', selectedFile.value)
    formData.append('title', uploadForm.value.title)
    formData.append('description', uploadForm.value.description)
    formData.append('category', uploadForm.value.category)
    formData.append('featured', uploadForm.value.featured)

    const response = await api.post('/admin/gallery', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

    images.value.unshift(response.data)
    await fetchCategoryStats()
    closeUploadModal()
    alert('Image uploadée avec succès!')
  } catch (error) {
    console.error('Failed to upload image:', error)
    alert('Erreur lors de l\'upload de l\'image')
  } finally {
    uploading.value = false
  }
}

function editImage(image) {
  editingImage.value = image
  editForm.value = {
    title: image.title,
    description: image.description || '',
    category: image.category,
    featured: image.featured || false
  }
}

function closeEditModal() {
  editingImage.value = null
  editForm.value = {
    title: '',
    description: '',
    category: '',
    featured: false
  }
}

async function saveImage() {
  saving.value = true
  try {
    await api.put(`/admin/gallery/${editingImage.value.id}`, editForm.value)
    
    // Update local image
    const index = images.value.findIndex(img => img.id === editingImage.value.id)
    if (index !== -1) {
      images.value[index] = { ...images.value[index], ...editForm.value }
    }
    
    closeEditModal()
    alert('Image mise à jour avec succès!')
  } catch (error) {
    console.error('Failed to update image:', error)
    alert('Erreur lors de la mise à jour de l\'image')
  } finally {
    saving.value = false
  }
}

function confirmDelete(image) {
  deletingImage.value = image
}

async function deleteImage() {
  deleting.value = true
  try {
    await api.delete(`/admin/gallery/${deletingImage.value.id}`)
    
    // Remove from local list
    images.value = images.value.filter(img => img.id !== deletingImage.value.id)
    
    deletingImage.value = null
    await fetchCategoryStats()
    alert('Image supprimée avec succès!')
  } catch (error) {
    console.error('Failed to delete image:', error)
    alert('Erreur lors de la suppression de l\'image')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchImages()
  fetchCategoryStats()
})
</script>

<style scoped>
.gallery-manager {
  padding: var(--spacing-3xl);
  max-width: 1400px;
  margin: 0 auto;
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-3xl);
}

.manager-header h1 {
  font-size: var(--font-size-3xl);
  color: var(--color-charcoal);
}

.header-actions {
  display: flex;
  gap: var(--spacing-md);
}

.scan-btn,
.upload-btn {
  padding: var(--spacing-md) var(--spacing-xl);
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: all var(--transition-base);
  color: white;
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.scan-btn {
  background: var(--color-gold);
}

.upload-btn {
  background: var(--color-black);
}

.scan-btn:hover:not(:disabled) {
  background: var(--color-gold-dark, #c4a962);
  transform: translateY(-2px);
}

.upload-btn:hover {
  background: #333;
  transform: translateY(-2px);
}

.scan-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-3xl);
}

.stat-card {
  background: white;
  padding: var(--spacing-xl);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  text-align: center;
}

.stat-card h3 {
  font-size: var(--font-size-sm);
  color: var(--color-gray);
  margin-bottom: var(--spacing-sm);
  font-weight: var(--font-weight-medium);
}

.stat-card .count {
  font-size: var(--font-size-4xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-gold);
  margin: 0;
}

.stat-card .label {
  font-size: var(--font-size-xs);
  color: var(--color-gray);
  margin: 0;
}

/* Controls */
.controls {
  display: flex;
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-2xl);
}

.search-box {
  flex: 1;
}

.search-input {
  width: 100%;
  padding: var(--spacing-md);
  border: 1px solid #ddd;
  border-radius: var(--radius-md);
  font-size: var(--font-size-base);
}

.category-filter {
  padding: var(--spacing-md);
  border: 1px solid #ddd;
  border-radius: var(--radius-md);
  font-size: var(--font-size-base);
  min-width: 200px;
}

/* Images Grid */
.images-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--spacing-xl);
}

.image-card {
  background: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-md);
  transition: transform var(--transition-base);
}

.image-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-xl);
}

.image-preview {
  position: relative;
  aspect-ratio: 4/3;
  overflow: hidden;
}

.image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-md);
  opacity: 0;
  transition: opacity var(--transition-base);
}

.image-card:hover .image-overlay {
  opacity: 1;
}

.edit-btn,
.delete-btn {
  padding: var(--spacing-sm) var(--spacing-md);
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  cursor: pointer;
  transition: all var(--transition-base);
}

.edit-btn {
  background: var(--color-gold);
  color: white;
}

.edit-btn:hover {
  background: var(--color-gold-dark, #c4a962);
}

.delete-btn {
  background: #dc3545;
  color: white;
}

.delete-btn:hover {
  background: #c82333;
}

.image-info {
  padding: var(--spacing-lg);
}

.image-info h4 {
  font-size: var(--font-size-lg);
  margin-bottom: var(--spacing-sm);
  color: var(--color-charcoal);
}

.category-badge {
  display: inline-block;
  padding: var(--spacing-xs) var(--spacing-sm);
  background: var(--color-gold);
  color: white;
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  margin-bottom: var(--spacing-sm);
}

.description {
  font-size: var(--font-size-sm);
  color: var(--color-gray);
  margin-bottom: var(--spacing-sm);
}

.meta {
  display: flex;
  gap: var(--spacing-sm);
}

.badge {
  padding: var(--spacing-xs) var(--spacing-sm);
  background: #e9ecef;
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
}

.badge.featured {
  background: #ffc107;
  color: #000;
}

/* Loading & Empty States */
.loading-state,
.empty-state {
  text-align: center;
  padding: var(--spacing-5xl) 0;
  color: var(--color-gray);
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-top-color: var(--color-gold);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto var(--spacing-md);
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--spacing-xl);
}

.modal-content {
  background: white;
  border-radius: var(--radius-xl);
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-content.small {
  max-width: 400px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-xl);
  border-bottom: 1px solid #e9ecef;
}

.modal-header h2 {
  margin: 0;
  font-size: var(--font-size-2xl);
}

.close-btn {
  background: none;
  border: none;
  font-size: var(--font-size-3xl);
  cursor: pointer;
  color: var(--color-gray);
  line-height: 1;
  padding: 0;
  width: 30px;
  height: 30px;
}

.close-btn:hover {
  color: var(--color-charcoal);
}

.modal-body {
  padding: var(--spacing-xl);
}

.form-group {
  margin-bottom: var(--spacing-lg);
}

.form-group label {
  display: block;
  margin-bottom: var(--spacing-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-charcoal);
}

.form-input,
.form-textarea,
.form-select {
  width: 100%;
  padding: var(--spacing-md);
  border: 1px solid #ddd;
  border-radius: var(--radius-md);
  font-size: var(--font-size-base);
  font-family: inherit;
}

.form-textarea {
  resize: vertical;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: auto;
  cursor: pointer;
}

.warning {
  color: #dc3545;
  font-weight: var(--font-weight-medium);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-md);
  padding: var(--spacing-xl);
  border-top: 1px solid #e9ecef;
}

.btn-primary,
.btn-secondary,
.btn-danger {
  padding: var(--spacing-md) var(--spacing-xl);
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: all var(--transition-base);
}

.btn-primary {
  background: var(--color-gold);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-gold-dark, #c4a962);
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background: #5a6268;
}

.btn-danger {
  background: #dc3545;
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background: #c82333;
}

.btn-primary:disabled,
.btn-danger:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Modal Transitions */
.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--transition-base);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

@media (max-width: 768px) {
  .gallery-manager {
    padding: var(--spacing-xl);
  }

  .manager-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-md);
  }

  .scan-btn {
    width: 100%;
  }

  .controls {
    flex-direction: column;
  }

  .category-filter {
    width: 100%;
  }

  .stats-grid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  }

  .images-grid {
    grid-template-columns: 1fr;
  }
}
</style>
