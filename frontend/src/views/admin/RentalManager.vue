<template>
  <div class="rental-manager">
    <div class="manager-header">
      <h1>{{ t('admin.rental.title') }}</h1>
      <button class="add-btn" @click="openAddModal">
        ➕ {{ t('admin.rental.add_item') }}
      </button>
    </div>

    <!-- Controls -->
    <div class="controls">
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          type="text" 
          :placeholder="t('admin.rental.search_placeholder')"
          class="search-input"
        />
      </div>
      <select v-model="filterCategoryId" class="category-filter">
        <option value="">{{ t('admin.rental.all_categories') }}</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </option>
      </select>
    </div>

    <!-- Items Grid -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>{{ t('admin.rental.loading') }}</p>
    </div>

    <div v-else class="items-grid">
      <div v-for="item in filteredItems" :key="item.id" class="item-card">
        <div class="item-preview">
          <img :src="item.image_url" :alt="item.title" />
          <div class="item-overlay">
            <button class="edit-btn" @click="editItem(item)">✏️ {{ t('admin.rental.edit') }}</button>
            <button class="delete-btn" @click="confirmDelete(item)">🗑️ {{ t('admin.rental.delete') }}</button>
          </div>
        </div>
        <div class="item-info">
          <div class="info-header">
            <h4>{{ item.title }}</h4>
            <span class="price">{{ formatPrice(item.price) }}</span>
          </div>
          <p class="category-badge">{{ getCategoryLabel(item.category_id) }}</p>
          <p v-if="item.description" class="description">{{ item.description }}</p>
          <div class="meta">
            <span v-if="item.featured" class="badge featured">⭐ {{ t('admin.rental.featured') }}</span>
            <span v-if="!item.available" class="badge unavailable">❌ {{ t('admin.rental.unavailable') }}</span>
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
              <h2>{{ editingItem ? t('admin.rental.modal_title_edit') : t('admin.rental.modal_title_new') }}</h2>
              <button class="close-btn" @click="closeModal">×</button>
            </div>
            <div class="modal-body">
              <form @submit.prevent="saveItem" class="item-form">
                <div class="form-group" v-if="!editingItem || changeImage">
                  <label>{{ t('admin.rental.form.image') }} *</label>
                  <input type="file" ref="fileInput" @change="handleFileChange" accept="image/*" class="form-input" :required="!editingItem" />
                </div>
                <div class="form-group" v-if="editingItem && !changeImage">
                  <label>{{ t('admin.rental.form.current_image') }}</label>
                  <div class="current-image">
                    <img :src="editingItem.image_url" height="100" />
                    <button type="button" class="btn-text" @click="changeImage = true">{{ t('admin.rental.form.change_image') }}</button>
                  </div>
                </div>

                <div class="form-row">
                  <div class="form-group">
                    <label>{{ t('admin.rental.form.title') }} *</label>
                    <input v-model="form.title" type="text" class="form-input" required />
                  </div>
                  <div class="form-group">
                    <label>{{ t('admin.rental.form.price') }} *</label>
                    <input v-model="form.price" type="number" step="0.01" class="form-input" required />
                  </div>
                </div>

                <div class="form-group">
                  <label>{{ t('admin.rental.form.category') }} *</label>
                  <select v-model="form.category_id" class="form-select" required>
                    <option value="" disabled>{{ t('admin.rental.form.choose_category') }}</option>
                    <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                      {{ cat.name }}
                    </option>
                  </select>
                </div>

                <div class="form-group">
                  <label>{{ t('admin.rental.form.description') }}</label>
                  <textarea v-model="form.description" class="form-textarea" rows="3"></textarea>
                </div>

                <div class="form-checkboxes">
                  <label class="checkbox-label">
                    <input v-model="form.featured" type="checkbox" />
                    <span>{{ t('admin.rental.form.make_featured') }}</span>
                  </label>
                  <label class="checkbox-label">
                    <input v-model="form.available" type="checkbox" />
                    <span>{{ t('admin.rental.form.available') }}</span>
                  </label>
                </div>

                <div class="modal-footer">
                  <button type="button" class="btn-secondary" @click="closeModal">{{ t('admin.rental.cancel') }}</button>
                  <button type="submit" class="btn-primary" :disabled="saving">
                    {{ saving ? t('admin.rental.saving') : t('admin.rental.save') }}
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
import { useI18n } from 'vue-i18n'
import api from '../../services/api'

const { t } = useI18n()

const items = ref([])
const loading = ref(true)
const saving = ref(false)
const searchQuery = ref('')
const filterCategoryId = ref('')
const showModal = ref(false)
const editingItem = ref(null)
const changeImage = ref(false)
const selectedFile = ref(null)

const form = ref({
  title: '',
  price: '',
  category_id: '',
  description: '',
  featured: false,
  available: true
})

const categories = ref([])

async function fetchCategories() {
  try {
    const res = await api.get('/admin/categories')
    categories.value = res.data.filter(c => c.type === 'rental')
  } catch (err) {
    console.error(t('admin.rental.category_error'), err)
  }
}

const filteredItems = computed(() => {
  let result = items.value
  if (filterCategoryId.value) {
    result = result.filter(i => i.category_id === filterCategoryId.value)
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

function getCategoryLabel(catId) {
  const cat = categories.value.find(c => c.id === catId)
  return cat ? cat.name : 'Inconnu'
}

async function fetchItems() {
  loading.value = true
  try {
    const res = await api.get('/admin/rentals')
    items.value = res.data
  } catch (err) {
    console.error(err)
    alert(t('admin.rental.load_error'))
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
    category_id: '',
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
    alert(t('admin.rental.image_required'))
    return
  }

  saving.value = true
  const formData = new FormData()
  formData.append('title', form.value.title)
  formData.append('price', form.value.price)
  formData.append('category_id', form.value.category_id)
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
    alert(t('admin.rental.save_error'))
  } finally {
    saving.value = false
  }
}

async function confirmDelete(item) {
  if (!confirm(t('admin.rental.delete_confirm'))) return
  try {
    await api.delete(`/admin/rentals/${item.id}`)
    items.value = items.value.filter(i => i.id !== item.id)
  } catch (err) {
    console.error(err)
    alert(t('admin.rental.delete_error'))
  }
}

onMounted(() => {
  fetchCategories()
  fetchItems()
})
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
