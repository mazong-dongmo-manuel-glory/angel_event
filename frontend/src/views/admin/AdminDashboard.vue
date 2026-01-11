<template>
  <div class="admin-dashboard">
    <!-- Stats Cards -->
    <div class="stats-grid">
      <div v-for="stat in stats" :key="stat.label" class="stat-card">
        <div class="stat-icon" :style="{ background: stat.color }">
          {{ stat.icon }}
        </div>
        <div class="stat-content">
          <h3>{{ stat.value }}</h3>
          <p>{{ stat.label }}</p>
        </div>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="dashboard-section">
      <h2>Activité Récente</h2>
      <Card>
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>Chargement...</p>
        </div>
        <div v-else class="activity-list">
          <div v-for="booking in recentBookings" :key="booking.id" class="activity-item">
            <div class="activity-icon">📅</div>
            <div class="activity-details">
              <strong>{{ booking.client?.name }}</strong>
              <p>{{ formatEventType(booking.event_type) }} - {{ formatDate(booking.event_date) }}</p>
            </div>
            <div class="activity-status" :class="booking.status">
              {{ formatStatus(booking.status) }}
            </div>
          </div>
          <div v-if="recentBookings.length === 0" class="empty-state">
            Aucune réservation récente
          </div>
        </div>
      </Card>
    </div>

    <!-- Quick Actions -->
    <div class="dashboard-section">
      <h2>Actions Rapides</h2>
      <div class="actions-grid">
        <button class="action-card" @click="$router.push('/admin/bookings')">
          <span class="action-icon">📅</span>
          <span>Gérer les réservations</span>
        </button>
        <button class="action-card" @click="$router.push('/admin/clients')">
          <span class="action-icon">👥</span>
          <span>Voir les clients</span>
        </button>
        <button class="action-card" @click="$router.push('/admin/gallery')">
          <span class="action-icon">🖼️</span>
          <span>Ajouter des photos</span>
        </button>
        <button class="action-card" @click="$router.push('/admin/newsletter')">
          <span class="action-icon">📧</span>
          <span>Envoyer newsletter</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Card from '../../components/ui/Card.vue'
import api from '../../services/api'

const loading = ref(true)
const stats = ref([
  { icon: '📅', label: 'Réservations totales', value: '0', color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
  { icon: '⏳', label: 'En attente', value: '0', color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
  { icon: '✓', label: 'Confirmées', value: '0', color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
  { icon: '👥', label: 'Clients', value: '0', color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' },
])

const recentBookings = ref([])

const eventTypeLabels = {
  proposal: 'Demande en mariage',
  wedding: 'Mariage',
  birthday: 'Anniversaire',
  baby_shower: 'Baby Shower',
  corporate: 'Événement corporatif',
  other: 'Autre'
}

const statusLabels = {
  pending: 'En attente',
  confirmed: 'Confirmée',
  paid: 'Payée',
  completed: 'Terminée',
  cancelled: 'Annulée'
}

function formatEventType(type) {
  return eventTypeLabels[type] || type
}

function formatStatus(status) {
  return statusLabels[status] || status
}

function formatDate(dateString) {
  const date = new Date(dateString)
  return date.toLocaleDateString('fr-FR', { year: 'numeric', month: 'long', day: 'numeric' })
}

async function fetchDashboardData() {
  loading.value = true
  try {
    const [statsResponse, bookingsResponse] = await Promise.all([
      api.get('/admin/dashboard/stats'),
      api.get('/admin/bookings', { params: { limit: 5 } })
    ])

    stats.value[0].value = statsResponse.data.total_bookings.toString()
    stats.value[1].value = statsResponse.data.pending_bookings.toString()
    stats.value[2].value = statsResponse.data.confirmed_bookings.toString()
    stats.value[3].value = statsResponse.data.total_clients.toString()

    recentBookings.value = bookingsResponse.data.slice(0, 5)
  } catch (error) {
    console.error('Failed to fetch dashboard data:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDashboardData()
})
</script>

<style scoped>
.admin-dashboard {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2xl);
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-xl);
}

.stat-card {
  background: var(--color-white);
  padding: var(--spacing-xl);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  transition: transform var(--transition-base);
}

.stat-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-lg);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-2xl);
  color: white;
}

.stat-content h3 {
  font-size: var(--font-size-3xl);
  color: var(--color-black-soft);
  margin-bottom: var(--spacing-xs);
}

.stat-content p {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
}

/* Dashboard Section */
.dashboard-section {
  background: var(--color-white);
  padding: var(--spacing-2xl);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}

.dashboard-section h2 {
  color: var(--color-gold);
  margin-bottom: var(--spacing-xl);
  font-size: var(--font-size-2xl);
}

/* Activity List */
.activity-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.activity-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  background: var(--color-gray-lighter);
}

.activity-icon {
  font-size: var(--font-size-2xl);
}

.activity-details {
  flex: 1;
}

.activity-details strong {
  color: var(--color-black-soft);
  display: block;
  margin-bottom: var(--spacing-xs);
}

.activity-details p {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
}

.activity-status {
  padding: var(--spacing-xs) var(--spacing-md);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
}

.activity-status.pending {
  background: rgba(247, 119, 0, 0.1);
  color: var(--color-warning);
}

.activity-status.confirmed,
.activity-status.paid {
  background: rgba(45, 122, 79, 0.1);
  color: var(--color-success);
}

/* Actions Grid */
.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-lg);
}

.action-card {
  background: linear-gradient(135deg, var(--color-gold) 0%, var(--color-gold-light) 100%);
  color: white;
  padding: var(--spacing-xl);
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  font-family: var(--font-sans);
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-medium);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-md);
  transition: transform var(--transition-base);
}

.action-card:hover {
  transform: translateY(-3px);
}

.action-icon {
  font-size: var(--font-size-3xl);
}

.loading-state,
.empty-state {
  text-align: center;
  padding: var(--spacing-2xl);
  color: var(--color-gray);
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
