<template>
  <div class="booking-manager-page">
    <div class="header-section">
      <h1>Gestion des Réservations</h1>
      
      <div class="view-toggles">
        <button 
          :class="['toggle-btn', { active: currentView === 'list' }]"
          @click="currentView = 'list'"
        >
          📋 Liste
        </button>
        <button 
          :class="['toggle-btn', { active: currentView === 'calendar' }]"
          @click="currentView = 'calendar'"
        >
          📅 Calendrier
        </button>
      </div>
    </div>

    <!-- Calendar View -->
    <div v-if="currentView === 'calendar'">
      <AvailabilityManager />
    </div>

    <!-- List View -->
    <div v-else class="bookings-list-view">
      <!-- Stats Cards -->
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">📅</div>
          <div class="stat-content">
            <h3>{{ stats.total }}</h3>
            <p>Total Réservations</p>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">⏳</div>
          <div class="stat-content">
            <h3>{{ stats.pending }}</h3>
            <p>En attente</p>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">✅</div>
          <div class="stat-content">
            <h3>{{ stats.confirmed }}</h3>
            <p>Confirmées</p>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">💰</div>
          <div class="stat-content">
            <h3>{{ formatPrice(stats.revenue) }}</h3>
            <p>Chiffre d'affaires</p>
          </div>
        </div>
      </div>

      <!-- Filters -->
      <div class="controls">
        <div class="filter-group">
          <select v-model="filterStatus" class="filter-select">
            <option value="all">Tous les statuts</option>
            <option value="pending">En attente</option>
            <option value="confirmed">Confirmées</option>
            <option value="paid">Payées</option>
            <option value="completed">Terminées</option>
            <option value="cancelled">Annulées</option>
          </select>
          
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Rechercher (client, email, type)..." 
            class="search-input"
          />
        </div>
      </div>

      <!-- Bookings Table -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Chargement des réservations...</p>
      </div>

      <div v-else-if="filteredBookings.length === 0" class="empty-state">
        <p>Aucune réservation trouvée</p>
      </div>

      <div v-else class="bookings-table-container">
        <table class="bookings-table">
          <thead>
            <tr>
              <th>Date</th>
              <th>Client</th>
              <th>Type</th>
              <th>Invités</th>
              <th>Montant</th>
              <th>Statut</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="booking in filteredBookings" :key="booking.id">
              <td>
                <div class="date-cell">
                  <strong>{{ formatDate(booking.event_date) }}</strong>
                  <span class="time-ago">{{ getTimeAgo(booking.event_date) }}</span>
                </div>
              </td>
              <td>
                <div class="client-cell">
                  <strong>{{ booking.client?.name || 'Client inconnu' }}</strong>
                  <span class="client-email">{{ booking.client?.email }}</span>
                </div>
              </td>
              <td>
                <span class="type-badge">{{ getEventTypeLabel(booking.event_type) }}</span>
              </td>
              <td>{{ booking.guest_count }}</td>
              <td>
                <div class="amount-cell">
                  <strong>{{ formatPrice(booking.total_amount) }}</strong>
                  <span class="deposit" v-if="booking.deposit_amount">
                    Acompte: {{ formatPrice(booking.deposit_amount) }}
                  </span>
                </div>
              </td>
              <td>
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
              </td>
              <td>
                <button @click="viewBookingDetails(booking)" class="btn-icon" title="Voir détails">
                  👁️
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Booking Details Modal -->
    <Teleport to="body">
      <div v-if="selectedBooking" class="modal-overlay" @click="selectedBooking = null">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>Détails de la réservation #{{ selectedBooking.id }}</h2>
            <button @click="selectedBooking = null" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <div class="details-grid">
              <div class="detail-section">
                <h3>Informations Client</h3>
                <p><strong>Nom:</strong> {{ selectedBooking.client?.name }}</p>
                <p><strong>Email:</strong> {{ selectedBooking.client?.email }}</p>
                <p><strong>Téléphone:</strong> {{ selectedBooking.client?.phone || '-' }}</p>
              </div>

              <div class="detail-section">
                <h3>Événement</h3>
                <p><strong>Date:</strong> {{ formatDate(selectedBooking.event_date) }}</p>
                <p><strong>Type:</strong> {{ getEventTypeLabel(selectedBooking.event_type) }}</p>
                <p><strong>Lieu:</strong> {{ selectedBooking.event_location || 'Non défini' }}</p>
                <p><strong>Invités:</strong> {{ selectedBooking.guest_count }} personnes</p>
              </div>

              <div class="detail-section full-width">
                <h3>Finances</h3>
                <div class="finance-grid">
                  <div class="finance-item">
                    <span>Budget Estimé</span>
                    <strong>{{ formatPrice(selectedBooking.budget) }}</strong>
                  </div>
                  <div class="finance-item">
                    <span>Montant Total</span>
                    <strong>{{ formatPrice(selectedBooking.total_amount) }}</strong>
                  </div>
                  <div class="finance-item">
                    <span>Acompte (30%)</span>
                    <strong>{{ formatPrice(selectedBooking.deposit_amount) }}</strong>
                  </div>
                </div>
              </div>

              <div class="detail-section full-width" v-if="selectedBooking.message">
                <h3>Message du client</h3>
                <p class="message-box">{{ selectedBooking.message }}</p>
              </div>

               <div class="detail-section full-width" v-if="selectedBooking.special_requests">
                <h3>Demandes spéciales</h3>
                <p class="message-box">{{ selectedBooking.special_requests }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AvailabilityManager from '../../components/admin/AvailabilityManager.vue'
import api from '../../services/api'

const currentView = ref('list')
const bookings = ref([])
const loading = ref(false)
const filterStatus = ref('all')
const searchQuery = ref('')
const selectedBooking = ref(null)

const eventTypeLabels = {
  proposal: '💍 Demande en mariage',
  wedding: '💐 Mariage',
  birthday: '🎂 Anniversaire',
  baby_shower: '👶 Baby Shower',
  corporate: '🏢 Corporatif',
  other: '✨ Autre'
}

const stats = computed(() => {
  const total = bookings.value.length
  const pending = bookings.value.filter(b => b.status === 'pending').length
  const confirmed = bookings.value.filter(b => b.status === 'confirmed').length
  const revenue = bookings.value
    .filter(b => ['paid', 'completed'].includes(b.status))
    .reduce((sum, b) => sum + (b.total_amount || 0), 0)
  
  return { total, pending, confirmed, revenue }
})

const filteredBookings = computed(() => {
  let filtered = bookings.value

  // Status filter
  if (filterStatus.value !== 'all') {
    filtered = filtered.filter(b => b.status === filterStatus.value)
  }

  // Search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(b => 
      (b.client?.name && b.client.name.toLowerCase().includes(query)) ||
      (b.client?.email && b.client.email.toLowerCase().includes(query)) ||
      b.event_type.toLowerCase().includes(query)
    )
  }

  return filtered
})

async function fetchBookings() {
  loading.value = true
  try {
    const response = await api.get('/admin/bookings')
    bookings.value = response.data || []
  } catch (error) {
    console.error('Failed to fetch bookings:', error)
  } finally {
    loading.value = false
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
  } catch (error) {
    console.error('Failed to update booking status:', error)
    alert('Erreur lors de la mise à jour du statut')
  }
}

function viewBookingDetails(booking) {
  selectedBooking.value = booking
}

function getEventTypeLabel(type) {
  return eventTypeLabels[type] || type
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('fr-FR', {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

function getTimeAgo(dateStr) {
  const date = new Date(dateStr)
  const now = new Date()
  const diffTime = date - now
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) return "Aujourd'hui"
  if (diffDays === 1) return "Demain"
  if (diffDays === -1) return "Hier"
  if (diffDays > 0) return `Dans ${diffDays} jours`
  return `Il y a ${Math.abs(diffDays)} jours`
}

function formatPrice(price) {
  if (!price) return '0 €'
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(price)
}

onMounted(() => {
  fetchBookings()
})
</script>

<style scoped>
.booking-manager-page {
  padding: var(--spacing-3xl);
  max-width: 1400px;
  margin: 0 auto;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-3xl);
}

.header-section h1 {
  color: var(--color-gold);
  font-size: var(--font-size-4xl);
}

.view-toggles {
  display: flex;
  background: var(--color-gray-lighter);
  padding: 4px;
  border-radius: var(--radius-md);
  gap: 4px;
}

.toggle-btn {
  padding: var(--spacing-sm) var(--spacing-lg);
  border: none;
  background: transparent;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-weight: var(--font-weight-semibold);
  transition: all var(--transition-base);
}

.toggle-btn.active {
  background: white;
  color: var(--color-gold);
  box-shadow: var(--shadow-sm);
}

/* Stats */
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
}

.stat-icon {
  font-size: var(--font-size-4xl);
}

.stat-content h3 {
  font-size: var(--font-size-2xl);
  color: var(--color-gold);
  margin-bottom: var(--spacing-xs);
}

.stat-content p {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
}

/* Controls */
.controls {
  margin-bottom: var(--spacing-xl);
}

.filter-group {
  display: flex;
  gap: var(--spacing-md);
  flex-wrap: wrap;
}

.filter-select,
.search-input {
  padding: var(--spacing-md);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
}

.search-input {
  flex: 1;
  min-width: 250px;
}

/* Table */
.bookings-table-container {
  background: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  overflow-x: auto;
}

.bookings-table {
  width: 100%;
  border-collapse: collapse;
}

.bookings-table th {
  text-align: left;
  padding: var(--spacing-lg);
  background: var(--color-gray-lighter);
  color: var(--color-text);
  font-weight: var(--font-weight-semibold);
  border-bottom: 2px solid var(--color-border);
}

.bookings-table td {
  padding: var(--spacing-lg);
  border-bottom: 1px solid var(--color-border);
  vertical-align: middle;
}

.date-cell, .client-cell, .amount-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.time-ago, .client-email {
  font-size: var(--font-size-xs);
  color: var(--color-gray);
}

.deposit {
  font-size: var(--font-size-xs);
  color: var(--color-gold);
}

.type-badge {
  background: var(--color-gray-lighter);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
}

.status-select {
  padding: 6px 12px;
  border-radius: 4px;
  border: 1px solid transparent;
  font-weight: bold;
  cursor: pointer;
}

.status-select.pending { background: #FFF3CD; color: #856404; }
.status-select.confirmed { background: #D4EDDA; color: #155724; }
.status-select.paid { background: #CCE5FF; color: #004085; }
.status-select.completed { background: #D1ECF1; color: #0C5460; }
.status-select.cancelled { background: #F8D7DA; color: #721C24; }

.btn-icon {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
}

.btn-icon:hover {
  background: var(--color-gray-lighter);
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: var(--radius-lg);
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: var(--shadow-xl);
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
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
}

.modal-body {
  padding: var(--spacing-xl);
}

.details-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-xl);
}

.detail-section h3 {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
  text-transform: uppercase;
  margin-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--color-border);
  padding-bottom: var(--spacing-xs);
}

.detail-section p {
  margin-bottom: var(--spacing-xs);
  line-height: 1.4;
}

.full-width {
  grid-column: span 2;
}

.finance-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--spacing-md);
  background: var(--color-gray-lighter);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
}

.finance-item {
  display: flex;
  flex-direction: column;
}

.finance-item span {
  font-size: var(--font-size-xs);
  color: var(--color-gray);
}

.finance-item strong {
  color: var(--color-gold);
  font-size: var(--font-size-lg);
}

.message-box {
  background: var(--color-gray-lighter);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  font-style: italic;
}.loading-state, .empty-state {
  text-align: center;
  padding: var(--spacing-3xl);
  color: var(--color-gray);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--color-gray-lighter);
  border-top-color: var(--color-gold);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin { to { transform: rotate(360deg); } }

@media (max-width: 768px) {
  .details-grid {
    grid-template-columns: 1fr;
  }
  
  .full-width {
    grid-column: span 1;
  }
  
  .finance-grid {
    grid-template-columns: 1fr;
  }
}
</style>
