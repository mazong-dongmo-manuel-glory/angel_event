<template>
  <div class="content-manager">
    <!-- Header -->
    <div class="manager-header">
      <div>
        <h1>Gestion du Contenu</h1>
        <p class="subtitle">Gérez le contenu du site web (textes, titres, descriptions)</p>
      </div>
      <button @click="showAddModal = true" class="btn-add">
        ✏️ Ajouter un contenu
      </button>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📝</div>
        <div class="stat-content">
          <h3>{{ stats.total }}</h3>
          <p>Total Contenus</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">🏠</div>
        <div class="stat-content">
          <h3>{{ stats.sections }}</h3>
          <p>Sections</p>
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

    <!-- Filters -->
    <div class="controls">
      <div class="filter-group">
        <label>Section:</label>
        <select v-model="filterSection" class="filter-select">
          <option value="">Toutes les sections</option>
          <option v-for="section in availableSections" :key="section" :value="section">
            {{ getSectionLabel(section) }}
          </option>
        </select>
      </div>
      <div class="filter-group">
        <label>Langue:</label>
        <select v-model="filterLanguage" class="filter-select">
          <option value="">Toutes les langues</option>
          <option value="fr">🇫🇷 Français</option>
          <option value="en">🇬🇧 English</option>
        </select>
      </div>
      <div class="search-group">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Rechercher par clé ou valeur..."
          class="search-input"
        />
      </div>
    </div>

    <!-- Content Grid -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement du contenu...</p>
    </div>

    <div v-else-if="filteredContent.length === 0" class="empty-state">
      <p>Aucun contenu trouvé</p>
    </div>

    <div v-else class="content-sections">
      <div v-for="section in groupedContent" :key="section.name" class="section-group">
        <h2 class="section-title">{{ getSectionLabel(section.name) }}</h2>
        <div class="content-grid">
          <div
            v-for="item in section.items"
            :key="item.id"
            class="content-card"
          >
            <div class="card-header">
              <div class="content-key">
                <code>{{ item.key }}</code>
                <span class="language-badge">{{ item.language.toUpperCase() }}</span>
              </div>
            </div>
            <div class="card-body">
              <p class="content-value">{{ item.value }}</p>
              <div class="meta">
                <span class="date">{{ formatDate(item.updated_at) }}</span>
              </div>
            </div>
            <div class="card-actions">
              <button @click="editContent(item)" class="btn-action btn-secondary">
                ✏️ Modifier
              </button>
              <button @click="confirmDelete(item)" class="btn-action btn-danger">
                🗑️ Supprimer
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit/Add Modal -->
    <Teleport to="body">
      <div v-if="showEditModal || showAddModal" class="modal-overlay" @click="closeModal">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>{{ showAddModal ? 'Ajouter un contenu' : 'Modifier le contenu' }}</h2>
            <button @click="closeModal" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveContent" class="edit-form">
              <div class="form-group">
                <label>Clé * <span class="hint">(ex: hero_title, about_description)</span></label>
                <input 
                  v-model="editForm.key" 
                  type="text" 
                  required 
                  :disabled="!showAddModal"
                  placeholder="hero_title"
                />
              </div>
              <div class="form-group">
                <label>Valeur *</label>
                <textarea 
                  v-model="editForm.value" 
                  rows="6" 
                  required
                  placeholder="Entrez le contenu ici..."
                ></textarea>
              </div>
              <div class="form-row">
                <div class="form-group">
                  <label>Section *</label>
                  <select v-model="editForm.section" required>
                    <option value="home">🏠 Accueil</option>
                    <option value="services">💼 Services</option>
                    <option value="about">ℹ️ À propos</option>
                    <option value="contact">📧 Contact</option>
                    <option value="gallery">🖼️ Galerie</option>
                    <option value="testimonials">💬 Témoignages</option>
                    <option value="footer">📄 Footer</option>
                  </select>
                </div>
                <div class="form-group">
                  <label>Langue *</label>
                  <select v-model="editForm.language" required>
                    <option value="fr">🇫🇷 Français</option>
                    <option value="en">🇬🇧 English</option>
                  </select>
                </div>
              </div>
              <div class="form-actions">
                <button type="button" @click="closeModal" class="btn-secondary">
                  Annuler
                </button>
                <button type="submit" class="btn-primary">
                  {{ showAddModal ? 'Ajouter' : 'Enregistrer' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <div v-if="deletingContent" class="modal-overlay" @click="deletingContent = null">
        <div class="modal-content modal-small" @click.stop>
          <div class="modal-header">
            <h2>Confirmer la suppression</h2>
            <button @click="deletingContent = null" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <p>Êtes-vous sûr de vouloir supprimer le contenu <strong>{{ deletingContent.key }}</strong> ?</p>
            <p class="warning-text">Cette action est irréversible.</p>
            <div class="form-actions">
              <button @click="deletingContent = null" class="btn-secondary">
                Annuler
              </button>
              <button @click="deleteContent" class="btn-danger">
                Supprimer
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../../services/api'

const content = ref([])
const loading = ref(false)
const filterSection = ref('')
const filterLanguage = ref('')
const searchQuery = ref('')
const showEditModal = ref(false)
const showAddModal = ref(false)
const editingContent = ref(null)
const deletingContent = ref(null)

const editForm = ref({
  key: '',
  value: '',
  section: 'home',
  language: 'fr'
})

const sectionLabels = {
  home: '🏠 Accueil',
  services: '💼 Services',
  about: 'ℹ️ À propos',
  contact: '📧 Contact',
  gallery: '🖼️ Galerie',
  testimonials: '💬 Témoignages',
  footer: '📄 Footer'
}

const stats = computed(() => {
  const total = content.value.length
  const sections = new Set(content.value.map(c => c.section)).size
  const languages = new Set(content.value.map(c => c.language)).size
  
  return { total, sections, languages }
})

const availableSections = computed(() => {
  return [...new Set(content.value.map(c => c.section))].sort()
})

const filteredContent = computed(() => {
  let filtered = content.value

  if (filterSection.value) {
    filtered = filtered.filter(c => c.section === filterSection.value)
  }

  if (filterLanguage.value) {
    filtered = filtered.filter(c => c.language === filterLanguage.value)
  }

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(c => 
      c.key.toLowerCase().includes(query) || 
      c.value.toLowerCase().includes(query)
    )
  }

  return filtered
})

const groupedContent = computed(() => {
  const groups = {}
  
  filteredContent.value.forEach(item => {
    if (!groups[item.section]) {
      groups[item.section] = {
        name: item.section,
        items: []
      }
    }
    groups[item.section].items.push(item)
  })

  return Object.values(groups).sort((a, b) => a.name.localeCompare(b.name))
})

async function fetchContent() {
  loading.value = true
  try {
    const response = await api.get('/admin/content')
    content.value = response.data || []
  } catch (error) {
    console.error('Failed to fetch content:', error)
    alert('Erreur lors du chargement du contenu')
  } finally {
    loading.value = false
  }
}

function getSectionLabel(section) {
  return sectionLabels[section] || section
}

function editContent(item) {
  editingContent.value = item
  editForm.value = {
    key: item.key,
    value: item.value,
    section: item.section,
    language: item.language
  }
  showEditModal.value = true
}

function closeModal() {
  showEditModal.value = false
  showAddModal.value = false
  editingContent.value = null
  editForm.value = {
    key: '',
    value: '',
    section: 'home',
    language: 'fr'
  }
}

async function saveContent() {
  try {
    await api.post('/admin/content', editForm.value)
    
    // Refresh content list
    await fetchContent()
    
    closeModal()
  } catch (error) {
    console.error('Failed to save content:', error)
    alert('Erreur lors de la sauvegarde')
  }
}

function confirmDelete(item) {
  deletingContent.value = item
}

async function deleteContent() {
  try {
    // Note: Backend doesn't have delete endpoint, so we'll update to empty or handle differently
    // For now, we'll just show an alert
    alert('La suppression de contenu n\'est pas encore implémentée dans le backend')
    deletingContent.value = null
  } catch (error) {
    console.error('Failed to delete content:', error)
    alert('Erreur lors de la suppression')
  }
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('fr-FR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

onMounted(() => {
  fetchContent()
})
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
}

.btn-add {
  padding: var(--spacing-md) var(--spacing-xl);
  background: var(--color-gold);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  font-weight: var(--font-weight-semibold);
  transition: all var(--transition-base);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-size: var(--font-size-sm);
}

.btn-add:hover {
  background: #C5A028;
  transform: translateY(-2px);
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-xl);
  margin-bottom: var(--spacing-3xl);
}

.stat-card {
  background: white;
  padding: var(--spacing-xl);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  transition: transform var(--transition-base);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.stat-icon {
  font-size: var(--font-size-4xl);
}

.stat-content h3 {
  font-size: var(--font-size-3xl);
  color: var(--color-gold);
  margin-bottom: var(--spacing-xs);
}

.stat-content p {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
}

/* Controls */
.controls {
  display: flex;
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-2xl);
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.filter-group label {
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.filter-select {
  padding: var(--spacing-sm) var(--spacing-md);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  background: white;
  cursor: pointer;
  transition: border-color var(--transition-base);
}

.filter-select:focus {
  outline: none;
  border-color: var(--color-gold);
}

.search-group {
  flex: 1;
  min-width: 250px;
}

.search-input {
  width: 100%;
  padding: var(--spacing-sm) var(--spacing-md);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: border-color var(--transition-base);
}

.search-input:focus {
  outline: none;
  border-color: var(--color-gold);
}

/* Content Sections */
.content-sections {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3xl);
}

.section-group {
  background: var(--color-gray-lighter);
  padding: var(--spacing-2xl);
  border-radius: var(--radius-xl);
}

.section-title {
  color: var(--color-gold);
  font-size: var(--font-size-2xl);
  margin-bottom: var(--spacing-xl);
  padding-bottom: var(--spacing-md);
  border-bottom: 2px solid var(--color-border);
}

.content-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: var(--spacing-xl);
}

.content-card {
  background: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  overflow: hidden;
  transition: transform var(--transition-base);
}

.content-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.card-header {
  padding: var(--spacing-lg);
  border-bottom: 1px solid var(--color-border);
}

.content-key {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-md);
}

.content-key code {
  background: var(--color-gray-lighter);
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-sm);
  font-family: 'Courier New', monospace;
  font-size: var(--font-size-sm);
  color: var(--color-gold);
}

.language-badge {
  padding: var(--spacing-xs) var(--spacing-sm);
  background: var(--color-gold);
  color: white;
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
}

.card-body {
  padding: var(--spacing-lg);
}

.content-value {
  color: var(--color-text);
  line-height: var(--line-height-relaxed);
  margin-bottom: var(--spacing-md);
  min-height: 60px;
}

.meta {
  display: flex;
  justify-content: flex-end;
  font-size: var(--font-size-sm);
  color: var(--color-gray);
}

.card-actions {
  padding: var(--spacing-md);
  background: var(--color-gray-lighter);
  display: flex;
  gap: var(--spacing-sm);
}

.btn-action {
  flex: 1;
  padding: var(--spacing-xs) var(--spacing-md);
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  transition: all var(--transition-base);
}

.btn-secondary {
  background: var(--color-gray-light);
  color: var(--color-text);
}

.btn-secondary:hover {
  background: var(--color-gray);
}

.btn-danger {
  background: #C1292E;
  color: white;
}

.btn-danger:hover {
  background: #A01F23;
}

/* Loading & Empty States */
.loading-state,
.empty-state {
  text-align: center;
  padding: var(--spacing-5xl);
  color: var(--color-gray);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--color-border);
  border-top-color: var(--color-gold);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto var(--spacing-md);
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Modals */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal);
  padding: var(--spacing-xl);
}

.modal-content {
  background: white;
  border-radius: var(--radius-xl);
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: var(--shadow-xl);
}

.modal-small {
  max-width: 400px;
}

.modal-header {
  padding: var(--spacing-xl);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h2 {
  color: var(--color-gold);
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
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  transition: background var(--transition-base);
}

.close-btn:hover {
  background: var(--color-gray-lighter);
}

.modal-body {
  padding: var(--spacing-xl);
}

/* Forms */
.edit-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-lg);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.form-group label {
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.hint {
  font-size: var(--font-size-xs);
  color: var(--color-gray);
  font-weight: normal;
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: var(--spacing-md);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-sans);
  font-size: var(--font-size-base);
  transition: border-color var(--transition-base);
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--color-gold);
}

.form-group input:disabled {
  background: var(--color-gray-lighter);
  cursor: not-allowed;
}

.form-actions {
  display: flex;
  gap: var(--spacing-md);
  justify-content: flex-end;
  margin-top: var(--spacing-lg);
}

.btn-primary {
  padding: var(--spacing-md) var(--spacing-xl);
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-semibold);
  cursor: pointer;
  transition: all var(--transition-base);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-size: var(--font-size-sm);
  background: var(--color-gold);
  color: white;
}

.btn-primary:hover {
  background: #C5A028;
}

.warning-text {
  color: #C1292E;
  font-size: var(--font-size-sm);
  margin-top: var(--spacing-sm);
}

@media (max-width: 768px) {
  .content-manager {
    padding: var(--spacing-xl);
  }

  .manager-header {
    flex-direction: column;
    gap: var(--spacing-lg);
  }

  .btn-add {
    width: 100%;
  }

  .content-grid {
    grid-template-columns: 1fr;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .controls {
    flex-direction: column;
  }

  .search-group {
    width: 100%;
  }
}
</style>
