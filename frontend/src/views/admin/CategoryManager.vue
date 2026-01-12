<template>
  <div class="category-manager">
    <div class="manager-header">
      <h1>Gestion des Catégories</h1>
      <button class="add-btn" @click="openAddModal">
        ➕ Nouvelle Catégorie
      </button>
    </div>

    <!-- Categories List -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement...</p>
    </div>

    <div v-else class="categories-list">
      <table class="data-table">
        <thead>
          <tr>
            <th>Nom</th>
            <th>Type</th>
            <th>Slug</th>
            <th>Description</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="cat in categories" :key="cat.id">
            <td><strong>{{ cat.name }}</strong></td>
            <td>
              <span class="badge" :class="cat.type">{{ cat.type }}</span>
            </td>
            <td><code>{{ cat.slug }}</code></td>
            <td>{{ cat.description }}</td>
            <td class="actions">
              <button class="edit-btn" @click="editCategory(cat)">✏️</button>
              <button class="delete-btn" @click="confirmDelete(cat)">🗑️</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showModal" class="modal-overlay" @click="closeModal">
          <div class="modal-content" @click.stop>
            <div class="modal-header">
              <h2>{{ editingCategory ? 'Éditer la catégorie' : 'Nouvelle catégorie' }}</h2>
              <button class="close-btn" @click="closeModal">×</button>
            </div>
            <div class="modal-body">
              <form @submit.prevent="saveCategory">
                <div class="form-group">
                  <label>Nom *</label>
                  <input v-model="form.name" type="text" class="form-input" required placeholder="ex: Animation" />
                </div>
                
                <div class="form-group">
                  <label>Slug (Identifiant unique) *</label>
                  <input v-model="form.slug" type="text" class="form-input" required placeholder="ex: animation" />
                  <small>Utilisé dans l'URL et le code. Sans espaces, minuscules.</small>
                </div>

                <div class="form-group">
                  <label>Type</label>
                  <select v-model="form.type" class="form-select">
                    <option value="rental">Location</option>
                    <option value="gallery">Galerie</option>
                  </select>
                </div>

                <div class="form-group">
                  <label>Description</label>
                  <textarea v-model="form.description" class="form-textarea" rows="3"></textarea>
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
import { ref, onMounted } from 'vue'
import api from '../../services/api'

const categories = ref([])
const loading = ref(true)
const saving = ref(false)
const showModal = ref(false)
const editingCategory = ref(null)

const form = ref({
  name: '',
  slug: '',
  description: '',
  type: 'rental'
})

async function fetchCategories() {
  loading.value = true
  try {
    const res = await api.get('/admin/categories')
    categories.value = res.data
  } catch (err) {
    console.error(err)
    alert('Erreur chargement catégories')
  } finally {
    loading.value = false
  }
}

function openAddModal() {
  editingCategory.value = null
  resetForm()
  showModal.value = true
}

function editCategory(cat) {
  editingCategory.value = cat
  form.value = { ...cat }
  showModal.value = true
}

function resetForm() {
  form.value = {
    name: '',
    slug: '',
    description: '',
    type: 'rental'
  }
}

function closeModal() {
  showModal.value = false
  editingCategory.value = null
}

async function saveCategory() {
  saving.value = true
  try {
    if (editingCategory.value) {
      const res = await api.put(`/admin/categories/${editingCategory.value.id}`, form.value)
      const idx = categories.value.findIndex(c => c.id === editingCategory.value.id)
      categories.value[idx] = res.data
    } else {
      const res = await api.post('/admin/categories', form.value)
      categories.value.push(res.data)
    }
    closeModal()
  } catch (err) {
    console.error(err)
    alert('Erreur sauvegarde')
  } finally {
    saving.value = false
  }
}

async function confirmDelete(cat) {
  if (!confirm(`Supprimer la catégorie "${cat.name}" ?`)) return
  try {
    await api.delete(`/admin/categories/${cat.id}`)
    categories.value = categories.value.filter(c => c.id !== cat.id)
  } catch (err) {
    console.error(err)
    alert('Erreur suppression')
  }
}

onMounted(fetchCategories)
</script>

<style scoped>
.category-manager {
  max-width: 1000px;
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

.data-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.data-table th, .data-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.data-table th {
  background: #f8f9fa;
  font-weight: 600;
}

.badge {
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-size: 0.8rem;
  text-transform: uppercase;
  font-weight: bold;
}

.badge.rental { background: #e3f2fd; color: #1976d2; }
.badge.gallery { background: #f3e5f5; color: #7b1fa2; }

.actions {
  display: flex;
  gap: 0.5rem;
}

.edit-btn, .delete-btn {
  padding: 0.4rem 0.8rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.edit-btn { background: #eee; }
.delete-btn { background: #fee; color: #d32f2f; }

/* Modal styles - similar to RentalManager */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5);
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
}

.modal-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-input, .form-select, .form-textarea {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  margin-top: 0.3rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

.btn-primary { background: var(--color-gold); color: white; padding: 0.6rem 1.2rem; border-radius: 6px; border:none; cursor: pointer; }
.btn-secondary { background: #eee; padding: 0.6rem 1.2rem; border-radius: 6px; border:none; cursor: pointer; }
</style>
