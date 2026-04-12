<template>
  <div class="booking-page">
    <Header />

    <div class="booking-container">
      <div class="container container-narrow">
        <div class="booking-header fade-in-up">
          <h1 class="font-script text-gold">{{ t('booking.title') }}</h1>
          <p class="subtitle">{{ t('booking.subtitle') }}</p>
        </div>

        <div v-if="hasRentalSelection" class="selected-rentals-panel fade-in-up">
          <div class="selected-rentals-header">
            <div>
              <p class="cart-kicker">{{ t('booking.cart.kicker') }}</p>
              <h2>{{ t('booking.cart.title') }}</h2>
              <p>{{ t('booking.cart.count', { count: cartCount }) }}</p>
            </div>

            <Button variant="outline" size="sm" @click="goToRentals">
              {{ t('booking.cart.add_more') }}
            </Button>
          </div>

          <div class="selected-rentals-list">
            <div v-for="item in cartItems" :key="item.id" class="selected-rental-item">
              <img
                :src="getImageWithFallback(item.image_url, 'rental', item.category?.slug || item.category_enum)"
                :alt="item.title"
                class="rental-img"
                loading="lazy"
              />

              <div class="rental-details">
                <h4>{{ item.title }}</h4>
                <p>{{ item.category?.name || t('booking.cart.generic_item') }}</p>
              </div>

              <strong class="rental-price">{{ formatPrice(item.price) }}</strong>

              <button class="remove-rental-btn" type="button" @click="removeRentalItem(item.id)">
                {{ t('booking.cart.remove') }}
              </button>
            </div>
          </div>

          <div class="selected-rentals-footer">
            <span>{{ t('booking.cart.total') }}</span>
            <strong>{{ formatPrice(cartTotal) }}</strong>
          </div>
        </div>

        <div v-if="!bookingComplete" class="booking-wizard">
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

          <div v-show="currentStep === 0" class="step-content fade-in">
            <CalendarPicker
              v-model="formData.event_date"
              :title="t('booking.step1.title')"
              :subtitle="t('booking.step1.subtitle')"
              @date-selected="handleDateSelected"
            />

            <div class="step-actions">
              <Button size="lg" :disabled="!formData.event_date" @click="nextStep">
                {{ t('booking.step1.continue') }}
              </Button>
            </div>
          </div>

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
                    <option value="baptism">🕊️ {{ t('testimonials.types.baptism') }}</option>
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
                  :min="hasRentalSelection ? Math.ceil(cartTotal) : 500"
                  step="100"
                  required
                  placeholder="Ex: 5000"
                />
                <small v-if="hasRentalSelection" class="budget-info">
                  {{ t('booking.cart.budget_hint', { amount: formatPrice(cartTotal) }) }}
                </small>
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
              <Button size="lg" :disabled="!formData.event_type || !formData.budget" @click="nextStep">
                {{ t('booking.step1.continue') }}
              </Button>
            </div>
          </div>

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
                  placeholder="+1 (819) 244-4702"
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
              <Button size="lg" :disabled="!formData.name || !formData.email" @click="nextStep">
                {{ t('booking.step1.continue') }}
              </Button>
            </div>
          </div>

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

              <div v-if="hasRentalSelection" class="summary-section rental-card-summary">
                <h4>{{ t('booking.cart.summary_title') }}</h4>

                <div v-for="item in cartItems" :key="item.id" class="selected-rental-item compact">
                  <img
                    :src="getImageWithFallback(item.image_url, 'rental', item.category?.slug || item.category_enum)"
                    :alt="item.title"
                    class="rental-img"
                    loading="lazy"
                  />

                  <div class="rental-details">
                    <h4>{{ item.title }}</h4>
                    <p>{{ item.category?.name || t('booking.cart.generic_item') }}</p>
                  </div>

                  <strong class="rental-price">{{ formatPrice(item.price) }}</strong>
                </div>
              </div>

              <div class="summary-section payment-section">
                <h4>{{ t('booking.step4.section_payment') }}</h4>
                <div class="payment-details">
                  <div class="payment-row">
                    <span>{{ t('booking.step4.total') }}</span>
                    <strong>{{ formatPrice(formData.budget) }}</strong>
                  </div>

                  <div v-if="hasRentalSelection" class="payment-row">
                    <span>{{ t('booking.cart.items_total') }}</span>
                    <strong>{{ formatPrice(cartTotal) }}</strong>
                  </div>

                  <div v-if="hasRentalSelection" class="payment-row">
                    <span>{{ t('booking.cart.count', { count: cartCount }) }}</span>
                    <strong>{{ cartCount }}</strong>
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
              <Button size="lg" :loading="loading" @click="handleSubmit">
                {{ loading ? t('booking.step4.processing') : t('booking.step4.confirm') }}
              </Button>
            </div>
          </div>
        </div>

        <div v-else class="booking-success fade-in">
          <div class="success-icon">✓</div>
          <h2>{{ t('booking.success.title') }}</h2>
          <p>{{ t('booking.success.msg1') }}</p>
          <p>{{ t('booking.success.msg2', { email: formData.email }) }}</p>
          <p>{{ t('booking.success.msg3') }}</p>
          <Button size="lg" @click="router.push('/')">{{ t('booking.success.home') }}</Button>
        </div>
      </div>
    </div>

    <Footer />
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Header from '../../components/Header.vue'
import Footer from '../../components/Footer.vue'
import Button from '../../components/ui/Button.vue'
import CalendarPicker from '../../components/CalendarPicker.vue'
import api from '../../services/api'
import { getImageWithFallback } from '../../config/defaultImages'
import { useRentalCartStore } from '../../stores/rentalCart'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const rentalCartStore = useRentalCartStore()
const { cartItems, cartCount, cartItemIds, cartTotal } = storeToRefs(rentalCartStore)

const currentStep = ref(0)
const loading = ref(false)
const error = ref(null)
const bookingComplete = ref(false)

const steps = computed(() => [
  { id: 1, label: t('booking.steps.date') },
  { id: 2, label: t('booking.steps.event') },
  { id: 3, label: t('booking.steps.details') },
  { id: 4, label: t('booking.steps.confirm') },
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
  special_requests: '',
})

const hasRentalSelection = computed(() => cartCount.value > 0)

const depositAmount = computed(() => {
  return formData.value.budget ? (formData.value.budget * 0.3).toFixed(2) : 0
})

onMounted(async () => {
  try {
    const response = await api.get('/public/rentals')
    const rentals = response.data || []
    rentalCartStore.syncCatalog(rentals)

    if (route.query.rental_item_id) {
      const requestedId = Number.parseInt(route.query.rental_item_id, 10)
      const selectedItem = rentals.find((item) => item.id === requestedId)

      if (selectedItem) {
        rentalCartStore.addItem(selectedItem)
      }
    }
  } catch (err) {
    console.error('Failed to load rental items', err)
  }

  if (hasRentalSelection.value && !formData.value.event_type) {
    formData.value.event_type = 'other'
  }
})

watch(
  cartTotal,
  (total) => {
    if (total > 0 && (!formData.value.budget || formData.value.budget < total)) {
      formData.value.budget = Number(total.toFixed(2))
    }
  },
  { immediate: true }
)

function formatPrice(price) {
  const activeLocale = locale.value === 'en' ? 'en-CA' : 'fr-CA'
  return new Intl.NumberFormat(activeLocale, { style: 'currency', currency: 'CAD' }).format(price || 0)
}

function handleDateSelected(data) {
  formData.value.event_date = data.date
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
  if (!date) {
    return ''
  }

  const activeLocale = locale.value === 'en' ? 'en-CA' : 'fr-CA'
  return new Date(date).toLocaleDateString(activeLocale, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

function getEventTypeLabel(type) {
  const types = {
    proposal: t('testimonials.types.proposal'),
    wedding: t('testimonials.types.wedding'),
    birthday: t('testimonials.types.birthday'),
    baby_shower: t('testimonials.types.baby_shower'),
    baptism: t('testimonials.types.baptism'),
    corporate: t('testimonials.types.corporate'),
    other: t('testimonials.types.other'),
  }

  return types[type] || type
}

function removeRentalItem(itemId) {
  rentalCartStore.removeItem(itemId)
}

function goToRentals() {
  router.push({ name: 'rentals' })
}

async function handleSubmit() {
  loading.value = true
  error.value = null

  try {
    const payload = {
      ...formData.value,
      language: locale.value,
    }

    if (cartItemIds.value.length > 0) {
      payload.rental_item_ids = cartItemIds.value
    }

    await api.post('/public/bookings', payload)
    rentalCartStore.clearCart()
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

.selected-rentals-panel {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.96) 0%, rgba(247, 231, 206, 0.92) 100%);
  border: 1px solid rgba(212, 175, 55, 0.25);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
  padding: var(--spacing-2xl);
  margin-bottom: var(--spacing-3xl);
}

.selected-rentals-header {
  display: flex;
  justify-content: space-between;
  gap: var(--spacing-xl);
  align-items: flex-start;
  margin-bottom: var(--spacing-xl);
}

.cart-kicker {
  margin-bottom: var(--spacing-xs);
  text-transform: uppercase;
  letter-spacing: 0.14em;
  font-size: var(--font-size-xs);
  color: var(--color-gold);
}

.selected-rentals-header h2 {
  margin-bottom: var(--spacing-xs);
  font-size: var(--font-size-2xl);
}

.selected-rentals-header p {
  color: var(--color-gray);
  margin: 0;
}

.selected-rentals-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.selected-rentals-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-md);
  margin-top: var(--spacing-lg);
  padding-top: var(--spacing-lg);
  border-top: 1px solid rgba(212, 175, 55, 0.2);
  color: var(--color-gray);
}

.selected-rentals-footer strong {
  color: var(--color-gold);
  font-size: var(--font-size-xl);
}

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
  gap: var(--spacing-md);
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
  background: rgba(255, 255, 255, 0.86);
  border-radius: var(--radius-lg);
  padding: var(--spacing-md);
}

.selected-rental-item.compact {
  padding: var(--spacing-sm);
}

.rental-img {
  width: 100px;
  height: 100px;
  object-fit: cover;
  border-radius: var(--radius-md);
}

.selected-rental-item.compact .rental-img {
  width: 72px;
  height: 72px;
}

.rental-details h4 {
  font-weight: bold;
  color: var(--color-charcoal);
  margin-bottom: 4px;
}

.rental-details p {
  margin: 0;
  color: var(--color-gray);
}

.rental-price {
  color: var(--color-gold);
  font-weight: bold;
  margin-left: auto;
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

.step-actions {
  display: flex;
  gap: var(--spacing-lg);
  justify-content: center;
  margin-top: var(--spacing-2xl);
}

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
  .selected-rentals-header {
    flex-direction: column;
  }

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

  .selected-rental-item {
    align-items: flex-start;
    flex-wrap: wrap;
  }

  .rental-price,
  .remove-rental-btn {
    margin-left: 0;
  }

  .payment-row,
  .selected-rentals-footer {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
