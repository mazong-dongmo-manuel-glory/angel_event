<template>
  <div class="booking-page">
    <Header />
    
    <div class="booking-container">
      <div class="container container-narrow">
        <div class="booking-header fade-in-up">
          <h1 class="font-script text-gold">{{ t('booking.title') }}</h1>
          <p class="subtitle">{{ t('booking.subtitle') }}</p>
        </div>

        <div v-if="!bookingComplete" class="booking-wizard">
          <!-- Step Indicator -->
          <div class="steps-indicator">
            <div
              v-for="(step, index) in steps"
              :key="step.id"
              class="step"
              :class="{ active: currentStep === index, completed: currentStep > index }"
            >
              <div class="step-number">
                <span v-if="currentStep > index">✓</span>
                <span v-else>{{ index + 1 }}</span>
              </div>
              <span class="step-label">{{ step.label }}</span>
            </div>
          </div>

          <!-- Step 1: Date Selection -->
          <div v-show="currentStep === 0" class="step-content fade-in">
            <CalendarPicker
              v-model="formData.event_date"
              :title="t('booking.step1.title')"
              :subtitle="t('booking.step1.subtitle')"
              @date-selected="handleDateSelected"
            />
            
            <div class="step-actions">
              <Button
                size="lg"
                :disabled="!formData.event_date"
                @click="nextStep"
              >
                {{ t('booking.step1.continue') }}
              </Button>
            </div>
          </div>

          <!-- Step 2: Event Details -->
          <div v-show="currentStep === 1" class="step-content fade-in">
            <div class="form-card">
              <h3>{{ t('booking.step2.title') }}</h3>
              
              <div class="form-row">
                <div class="form-group">
                  <label for="event-type">{{ t('booking.step2.type') }}</label>
                  <select id="event-type" v-model="formData.event_type" required>
                    <option value="">{{ t('booking.step2.select') }}</option>
                    <option value="proposal">💍 {{ t('testimonials.types.proposal') }}</option>
                    <option value="wedding">💐 {{ t('testimonials.types.wedding') }}</option>
                    <option value="birthday">🎂 {{ t('testimonials.types.birthday') }}</option>
                    <option value="baby_shower">👶 {{ t('testimonials.types.baby_shower') }}</option>
                    <option value="corporate">🏢 {{ t('testimonials.types.corporate') }}</option>
                    <option value="other">✨ {{ t('testimonials.types.other') }}</option>
                  </select>
                </div>
                <div class="form-group">
                  <label for="guests">{{ t('booking.step2.guests') }}</label>
                  <input
                    id="guests"
                    v-model.number="formData.guest_count"
                    type="number"
                    min="1"
                    placeholder="Ex: 50"
                  />
                </div>
              </div>

              <div class="form-group">
                <label for="location">{{ t('booking.step2.location') }}</label>
                <input
                  id="location"
                  v-model="formData.event_location"
                  type="text"
                  :placeholder="t('booking.step2.location_ph')"
                />
              </div>

              <div class="form-group">
                <label for="budget">{{ t('booking.step2.budget') }}</label>
                <input
                  id="budget"
                  v-model.number="formData.budget"
                  type="number"
                  min="500"
                  step="100"
                  required
                  placeholder="Ex: 5000"
                />
                <small v-if="depositAmount > 0" class="budget-info">
                  {{ t('booking.step2.deposit', { amount: depositAmount }) }}
                </small>
              </div>

              <div class="form-group">
                <label for="message">{{ t('booking.step2.vision') }}</label>
                <textarea
                  id="message"
                  v-model="formData.message"
                  rows="4"
                  :placeholder="t('booking.step2.vision_ph')"
                ></textarea>
              </div>
            </div>

            <div class="step-actions">
              <Button variant="ghost" @click="prevStep">
                {{ t('gallery.prev') || 'Retour' }}
              </Button>
              <Button
                size="lg"
                :disabled="!formData.event_type || !formData.budget"
                @click="nextStep"
              >
                {{ t('booking.step1.continue') }}
              </Button>
            </div>
          </div>

          <!-- Step 3: Personal Information -->
          <div v-show="currentStep === 2" class="step-content fade-in">
            <div class="form-card">
              <h3>{{ t('booking.step3.title') }}</h3>
              
              <div class="form-row">
                <div class="form-group">
                  <label for="name">{{ t('booking.step3.name') }}</label>
                  <input
                    id="name"
                    v-model="formData.name"
                    type="text"
                    required
                    placeholder="Jean Dupont"
                  />
                </div>
                <div class="form-group">
                  <label for="email">{{ t('booking.step3.email') }}</label>
                  <input
                    id="email"
                    v-model="formData.email"
                    type="email"
                    required
                    placeholder="jean@example.com"
                  />
                </div>
              </div>

              <div class="form-group">
                <label for="phone">{{ t('booking.step3.phone') }}</label>
                <input
                  id="phone"
                  v-model="formData.phone"
                  type="tel"
                  placeholder="+1 (555) 123-4567"
                />
              </div>

              <div class="form-group">
                <label for="special-requests">{{ t('booking.step3.special') }}</label>
                <textarea
                  id="special-requests"
                  v-model="formData.special_requests"
                  rows="3"
                  :placeholder="t('booking.step3.special_ph')"
                ></textarea>
              </div>
            </div>

            <div class="step-actions">
              <Button variant="ghost" @click="prevStep">
                {{ t('gallery.prev') || 'Retour' }}
              </Button>
              <Button
                size="lg"
                :disabled="!formData.name || !formData.email"
                @click="nextStep"
              >
                {{ t('booking.step1.continue') }}
              </Button>
            </div>
          </div>

          <!-- Step 4: Confirmation & Payment -->
          <div v-show="currentStep === 3" class="step-content fade-in">
            <div class="booking-summary">
              <h3>{{ t('booking.step4.title') }}</h3>
              
              <div class="summary-section">
                <h4>{{ t('booking.step4.section_date') }}</h4>
                <p><strong>{{ t('booking.steps.date') }}:</strong> {{ formatDate(formData.event_date) }}</p>
                <p><strong>{{ t('booking.steps.event') }}:</strong> {{ getEventTypeLabel(formData.event_type) }}</p>
                <p v-if="formData.event_location"><strong>{{ t('booking.step2.location') }}:</strong> {{ formData.event_location }}</p>
                <p v-if="formData.guest_count"><strong>{{ t('booking.step2.guests') }}:</strong> {{ formData.guest_count }}</p>
              </div>

              <div class="summary-section">
                <h4>{{ t('booking.step4.section_contact') }}</h4>
                <p><strong>{{ t('booking.step3.name') }}:</strong> {{ formData.name }}</p>
                <p><strong>{{ t('booking.step3.email') }}:</strong> {{ formData.email }}</p>
                <p v-if="formData.phone"><strong>{{ t('booking.step3.phone') }}:</strong> {{ formData.phone }}</p>
              </div>

              <div class="summary-section payment-section">
                <h4>{{ t('booking.step4.section_payment') }}</h4>
                <div class="payment-details">
                  <div class="payment-row">
                    <span>{{ t('booking.step4.total') }}</span>
                    <strong>{{ formData.budget }} CAD</strong>
                  </div>

                  
                  <div v-if="selectedRentalItem" class="payment-row">
                    <span>{{ t('booking.step4.item_prefix') }} ({{ selectedRentalItem.title }}):</span>
                    <strong>{{ formatPrice(selectedRentalItem.price) }}</strong>
                  </div>
                </div>
                <p class="payment-note">
                  {{ t('booking.step4.note') }}
                </p>
              </div>
            </div>

            <div v-if="error" class="error-message">
              {{ error }}
            </div>

            <div class="step-actions">
              <Button variant="ghost" @click="prevStep">
                {{ t('gallery.prev') || 'Retour' }}
              </Button>
              <Button
                size="lg"
                :loading="loading"
                @click="handleSubmit"
              >
                {{ loading ? t('booking.step4.processing') : t('booking.step4.confirm') }}
              </Button>
            </div>
          </div>
        </div>

        <!-- Success Message -->
        <div v-else class="booking-success fade-in">
          <div class="success-icon">✓</div>
          <h2>{{ t('booking.success.title') }}</h2>
          <p>{{ t('booking.success.msg1') }}</p>
          <p>{{ t('booking.success.msg2', { email: formData.email }) }}</p>
          <p>{{ t('booking.success.msg3') }}</p>
          <Button size="lg" @click="$router.push('/')">{{ t('booking.success.home') }}</Button>
        </div>
      </div>
    </div>

    <Footer />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Header from '../../components/Header.vue'
import Footer from '../../components/Footer.vue'
import Button from '../../components/ui/Button.vue'
import CalendarPicker from '../../components/CalendarPicker.vue'
import api from '../../services/api'

const { t } = useI18n()
const currentStep = ref(0)
const loading = ref(false)
const error = ref(null)
const bookingComplete = ref(false)
const selectedRentalItem = ref(null)

const steps = computed(() => [
  { id: 1, label: t('booking.steps.date') },
  { id: 2, label: t('booking.steps.event') },
  { id: 3, label: t('booking.steps.details') },
  { id: 4, label: t('booking.steps.confirm') }
])

const formData = ref({
  name: '',
  email: '',
  phone: '',
  event_type: '',
  event_date: '',
  event_location: '',
  guest_count: null,
  budget: null,
  message: '',
  special_requests: ''
})

const depositAmount = computed(() => {
  return formData.value.budget ? (formData.value.budget * 0.3).toFixed(2) : 0
})

const route = useRoute()

onMounted(async () => {
  if (route.query.rental_item_id) {
    try {
      // In a real app we'd fetch the item. For now, we'll need an endpoint or store.
      // Since we don't have a public "get single item" endpoint, we'll fetch all and find it.
      // Optimization: Add GET /public/rentals/:id later
      const res = await api.get('/public/rentals')
      selectedRentalItem.value = res.data.find(i => i.id === parseInt(route.query.rental_item_id))
    } catch (err) {
      console.error('Failed to load rental item', err)
    }
  }
})

function formatPrice(price) {
  return new Intl.NumberFormat(t.locale || 'fr-CA', { style: 'currency', currency: 'CAD' }).format(price)
}

function handleDateSelected(data) {
  console.log('Date selected:', data)
  formData.value.event_date = data.date
  console.log('formData.event_date:', formData.value.event_date)
}

function nextStep() {
  if (currentStep.value < steps.value.length - 1) {
    currentStep.value++
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

function prevStep() {
  if (currentStep.value > 0) {
    currentStep.value--
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

function formatDate(date) {
  if (!date) return ''
  return new Date(date).toLocaleDateString(t.locale || 'fr-CA', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

function getEventTypeLabel(type) {
  const types = {
    proposal: t('testimonials.types.proposal'),
    wedding: t('testimonials.types.wedding'),
    birthday: t('testimonials.types.birthday'),
    baby_shower: t('testimonials.types.baby_shower'),
    corporate: t('testimonials.types.corporate'),
    other: t('testimonials.types.other')
  }
  return types[type] || type
}

async function handleSubmit() {
  loading.value = true
  error.value = null

  try {
    // Create booking
    const payload = { ...formData.value, language: t.locale.value }
    if (selectedRentalItem.value) {
      payload.rental_item_ids = [selectedRentalItem.value.id]
    }
    await api.post('/public/bookings', payload)
    bookingComplete.value = true
  } catch (err) {
    error.value = err.response?.data?.error || 'Une erreur est survenue. Veuillez réessayer.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.booking-page {
  min-height: 100vh;
  background: var(--color-ivory);
}

.booking-container {
  padding: calc(80px + var(--spacing-5xl)) 0 var(--spacing-5xl);
}

.booking-header {
  text-align: center;
  margin-bottom: var(--spacing-4xl);
}

.booking-header h1 {
  font-size: var(--font-size-5xl);
  margin-bottom: var(--spacing-md);
}

.subtitle {
  font-size: var(--font-size-xl);
  color: var(--color-gray);
}

/* Steps Indicator */
.steps-indicator {
  display: flex;
  justify-content: space-between;
  margin-bottom: var(--spacing-4xl);
  position: relative;
}

.steps-indicator::before {
  content: '';
  position: absolute;
  top: 20px;
  left: 10%;
  right: 10%;
  height: 2px;
  background: var(--color-border);
  z-index: 0;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-sm);
  position: relative;
  z-index: 1;
  flex: 1;
}

.step-number {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--color-white);
  border: 2px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-semibold);
  color: var(--color-gray);
  transition: all var(--transition-base);
}

.step.active .step-number {
  background: var(--color-gold);
  border-color: var(--color-gold);
  color: white;
  transform: scale(1.1);
}

.step.completed .step-number {
  background: var(--color-success);
  border-color: var(--color-success);
  color: white;
}

.step-label {
  font-size: var(--font-size-sm);
  color: var(--color-gray);
  font-weight: var(--font-weight-medium);
}

.step.active .step-label {
  color: var(--color-gold);
}

/* Step Content */
.step-content {
  animation: fadeIn 0.3s ease-in;
}

.form-card {
  background: white;
  padding: var(--spacing-3xl);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
  margin-bottom: var(--spacing-2xl);
}

.form-card h3 {
  color: var(--color-gold);
  font-size: var(--font-size-2xl);
  margin-bottom: var(--spacing-xl);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.form-group label {
  font-weight: var(--font-weight-medium);
  color: var(--color-text);
  font-size: var(--font-size-sm);
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

.form-group small {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
}

.budget-info {
  padding: var(--spacing-sm);
  background: rgba(212, 175, 55, 0.1);
  border-radius: var(--radius-sm);
  color: var(--color-gold);
}

/* Booking Summary */
.booking-summary {
  background: white;
  padding: var(--spacing-3xl);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
  margin-bottom: var(--spacing-2xl);
}

.booking-summary h3 {
  color: var(--color-gold);
  font-size: var(--font-size-2xl);
  margin-bottom: var(--spacing-2xl);
  text-align: center;
}

.summary-section {
  padding: var(--spacing-xl);
  background: var(--color-gray-lighter);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-lg);
}

.summary-section h4 {
  color: var(--color-gold);
  margin-bottom: var(--spacing-md);
  font-size: var(--font-size-lg);
}

.summary-section p {
  color: var(--color-text);
  margin-bottom: var(--spacing-sm);
  line-height: var(--line-height-relaxed);
}

.payment-section {
  background: linear-gradient(135deg, rgba(212, 175, 55, 0.1) 0%, rgba(229, 193, 88, 0.1) 100%);
  border: 2px solid var(--color-gold);
}

.payment-details {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.payment-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) 0;
}

.payment-row.highlight {
  padding: var(--spacing-md);
  background: white;
  border-radius: var(--radius-md);
  font-size: var(--font-size-lg);
}

.payment-note {
  font-size: var(--font-size-sm);
  color: var(--color-gray);
  text-align: center;
  margin-top: var(--spacing-md);
}

.rental-card-summary {
  background: var(--color-white);
  border: 2px solid var(--color-gold);
}

.selected-rental-item {
  display: flex;
  gap: var(--spacing-lg);
  align-items: center;
}

.rental-img {
  width: 100px;
  height: 100px;
  object-fit: cover;
  border-radius: var(--radius-md);
}

.rental-details h4 {
  font-weight: bold;
  color: var(--color-charcoal);
}

.rental-price {
  color: var(--color-gold);
  font-weight: bold;
}

.remove-rental-btn {
  background: none;
  border: 1px solid var(--color-error);
  color: var(--color-error);
  padding: 0.5rem 1rem;
  border-radius: var(--radius-sm);
  cursor: pointer;
  margin-left: auto;
}

/* Step Actions */
.step-actions {
  display: flex;
  gap: var(--spacing-lg);
  justify-content: center;
  margin-top: var(--spacing-2xl);
}

/* Success State */
.booking-success {
  text-align: center;
  background: white;
  padding: var(--spacing-5xl) var(--spacing-3xl);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
}

.success-icon {
  width: 100px;
  height: 100px;
  background: var(--color-success);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-5xl);
  margin: 0 auto var(--spacing-2xl);
  animation: scaleIn 0.5s ease-out;
}

.booking-success h2 {
  color: var(--color-gold);
  margin-bottom: var(--spacing-lg);
  font-size: var(--font-size-3xl);
}

.booking-success p {
  color: var(--color-gray);
  margin-bottom: var(--spacing-md);
  line-height: var(--line-height-relaxed);
  font-size: var(--font-size-lg);
}

.error-message {
  background: rgba(193, 41, 46, 0.1);
  color: var(--color-error);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-lg);
  text-align: center;
}

@media (max-width: 768px) {
  .steps-indicator {
    flex-wrap: wrap;
  }
  
  .step-label {
    font-size: var(--font-size-xs);
  }
  
  .form-card {
    padding: var(--spacing-xl);
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .step-actions {
    flex-direction: column;
  }
}
</style>
