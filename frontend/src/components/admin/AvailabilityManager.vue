<template>
  <div class="availability-manager">
    <div class="manager-header">
      <h2>Gestion des disponibilités</h2>
      <p>Définissez les dates où vous n'êtes pas disponible pour des événements</p>
    </div>

    <div class="manager-content">
      <!-- Calendar for selecting dates -->
      <div class="calendar-section">
        <h3>Calendrier</h3>
        <p class="section-subtitle">Cliquez sur les dates pour les marquer comme indisponibles</p>
        
        <VCalendar
          v-model="selectedDates"
          :attributes="calendarAttributes"
          :min-date="new Date()"
          mode="multiple"
          expanded
          borderless
        >
          <template #day-content="{ day, attributes }">
            <div class="calendar-day" :class="getDayClass(attributes)">
              <span class="day-number">{{ day.day }}</span>
            </div>
          </template>
        </VCalendar>
      </div>

      <!-- Availability List -->
      <div class="availability-list">
        <div class="list-header">
          <h3>Dates indisponibles</h3>
          <Button size="sm" @click="showAddModal = true">
            + Ajouter une plage
          </Button>
        </div>

        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>Chargement...</p>
        </div>

        <div v-else-if="unavailableDates.length === 0" class="empty-state">
          <p>Aucune date indisponible définie</p>
          <p class="hint">Toutes les dates futures sont disponibles par défaut</p>
        </div>

        <div v-else class="dates-list">
          <div
            v-for="avail in unavailableDates"
            :key="avail.id"
            class="date-item"
          >
            <div class="date-info">
              <span class="date-icon">📅</span>
              <div>
                <strong>{{ formatDate(avail.date) }}</strong>
                <p v-if="avail.reason">{{ avail.reason }}</p>
              </div>
            </div>
            <button class="delete-btn" @click="deleteAvailability(avail.id)">
              🗑️
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Availability Modal -->
    <Modal v-model="showAddModal">
      <template #header>
        <h2>Ajouter une indisponibilité</h2>
      </template>

      <form @submit.prevent="handleAddAvailability" class="add-form">
        <div class="form-group">
          <label>Type</label>
          <div class="radio-group">
            <label class="radio-label">
              <input v-model="newAvailability.type" type="radio" value="single" />
              <span>Date unique</span>
            </label>
            <label class="radio-label">
              <input v-model="newAvailability.type" type="radio" value="range" />
              <span>Plage de dates</span>
            </label>
          </div>
        </div>

        <div v-if="newAvailability.type === 'single'" class="form-group">
          <label for="single-date">Date</label>
          <input
            id="single-date"
            v-model="newAvailability.date"
            type="date"
            required
            :min="minDate"
          />
        </div>

        <div v-else class="form-row">
          <div class="form-group">
            <label for="start-date">Date de début</label>
            <input
              id="start-date"
              v-model="newAvailability.startDate"
              type="date"
              required
              :min="minDate"
            />
          </div>
          <div class="form-group">
            <label for="end-date">Date de fin</label>
            <input
              id="end-date"
              v-model="newAvailability.endDate"
              type="date"
              required
              :min="newAvailability.startDate"
            />
          </div>
        </div>

        <div class="form-group">
          <label for="reason">Raison (optionnel)</label>
          <input
            id="reason"
            v-model="newAvailability.reason"
            type="text"
            placeholder="Ex: Vacances, autre événement..."
          />
        </div>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>
      </form>

      <template #footer>
        <Button variant="ghost" @click="showAddModal = false">
          Annuler
        </Button>
        <Button :loading="submitting" @click="handleAddAvailability">
          Ajouter
        </Button>
      </template>
    </Modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Calendar as VCalendar } from 'v-calendar'
import 'v-calendar/style.css'
import Button from '../ui/Button.vue'
import Modal from '../ui/Modal.vue'
import api from '../../services/api'

const selectedDates = ref([])
const unavailableDates = ref([])
const loading = ref(false)
const submitting = ref(false)
const error = ref(null)
const showAddModal = ref(false)

const newAvailability = ref({
  type: 'single',
  date: '',
  startDate: '',
  endDate: '',
  reason: ''
})

const minDate = computed(() => {
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  return tomorrow.toISOString().split('T')[0]
})

const calendarAttributes = computed(() => {
  const attrs = []
  
  // Mark unavailable dates in red
  unavailableDates.value.forEach(avail => {
    attrs.push({
      key: `unavailable-${avail.id}`,
      dates: new Date(avail.date),
      dot: {
        color: 'red',
        class: 'unavailable-dot'
      },
      popover: {
        label: avail.reason || 'Indisponible'
      }
    })
  })
  
  return attrs
})

function getDayClass(attributes) {
  const isUnavailable = attributes.some(attr => 
    attr.key.startsWith('unavailable-')
  )
  
  return isUnavailable ? 'unavailable' : 'available'
}

function formatDate(dateStr) {
  const date = new Date(dateStr)
  return date.toLocaleDateString('fr-FR', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

async function fetchAvailabilities() {
  loading.value = true
  try {
    const response = await api.get('/admin/availabilities')
    unavailableDates.value = response.data.filter(a => !a.available)
  } catch (err) {
    console.error('Failed to fetch availabilities:', err)
  } finally {
    loading.value = false
  }
}

async function handleAddAvailability() {
  submitting.value = true
  error.value = null

  try {
    if (newAvailability.value.type === 'single') {
      // Add single date
      await api.post('/admin/availabilities', {
        date: newAvailability.value.date,
        available: false,
        reason: newAvailability.value.reason
      })
    } else {
      // Add date range
      const start = new Date(newAvailability.value.startDate)
      const end = new Date(newAvailability.value.endDate)
      
      const promises = []
      for (let d = new Date(start); d <= end; d.setDate(d.getDate() + 1)) {
        const dateStr = d.toISOString().split('T')[0]
        promises.push(
          api.post('/admin/availabilities', {
            date: dateStr,
            available: false,
            reason: newAvailability.value.reason
          })
        )
      }
      
      await Promise.all(promises)
    }

    showAddModal.value = false
    newAvailability.value = {
      type: 'single',
      date: '',
      startDate: '',
      endDate: '',
      reason: ''
    }
    
    await fetchAvailabilities()
  } catch (err) {
    error.value = err.response?.data?.error || 'Une erreur est survenue'
  } finally {
    submitting.value = false
  }
}

async function deleteAvailability(id) {
  if (!confirm('Êtes-vous sûr de vouloir supprimer cette indisponibilité?')) {
    return
  }

  try {
    await api.delete(`/admin/availabilities/${id}`)
    await fetchAvailabilities()
  } catch (err) {
    console.error('Failed to delete availability:', err)
  }
}

onMounted(() => {
  fetchAvailabilities()
})
</script>

<style scoped>
.availability-manager {
  background: white;
  border-radius: var(--radius-lg);
  padding: var(--spacing-2xl);
}

.manager-header {
  margin-bottom: var(--spacing-3xl);
}

.manager-header h2 {
  color: var(--color-gold);
  font-size: var(--font-size-2xl);
  margin-bottom: var(--spacing-sm);
}

.manager-header p {
  color: var(--color-gray);
}

.manager-content {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: var(--spacing-3xl);
}

.calendar-section h3,
.availability-list h3 {
  color: var(--color-black-soft);
  margin-bottom: var(--spacing-md);
}

.section-subtitle {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
  margin-bottom: var(--spacing-lg);
}

/* Calendar Styles */
:deep(.vc-container) {
  border: 2px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
}

.calendar-day {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  transition: all var(--transition-base);
}

.calendar-day.available:hover {
  background: rgba(45, 122, 79, 0.1);
  cursor: pointer;
}

.calendar-day.unavailable {
  background: rgba(193, 41, 46, 0.1);
  color: var(--color-error);
}

/* Availability List */
.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.dates-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  max-height: 500px;
  overflow-y: auto;
}

.date-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  background: var(--color-gray-lighter);
  border-radius: var(--radius-md);
  transition: background var(--transition-base);
}

.date-item:hover {
  background: rgba(193, 41, 46, 0.1);
}

.date-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.date-icon {
  font-size: var(--font-size-xl);
}

.date-info strong {
  display: block;
  color: var(--color-black-soft);
  margin-bottom: var(--spacing-xs);
}

.date-info p {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
}

.delete-btn {
  background: none;
  border: none;
  font-size: var(--font-size-xl);
  cursor: pointer;
  padding: var(--spacing-sm);
  transition: transform var(--transition-base);
}

.delete-btn:hover {
  transform: scale(1.2);
}

/* Add Form */
.add-form {
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
  font-weight: var(--font-weight-medium);
  color: var(--color-text);
}

.form-group input {
  padding: var(--spacing-md);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-sans);
  font-size: var(--font-size-base);
}

.form-group input:focus {
  outline: none;
  border-color: var(--color-gold);
}

.radio-group {
  display: flex;
  gap: var(--spacing-lg);
}

.radio-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
}

.radio-label input[type="radio"] {
  width: auto;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: var(--spacing-3xl);
  color: var(--color-gray);
}

.hint {
  font-size: var(--font-size-sm);
  margin-top: var(--spacing-sm);
}

.error-message {
  background: rgba(193, 41, 46, 0.1);
  color: var(--color-error);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  text-align: center;
}

@media (max-width: 968px) {
  .manager-content {
    grid-template-columns: 1fr;
  }
}
</style>
