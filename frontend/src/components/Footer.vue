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
            Créer l'instant parfait avec élégance et passion.
          </p>

        </div>

        <!-- Navigation -->
        <div class="footer-section nav-section">
          <h4>EXPLORER</h4>
          <RouterLink to="/services">SERVICES</RouterLink>
          <RouterLink to="/galerie">GALERIE</RouterLink>
          <RouterLink to="/temoignages">TÉMOIGNAGES</RouterLink>
          <RouterLink to="/a-propos">À PROPOS</RouterLink>
        </div>

        <!-- Contact -->
        <div class="footer-section contact-section">
          <h4>CONTACT</h4>
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
          <h4>NEWSLETTER</h4>
          <p>Restez inspiré</p>
          <form @submit.prevent="handleNewsletter" class="newsletter-form">
            <input v-model="email" type="email" placeholder="Votre email" required />
            <button type="submit" aria-label="S'inscrire" class="newsletter-btn">
              <ArrowRight class="icon-sm" />
            </button>
          </form>
        </div>
      </div>

      <div class="footer-bottom">
        <p>&copy; {{ new Date().getFullYear() }} Angel Event. Tous droits réservés.</p>
        <div class="legal-links">
          <a href="#">Mentions légales</a>
          <a href="#">Confidentialité</a>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Instagram, Facebook, Twitter, Mail, Phone, MapPin, ArrowRight } from 'lucide-vue-next'
import api from '../services/api'

const email = ref('')

async function handleNewsletter() {
  try {
    await api.post('/public/newsletter/subscribe', { email: email.value, language: 'fr' })
    email.value = ''
    alert('Merci de vous être abonné!')
  } catch (error) {
    console.error('Newsletter subscription failed:', error)
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
</style>
