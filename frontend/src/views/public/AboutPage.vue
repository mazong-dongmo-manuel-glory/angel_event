<template>
  <div class="about-page">
    <Header />
    <!-- Hero Section -->
    <section class="about-hero">
      <div class="container text-center">
        <h1 class="font-script text-gold fade-in-up">À Propos d'Angel Event</h1>
        <p class="hero-subtitle fade-in-up">L'art de sublimer vos moments précieux</p>
      </div>
    </section>

    <!-- Story Section -->
    <section class="story-section">
      <div class="container">
        <div class="story-content">
          <div class="story-text fade-in-up">
            <h2>Notre Histoire</h2>
            <p>
              Angel Event est née d'une passion profonde pour la création d'expériences inoubliables.
              Chaque événement est une opportunité unique de transformer des rêves en réalité, de créer
              des souvenirs qui dureront toute une vie.
            </p>
            <p>
              Spécialisés dans les demandes en mariage, les mariages et les événements haut de gamme,
              nous mettons notre expertise et notre créativité au service de vos moments les plus précieux.
            </p>
            <p>
              Notre approche personnalisée garantit que chaque détail reflète votre vision unique et
              dépasse vos attentes.
            </p>
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
        <h2 class="section-title text-center fade-in-up">Nos Valeurs</h2>
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
        <h2 class="section-title text-center fade-in-up">Pourquoi Nous Choisir</h2>
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
        <h2 class="font-script fade-in-up">Créons Ensemble Votre Moment Parfait</h2>
        <p class="fade-in-up">Contactez-nous pour une consultation gratuite</p>
        <div class="cta-buttons fade-in-up">
          <Button size="lg" @click="$router.push('/reserver')">Réserver maintenant</Button>
          <Button size="lg" variant="white" @click="$router.push('/contact')">Nous contacter</Button>
        </div>
      </div>
    </section>

    <Footer />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Header from '../../components/Header.vue'
import Footer from '../../components/Footer.vue'
import Button from '../../components/ui/Button.vue'

const values = ref([
  {
    id: 1,
    icon: '✨',
    title: 'Excellence',
    description: 'Nous visons la perfection dans chaque détail, sans compromis sur la qualité.'
  },
  {
    id: 2,
    icon: '💝',
    title: 'Passion',
    description: 'Notre amour pour notre métier se reflète dans chaque événement que nous créons.'
  },
  {
    id: 3,
    icon: '🎯',
    title: 'Personnalisation',
    description: 'Chaque événement est unique et conçu sur mesure selon vos désirs.'
  },
  {
    id: 4,
    icon: '🤝',
    title: 'Confiance',
    description: 'Nous bâtissons des relations durables basées sur la transparence et l\'intégrité.'
  }
])

const reasons = ref([
  {
    id: 1,
    title: 'Expertise Reconnue',
    description: 'Des années d\'expérience dans la création d\'événements haut de gamme et mémorables.'
  },
  {
    id: 2,
    title: 'Service Clé en Main',
    description: 'De la conception à la réalisation, nous gérons tous les aspects de votre événement.'
  },
  {
    id: 3,
    title: 'Attention aux Détails',
    description: 'Chaque élément est soigneusement pensé et exécuté avec précision.'
  },
  {
    id: 4,
    title: 'Réseau de Partenaires',
    description: 'Accès à un réseau exclusif de fournisseurs et lieux prestigieux.'
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
