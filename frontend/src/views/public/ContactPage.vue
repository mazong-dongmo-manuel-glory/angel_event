<template>
  <div class="contact-page">
    <Header />
    <div class="container container-narrow contact-container">
      <div class="contact-header">
        <h1 class="font-script text-gold">{{ t('contact.title') }}</h1>
        <p class="subtitle">{{ t('contact.subtitle') }}</p>
      </div>

      <div class="contact-content">
        <div class="contact-info">
          <div class="info-card">
            <div class="info-icon">📧</div>
            <h3>{{ t('contact.info.email') }}</h3>
            <p><a href="mailto:contact@angelevent.com">contact@angelevent.com</a></p>
          </div>
          <div class="info-card">
            <div class="info-icon">📱</div>
            <h3>{{ t('contact.info.phone') }}</h3>
            <p><a href="tel:+15551234567">+1 (555) 123-4567</a></p>
          </div>
          <div class="info-card">
            <div class="info-icon">💬</div>
            <h3>{{ t('contact.info.whatsapp') }}</h3>
            <p><a href="https://wa.me/15551234567" target="_blank">{{ t('contact.info.chat') }}</a></p>
          </div>
        </div>

        <div class="contact-form-container">
          <form v-if="!submitted" @submit.prevent="handleSubmit" class="contact-form">
            <div class="form-group">
              <label for="name">{{ t('contact.form.name') }}</label>
              <input
                id="name"
                v-model="formData.name"
                type="text"
                required
                :placeholder="t('contact.form.name_ph')"
              />
            </div>

            <div class="form-group">
              <label for="email">{{ t('contact.form.email') }}</label>
              <input
                id="email"
                v-model="formData.email"
                type="email"
                required
                :placeholder="t('contact.form.email_ph')"
              />
            </div>

            <div class="form-group">
              <label for="phone">{{ t('contact.form.phone') }}</label>
              <input
                id="phone"
                v-model="formData.phone"
                type="tel"
                :placeholder="t('contact.form.phone_ph')"
              />
            </div>

            <div class="form-group">
              <label for="message">{{ t('contact.form.message') }}</label>
              <textarea
                id="message"
                v-model="formData.message"
                rows="6"
                required
                :placeholder="t('contact.form.message_ph')"
              ></textarea>
            </div>

            <div v-if="error" class="error-message">
              {{ error }}
            </div>

            <Button type="submit" size="lg" :loading="loading" block>
              {{ loading ? t('contact.form.sending') : t('contact.form.submit') }}
            </Button>
          </form>

          <div v-else class="success-message">
            <div class="success-icon">✓</div>
            <h3>{{ t('contact.form.success_title') }}</h3>
            <p>{{ t('contact.form.success_msg') }}</p>
            <Button @click="resetForm">{{ t('contact.form.reset') }}</Button>
          </div>
        </div>

        <!-- Newsletter Signup -->
        <div class="newsletter-section">
          <h3>{{ t('contact.newsletter.title') }}</h3>
          <p>{{ t('contact.newsletter.text') }}</p>
          <form v-if="!newsletterSubmitted" @submit.prevent="handleNewsletterSubmit" class="newsletter-form">
            <input
              v-model="newsletterEmail"
              type="email"
              :placeholder="t('contact.newsletter.placeholder')"
              required
            />
            <Button type="submit" :loading="newsletterLoading">
              {{ t('contact.newsletter.submit') }}
            </Button>
          </form>
          <div v-else class="newsletter-success">
            {{ t('contact.newsletter.success') }}
          </div>
        </div>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Header from '../../components/Header.vue'
import Footer from '../../components/Footer.vue'
import Button from '../../components/ui/Button.vue'
import api from '../../services/api'

const { t } = useI18n()

const formData = ref({
  name: '',
  email: '',
  phone: '',
  message: ''
})

const loading = ref(false)
const error = ref(null)
const submitted = ref(false)

const newsletterEmail = ref('')
const newsletterLoading = ref(false)
const newsletterSubmitted = ref(false)

const route = useRoute()

// Watch for query params changes (immediate: true handles initial load)
watch(() => route.query, (newQuery) => {
  if (newQuery.rental_item_name) {
    formData.value.message = t('contact.form.rental_interest', { item: newQuery.rental_item_name })
  }
}, { immediate: true })

async function handleSubmit() {
  loading.value = true
  error.value = null

  try {
    await api.post('/public/contact', formData.value)
    submitted.value = true
  } catch (err) {
    error.value = err.response?.data?.error || 'Une erreur est survenue. Veuillez réessayer.'
  } finally {
    loading.value = false
  }
}

function resetForm() {
  formData.value = {
    name: '',
    email: '',
    phone: '',
    message: ''
  }
  submitted.value = false
  error.value = null
}

async function handleNewsletterSubmit() {
  newsletterLoading.value = true

  try {
    await api.post('/public/newsletter/subscribe', {
      email: newsletterEmail.value,
      language: t.locale.value
    })
    newsletterSubmitted.value = true
    newsletterEmail.value = ''
  } catch (err) {
    console.error('Newsletter subscription failed:', err)
  } finally {
    newsletterLoading.value = false
  }
}
</script>

<style scoped>
.contact-page {
  min-height: 100vh;
  background: var(--color-ivory);
}

.contact-container {
  padding: calc(80px + var(--spacing-5xl)) 0 var(--spacing-5xl);
}

.contact-header {
  text-align: center;
  margin-bottom: var(--spacing-4xl);
}

.contact-header h1 {
  font-size: var(--font-size-5xl);
  margin-bottom: var(--spacing-md);
}

.subtitle {
  font-size: var(--font-size-xl);
  color: var(--color-gray);
}

.contact-content {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3xl);
}

.contact-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-xl);
}

.info-card {
  background: var(--color-white);
  padding: var(--spacing-2xl);
  border-radius: var(--radius-lg);
  text-align: center;
  box-shadow: var(--shadow-md);
  transition: transform var(--transition-base);
}

.info-card:hover {
  transform: translateY(-5px);
}

.info-icon {
  font-size: var(--font-size-4xl);
  margin-bottom: var(--spacing-md);
}

.info-card h3 {
  color: var(--color-gold);
  margin-bottom: var(--spacing-sm);
}

.info-card a {
  color: var(--color-text);
  text-decoration: none;
  transition: color var(--transition-base);
}

.info-card a:hover {
  color: var(--color-gold);
}

.contact-form-container {
  background: var(--color-white);
  padding: var(--spacing-3xl);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
}

.contact-form {
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

.error-message {
  background: rgba(193, 41, 46, 0.1);
  color: var(--color-error);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  text-align: center;
}

.success-message {
  text-align: center;
  padding: var(--spacing-2xl);
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

.success-message h3 {
  color: var(--color-gold);
  margin-bottom: var(--spacing-md);
}

.success-message p {
  color: var(--color-gray);
  margin-bottom: var(--spacing-xl);
}

/* Newsletter Section */
.newsletter-section {
  background: linear-gradient(135deg, var(--color-gold) 0%, var(--color-gold-light) 100%);
  color: var(--color-white);
  padding: var(--spacing-3xl);
  border-radius: var(--radius-xl);
  text-align: center;
}

.newsletter-section h3 {
  color: var(--color-white);
  margin-bottom: var(--spacing-md);
}

.newsletter-section p {
  margin-bottom: var(--spacing-xl);
  opacity: 0.9;
}

.newsletter-form {
  display: flex;
  gap: var(--spacing-md);
  max-width: 500px;
  margin: 0 auto;
}

.newsletter-form input {
  flex: 1;
  padding: var(--spacing-md);
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--font-size-base);
}

.newsletter-success {
  background: rgba(255, 255, 255, 0.2);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-medium);
}

@media (max-width: 768px) {
  .contact-info {
    grid-template-columns: 1fr;
  }
  
  .contact-form-container {
    padding: var(--spacing-xl);
  }
  
  .newsletter-form {
    flex-direction: column;
  }
}
</style>
