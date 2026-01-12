<template>
  <div class="newsletter-manager">
    <!-- Header -->
    <div class="manager-header">
      <div>
        <h1>Gestion de la Newsletter</h1>
        <p class="subtitle">Gérez les abonnés et envoyez des campagnes d'emailing</p>
      </div>
      <button @click="showSendModal = true" class="btn-send">
        📧 Envoyer une Newsletter
      </button>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">👥</div>
        <div class="stat-content">
          <h3>{{ stats.total }}</h3>
          <p>Total Abonnés</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">✅</div>
        <div class="stat-content">
          <h3>{{ stats.active }}</h3>
          <p>Actifs</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">🇫🇷</div>
        <div class="stat-content">
          <h3>{{ stats.french }}</h3>
          <p>Français</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">🇬🇧</div>
        <div class="stat-content">
          <h3>{{ stats.english }}</h3>
          <p>English</p>
        </div>
      </div>
    </div>

    <!-- Controls -->
    <div class="controls">
      <div class="filter-buttons">
        <button
          :class="['filter-btn', { active: filter === 'all' }]"
          @click="filter = 'all'"
        >
          Tous
        </button>
        <button
          :class="['filter-btn', { active: filter === 'active' }]"
          @click="filter = 'active'"
        >
          Actifs
        </button>
        <button
          :class="['filter-btn', { active: filter === 'inactive' }]"
          @click="filter = 'inactive'"
        >
          Inactifs
        </button>
      </div>
      <div class="search-group">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Rechercher par email ou nom..."
          class="search-input"
        />
      </div>
    </div>

    <!-- Subscribers List -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement des abonnés...</p>
    </div>

    <div v-else-if="filteredSubscribers.length === 0" class="empty-state">
      <p>Aucun abonné trouvé</p>
    </div>

    <div v-else class="subscribers-table">
      <table>
        <thead>
          <tr>
            <th>Email</th>
            <th>Nom</th>
            <th>Langue</th>
            <th>Statut</th>
            <th>Date d'inscription</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="subscriber in filteredSubscribers" :key="subscriber.id">
            <td>
              <strong>{{ subscriber.email }}</strong>
            </td>
            <td>{{ subscriber.name || '-' }}</td>
            <td>
              <span class="language-badge">
                {{ subscriber.language === 'fr' ? '🇫🇷 FR' : '🇬🇧 EN' }}
              </span>
            </td>
            <td>
              <span 
                class="status-badge" 
                :class="subscriber.active ? 'status-active' : 'status-inactive'"
              >
                {{ subscriber.active ? '✓ Actif' : '✗ Inactif' }}
              </span>
            </td>
            <td>{{ formatDate(subscriber.created_at) }}</td>
            <td>
              <div class="action-buttons">
                <button
                  v-if="subscriber.active"
                  @click="toggleActive(subscriber)"
                  class="btn-action btn-warning"
                  title="Désactiver"
                >
                  🔕
                </button>
                <button
                  v-else
                  @click="toggleActive(subscriber)"
                  class="btn-action btn-success"
                  title="Activer"
                >
                  🔔
                </button>
                <button
                  @click="confirmDelete(subscriber)"
                  class="btn-action btn-danger"
                  title="Supprimer"
                >
                  🗑️
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Send Newsletter Modal -->
    <Teleport to="body">
      <div v-if="showSendModal" class="modal-overlay" @click="showSendModal = false">
        <div class="modal-content modal-large" @click.stop>
          <div class="modal-header">
            <h2>Envoyer une Newsletter</h2>
            <button @click="showSendModal = false" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="sendNewsletter" class="send-form">
              <div class="form-group">
                <label>Destinataires *</label>
                <select v-model="sendForm.language" required>
                  <option value="">Tous les abonnés actifs</option>
                  <option value="fr">🇫🇷 Abonnés français uniquement</option>
                  <option value="en">🇬🇧 Abonnés anglais uniquement</option>
                </select>
                <p class="recipient-count">
                  {{ getRecipientCount() }} destinataire(s)
                </p>
              </div>
              <div class="form-group">
                <label>Sujet *</label>
                <input 
                  v-model="sendForm.subject" 
                  type="text" 
                  required
                  placeholder="Nouveautés Angel Event - Janvier 2026"
                />
              </div>
              <div class="form-group">
                <label>Contenu * <span class="hint">(HTML supporté)</span></label>
                <textarea 
                  v-model="sendForm.content" 
                  rows="12" 
                  required
                  placeholder="<h2>Bonjour,</h2><p>Découvrez nos nouveautés...</p>"
                ></textarea>
              </div>
              <div class="preview-section">
                <h3>Aperçu</h3>
                <div class="preview-box" v-html="sendForm.content"></div>
              </div>
              <div class="form-actions">
                <button type="button" @click="showSendModal = false" class="btn-secondary">
                  Annuler
                </button>
                <button type="submit" class="btn-primary" :disabled="sending">
                  {{ sending ? 'Envoi en cours...' : '📧 Envoyer' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <div v-if="deletingSubscriber" class="modal-overlay" @click="deletingSubscriber = null">
        <div class="modal-content modal-small" @click.stop>
          <div class="modal-header">
            <h2>Confirmer la suppression</h2>
            <button @click="deletingSubscriber = null" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <p>Êtes-vous sûr de vouloir supprimer l'abonné <strong>{{ deletingSubscriber.email }}</strong> ?</p>
            <p class="warning-text">Cette action est irréversible.</p>
            <div class="form-actions">
              <button @click="deletingSubscriber = null" class="btn-secondary">
                Annuler
              </button>
              <button @click="deleteSubscriber" class="btn-danger">
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

const subscribers = ref([])
const loading = ref(false)
const sending = ref(false)
const filter = ref('all')
const searchQuery = ref('')
const showSendModal = ref(false)
const deletingSubscriber = ref(null)

const sendForm = ref({
  language: '',
  subject: '',
  content: ''
})

const stats = computed(() => {
  const total = subscribers.value.length
  const active = subscribers.value.filter(s => s.active).length
  const french = subscribers.value.filter(s => s.language === 'fr').length
  const english = subscribers.value.filter(s => s.language === 'en').length
  
  return { total, active, french, english }
})

const filteredSubscribers = computed(() => {
  let filtered = subscribers.value

  if (filter.value === 'active') {
    filtered = filtered.filter(s => s.active)
  } else if (filter.value === 'inactive') {
    filtered = filtered.filter(s => !s.active)
  }

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(s => 
      s.email.toLowerCase().includes(query) || 
      (s.name && s.name.toLowerCase().includes(query))
    )
  }

  return filtered.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
})

async function fetchSubscribers() {
  loading.value = true
  try {
    const response = await api.get('/admin/newsletter/subscribers')
    subscribers.value = response.data || []
  } catch (error) {
    console.error('Failed to fetch subscribers:', error)
    alert('Erreur lors du chargement des abonnés')
  } finally {
    loading.value = false
  }
}

async function toggleActive(subscriber) {
  try {
    const updatedActive = !subscriber.active
    await api.put(`/admin/newsletter/subscribers/${subscriber.id}`, {
      active: updatedActive
    })
    
    // Update local data
    subscriber.active = updatedActive
  } catch (error) {
    console.error('Failed to toggle active status:', error)
    alert('Erreur lors de la modification du statut')
  }
}


function confirmDelete(subscriber) {
  deletingSubscriber.value = subscriber
}

async function deleteSubscriber() {
  if (!deletingSubscriber.value) return
  
  try {
    await api.delete(`/admin/newsletter/subscribers/${deletingSubscriber.value.id}`)
    
    // Remove from local list
    subscribers.value = subscribers.value.filter(s => s.id !== deletingSubscriber.value.id)
    deletingSubscriber.value = null
  } catch (error) {
    console.error('Failed to delete subscriber:', error)
    alert('Erreur lors de la suppression')
    deletingSubscriber.value = null
  }
}


function getRecipientCount() {
  let count = subscribers.value.filter(s => s.active)
  
  if (sendForm.value.language) {
    count = count.filter(s => s.language === sendForm.value.language)
  }
  
  return count.length
}

async function sendNewsletter() {
  if (!confirm(`Êtes-vous sûr de vouloir envoyer cette newsletter à ${getRecipientCount()} destinataire(s) ?`)) {
    return
  }

  sending.value = true
  try {
    await api.post('/admin/newsletter/send', sendForm.value)
    
    alert('Newsletter envoyée avec succès!')
    showSendModal.value = false
    sendForm.value = {
      language: '',
      subject: '',
      content: ''
    }
  } catch (error) {
    console.error('Failed to send newsletter:', error)
    alert('Erreur lors de l\'envoi de la newsletter')
  } finally {
    sending.value = false
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
  fetchSubscribers()
})
</script>

<style scoped>
.newsletter-manager {
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

.btn-send {
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

.btn-send:hover {
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

.filter-buttons {
  display: flex;
  gap: var(--spacing-md);
}

.filter-btn {
  padding: var(--spacing-md) var(--spacing-xl);
  border: 2px solid var(--color-border);
  background: white;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-base);
  font-weight: var(--font-weight-semibold);
  text-transform: uppercase;
  font-size: var(--font-size-sm);
  letter-spacing: 0.05em;
}

.filter-btn:hover {
  border-color: var(--color-gold);
}

.filter-btn.active {
  background: var(--color-gold);
  border-color: var(--color-gold);
  color: white;
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

/* Table */
.subscribers-table {
  background: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  overflow: hidden;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  background: var(--color-gray-lighter);
}

th {
  padding: var(--spacing-lg);
  text-align: left;
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
  border-bottom: 2px solid var(--color-border);
}

td {
  padding: var(--spacing-lg);
  border-bottom: 1px solid var(--color-border);
}

tbody tr:hover {
  background: var(--color-gray-lighter);
}

tbody tr:last-child td {
  border-bottom: none;
}

.language-badge {
  padding: var(--spacing-xs) var(--spacing-sm);
  background: var(--color-gold);
  color: white;
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  display: inline-block;
}

.status-badge {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  display: inline-block;
}

.status-active {
  background: #2D7A4F;
  color: white;
}

.status-inactive {
  background: #C1292E;
  color: white;
}

.action-buttons {
  display: flex;
  gap: var(--spacing-xs);
}

.btn-action {
  padding: var(--spacing-xs) var(--spacing-sm);
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-base);
  transition: all var(--transition-base);
}

.btn-success {
  background: #2D7A4F;
  color: white;
}

.btn-success:hover {
  background: #246339;
}

.btn-warning {
  background: #FFA500;
  color: white;
}

.btn-warning:hover {
  background: #FF8C00;
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
  background: white;
  border-radius: var(--radius-lg);
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

.modal-large {
  max-width: 800px;
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
.send-form {
  display: flex;
  flex-direction: column;
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

.recipient-count {
  font-size: var(--font-size-sm);
  color: var(--color-gold);
  font-weight: var(--font-weight-semibold);
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

.preview-section {
  margin-top: var(--spacing-lg);
}

.preview-section h3 {
  color: var(--color-gold);
  margin-bottom: var(--spacing-md);
}

.preview-box {
  padding: var(--spacing-lg);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-gray-lighter);
  min-height: 200px;
  max-height: 400px;
  overflow-y: auto;
}

.form-actions {
  display: flex;
  gap: var(--spacing-md);
  justify-content: flex-end;
  margin-top: var(--spacing-lg);
}

.btn-secondary {
  padding: var(--spacing-md) var(--spacing-xl);
  border: 2px solid var(--color-border);
  background: white;
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-semibold);
  cursor: pointer;
  transition: all var(--transition-base);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-size: var(--font-size-sm);
}

.btn-secondary:hover {
  border-color: var(--color-gold);
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

.btn-primary:hover:not(:disabled) {
  background: #C5A028;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.warning-text {
  color: #C1292E;
  font-size: var(--font-size-sm);
  margin-top: var(--spacing-sm);
}

@media (max-width: 768px) {
  .newsletter-manager {
    padding: var(--spacing-xl);
  }

  .manager-header {
    flex-direction: column;
    gap: var(--spacing-lg);
  }

  .btn-send {
    width: 100%;
  }

  .controls {
    flex-direction: column;
  }

  .search-group {
    width: 100%;
  }

  .subscribers-table {
    overflow-x: auto;
  }

  table {
    min-width: 800px;
  }
}
</style>
