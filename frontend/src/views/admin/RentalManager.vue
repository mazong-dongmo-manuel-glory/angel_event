<template>
  <div class="rental-manager">
    <div class="manager-header">
      <h1>Gestion des Locations</h1>
      <button class="add-btn" @click="openAddModal">
        ➕ Ajouter un article
      </button>
    </div>

    <!-- Controls -->
    <div class="controls">
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Rechercher par titre..."
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

    <!-- Items Grid -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement des articles...</p>
    </div>

    <div v-else class="items-grid">
      <div v-for="item in filteredItems" :key="item.id" class="item-card">
        <div class="item-preview">
          <img :src="item.image_url" :alt="item.title" />
          <div class="item-overlay">
            <button class="edit-btn" @click="editItem(item)">✏️ Éditer</button>
            <button class="delete-btn" @click="confirmDelete(item)">🗑️ Supprimer</button>
          </div>
        </div>
        <div class="item-info">
          <div class="info-header">
            <h4>{{ item.title }}</h4>
            <span class="price">{{ formatPrice(item.price) }}</span>
          </div>
          <p class="category-badge">{{ getCategoryLabel(item.category) }}</p>
          <p v-if="item.description" class="description">{{ item.description }}</p>
          <div class="meta">
            <span v-if="item.featured" class="badge featured">⭐ Vedette</span>
            <span v-if="!item.available" class="badge unavailable">❌ Indisponible</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showModal" class="modal-overlay" @click="closeModal">
          <div class="modal-content" @click.stop>
            <div class="modal-header">
              <h2>{{ editingItem ? 'Éditer l\'article' : 'Nouvel article' }}</h2>
              <button class="close-btn" @click="closeModal">×</button>
            </div>
            <div class="modal-body">
              <form @submit.prevent="saveItem" class="item-form">
                <div class="form-group" v-if="!editingItem || changeImage">
                  <label>Image *</label>
                  <input type="file" ref="fileInput" @change="handleFileChange" accept="image/*" class="form-input" :required="!editingItem" />
                </div>
                <div class="form-group" v-if="editingItem && !changeImage">
                  <label>Image actuelle</label>
                  <div class="current-image">
                    <img :src="editingItem.image_url" height="100" />
                    <button type="button" class="btn-text" @click="changeImage = true">Changer l'image</button>
                  </div>
                </div>

                <div class="form-row">
                  <div class="form-group">
                    <label>Titre *</label>
                    <input v-model="form.title" type="text" class="form-input" required />
                  </div>
                  <div class="form-group">
                    <label>Prix ($) *</label>
                    <input v-model="form.price" type="number" step="0.01" class="form-input" required />
                  </div>
                </div>

                <div class="form-group">
                  <label>Catégorie *</label>
                  <select v-model="form.category" class="form-select" required>
                    <option value="" disabled>Choisir une catégorie</option>
                    <option v-for="cat in categories" :key="cat.value" :value="cat.value">
                      {{ cat.label }}
                    </option>
                  </select>
                </div>

                <div class="form-group">
                  <label>Description</label>
                  <textarea v-model="form.description" class="form-textarea" rows="3"></textarea>
                </div>

                <div class="form-checkboxes">
                  <label class="checkbox-label">
                    <input v-model="form.featured" type="checkbox" />
                    <span>Mettre en vedette</span>
                  </label>
                  <label class="checkbox-label">
                    <input v-model="form.available" type="checkbox" />
                    <span>Disponible</span>
                  </label>
                </div>

                <div class="modal-footer">
                  <button type="button" class="btn-secondary" @click="closeModal">Annuler</button>
                  <button type="submit" class="btn-primary" :disabled="saving">
                    {{ saving ? 'Enregistrement...' : 'Enregistrer' }}
                  </button>
                </div>
              </form>
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

const items = ref([])
const loading = ref(true)
const saving = ref(false)
const searchQuery = ref('')
const filterCategory = ref('')
const showModal = ref(false)
const editingItem = ref(null)
const changeImage = ref(false)
const selectedFile = ref(null)

const form = ref({
  title: '',
  price: '',
  category: '',
  description: '',
  featured: false,
  available: true
})

const categories = [
  { label: 'Centre de table', value: 'centerpiece' },
  { label: 'Backdrop / Panneau', value: 'backdrop' },
  { label: 'Fleur', value: 'flower' },
  { label: 'Autre article', value: 'other' }
]

const filteredItems = computed(() => {
  let result = items.value
  if (filterCategory.value) {
    result = result.filter(i => i.category === filterCategory.value)
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(i => i.title.toLowerCase().includes(q))
  }
  return result
})

function formatPrice(price) {
  return new Intl.NumberFormat('fr-CA', { style: 'currency', currency: 'CAD' }).format(price)
}

function getCategoryLabel(val) {
  const cat = categories.find(c => c.value === val)
  return cat ? cat.label : val
}

async function fetchItems() {
  loading.value = true
  try {
    const res = await api.get('/admin/rentals')
    items.value = res.data
  } catch (err) {
    console.error(err)
    alert('Erreur chargement articles')
  } finally {
    loading.value = false
  }
}

function openAddModal() {
  editingItem.value = null
  resetForm()
  showModal.value = true
}

function editItem(item) {
  editingItem.value = item
  form.value = { ...item }
  changeImage.value = false
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingItem.value = null
  selectedFile.value = null
}

function resetForm() {
  form.value = {
    title: '',
    price: '',
    category: '',
    description: '',
    featured: false,
    available: true
  }
  selectedFile.value = null
  changeImage.value = false
}

function handleFileChange(e) {
  selectedFile.value = e.target.files[0]
}

async function saveItem() {
  if (!editingItem.value && !selectedFile.value) {
    alert('Image requise')
    return
  }

  saving.value = true
  const formData = new FormData()
  formData.append('title', form.value.title)
  formData.append('price', form.value.price)
  formData.append('category', form.value.category)
  formData.append('description', form.value.description)
  formData.append('featured', form.value.featured)
  formData.append('available', form.value.available)
  
  if (selectedFile.value) {
    formData.append('image', selectedFile.value)
  }

  try {
    if (editingItem.value) {
      const res = await api.put(`/admin/rentals/${editingItem.value.id}`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
      const idx = items.value.findIndex(i => i.id === editingItem.value.id)
      items.value[idx] = res.data
    } else {
      const res = await api.post('/admin/rentals', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
      items.value.unshift(res.data)
    }
    closeModal()
  } catch (err) {
    console.error(err)
    alert('Erreur sauvegarde')
  } finally {
    saving.value = false
  }
}

async function confirmDelete(item) {
  if (!confirm('Voulez-vous vraiment supprimer cet article ?')) return
  try {
    await api.delete(`/admin/rentals/${item.id}`)
    items.value = items.value.filter(i => i.id !== item.id)
  } catch (err) {
    console.error(err)
    alert('Erreur suppression')
  }
}

onMounted(fetchItems)
</script>

<style scoped>
.rental-manager {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.add-btn {
  background: var(--color-gold);
  color: white;
  border: none;
  padding: 0.8rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  font-weight: bold;
}

.controls {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.search-input, .category-filter {
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
}

.items-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
}

.item-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.item-preview {
  position: relative;
  aspect-ratio: 4/3;
}

.item-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item-overlay {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.5rem;
  opacity: 0;
  transition: opacity 0.2s;
}

.item-card:hover .item-overlay {
  opacity: 1;
}

.item-info {
  padding: 1rem;
}

.info-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.price {
  font-weight: bold;
  color: var(--color-gold);
}

.category-badge {
  background: #f0f0f0;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-size: 0.8rem;
  display: inline-block;
}

/* Modal form styles */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-row {
  display: flex;
  gap: 1rem;
}

.form-input, .form-select, .form-textarea {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  margin-top: 0.3rem;
}

.btn-primary, .btn-secondary, .edit-btn, .delete-btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.btn-primary { background: var(--color-gold); color: white; }
.btn-secondary { background: #eee; }
.edit-btn { background: white; }
.delete-btn { background: #ff4444; color: white; }
</style>
