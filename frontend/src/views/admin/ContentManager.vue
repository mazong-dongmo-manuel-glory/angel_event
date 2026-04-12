<template>
  <div class="content-manager">
    <div class="manager-header">
      <div>
        <h1>Gestion du Contenu</h1>
        <p class="subtitle">Modifiez les textes, logos, photos et visuels des pages vitrines depuis ce module.</p>
      </div>
      <button @click="openCreateModal" class="btn-add">
        + Ajouter un contenu
      </button>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📝</div>
        <div class="stat-content">
          <h3>{{ stats.total }}</h3>
          <p>Total contenus</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">🖼️</div>
        <div class="stat-content">
          <h3>{{ stats.images }}</h3>
          <p>Visuels</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">🌐</div>
        <div class="stat-content">
          <h3>{{ stats.languages }}</h3>
          <p>Langues</p>
        </div>
      </div>
    </div>

    <div class="controls">
      <div class="filter-group">
        <label>Section</label>
        <select v-model="filterSection" class="filter-select">
          <option value="">Toutes</option>
          <option v-for="section in availableSections" :key="section" :value="section">
            {{ getSectionLabel(section) }}
          </option>
        </select>
      </div>

      <div class="filter-group">
        <label>Langue</label>
        <select v-model="filterLanguage" class="filter-select">
          <option value="">Toutes</option>
          <option value="fr">FR</option>
          <option value="en">EN</option>
        </select>
      </div>

      <div class="filter-group">
        <label>Type</label>
        <select v-model="filterType" class="filter-select">
          <option value="">Tous</option>
          <option value="text">Texte</option>
          <option value="image">Image</option>
        </select>
      </div>

      <input
        v-model="searchQuery"
        type="text"
        class="search-input"
        placeholder="Rechercher par clé, valeur ou section..."
      />
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement du contenu...</p>
    </div>

    <div v-else-if="groupedContent.length === 0" class="empty-state">
      <p>Aucun contenu trouvé.</p>
    </div>

    <div v-else class="content-sections">
      <section v-for="group in groupedContent" :key="group.name" class="section-group">
        <div class="section-header">
          <h2>{{ getSectionLabel(group.name) }}</h2>
          <span>{{ group.items.length }} entrée(s)</span>
        </div>

        <div class="content-grid">
          <article v-for="item in group.items" :key="item.id" class="content-card">
            <div class="card-header">
              <code>{{ item.key }}</code>
              <div class="badges">
                <span class="badge language">{{ item.language.toUpperCase() }}</span>
                <span class="badge" :class="item.type === 'image' ? 'image' : 'text'">
                  {{ item.type || 'text' }}
                </span>
              </div>
            </div>

            <div class="card-body">
              <img v-if="item.type === 'image'" :src="item.value" :alt="item.key" class="content-preview-image" />
              <p v-else class="content-value">{{ item.value }}</p>
              <div class="meta">
                <span>{{ formatDate(item.updated_at) }}</span>
              </div>
            </div>

            <div class="card-actions">
              <button class="btn-secondary" @click="editContent(item)">Modifier</button>
              <button class="btn-danger" @click="confirmDelete(item)">Supprimer</button>
            </div>
          </article>
        </div>
      </section>
    </div>

    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click="closeModal">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>{{ editingContent ? 'Modifier le contenu' : 'Ajouter un contenu' }}</h2>
            <button @click="closeModal" class="close-btn">×</button>
          </div>

          <form class="edit-form" @submit.prevent="saveContent">
            <div class="form-row">
              <div class="form-group">
                <label>Clé</label>
                <input v-model="editForm.key" type="text" required :disabled="!!editingContent" placeholder="about_story_image" />
              </div>

              <div class="form-group">
                <label>Type</label>
                <select v-model="editForm.type" required>
                  <option value="text">Texte</option>
                  <option value="image">Image</option>
                </select>
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>Section</label>
                <select v-model="editForm.section" required>
                  <option value="global">Global</option>
                  <option value="home">Accueil</option>
                  <option value="services">Services</option>
                  <option value="about">À propos</option>
                  <option value="contact">Contact</option>
                  <option value="gallery">Galerie</option>
                  <option value="testimonials">Témoignages</option>
                  <option value="footer">Footer</option>
                </select>
              </div>

              <div class="form-group">
                <label>Langue</label>
                <select v-model="editForm.language" required>
                  <option value="fr">Français</option>
                  <option value="en">English</option>
                </select>
              </div>
            </div>

            <div v-if="editForm.type === 'image'" class="form-group">
              <label>URL de l'image</label>
              <input
                v-model="editForm.value"
                type="text"
                placeholder="/uploads/content/mon-image.jpg ou https://..."
              />
            </div>

            <div v-if="editForm.type === 'image'" class="form-group">
              <label>Uploader une image</label>
              <input type="file" accept=".jpg,.jpeg,.png,.gif,.webp,.svg" @change="handleFileChange" />
              <small>Si un fichier est sélectionné, il remplacera l'URL ci-dessus.</small>
            </div>

            <div v-if="imagePreview" class="image-preview-box">
              <img :src="imagePreview" alt="Preview" />
            </div>

            <div v-if="editForm.type !== 'image'" class="form-group">
              <label>Valeur</label>
              <textarea v-model="editForm.value" rows="6" required placeholder="Entrez votre contenu ici..."></textarea>
            </div>

            <div class="form-actions">
              <button type="button" class="btn-secondary" @click="closeModal">Annuler</button>
              <button type="submit" class="btn-primary" :disabled="saving">
                {{ saving ? 'Enregistrement...' : 'Enregistrer' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="deletingContent" class="modal-overlay" @click="deletingContent = null">
        <div class="modal-content modal-small" @click.stop>
          <div class="modal-header">
            <h2>Supprimer ce contenu ?</h2>
            <button @click="deletingContent = null" class="close-btn">×</button>
          </div>

          <div class="modal-body">
            <p>La clé <strong>{{ deletingContent.key }}</strong> sera supprimée définitivement.</p>

            <div class="form-actions">
              <button class="btn-secondary" @click="deletingContent = null">Annuler</button>
              <button class="btn-danger" @click="deleteContent">Supprimer</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import api from '../../services/api'

const content = ref([])
const loading = ref(false)
const saving = ref(false)
const filterSection = ref('')
const filterLanguage = ref('')
const filterType = ref('')
const searchQuery = ref('')
const showModal = ref(false)
const editingContent = ref(null)
const deletingContent = ref(null)
const imageFile = ref(null)
const imagePreview = ref('')

const editForm = ref({
  key: '',
  value: '',
  section: 'home',
  language: 'fr',
  type: 'text',
})

const sectionLabels = {
  global: 'Global',
  home: 'Accueil',
  services: 'Services',
  about: 'À propos',
  contact: 'Contact',
  gallery: 'Galerie',
  testimonials: 'Témoignages',
  footer: 'Footer',
}

const stats = computed(() => ({
  total: content.value.length,
  images: content.value.filter((item) => item.type === 'image').length,
  languages: new Set(content.value.map((item) => item.language)).size,
}))

const availableSections = computed(() => {
  return [...new Set(content.value.map((item) => item.section))].filter(Boolean).sort()
})

const filteredContent = computed(() => {
  return content.value.filter((item) => {
    const matchesSection = !filterSection.value || item.section === filterSection.value
    const matchesLanguage = !filterLanguage.value || item.language === filterLanguage.value
    const matchesType = !filterType.value || (item.type || 'text') === filterType.value
    const query = searchQuery.value.trim().toLowerCase()
    const matchesSearch =
      !query ||
      item.key.toLowerCase().includes(query) ||
      item.value.toLowerCase().includes(query) ||
      item.section.toLowerCase().includes(query)

    return matchesSection && matchesLanguage && matchesType && matchesSearch
  })
})

const groupedContent = computed(() => {
  const groups = {}

  for (const item of filteredContent.value) {
    if (!groups[item.section]) {
      groups[item.section] = {
        name: item.section,
        items: [],
      }
    }

    groups[item.section].items.push(item)
  }

  return Object.values(groups)
    .map((group) => ({
      ...group,
      items: [...group.items].sort((a, b) => a.key.localeCompare(b.key)),
    }))
    .sort((a, b) => a.name.localeCompare(b.name))
})

function getSectionLabel(section) {
  return sectionLabels[section] || section
}

function resetForm() {
  editForm.value = {
    key: '',
    value: '',
    section: 'home',
    language: 'fr',
    type: 'text',
  }
  imageFile.value = null
  imagePreview.value = ''
}

async function fetchContent() {
  loading.value = true

  try {
    const response = await api.get('/admin/content')
    content.value = response.data || []
  } catch (error) {
    console.error('Failed to fetch content:', error)
    alert('Erreur lors du chargement du contenu.')
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingContent.value = null
  resetForm()
  showModal.value = true
}

function editContent(item) {
  editingContent.value = item
  editForm.value = {
    key: item.key,
    value: item.value,
    section: item.section,
    language: item.language,
    type: item.type || 'text',
  }
  imageFile.value = null
  imagePreview.value = item.type === 'image' ? item.value : ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingContent.value = null
  resetForm()
}

function handleFileChange(event) {
  const file = event.target.files?.[0]
  imageFile.value = file || null
  imagePreview.value = file ? URL.createObjectURL(file) : editForm.value.value
}

async function saveContent() {
  saving.value = true

  try {
    if (editForm.value.type === 'image' && imageFile.value) {
      const formData = new FormData()
      formData.append('key', editForm.value.key)
      formData.append('section', editForm.value.section)
      formData.append('language', editForm.value.language)
      formData.append('type', 'image')
      formData.append('image', imageFile.value)

      await api.post('/admin/content/upload', formData)
    } else {
      await api.post('/admin/content', editForm.value)
    }

    await fetchContent()
    closeModal()
  } catch (error) {
    console.error('Failed to save content:', error)
    alert(error.response?.data?.error || 'Erreur lors de la sauvegarde.')
  } finally {
    saving.value = false
  }
}

function confirmDelete(item) {
  deletingContent.value = item
}

async function deleteContent() {
  try {
    await api.delete(`/admin/content/${deletingContent.value.id}`)
    deletingContent.value = null
    await fetchContent()
  } catch (error) {
    console.error('Failed to delete content:', error)
    alert(error.response?.data?.error || 'Erreur lors de la suppression.')
  }
}

function formatDate(dateStr) {
  if (!dateStr) {
    return '-'
  }

  return new Date(dateStr).toLocaleDateString('fr-FR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

onMounted(fetchContent)
</script>

<style scoped>
.content-manager {
  padding: var(--spacing-3xl);
  max-width: 1400px;
  margin: 0 auto;
}

.manager-header {
  display: flex;
  justify-content: space-between;
  gap: var(--spacing-xl);
  align-items: flex-start;
  margin-bottom: var(--spacing-3xl);
}

.manager-header h1 {
  color: var(--color-gold);
  font-size: var(--font-size-4xl);
  margin-bottom: var(--spacing-sm);
}

.subtitle {
  color: var(--color-gray);
  font-size: var(--font-size-lg);
  max-width: 760px;
}

.btn-add,
.btn-primary,
.btn-secondary,
.btn-danger {
  border: none;
  border-radius: var(--radius-md);
  padding: var(--spacing-md) var(--spacing-xl);
  cursor: pointer;
  font-weight: var(--font-weight-semibold);
  transition: all var(--transition-base);
}

.btn-add,
.btn-primary {
  background: var(--color-gold);
  color: white;
}

.btn-add:hover,
.btn-primary:hover {
  transform: translateY(-2px);
  background: #c59f2c;
}

.btn-secondary {
  background: var(--color-gray-lighter);
  color: var(--color-text);
}

.btn-danger {
  background: rgba(193, 41, 46, 0.12);
  color: var(--color-error);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: var(--spacing-xl);
  margin-bottom: var(--spacing-3xl);
}

.stat-card {
  background: white;
  border-radius: var(--radius-lg);
  padding: var(--spacing-xl);
  box-shadow: var(--shadow-md);
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
}

.stat-icon {
  font-size: var(--font-size-4xl);
}

.stat-content h3 {
  font-size: var(--font-size-2xl);
  color: var(--color-gold);
}

.stat-content p {
  color: var(--color-gray);
}

.controls {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-3xl);
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.filter-select,
.search-input,
.edit-form input,
.edit-form select,
.edit-form textarea {
  width: 100%;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--spacing-md);
  font-family: var(--font-sans);
}

.content-sections {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3xl);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.section-header h2 {
  color: var(--color-gold);
}

.section-header span {
  color: var(--color-gray);
}

.content-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--spacing-xl);
}

.content-card {
  background: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  gap: var(--spacing-md);
  align-items: center;
  padding: var(--spacing-lg);
  border-bottom: 1px solid var(--color-border);
}

.card-header code {
  color: var(--color-gold);
  font-weight: var(--font-weight-semibold);
}

.badges {
  display: flex;
  gap: var(--spacing-sm);
}

.badge {
  padding: 4px 10px;
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  text-transform: uppercase;
  font-weight: var(--font-weight-semibold);
}

.badge.language {
  background: rgba(0, 0, 0, 0.06);
}

.badge.text {
  background: rgba(212, 175, 55, 0.14);
  color: var(--color-gold);
}

.badge.image {
  background: rgba(79, 70, 229, 0.12);
  color: #4338ca;
}

.card-body {
  padding: var(--spacing-lg);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  flex: 1;
}

.content-preview-image {
  width: 100%;
  height: 220px;
  object-fit: cover;
  border-radius: var(--radius-md);
  background: var(--color-gray-lighter);
}

.content-value {
  color: var(--color-text);
  line-height: 1.7;
  white-space: pre-wrap;
}

.meta {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
}

.card-actions,
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-md);
}

.card-actions {
  padding: 0 var(--spacing-lg) var(--spacing-lg);
}

.loading-state,
.empty-state {
  text-align: center;
  padding: var(--spacing-4xl);
  color: var(--color-gray);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--color-gray-lighter);
  border-top-color: var(--color-gold);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto var(--spacing-md);
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.55);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--spacing-xl);
}

.modal-content {
  width: min(760px, 100%);
  background: white;
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-2xl);
  overflow: hidden;
}

.modal-small {
  width: min(420px, 100%);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-xl);
  border-bottom: 1px solid var(--color-border);
}

.modal-header h2 {
  color: var(--color-gold);
}

.close-btn {
  border: none;
  background: transparent;
  font-size: 2rem;
  line-height: 1;
  cursor: pointer;
}

.edit-form,
.modal-body {
  padding: var(--spacing-xl);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--spacing-lg);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-lg);
}

.image-preview-box {
  margin-bottom: var(--spacing-lg);
}

.image-preview-box img {
  width: 100%;
  max-height: 280px;
  object-fit: cover;
  border-radius: var(--radius-lg);
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 900px) {
  .controls {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .content-manager {
    padding: var(--spacing-xl);
  }

  .manager-header,
  .section-header,
  .card-actions,
  .form-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .controls,
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
