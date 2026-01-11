<template>
  <div class="home-page">
    <Header />

    <!-- Luxury Hero Section -->
    <section class="hero-luxury">
      <div class="hero-image-container">
        <!-- Using a high-fashion, clean event image -->
        <img src="https://images.unsplash.com/photo-1519741497674-611481863552?w=1920" alt="Luxury Event" class="hero-img" />
        <div class="hero-overlay"></div>
      </div>
      
      <div class="hero-content text-center" v-motion-fade-visible>
        <span class="hero-pretitle">EST. 2024</span>
        <h1 class="hero-title">CRÉATEUR D'ÉVÉNEMENTS</h1>
        <div class="hero-divider"></div>
        <p class="hero-subtitle">
          "En travaillant ensemble et en réunissant nos richesses, nous pouvons accomplir de grandes choses."
        </p>
        <Button variant="white" size="lg" @click="$router.push('/services')">
          DÉCOUVRIR L'EXPÉRIENCE
        </Button>
      </div>
    </section>

    <!-- Welcome Section (Clean White) -->
    <section class="section-clean">
      <div class="container container-narrow text-center">
        <h2 class="section-title">BIENVENUE CHEZ ANGEL EVENT</h2>
        <p class="section-text">
          Nous sommes une agence de planification d'événements haut de gamme, spécialisée dans la création de moments inoubliables. 
          Notre approche minimaliste et sophistiquée met en valeur l'essence de chaque célébration.
        </p>
        <div class="signature">Angel Event Team</div>
      </div>
    </section>

    <!-- Services Trio (Widget Style) -->
    <section class="section-gray">
      <div class="container">
        <div class="services-grid">
          <div class="service-widget" v-for="service in services" :key="service.id">
            <div class="widget-icon">
              <!-- Dynamic Icon Component -->
              <component :is="service.icon" stroke-width="1" />
            </div>
            <h3>{{ service.title }}</h3>
            <p>{{ service.description }}</p>
            <a @click="$router.push('/services')" class="link-luxury">EN SAVOIR PLUS</a>
          </div>
        </div>
      </div>
    </section>

    <!-- Image + Quote Section (Parallax feel) -->
    <section class="quote-section">
      <div class="quote-content">
        <Quote class="quote-icon" />
        <blockquote>
          L'élégance n'est pas de se faire remarquer, mais de se faire souvenir.
        </blockquote>
        <cite>Giorgio Armani</cite>
      </div>
    </section>

    <!-- Gallery Strip -->
    <section class="gallery-strip">
      <div class="gallery-item" v-for="(img, idx) in galleryImages" :key="idx">
        <img :src="img" alt="Gallery" />
        <div class="gallery-overlay">
          <Instagram class="icon-white" />
        </div>
      </div>
    </section>

    <Footer />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Crown, Heart, PartyPopper, Quote, Instagram } from 'lucide-vue-next'
import Header from '../../components/Header.vue'
import Footer from '../../components/Footer.vue'
import Button from '../../components/ui/Button.vue'

const services = ref([
  {
    id: 1,
    icon: Heart,
    title: 'MARIAGES',
    description: 'Une coordination parfaite pour le plus beau jour de votre vie.'
  },
  {
    id: 2,
    icon: Crown,
    title: 'CORPORATIF',
    description: 'Des événements professionnels qui marquent les esprits.'
  },
  {
    id: 3,
    icon: PartyPopper,
    title: 'CÉLÉBRATIONS',
    description: 'Anniversaires, fiançailles et moments précieux.'
  }
])

const galleryImages = ref([
  'https://images.unsplash.com/photo-1511795409834-ef04bbd61622?w=600',
  'https://images.unsplash.com/photo-1519225421980-715cb0215aed?w=600',
  'https://images.unsplash.com/photo-1469334031218-e382a71b716b?w=600',
  'https://images.unsplash.com/photo-1519741497674-611481863552?w=600'
])
</script>

<style scoped>
/* Luxury Hero */
.hero-luxury {
  position: relative;
  height: 90vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--color-black);
}

.hero-image-container {
  position: absolute;
  inset: 0;
  opacity: 0.8; /* Slightly darkened for text contrast */
}

.hero-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hero-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.3);
}

.hero-content {
  position: relative;
  z-index: 1;
  color: var(--color-white);
  padding: 0 var(--spacing-xl);
  max-width: 900px;
}

.hero-pretitle {
  display: block;
  font-size: var(--font-size-sm);
  letter-spacing: 0.3em;
  margin-bottom: var(--spacing-md);
}

.hero-title {
  font-family: var(--font-serif);
  font-size: clamp(2.5rem, 5vw, 5rem);
  line-height: 1.1;
  letter-spacing: 0.05em;
  margin-bottom: var(--spacing-lg);
}

.hero-divider {
  width: 80px;
  height: 2px;
  background: var(--color-white);
  margin: 0 auto var(--spacing-lg);
}

.hero-subtitle {
  font-family: var(--font-serif);
  font-size: 1.2rem;
  font-style: italic;
  opacity: 0.9;
  margin-bottom: var(--spacing-2xl);
  line-height: 1.8;
}

/* Sections */
.section-clean {
  padding: var(--spacing-5xl) 0;
  background: var(--color-white);
}

.section-title {
  font-size: 2.5rem;
  letter-spacing: 0.1em;
  margin-bottom: var(--spacing-lg);
  color: var(--color-black);
}

.section-text {
  color: var(--color-text-light);
  font-size: 1.1rem;
  line-height: 2;
  margin-bottom: var(--spacing-xl);
}

.signature {
  font-family: var(--font-script);
  font-size: 3rem;
  color: var(--color-black);
}

/* Services Widget */
.section-gray {
  padding: var(--spacing-5xl) 0;
  background: var(--color-surface);
}

.services-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: var(--spacing-xl);
}

.service-widget {
  background: var(--color-white);
  padding: var(--spacing-2xl) var(--spacing-xl);
  text-align: center;
  transition: transform 0.3s;
  border: 1px solid transparent;
}

.service-widget:hover {
  transform: translateY(-10px);
  border-color: var(--color-border);
  box-shadow: var(--shadow-lg);
}

.widget-icon {
  color: var(--color-gold); /* Keep gold accent ONLY on icons for luxury touch */
  width: 48px;
  height: 48px;
  margin: 0 auto var(--spacing-lg);
}

.widget-icon svg {
  width: 100%;
  height: 100%;
}

.service-widget h3 {
  font-size: 1.2rem;
  letter-spacing: 0.2em;
  margin-bottom: var(--spacing-md);
  color: var(--color-black);
}

.service-widget p {
  color: var(--color-text-light);
  margin-bottom: var(--spacing-xl);
}

.link-luxury {
  font-size: 0.8rem;
  letter-spacing: 0.2em;
  border-bottom: 1px solid var(--color-border);
  padding-bottom: 5px;
  cursor: pointer;
  transition: all 0.3s;
}

.link-luxury:hover {
  border-color: var(--color-black);
  color: var(--color-black);
}

/* Quote Section */
.quote-section {
  padding: var(--spacing-5xl) 0;
  background: var(--color-black);
  color: var(--color-white);
  text-align: center;
  position: relative;
  overflow: hidden;
}

.quote-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 var(--spacing-xl);
}

.quote-icon {
  width: 40px;
  height: 40px;
  opacity: 0.3;
  margin-bottom: var(--spacing-xl);
}

blockquote {
  font-family: var(--font-serif);
  font-size: 2rem;
  line-height: 1.4;
  margin-bottom: var(--spacing-lg);
}

cite {
  font-style: normal;
  letter-spacing: 0.2em;
  font-size: 0.9rem;
  opacity: 0.7;
}

/* Gallery Strip */
.gallery-strip {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  height: 300px;
}

.gallery-item {
  position: relative;
  height: 100%;
  overflow: hidden;
}

.gallery-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.gallery-item:hover img {
  transform: scale(1.1);
}

.gallery-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.gallery-item:hover .gallery-overlay {
  opacity: 1;
}

.icon-white {
  color: white;
  width: 32px;
  height: 32px;
}

@media (max-width: 768px) {
  .gallery-strip {
    grid-template-columns: repeat(2, 1fr);
    height: auto;
  }
  .gallery-item {
    height: 250px;
  }
  .hero-title {
    font-size: 2.5rem;
  }
}
</style>
