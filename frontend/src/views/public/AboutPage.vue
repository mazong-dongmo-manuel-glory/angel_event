<template>
  <div class="about-page">
    <Header />
    <!-- Hero Section -->
    <section class="about-hero">
      <div class="container text-center">
        <h1 class="font-script text-gold fade-in-up">{{ t('about.hero.title') }}</h1>
        <p class="hero-subtitle fade-in-up">{{ t('about.hero.subtitle') }}</p>
      </div>
    </section>

    <!-- Story Section -->
    <section class="story-section">
      <div class="container">
        <div class="story-content">
          <div class="story-text fade-in-up">
            <h2>{{ t('about.story.title') }}</h2>
            <p>{{ t('about.story.p1') }}</p>
            <p>{{ t('about.story.p2') }}</p>
            <p>{{ t('about.story.p3') }}</p>
          </div>
          <div class="story-image fade-in-up">
            <img src="https://images.unsplash.com/photo-1511285560929-80b456fea0bc?w=800" alt="Angel Event" />
          </div>
        </div>
      </div>
    </section>

    <!-- Values Section -->
    <section class="values-section">
      <div class="container">
        <h2 class="section-title text-center fade-in-up">{{ t('about.values.title') }}</h2>
        <div class="values-grid">
          <div v-for="value in values" :key="value.id" class="value-card fade-in-up">
            <div class="value-icon">{{ value.icon }}</div>
            <h3>{{ value.title }}</h3>
            <p>{{ value.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Why Choose Us -->
    <section class="why-section">
      <div class="container">
        <h2 class="section-title text-center fade-in-up">{{ t('about.why.title') }}</h2>
        <div class="why-grid">
          <div v-for="reason in reasons" :key="reason.id" class="why-item fade-in-up">
            <div class="why-number">{{ reason.id }}</div>
            <h3>{{ reason.title }}</h3>
            <p>{{ reason.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta-section">
      <div class="container text-center">
        <h2 class="font-script fade-in-up">{{ t('about.cta.title') }}</h2>
        <p class="fade-in-up">{{ t('about.cta.subtitle') }}</p>
        <div class="cta-buttons fade-in-up">
          <Button size="lg" @click="$router.push('/reserver')">{{ t('about.cta.book') }}</Button>
          <Button size="lg" variant="white" @click="$router.push('/contact')">{{ t('about.cta.contact') }}</Button>
        </div>
      </div>
    </section>

    <Footer />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Header from '../../components/Header.vue'
import Footer from '../../components/Footer.vue'
import Button from '../../components/ui/Button.vue'

const { t } = useI18n()

const values = computed(() => [
  {
    id: 1,
    icon: '✨',
    title: t('about.values.excellence.title'),
    description: t('about.values.excellence.desc')
  },
  {
    id: 2,
    icon: '💝',
    title: t('about.values.passion.title'),
    description: t('about.values.passion.desc')
  },
  {
    id: 3,
    icon: '🎯',
    title: t('about.values.custom.title'),
    description: t('about.values.custom.desc')
  },
  {
    id: 4,
    icon: '🤝',
    title: t('about.values.trust.title'),
    description: t('about.values.trust.desc')
  }
])

const reasons = computed(() => [
  {
    id: 1,
    title: t('about.why.exp.title'),
    description: t('about.why.exp.desc')
  },
  {
    id: 2,
    title: t('about.why.turnkey.title'),
    description: t('about.why.turnkey.desc')
  },
  {
    id: 3,
    title: t('about.why.detail.title'),
    description: t('about.why.detail.desc')
  },
  {
    id: 4,
    title: t('about.why.network.title'),
    description: t('about.why.network.desc')
  }
])
</script>

<style scoped>
.about-page {
  min-height: 100vh;
}

/* Hero Section */
.about-hero {
  padding: calc(80px + var(--spacing-5xl)) 0 var(--spacing-4xl);
  background: linear-gradient(135deg, var(--color-ivory) 0%, var(--color-champagne) 100%);
}

.about-hero h1 {
  font-size: var(--font-size-6xl);
  margin-bottom: var(--spacing-lg);
}

.hero-subtitle {
  font-size: var(--font-size-xl);
  color: var(--color-gray);
}

/* Story Section */
.story-section {
  padding: var(--spacing-5xl) 0;
}

.story-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-4xl);
  align-items: center;
}

.story-text h2 {
  color: var(--color-gold);
  margin-bottom: var(--spacing-xl);
  font-size: var(--font-size-3xl);
}

.story-text p {
  color: var(--color-gray);
  line-height: var(--line-height-relaxed);
  margin-bottom: var(--spacing-lg);
  font-size: var(--font-size-lg);
}

.story-image {
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: var(--shadow-xl);
}

.story-image img {
  width: 100%;
  height: 500px;
  object-fit: cover;
}

/* Values Section */
.values-section {
  padding: var(--spacing-5xl) 0;
  background: var(--color-gray-lighter);
}

.section-title {
  font-size: var(--font-size-4xl);
  margin-bottom: var(--spacing-4xl);
  color: var(--color-black-soft);
}

.values-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-2xl);
}

.value-card {
  background: var(--color-white);
  padding: var(--spacing-2xl);
  border-radius: var(--radius-lg);
  text-align: center;
  box-shadow: var(--shadow-md);
  transition: transform var(--transition-base);
}

.value-card:hover {
  transform: translateY(-5px);
}

.value-icon {
  font-size: var(--font-size-5xl);
  margin-bottom: var(--spacing-lg);
}

.value-card h3 {
  color: var(--color-gold);
  margin-bottom: var(--spacing-md);
  font-size: var(--font-size-xl);
}

.value-card p {
  color: var(--color-gray);
  line-height: var(--line-height-relaxed);
}

/* Why Section */
.why-section {
  padding: var(--spacing-5xl) 0;
}

.why-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: var(--spacing-3xl);
}

.why-item {
  text-align: center;
}

.why-number {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, var(--color-gold) 0%, var(--color-gold-light) 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-2xl);
  font-weight: bold;
  margin: 0 auto var(--spacing-lg);
}

.why-item h3 {
  color: var(--color-gold);
  margin-bottom: var(--spacing-md);
}

.why-item p {
  color: var(--color-gray);
  line-height: var(--line-height-relaxed);
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
  .story-content {
    grid-template-columns: 1fr;
  }
  
  .values-grid,
  .why-grid {
    grid-template-columns: 1fr;
  }
  
  .cta-buttons {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
