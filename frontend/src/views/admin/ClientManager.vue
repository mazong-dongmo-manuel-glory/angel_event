<template>
  <div class="client-manager">
    <!-- Header -->
    <div class="manager-header">
      <h1>Gestion des Clients</h1>
      <p class="subtitle">Gérez vos clients et consultez leurs réservations</p>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">👥</div>
        <div class="stat-content">
          <h3>{{ stats.total }}</h3>
          <p>Total Clients</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">📅</div>
        <div class="stat-content">
          <h3>{{ stats.withBookings }}</h3>
          <p>Avec Réservations</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">✨</div>
        <div class="stat-content">
          <h3>{{ stats.newThisMonth }}</h3>
          <p>Nouveaux ce mois</p>
        </div>
      </div>
    </div>

    <!-- Search and Filters -->
    <div class="controls">
      <div class="search-box">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Rechercher par nom ou email..."
          @input="handleSearch"
        />
      </div>
    </div>

    <!-- Clients Table -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement des clients...</p>
    </div>

    <div v-else-if="filteredClients.length === 0" class="empty-state">
      <p>Aucun client trouvé</p>
    </div>

    <div v-else class="clients-table">
      <table>
        <thead>
          <tr>
            <th>Nom</th>
            <th>Email</th>
            <th>Téléphone</th>
            <th>Réservations</th>
            <th>Créé le</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="client in filteredClients" :key="client.id">
            <td>
              <strong>{{ client.name }}</strong>
            </td>
            <td>{{ client.email }}</td>
            <td>{{ client.phone || '-' }}</td>
            <td>
              <span class="badge">{{ client.bookings?.length || 0 }}</span>
            </td>
            <td>{{ formatDate(client.created_at) }}</td>
            <td>
              <div class="action-buttons">
                <button @click="viewClient(client)" class="btn-icon" title="Voir détails">
                  👁️
                </button>
                <button @click="editClient(client)" class="btn-icon" title="Modifier">
                  ✏️
                </button>
                <button @click="confirmDelete(client)" class="btn-icon btn-danger" title="Supprimer">
                  🗑️
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- View Client Modal -->
    <Teleport to="body">
      <div v-if="viewingClient" class="modal-overlay" @click="viewingClient = null">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>Détails du Client</h2>
            <button @click="viewingClient = null" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <div class="client-details">
              <div class="detail-row">
                <strong>Nom:</strong>
                <span>{{ viewingClient.name }}</span>
              </div>
              <div class="detail-row">
                <strong>Email:</strong>
                <span>{{ viewingClient.email }}</span>
              </div>
              <div class="detail-row">
                <strong>Téléphone:</strong>
                <span>{{ viewingClient.phone || 'Non renseigné' }}</span>
              </div>
              <div class="detail-row">
                <strong>Notes:</strong>
                <span>{{ viewingClient.notes || 'Aucune note' }}</span>
              </div>
              <div class="detail-row">
                <strong>Créé le:</strong>
                <span>{{ formatDate(viewingClient.created_at) }}</span>
              </div>
            </div>

            <div v-if="viewingClient.bookings && viewingClient.bookings.length > 0" class="bookings-section">
              <h3>Réservations ({{ viewingClient.bookings.length }})</h3>
              <div class="bookings-list">
                <div v-for="booking in viewingClient.bookings" :key="booking.id" class="booking-card">
                  <div class="booking-info">
                    <strong>{{ booking.event_type }}</strong>
                    <span>{{ formatDate(booking.event_date) }}</span>
                    <span class="booking-amount">{{ formatPrice(booking.total_amount) }}</span>
                  </div>
                  
                  <div class="booking-actions">
                    <select 
                      :value="booking.status" 
                      @change="updateStatus(booking, $event.target.value)"
                      :class="['status-select', booking.status]"
                    >
                      <option value="pending">En attente</option>
                      <option value="confirmed">Confirmé</option>
                      <option value="paid">Payé</option>
                      <option value="completed">Terminé</option>
                      <option value="cancelled">Annulé</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Edit Client Modal -->
    <Teleport to="body">
      <div v-if="editingClient" class="modal-overlay" @click="editingClient = null">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>Modifier le Client</h2>
            <button @click="editingClient = null" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveClient" class="edit-form">
              <div class="form-group">
                <label>Nom *</label>
                <input v-model="editForm.name" type="text" required />
              </div>
              <div class="form-group">
                <label>Email *</label>
                <input v-model="editForm.email" type="email" required />
              </div>
              <div class="form-group">
                <label>Téléphone</label>
                <input v-model="editForm.phone" type="tel" />
              </div>
              <div class="form-group">
                <label>Notes</label>
                <textarea v-model="editForm.notes" rows="4"></textarea>
              </div>
              <div class="form-actions">
                <button type="button" @click="editingClient = null" class="btn-secondary">
                  Annuler
                </button>
                <button type="submit" class="btn-primary">
                  Enregistrer
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <div v-if="deletingClient" class="modal-overlay" @click="deletingClient = null">
        <div class="modal-content modal-small" @click.stop>
          <div class="modal-header">
            <h2>Confirmer la suppression</h2>
            <button @click="deletingClient = null" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <p>Êtes-vous sûr de vouloir supprimer le client <strong>{{ deletingClient.name }}</strong> ?</p>
            <p class="warning-text">Cette action est irréversible.</p>
            <div class="form-actions">
              <button @click="deletingClient = null" class="btn-secondary">
                Annuler
              </button>
              <button @click="deleteClient" class="btn-danger">
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

const clients = ref([])
const loading = ref(false)
const searchQuery = ref('')
const viewingClient = ref(null)
const editingClient = ref(null)
const deletingClient = ref(null)

const editForm = ref({
  name: '',
  email: '',
  phone: '',
  notes: ''
})

const stats = computed(() => {
  const total = clients.value.length
  const withBookings = clients.value.filter(c => c.bookings && c.bookings.length > 0).length
  
  const now = new Date()
  const firstDayOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
  const newThisMonth = clients.value.filter(c => new Date(c.created_at) >= firstDayOfMonth).length
  
  return { total, withBookings, newThisMonth }
})

const filteredClients = computed(() => {
  if (!searchQuery.value) return clients.value
  
  const query = searchQuery.value.toLowerCase()
  return clients.value.filter(client =>
    client.name.toLowerCase().includes(query) ||
    client.email.toLowerCase().includes(query)
  )
})

async function fetchClients() {
  loading.value = true
  try {
    const response = await api.get('/admin/clients')
    clients.value = response.data
  } catch (error) {
    console.error('Failed to fetch clients:', error)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  // Search is handled by computed property
}

function viewClient(client) {
  viewingClient.value = client
}

function editClient(client) {
  editingClient.value = client
  editForm.value = {
    name: client.name,
    email: client.email,
    phone: client.phone || '',
    notes: client.notes || ''
  }
}

async function updateStatus(booking, newStatus) {
  try {
    const response = await api.put(`/admin/bookings/${booking.id}/status`, {
      status: newStatus,
      admin_notes: booking.admin_notes
    })
    
    // Update local data
    booking.status = newStatus
    
    // Also update in the main list
    const client = clients.value.find(c => c.id === booking.client_id)
    if (client) {
      const b = client.bookings.find(b => b.id === booking.id)
      if (b) b.status = newStatus
    }
  } catch (error) {
    console.error('Failed to update booking status:', error)
    alert('Erreur lors de la mise à jour du statut')
  }
}

async function saveClient() {
  try {
    await api.put(`/admin/clients/${editingClient.value.id}`, editForm.value)
    
    // Update local data
    const index = clients.value.findIndex(c => c.id === editingClient.value.id)
    if (index !== -1) {
      clients.value[index] = { ...clients.value[index], ...editForm.value }
    }
    
    editingClient.value = null
  } catch (error) {
    console.error('Failed to update client:', error)
    alert('Erreur lors de la mise à jour du client')
  }
}

function confirmDelete(client) {
  deletingClient.value = client
}

async function deleteClient() {
  try {
    await api.delete(`/admin/clients/${deletingClient.value.id}`)
    
    // Remove from local data
    clients.value = clients.value.filter(c => c.id !== deletingClient.value.id)
    
    deletingClient.value = null
  } catch (error) {
    console.error('Failed to delete client:', error)
    alert('Erreur lors de la suppression du client')
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

function formatPrice(price) {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(price)
}

onMounted(() => {
  fetchClients()
})
</script>

<style scoped>
.client-manager {
  padding: var(--spacing-3xl);
  max-width: 1400px;
  margin: 0 auto;
}

.manager-header {
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

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
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
  margin-bottom: var(--spacing-2xl);
}

.search-box input {
  width: 100%;
  max-width: 500px;
  padding: var(--spacing-md) var(--spacing-lg);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--font-size-base);
  transition: border-color var(--transition-base);
}

.search-box input:focus {
  outline: none;
  border-color: var(--color-gold);
}

/* Table */
.clients-table {
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
  font-size: var(--font-size-sm);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

td {
  padding: var(--spacing-lg);
  border-top: 1px solid var(--color-border);
}

tbody tr:hover {
  background: var(--color-gray-lighter);
}

.badge {
  display: inline-block;
  padding: var(--spacing-xs) var(--spacing-sm);
  background: var(--color-gold);
  color: white;
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
}

.action-buttons {
  display: flex;
  gap: var(--spacing-sm);
}

.btn-icon {
  background: none;
  border: none;
  font-size: var(--font-size-xl);
  cursor: pointer;
  padding: var(--spacing-xs);
  border-radius: var(--radius-sm);
  transition: background var(--transition-base);
}

.btn-icon:hover {
  background: var(--color-gray-lighter);
}

.btn-icon.btn-danger:hover {
  background: rgba(193, 41, 46, 0.1);
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

.client-details {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-2xl);
}

.detail-row {
  display: flex;
  gap: var(--spacing-md);
}

.detail-row strong {
  min-width: 120px;
  color: var(--color-gray);
}

.bookings-section {
  margin-top: var(--spacing-2xl);
  padding-top: var(--spacing-2xl);
  border-top: 1px solid var(--color-border);
}

.bookings-section h3 {
  color: var(--color-gold);
  margin-bottom: var(--spacing-lg);
}

.bookings-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.booking-card {
  padding: var(--spacing-md);
  background: var(--color-gray-lighter);
  border-radius: var(--radius-md);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-md);
}

.booking-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  flex: 1;
}

.booking-amount {
  font-weight: var(--font-weight-semibold);
  color: var(--color-gold);
}

.status-badge {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  text-transform: uppercase;
}

.status-select {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  text-transform: uppercase;
  border: 1px solid transparent;
  cursor: pointer;
  background-color: white;
}

.status-select:focus {
  outline: none;
  border-color: rgba(0,0,0,0.2);
}

.status-select.pending {
  background-color: #FFF3CD;
  color: #856404;
}

.status-select.confirmed {
  background-color: #D4EDDA;
  color: #155724;
}

.status-select.paid {
  background-color: #CCE5FF;
  color: #004085;
}

.status-select.completed {
  background-color: #D1ECF1;
  color: #0C5460;
}

.status-select.cancelled {
  background-color: #F8D7DA;
  color: #721C24;
}

/* Forms */
.edit-form {
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

.form-group input,
.form-group textarea {
  padding: var(--spacing-md);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-sans);
  font-size: var(--font-size-base);
  transition: border-color var(--transition-base);
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--color-gold);
}

.form-actions {
  display: flex;
  gap: var(--spacing-md);
  justify-content: flex-end;
  margin-top: var(--spacing-lg);
}

.btn-primary,
.btn-secondary,
.btn-danger {
  padding: var(--spacing-md) var(--spacing-xl);
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-semibold);
  cursor: pointer;
  transition: all var(--transition-base);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-size: var(--font-size-sm);
}

.btn-primary {
  background: var(--color-gold);
  color: white;
}

.btn-primary:hover {
  background: #C5A028;
}

.btn-secondary {
  background: var(--color-gray-lighter);
  color: var(--color-text);
}

.btn-secondary:hover {
  background: var(--color-gray-light);
}

.btn-danger {
  background: #C1292E;
  color: white;
}

.btn-danger:hover {
  background: #A01F23;
}

.warning-text {
  color: #C1292E;
  font-size: var(--font-size-sm);
  margin-top: var(--spacing-sm);
}

@media (max-width: 768px) {
  .client-manager {
    padding: var(--spacing-xl);
  }

  .clients-table {
    overflow-x: auto;
  }

  table {
    min-width: 800px;
  }
}
</style>
