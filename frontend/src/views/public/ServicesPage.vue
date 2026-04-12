<template>
  <div class="services-page">
    <Header />
    <!-- Hero Section -->
    <section class="services-hero">
      <div class="container text-center">
        <h1 class="font-script text-gold fade-in-up">{{ getContent('services_hero_title', t('services_page.hero.title')) }}</h1>
        <p class="hero-subtitle fade-in-up">{{ getContent('services_hero_description', t('services_page.hero.subtitle')) }}</p>
      </div>
    </section>

    <!-- Services -->
    <section class="services-section">
      <div class="container">
        <div v-for="(service, index) in services" :key="service.id" class="service-detail" :class="{ reverse: index % 2 === 1 }">
          <div class="service-image fade-in-up">
            <img :src="service.image" :alt="service.title" />
          </div>
          <div class="service-content fade-in-up">
            <div class="service-icon">{{ service.icon }}</div>
            <h2>{{ service.title }}</h2>
            <p class="service-description">{{ service.description }}</p>
            <ul class="service-features">
              <li v-for="feature in service.features" :key="feature">
                <span class="check-icon">✓</span> {{ feature }}
              </li>
            </ul>
            <Button @click="$router.push('/reserver')">{{ getContent('services_list_cta', t('services_page.list.cta')) }}</Button>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta-section">
      <div class="container text-center">
        <h2 class="font-script fade-in-up">{{ getContent('services_cta_title', t('services_page.cta.title')) }}</h2>
        <p class="fade-in-up">{{ getContent('services_cta_subtitle', t('services_page.cta.subtitle')) }}</p>
        <div class="cta-buttons fade-in-up">
          <Button size="lg" @click="$router.push('/reserver')">{{ getContent('services_cta_quote', t('services_page.cta.quote')) }}</Button>
          <Button size="lg" variant="white" @click="$router.push('/contact')">{{ getContent('services_cta_contact', t('services_page.cta.contact')) }}</Button>
        </div>
      </div>
    </section>

    <Footer />
  </div>
</template>

<script setup>
import { onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Header from '../../components/Header.vue'
import Footer from '../../components/Footer.vue'
import Button from '../../components/ui/Button.vue'
import { useSiteContent } from '../../composables/useSiteContent'

const { t } = useI18n()
const { getContent, getImage } = useSiteContent('services')

const serviceDefinitions = [
  {
    id: 'wedding',
    icon: '💍',
    imageFallback: '/storage/wedding/photo_2026-01-10%2001.29.11.jpeg'
  },
  {
    id: 'proposal',
    icon: '💌',
    imageFallback: '/storage/marryme/photo_2026-01-10%2001.22.12.jpeg'
  },
  {
    id: 'baptism',
    icon: '🕊️',
    imageFallback: '/storage/bapteme/photo_2026-01-10%2001.31.48.jpeg'
  },
  {
    id: 'birthday',
    icon: '🎂',
    imageFallback: '/storage/birthday/photo_2026-01-10%2001.19.03.jpeg'
  },
  {
    id: 'baby_shower',
    icon: '🧸',
    imageFallback: '/storage/baby_shower/photo_2026-01-10%2001.34.36.jpeg'
  }
]

const services = computed(() => serviceDefinitions.map((service) => ({
  id: service.id,
  icon: service.icon,
  title: getContent(`services_${service.id}_title`, t(`services_page.list.${service.id}.title`)),
  description: getContent(`services_${service.id}_desc`, t(`services_page.list.${service.id}.desc`)),
  image: getImage(`services_${service.id}_image`, service.imageFallback),
  features: Array.from({ length: 8 }, (_, index) => (
    getContent(`services_${service.id}_feature_${index + 1}`, t(`services_page.list.${service.id}.f${index + 1}`))
  ))
})))

onMounted(() => {
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.classList.add('revealed')
      }
    })
  }, { threshold: 0.1 })

  document.querySelectorAll('.fade-in-up').forEach(el => {
    el.classList.add('scroll-reveal')
    observer.observe(el)
  })
})
</script>

<style scoped>
.services-page {
  min-height: 100vh;
  padding-top: 80px;
}

/* Hero Section */
.services-hero {
  padding: var(--spacing-5xl) 0 var(--spacing-4xl);
  background: linear-gradient(135deg, var(--color-ivory) 0%, var(--color-champagne) 100%);
}

.services-hero h1 {
  font-size: var(--font-size-6xl);
  margin-bottom: var(--spacing-lg);
}

.hero-subtitle {
  font-size: var(--font-size-xl);
  color: var(--color-gray);
  max-width: 600px;
  margin: 0 auto;
}

/* Services Section */
.services-section {
  padding: var(--spacing-5xl) 0;
}

.service-detail {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-4xl);
  align-items: center;
  margin-bottom: var(--spacing-5xl);
}

.service-detail.reverse {
  direction: rtl;
}

.service-detail.reverse > * {
  direction: ltr;
}

.service-image {
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: var(--shadow-xl);
}

.service-image img {
  width: 100%;
  height: 400px;
  object-fit: cover;
  transition: transform var(--transition-slow);
}

.service-image:hover img {
  transform: scale(1.05);
}

.service-content {
  padding: var(--spacing-xl);
}

.service-icon {
  font-size: var(--font-size-5xl);
  margin-bottom: var(--spacing-lg);
}

.service-content h2 {
  color: var(--color-gold);
  margin-bottom: var(--spacing-md);
  font-size: var(--font-size-3xl);
}

.service-description {
  color: var(--color-gray);
  margin-bottom: var(--spacing-xl);
  line-height: var(--line-height-relaxed);
  font-size: var(--font-size-lg);
}

.service-features {
  list-style: none;
  margin-bottom: var(--spacing-2xl);
}

.service-features li {
  padding: var(--spacing-sm) 0;
  color: var(--color-text);
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.check-icon {
  color: var(--color-gold);
  font-weight: bold;
  font-size: var(--font-size-lg);
}

/* CTA Section */
.cta-section {
  padding: var(--spacing-5xl) 0;
  background: var(--color-black-soft);
  color: var(--color-white);
}

.cta-section h2 {
  font-size: var(--font-size-5xl);
  color: var(--color-gold);
  margin-bottom: var(--spacing-lg);
}

.cta-section p {
  font-size: var(--font-size-xl);
  margin-bottom: var(--spacing-2xl);
  color: var(--color-champagne);
}

.cta-buttons {
  display: flex;
  gap: var(--spacing-lg);
  justify-content: center;
  flex-wrap: wrap;
}

@media (max-width: 768px) {
  .service-detail {
    grid-template-columns: 1fr;
  }
  
  .service-detail.reverse {
    direction: ltr;
  }

  .cta-buttons {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
