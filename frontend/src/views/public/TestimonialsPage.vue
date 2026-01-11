<template>
  <div class="testimonials-page">
    <Header />
    <div class="container testimonials-container">
      <div class="testimonials-header fade-in-up">
        <h1 class="font-script text-gold">Témoignages</h1>
        <p class="subtitle">Ce que nos clients disent de nous</p>
      </div>

      <!-- Submit Testimonial CTA -->
      <div class="submit-cta fade-in-up">
        <p>Vous avez travaillé avec nous?</p>
        <Button @click="showSubmitForm = true">Partager votre expérience</Button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Chargement des témoignages...</p>
      </div>

      <!-- Testimonials Grid -->
      <div v-else class="testimonials-grid">
        <Card v-for="testimonial in testimonials" :key="testimonial.id" class="testimonial-card fade-in-up" glass>
          <div class="testimonial-rating">
            <span v-for="n in testimonial.rating" :key="n" class="star filled">★</span>
            <span v-for="n in (5 - testimonial.rating)" :key="`empty-${n}`" class="star">★</span>
          </div>
          <p class="testimonial-content">"{{ testimonial.content }}"</p>
          <div class="testimonial-author">
            <strong>{{ testimonial.name }}</strong>
            <span class="event-type">{{ formatEventType(testimonial.event_type) }}</span>
          </div>
        </Card>
      </div>

      <!-- Submit Form Modal -->
      <Modal v-model="showSubmitForm">
        <template #header>
          <h2>Partager votre témoignage</h2>
        </template>

        <form v-if="!submitted" @submit.prevent="handleSubmit" class="testimonial-form">
          <div class="form-group">
            <label for="name">Votre nom *</label>
            <input
              id="name"
              v-model="formData.name"
              type="text"
              required
              placeholder="Jean Dupont"
            />
          </div>

          <div class="form-group">
            <label for="event-type">Type d'événement</label>
            <select id="event-type" v-model="formData.event_type">
              <option value="">Sélectionnez...</option>
              <option value="proposal">Demande en mariage</option>
              <option value="wedding">Mariage</option>
              <option value="birthday">Anniversaire</option>
              <option value="baby_shower">Baby Shower</option>
              <option value="other">Autre</option>
            </select>
          </div>

          <div class="form-group">
            <label>Note *</label>
            <div class="rating-input">
              <span
                v-for="n in 5"
                :key="n"
                class="star-input"
                :class="{ filled: formData.rating >= n }"
                @click="formData.rating = n"
              >
                ★
              </span>
            </div>
          </div>

          <div class="form-group">
            <label for="content">Votre témoignage *</label>
            <textarea
              id="content"
              v-model="formData.content"
              rows="5"
              required
              placeholder="Partagez votre expérience avec Angel Event..."
            ></textarea>
          </div>

          <div v-if="error" class="error-message">
            {{ error }}
          </div>
        </form>

        <div v-else class="success-message">
          <div class="success-icon">✓</div>
          <p>Merci pour votre témoignage! Il sera publié après modération.</p>
        </div>

        <template #footer>
          <Button v-if="!submitted" variant="ghost" @click="showSubmitForm = false">
            Annuler
          </Button>
          <Button v-if="!submitted" type="submit" :loading="submitting" @click="handleSubmit">
            Envoyer
          </Button>
          <Button v-else @click="showSubmitForm = false; submitted = false">
            Fermer
          </Button>
        </template>
      </Modal>
    </div>
    <Footer />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Header from '../../components/Header.vue'
import Footer from '../../components/Footer.vue'
import Button from '../../components/ui/Button.vue'
import Card from '../../components/ui/Card.vue'
import Modal from '../../components/ui/Modal.vue'
import api from '../../services/api'

const testimonials = ref([])
const loading = ref(true)
const showSubmitForm = ref(false)
const submitting = ref(false)
const submitted = ref(false)
const error = ref(null)

const formData = ref({
  name: '',
  event_type: '',
  rating: 5,
  content: ''
})

const eventTypeLabels = {
  proposal: 'Demande en mariage',
  wedding: 'Mariage',
  birthday: 'Anniversaire',
  baby_shower: 'Baby Shower',
  corporate: 'Événement corporatif',
  other: 'Autre'
}

function formatEventType(type) {
  return eventTypeLabels[type] || type
}

async function fetchTestimonials() {
  loading.value = true
  try {
    const response = await api.get('/public/testimonials')
    testimonials.value = response.data
  } catch (err) {
    console.error('Failed to fetch testimonials:', err)
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  submitting.value = true
  error.value = null

  try {
    await api.post('/public/testimonials', formData.value)
    submitted.value = true
    formData.value = {
      name: '',
      event_type: '',
      rating: 5,
      content: ''
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Une erreur est survenue'
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchTestimonials()
})
</script>

<style scoped>
.testimonials-page {
  min-height: 100vh;
  background: linear-gradient(135deg, var(--color-ivory) 0%, var(--color-champagne) 100%);
}

.testimonials-container {
  padding: calc(80px + var(--spacing-5xl)) 0 var(--spacing-5xl);
}

.testimonials-header {
  text-align: center;
  margin-bottom: var(--spacing-3xl);
}

.testimonials-header h1 {
  font-size: var(--font-size-5xl);
  margin-bottom: var(--spacing-md);
}

.subtitle {
  font-size: var(--font-size-xl);
  color: var(--color-gray);
}

.submit-cta {
  text-align: center;
  background: var(--color-white);
  padding: var(--spacing-2xl);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-3xl);
  box-shadow: var(--shadow-md);
}

.submit-cta p {
  margin-bottom: var(--spacing-md);
  color: var(--color-gray);
}

.testimonials-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: var(--spacing-xl);
}

.testimonial-card {
  padding: var(--spacing-2xl);
  text-align: center;
}

.testimonial-rating {
  font-size: var(--font-size-xl);
  margin-bottom: var(--spacing-lg);
}

.star {
  color: var(--color-gray-light);
}

.star.filled {
  color: var(--color-gold);
}

.testimonial-content {
  font-style: italic;
  color: var(--color-text);
  line-height: var(--line-height-relaxed);
  margin-bottom: var(--spacing-lg);
  font-size: var(--font-size-lg);
}

.testimonial-author {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.testimonial-author strong {
  color: var(--color-gold);
  font-size: var(--font-size-lg);
}

.event-type {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
}

/* Form */
.testimonial-form {
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
  font-weight: var(--font-weight-medium);
  color: var(--color-text);
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: var(--spacing-md);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-sans);
  font-size: var(--font-size-base);
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--color-gold);
}

.rating-input {
  display: flex;
  gap: var(--spacing-sm);
}

.star-input {
  font-size: var(--font-size-3xl);
  color: var(--color-gray-light);
  cursor: pointer;
  transition: color var(--transition-base);
}

.star-input.filled {
  color: var(--color-gold);
}

.star-input:hover {
  color: var(--color-gold-light);
}

.error-message {
  background: rgba(193, 41, 46, 0.1);
  color: var(--color-error);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  text-align: center;
}

.success-message {
  text-align: center;
  padding: var(--spacing-xl);
}

.success-icon {
  width: 60px;
  height: 60px;
  background: var(--color-success);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-3xl);
  margin: 0 auto var(--spacing-lg);
}

.loading-state {
  text-align: center;
  padding: var(--spacing-5xl) 0;
  color: var(--color-gray);
}

@media (max-width: 768px) {
  .testimonials-grid {
    grid-template-columns: 1fr;
  }
}
</style>
