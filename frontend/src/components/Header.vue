<template>
  <header class="header" :class="{ scrolled: isScrolled }">
    <!-- Top Bar -->
    <div class="top-bar">
      <div class="container top-bar-content">
        <nav class="top-nav">
          <RouterLink to="/a-propos">{{ t('nav.about') }}</RouterLink>
          <RouterLink to="/services">{{ t('nav.services') }}</RouterLink>
          <RouterLink to="/galerie">{{ t('nav.gallery') }}</RouterLink>
          <RouterLink to="/location">{{ t('nav.rentals') }}</RouterLink>
          <RouterLink to="/contact">{{ t('nav.contact') }}</RouterLink>
        </nav>


        <div class="top-socials">
          <button @click="toggleLanguage" class="lang-btn" aria-label="Toggle Language">
            <Globe class="icon-sm" />
            <span>{{ locale === 'fr' ? 'EN' : 'FR' }}</span>
          </button>
          <div class="divider"></div>
          <a href="https://www.instagram.com/angel_eventt/" target="_blank" rel="noopener noreferrer" aria-label="Instagram">
            <Instagram class="icon-sm" />
          </a>
          <a href="https://www.tiktok.com/@angel_eventt" target="_blank" rel="noopener noreferrer" aria-label="TikTok">
            <TikTok class="icon-sm" />
          </a>
        </div>
      </div>
    </div>

    <!-- Main Branding Area -->
    <div class="brand-area" :class="{ compact: isScrolled }">
      <div class="container brand-content">
        <RouterLink to="/" class="brand-logo">
          <div class="logo-icon-wrapper">
             <img src="/logo.jpeg" alt="Angel Event Logo" class="logo-image" />
          </div>
          <div class="brand-text">
            <span class="brand-name">ANGEL EVENT</span>
            <span class="brand-tagline">{{ t('nav.tagline') }}</span>
          </div>
        </RouterLink>
        
        <div class="header-actions">
           <Button @click="goToBooking" variant="outline" size="sm">
             {{ t('nav.book_now') }}
           </Button>
        </div>
        
        <button class="mobile-menu-btn" @click="mobileMenuOpen = !mobileMenuOpen">
          <Menu v-if="!mobileMenuOpen" />
          <X v-else />
        </button>
      </div>
    </div>

    <!-- Mobile Menu Overlay -->
    <div class="mobile-menu" :class="{ open: mobileMenuOpen }">
      <nav class="mobile-nav">
        <RouterLink to="/" @click="mobileMenuOpen = false">{{ t('nav.home') }}</RouterLink>
        <RouterLink to="/a-propos" @click="mobileMenuOpen = false">{{ t('nav.about') }}</RouterLink>
        <RouterLink to="/services" @click="mobileMenuOpen = false">{{ t('nav.services') }}</RouterLink>
        <RouterLink to="/galerie" @click="mobileMenuOpen = false">{{ t('nav.gallery') }}</RouterLink>
        <RouterLink to="/temoignages" @click="mobileMenuOpen = false">{{ t('nav.testimonials') }}</RouterLink>
        <RouterLink to="/contact" @click="mobileMenuOpen = false">{{ t('nav.contact') }}</RouterLink>
        <div class="mobile-actions">
           <Button @click="goToBooking" class="w-full">
             {{ t('nav.book') }}
           </Button>
        </div>
      </nav>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Search, Instagram, Menu, X, Globe } from 'lucide-vue-next'
import TikTok from './icons/TikTok.vue'
import Button from './ui/Button.vue'

const router = useRouter()
const { t, locale } = useI18n()
const isScrolled = ref(false)
const mobileMenuOpen = ref(false)

function toggleLanguage() {
  locale.value = locale.value === 'fr' ? 'en' : 'fr'
}

function handleScroll() {
  isScrolled.value = window.scrollY > 50
}

function goToBooking() {
  router.push('/reserver')
  mobileMenuOpen.value = false
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: var(--z-fixed);
  background: var(--color-white);
  transition: all var(--transition-base);
  box-shadow: var(--shadow-sm);
}

/* Top Bar */
.top-bar {
  border-bottom: 1px solid var(--color-gray-light);
  padding: var(--spacing-xs) 0;
  font-size: var(--font-size-xs);
  background: var(--color-white);
  display: block;
}

.header.scrolled .top-bar {
  display: none; /* Hide top bar on scroll for cleaner look */
}

.top-bar-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.top-nav {
  display: flex;
  gap: var(--spacing-xl);
  flex: 1;
  justify-content: center;
}

.top-nav a {
  text-decoration: none;
  color: var(--color-text);
  font-weight: 500;
  letter-spacing: 0.1em;
  transition: color 0.3s;
  font-size: 0.7rem;
}

.top-nav a:hover {
  color: var(--color-gold);
}

.top-socials {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-md);
}

.top-socials a {
  color: var(--color-text);
  transition: color 0.3s;
}

.top-socials a:hover {
  color: var(--color-gold);
}

.lang-btn {
  background: none;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  color: var(--color-text);
  font-weight: 500;
  font-size: 0.7rem;
  padding: 0;
  transition: color 0.3s;
}

.lang-btn:hover {
  color: var(--color-gold);
}

.divider {
  width: 1px;
  background-color: var(--color-gray-light);
  height: 14px;
}

.icon-sm {
  width: 14px;
  height: 14px;
}

/* Brand Area */
.brand-area {
  padding: var(--spacing-lg) 0;
  background: var(--color-white);
  transition: all var(--transition-base);
}

.brand-area.compact {
  padding: var(--spacing-sm) 0;
}

.brand-content {
  display: flex;
  align-items: center;
  justify-content: center; /* Centered by default due to reference image */
  position: relative;
}

.brand-logo {
  text-decoration: none;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  color: var(--color-black);
  transition: transform 0.3s;
}

.brand-area.compact .brand-logo {
  flex-direction: row;
  gap: var(--spacing-md);
}

.brand-area.compact .brand-logo .brand-text {
  text-align: left;
}

.logo-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-image {
  height: 60px;
  width: auto;
  object-fit: contain;
  transition: all var(--transition-base);
  margin-bottom: var(--spacing-xs);
}

.brand-area.compact .logo-image {
  height: 40px;
  margin-bottom: 0;
}

.brand-name {
  font-family: var(--font-serif);
  font-size: 2rem;
  letter-spacing: 0.2em;
  line-height: 1;
  text-transform: uppercase;
}

.brand-area.compact .brand-name {
  font-size: 1.5rem;
}

.brand-tagline {
  font-size: 0.6rem;
  letter-spacing: 0.4em;
  color: var(--color-text-light);
  margin-top: 5px;
  text-transform: uppercase;
}

.brand-area.compact .brand-tagline {
  display: none;
}

.header-actions {
  position: absolute;
  right: var(--spacing-xl);
  top: 50%;
  transform: translateY(-50%);
}

/* Mobile Menu */
.mobile-menu-btn {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  position: absolute;
  left: var(--spacing-xl);
  top: 50%;
  transform: translateY(-50%);
}

.mobile-menu {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: var(--color-white);
  z-index: var(--z-fixed);
  padding: 100px var(--spacing-xl);
  transform: translateX(-100%);
  transition: transform 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
}

.mobile-menu.open {
  transform: translateX(0);
}

.mobile-nav {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  align-items: center;
}

.mobile-nav a {
  font-size: 1.2rem;
  text-decoration: none;
  color: var(--color-black);
  letter-spacing: 0.1em;
  text-transform: uppercase;
}

.mobile-actions {
  margin-top: var(--spacing-xl);
  width: 100%;
  max-width: 300px;
}

/* Responsive */
@media (max-width: 1024px) {
  .top-bar {
    display: none;
  }
  
  .header-actions {
    display: none;
  }
  
  .mobile-menu-btn {
    display: block;
  }
  
  .brand-logo {
    transform: scale(0.8);
  }
}
</style>
