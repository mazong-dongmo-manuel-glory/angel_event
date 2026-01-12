<template>
  <footer class="footer">
    <div class="container bordered-top">
      <div class="footer-content">
        <!-- Brand Section -->
        <div class="footer-section brand-section">
          <div class="footer-brand">
            <img src="/logo.jpeg" alt="Angel Event Logo" class="footer-logo-img" />
            <h3 class="footer-logo">ANGEL EVENT</h3>
          </div>

          <p class="footer-desc">
            {{ $t('footer.desc') }}
          </p>
          <div class="social-links">
            <a href="https://www.instagram.com/angel_eventt/" target="_blank" rel="noopener noreferrer" aria-label="Instagram">
              <Instagram class="icon" />
            </a>
            <a href="https://www.tiktok.com/@angel_eventt" target="_blank" rel="noopener noreferrer" aria-label="TikTok">
              <TikTok class="icon" />
            </a>
          </div>

        </div>

        <!-- Navigation -->
        <div class="footer-section nav-section">
          <h4>{{ $t('footer.explore') }}</h4>
          <RouterLink to="/services">{{ $t('nav.services') }}</RouterLink>
          <RouterLink to="/galerie">{{ $t('nav.gallery') }}</RouterLink>
          <RouterLink to="/temoignages">{{ $t('nav.testimonials') }}</RouterLink>
          <RouterLink to="/a-propos">{{ $t('nav.about') }}</RouterLink>
        </div>

        <!-- Contact -->
        <div class="footer-section contact-section">
          <h4>{{ $t('footer.contact') }}</h4>
          <div class="contact-item">
            <Mail class="icon-xs" />
            <span>contact@angelevent.com</span>
          </div>
          <div class="contact-item">
            <Phone class="icon-xs" />
            <span>+1 (819) 244-4702</span>
          </div>
          <div class="contact-item">
            <MapPin class="icon-xs" />
            <span>Trois-Rivieres, QC</span>
          </div>
        </div>

        <!-- Newsletter -->
        <div class="footer-section newsletter-section">
          <h4>{{ $t('footer.newsletter') }}</h4>
          <p>{{ $t('footer.newsletter_desc') }}</p>
          <form @submit.prevent="handleNewsletter" class="newsletter-form">
            <input v-model="email" type="email" :placeholder="$t('footer.email_placeholder')" required />
            <button type="submit" :aria-label="$t('footer.subscribe')" class="newsletter-btn">
              <ArrowRight class="icon-sm" />
            </button>
          </form>
        </div>
      </div>

      <div class="footer-bottom">
        <div class="copyright">
          <p>&copy; {{ new Date().getFullYear() }} Angel Event. {{ $t('footer.rights_reserved') }}</p>
          <p class="mazong-credit">
            {{ $t('footer.made_by') }} <span class="mazong-highlight">Mazong</span>
          </p>
        </div>
        <div class="legal-links">
          <a href="#">{{ $t('footer.legal') }}</a>
          <a href="#">{{ $t('footer.privacy') }}</a>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Instagram, Mail, Phone, MapPin, ArrowRight } from 'lucide-vue-next'
import TikTok from './icons/TikTok.vue'
import api from '../services/api'
import i18n from '../i18n'

const { t } = useI18n()
const email = ref('')

async function handleNewsletter() {
  try {
    await api.post('/public/newsletter/subscribe', { 
      email: email.value, 
      language: i18n.global.locale.value 
    })
    email.value = ''
    alert(t('footer.subscribe_success'))
  } catch (error) {
    console.error('Newsletter subscription failed:', error)
    alert(t('footer.subscribe_error'))
  }
}
</script>

<style scoped>
.footer {
  background: var(--color-black-rich);
  color: var(--color-white);
  padding: var(--spacing-5xl) 0 var(--spacing-xl);
}

.bordered-top {
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  padding-top: var(--spacing-3xl);
}

.footer-content {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1.5fr;
  gap: var(--spacing-3xl);
  margin-bottom: var(--spacing-3xl);
}

.footer-section h3.footer-logo {
  font-family: var(--font-serif);
  font-size: 1.5rem;
  letter-spacing: 0.1em;
  margin-bottom: var(--spacing-md);
  color: var(--color-white);
}

.footer-brand {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.footer-logo-img {
  height: 50px;
  width: auto;
  object-fit: contain;
}

.footer-section h4 {
  color: var(--color-white);
  font-size: 0.8rem;
  letter-spacing: 0.2em;
  margin-bottom: var(--spacing-lg);
  opacity: 0.7;
}

.footer-desc {
  color: var(--color-gray-medium);
  margin-bottom: var(--spacing-lg);
  max-width: 300px;
}

.nav-section a {
  display: block;
  color: var(--color-gray-medium);
  text-decoration: none;
  margin-bottom: var(--spacing-sm);
  transition: color var(--transition-base);
  font-size: 0.9rem;
}

.nav-section a:hover {
  color: var(--color-white);
}

.contact-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  color: var(--color-gray-medium);
  margin-bottom: var(--spacing-sm);
  font-size: 0.9rem;
}

.newsletter-form {
  display: flex;
  gap: var(--spacing-sm);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  padding-bottom: var(--spacing-xs);
}

.newsletter-form input {
  background: none;
  border: none;
  color: white;
  padding: var(--spacing-sm) 0;
  width: 100%;
  outline: none;
}

.newsletter-form input::placeholder {
  color: rgba(255, 255, 255, 0.6);
}

.newsletter-btn {
  background: white;
  border: none;
  color: var(--color-black);
  cursor: pointer;
  padding: var(--spacing-sm);
  border-radius: var(--radius-sm);
  transition: all var(--transition-base);
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 40px;
}

.newsletter-btn:hover {
  background: var(--color-gold);
  color: white;
  transform: translateX(3px);
}

.social-links {
  display: flex;
  gap: var(--spacing-md);
}

.social-links a {
  color: var(--color-white);
  transition: color 0.3s;
}

.social-links a:hover {
  color: var(--color-gray-medium);
}

.icon { width: 20px; height: 20px; }
.icon-sm { width: 16px; height: 16px; }
.icon-xs { width: 14px; height: 14px; }

.footer-bottom {
  display: flex;
  justify-content: space-between;
  padding-top: var(--spacing-xl);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  color: var(--color-gray-dark);
  font-size: 0.8rem;
}

.legal-links a {
  color: var(--color-gray-dark);
  margin-left: var(--spacing-lg);
  text-decoration: none;
}

@media (max-width: 1024px) {
  .footer-content {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 768px) {
  .footer-content {
    grid-template-columns: 1fr;
  }
  .footer-bottom {
    flex-direction: column;
    gap: var(--spacing-md);
    text-align: center;
  }
  .legal-links {
    margin-top: var(--spacing-sm);
  }
  .legal-links a {
    margin: 0 var(--spacing-md);
  }
}

.copyright {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.mazong-credit {
  font-size: 0.8rem;
  opacity: 0.8;
}

.mazong-highlight {
  color: var(--color-gold);
  font-weight: 600;
  letter-spacing: 0.05em;
}
</style>
