<template>
  <div class="services-page">
    <Header />
    <!-- Hero Section -->
    <section class="services-hero">
      <div class="container text-center">
        <h1 class="font-script text-gold fade-in-up">Nos Services</h1>
        <p class="hero-subtitle fade-in-up">Des émotions mises en scène avec passion et excellence</p>
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
            <Button @click="$router.push('/reserver')">Réserver ce service</Button>
          </div>
        </div>
      </div>
    </section>

    <!-- Process Section -->
    <section class="process-section">
      <div class="container">
        <h2 class="section-title text-center fade-in-up">Notre Processus</h2>
        <p class="section-subtitle text-center fade-in-up">Comment nous créons votre événement parfait</p>
        
        <div class="process-steps">
          <div v-for="(step, index) in processSteps" :key="step.id" class="process-step fade-in-up">
            <div class="step-number">{{ index + 1 }}</div>
            <h3>{{ step.title }}</h3>
            <p>{{ step.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta-section">
      <div class="container text-center">
        <h2 class="font-script fade-in-up">Prêt à commencer?</h2>
        <p class="fade-in-up">Contactez-nous pour une consultation gratuite</p>
        <div class="cta-buttons fade-in-up">
          <Button size="lg" @click="$router.push('/reserver')">Demander une soumission</Button>
          <Button size="lg" variant="white" @click="$router.push('/contact')">Nous contacter</Button>
        </div>
      </div>
    </section>

    <Footer />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Header from '../../components/Header.vue'
import Footer from '../../components/Footer.vue'
import Button from '../../components/ui/Button.vue'

const services = ref([
  {
    id: 1,
    icon: '💍',
    title: 'Demandes en mariage',
    description: 'Créez le moment le plus mémorable de votre vie avec une mise en scène romantique et inoubliable.',
    image: 'https://images.unsplash.com/photo-1519741497674-611481863552?w=800',
    features: [
      'Décors "Marry Me" personnalisés',
      'Mise en scène romantique sur mesure',
      'Bougies, fleurs et arches élégantes',
      'Lettres lumineuses et néons',
      'Coordination complète le jour J',
      'Photographie professionnelle disponible'
    ]
  },
  {
    id: 2,
    icon: '🎨',
    title: 'Décoration événementielle',
    description: 'Transformez vos événements avec des décors élégants, raffinés et entièrement personnalisés selon vos goûts.',
    image: 'https://images.unsplash.com/photo-1464366400600-7168b8af9bc3?w=800',
    features: [
      'Anniversaires mémorables',
      'Baby showers élégants',
      'Événements corporatifs',
      'Fêtes privées haut de gamme',
      'Thèmes personnalisés',
      'Installation et démontage inclus'
    ]
  },
  {
    id: 3,
    icon: '📋',
    title: 'Planification complète',
    description: 'Laissez-nous gérer chaque détail de votre événement pour une expérience sans stress et parfaitement orchestrée.',
    image: 'https://images.unsplash.com/photo-1511795409834-ef04bbd61622?w=800',
    features: [
      'Accompagnement personnalisé',
      'Gestion du budget',
      'Coordination des fournisseurs',
      'Timeline détaillée',
      'Gestion le jour J',
      'Support 24/7 avant l\'événement'
    ]
  }
])

const processSteps = ref([
  {
    id: 1,
    title: 'Consultation',
    description: 'Rencontre initiale pour comprendre votre vision et vos besoins'
  },
  {
    id: 2,
    title: 'Planification',
    description: 'Création d\'un plan détaillé et personnalisé pour votre événement'
  },
  {
    id: 3,
    title: 'Préparation',
    description: 'Coordination de tous les détails et préparation minutieuse'
  },
  {
    id: 4,
    title: 'Réalisation',
    description: 'Mise en place et gestion complète le jour de l\'événement'
  }
])

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

/* Process Section */
.process-section {
  padding: var(--spacing-5xl) 0;
  background: var(--color-gray-lighter);
}

.section-title {
  font-size: var(--font-size-4xl);
  margin-bottom: var(--spacing-md);
}

.section-subtitle {
  font-size: var(--font-size-lg);
  color: var(--color-gray);
  margin-bottom: var(--spacing-4xl);
}

.process-steps {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-2xl);
  margin-top: var(--spacing-3xl);
}

.process-step {
  background: var(--color-white);
  padding: var(--spacing-2xl);
  border-radius: var(--radius-lg);
  text-align: center;
  box-shadow: var(--shadow-md);
  transition: transform var(--transition-base);
}

.process-step:hover {
  transform: translateY(-5px);
}

.step-number {
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

.process-step h3 {
  color: var(--color-gold);
  margin-bottom: var(--spacing-md);
}

.process-step p {
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
  .service-detail {
    grid-template-columns: 1fr;
  }
  
  .service-detail.reverse {
    direction: ltr;
  }
  
  .process-steps {
    grid-template-columns: 1fr;
  }
  
  .cta-buttons {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
