<template>
  <div class="calendar-picker">
    <div class="calendar-header">
      <h3>{{ title }}</h3>
      <p v-if="subtitle">{{ subtitle }}</p>
    </div>

    <VCalendar
      v-model="selectedDate"
      :attributes="calendarAttributes"
      :min-date="minDate"
      @dayclick="handleDayClick"
      expanded
      borderless
      :is-dark="false"
    />

    <div v-if="selectedDate" class="selected-date-info">
      <div class="info-icon">📅</div>
      <div>
        <strong>Date sélectionnée</strong>
        <p>{{ formatSelectedDate }}</p>
      </div>
    </div>

    <div v-if="loading" class="calendar-loading">
      <div class="spinner"></div>
      <p>Vérification de la disponibilité...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { Calendar as VCalendar } from 'v-calendar'
import 'v-calendar/style.css'
import api from '../services/api'

const props = defineProps({
  title: {
    type: String,
    default: 'Sélectionnez une date'
  },
  subtitle: {
    type: String,
    default: 'Les dates en vert sont disponibles'
  },
  modelValue: {
    type: [Date, String],
    default: null
  }
})

const emit = defineEmits(['update:modelValue', 'date-selected'])

const selectedDate = ref(props.modelValue ? new Date(props.modelValue) : null)
const availabilities = ref([])
const loading = ref(false)

const minDate = computed(() => {
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  return tomorrow
})

const calendarAttributes = computed(() => {
  const attrs = []
  
  // Available dates (green dots)
  availabilities.value.forEach(avail => {
    if (avail.available) {
      attrs.push({
        key: `available-${avail.date}`,
        dates: new Date(avail.date),
        dot: {
          color: 'green',
          class: 'available-dot'
        }
      })
    }
  })
  
  // Selected date highlight
  if (selectedDate.value) {
    attrs.push({
      key: 'selected',
      dates: selectedDate.value,
      highlight: {
        color: 'blue',
        fillMode: 'solid'
      }
    })
  }
  
  return attrs
})

const formatSelectedDate = computed(() => {
  if (!selectedDate.value) return ''
  return new Date(selectedDate.value).toLocaleDateString('fr-FR', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

function formatDateForAPI(date) {
  const d = new Date(date)
  return d.toISOString().split('T')[0]
}

function handleDayClick(day) {
  const clickedDate = day.date
  const dateStr = formatDateForAPI(clickedDate)
  
  // Check if date is in the past
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  if (clickedDate < today) {
    return
  }
  
  // Set the selected date
  selectedDate.value = clickedDate
  emit('update:modelValue', dateStr)
  emit('date-selected', {
    date: dateStr,
    formatted: formatSelectedDate.value,
    available: true
  })
}

async function fetchAvailabilities(month) {
  try {
    const monthStr = month.toISOString().slice(0, 7)
    const response = await api.get('/public/availabilities', {
      params: { month: monthStr }
    })
    
    if (response.data.length === 0) {
      // Generate available dates for the next 3 months
      const dates = []
      const today = new Date()
      for (let i = 1; i <= 90; i++) {
        const date = new Date(today)
        date.setDate(date.getDate() + i)
        dates.push({
          date: formatDateForAPI(date),
          available: true,
          max_events: 1,
          booked: 0
        })
      }
      availabilities.value = dates
    } else {
      availabilities.value = response.data
    }
  } catch (error) {
    console.error('Failed to fetch availabilities:', error)
    // Fallback: make all future dates available
    const dates = []
    const today = new Date()
    for (let i = 1; i <= 90; i++) {
      const date = new Date(today)
      date.setDate(date.getDate() + i)
      dates.push({
        date: formatDateForAPI(date),
        available: true,
        max_events: 1,
        booked: 0
      })
    }
    availabilities.value = dates
  }
}

watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    selectedDate.value = new Date(newVal)
  }
})

watch(selectedDate, (newVal) => {
  if (newVal) {
    const dateStr = formatDateForAPI(newVal)
    emit('update:modelValue', dateStr)
  }
})

onMounted(() => {
  fetchAvailabilities(new Date())
})
</script>

<style scoped>
.calendar-picker {
  background: white;
  border-radius: var(--radius-xl);
  padding: var(--spacing-2xl);
  box-shadow: var(--shadow-lg);
}

.calendar-header {
  text-align: center;
  margin-bottom: var(--spacing-xl);
}

.calendar-header h3 {
  color: var(--color-gold);
  font-size: var(--font-size-2xl);
  margin-bottom: var(--spacing-sm);
}

.calendar-header p {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
}

/* V-Calendar Custom Styles */
:deep(.vc-container) {
  border: none;
  font-family: var(--font-sans);
}

:deep(.vc-header) {
  padding: var(--spacing-lg) 0;
  margin-bottom: var(--spacing-md);
}

:deep(.vc-title) {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-black-soft);
}

:deep(.vc-arrow) {
  color: var(--color-gold);
  border-radius: var(--radius-sm);
  transition: background var(--transition-base);
}

:deep(.vc-arrow:hover) {
  background: rgba(212, 175, 55, 0.1);
}

:deep(.vc-weeks) {
  padding: 0;
}

:deep(.vc-weekday) {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  padding: var(--spacing-sm);
}

:deep(.vc-day) {
  min-height: 60px;
  position: relative;
}

:deep(.vc-day-content) {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.calendar-day {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  transition: all var(--transition-base);
  position: relative;
  padding: var(--spacing-sm);
}

.day-number {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-medium);
}

.booking-indicator {
  position: absolute;
  bottom: 4px;
  font-size: var(--font-size-xs);
  color: var(--color-gold);
}

.calendar-day.available {
  background: rgba(45, 122, 79, 0.1);
  color: var(--color-success);
  cursor: pointer;
}

.calendar-day.available:hover {
  background: rgba(45, 122, 79, 0.2);
  transform: scale(1.05);
}

.calendar-day.unavailable {
  background: rgba(193, 41, 46, 0.1);
  color: var(--color-error);
  cursor: not-allowed;
  opacity: 0.5;
}

:deep(.vc-day.is-disabled) {
  opacity: 0.3;
  pointer-events: none;
}

:deep(.vc-highlight) {
  border-radius: var(--radius-md);
}

.selected-date-info {
  margin-top: var(--spacing-xl);
  padding: var(--spacing-lg);
  background: linear-gradient(135deg, rgba(212, 175, 55, 0.1) 0%, rgba(229, 193, 88, 0.1) 100%);
  border-radius: var(--radius-lg);
  border: 2px solid var(--color-gold);
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.info-icon {
  font-size: var(--font-size-3xl);
}

.selected-date-info strong {
  display: block;
  color: var(--color-gold);
  margin-bottom: var(--spacing-xs);
  font-size: var(--font-size-lg);
}

.selected-date-info p {
  color: var(--color-text);
  font-size: var(--font-size-base);
}

.calendar-loading {
  text-align: center;
  padding: var(--spacing-lg);
  color: var(--color-gray);
}

@media (max-width: 768px) {
  .calendar-picker {
    padding: var(--spacing-lg);
  }
  
  :deep(.vc-day) {
    min-height: 50px;
  }
}
</style>
