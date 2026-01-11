<template>
  <div class="testimonial-manager">
    <!-- Header -->
    <div class="manager-header">
      <h1>Gestion des Témoignages</h1>
      <p class="subtitle">Gérez les témoignages clients et leur approbation</p>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">💬</div>
        <div class="stat-content">
          <h3>{{ stats.total }}</h3>
          <p>Total Témoignages</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">✅</div>
        <div class="stat-content">
          <h3>{{ stats.approved }}</h3>
          <p>Approuvés</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">⏳</div>
        <div class="stat-content">
          <h3>{{ stats.pending }}</h3>
          <p>En attente</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">⭐</div>
        <div class="stat-content">
          <h3>{{ stats.featured }}</h3>
          <p>En vedette</p>
        </div>
      </div>
    </div>

    <!-- Controls -->
    <div class="controls">
      <div class="filter-buttons">
        <button
          :class="['filter-btn', { active: filter === 'all' }]"
          @click="filter = 'all'"
        >
          Tous
        </button>
        <button
          :class="['filter-btn', { active: filter === 'pending' }]"
          @click="filter = 'pending'"
        >
          En attente
        </button>
        <button
          :class="['filter-btn', { active: filter === 'approved' }]"
          @click="filter = 'approved'"
        >
          Approuvés
        </button>
      </div>
    </div>

    <!-- Testimonials Grid -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement des témoignages...</p>
    </div>

    <div v-else-if="filteredTestimonials.length === 0" class="empty-state">
      <p>Aucun témoignage trouvé</p>
    </div>

    <div v-else class="testimonials-grid">
      <div
        v-for="testimonial in filteredTestimonials"
        :key="testimonial.id"
        class="testimonial-card"
      >
        <div class="card-header">
          <div class="author-info">
            <strong>{{ testimonial.name }}</strong>
            <div class="rating">
              <span v-for="i in 5" :key="i" class="star" :class="{ filled: i <= testimonial.rating }">
                ⭐
              </span>
            </div>
          </div>
          <div class="badges">
            <span v-if="testimonial.approved" class="badge badge-success">Approuvé</span>
            <span v-else class="badge badge-warning">En attente</span>
            <span v-if="testimonial.featured" class="badge badge-featured">Vedette</span>
          </div>
        </div>

        <div class="card-body">
          <p class="content">{{ testimonial.content }}</p>
          <div class="meta">
            <span class="event-type">{{ getEventTypeLabel(testimonial.event_type) }}</span>
            <span class="date">{{ formatDate(testimonial.created_at) }}</span>
          </div>
        </div>

        <div class="card-actions">
          <button
            v-if="!testimonial.approved"
            @click="approveTestimonial(testimonial)"
            class="btn-action btn-success"
          >
            ✓ Approuver
          </button>
          <button
            @click="toggleFeatured(testimonial)"
            class="btn-action"
            :class="testimonial.featured ? 'btn-warning' : 'btn-secondary'"
          >
            {{ testimonial.featured ? '★ Retirer vedette' : '☆ Mettre en vedette' }}
          </button>
          <button @click="editTestimonial(testimonial)" class="btn-action btn-secondary">
            ✏️ Modifier
          </button>
          <button @click="confirmDelete(testimonial)" class="btn-action btn-danger">
            🗑️ Supprimer
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Testimonial Modal -->
    <Teleport to="body">
      <div v-if="editingTestimonial" class="modal-overlay" @click="editingTestimonial = null">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>Modifier le Témoignage</h2>
            <button @click="editingTestimonial = null" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveTestimonial" class="edit-form">
              <div class="form-group">
                <label>Nom *</label>
                <input v-model="editForm.name" type="text" required />
              </div>
              <div class="form-group">
                <label>Contenu *</label>
                <textarea v-model="editForm.content" rows="6" required></textarea>
              </div>
              <div class="form-row">
                <div class="form-group">
                  <label>Note</label>
                  <select v-model.number="editForm.rating">
                    <option :value="5">5 étoiles</option>
                    <option :value="4">4 étoiles</option>
                    <option :value="3">3 étoiles</option>
                    <option :value="2">2 étoiles</option>
                    <option :value="1">1 étoile</option>
                  </select>
                </div>
                <div class="form-group">
                  <label>Type d'événement</label>
                  <select v-model="editForm.event_type">
                    <option value="proposal">Demande en mariage</option>
                    <option value="wedding">Mariage</option>
                    <option value="birthday">Anniversaire</option>
                    <option value="baby_shower">Baby Shower</option>
                    <option value="corporate">Corporatif</option>
                    <option value="other">Autre</option>
                  </select>
                </div>
              </div>
              <div class="form-group">
                <label class="checkbox-label">
                  <input v-model="editForm.approved" type="checkbox" />
                  Approuvé
                </label>
              </div>
              <div class="form-group">
                <label class="checkbox-label">
                  <input v-model="editForm.featured" type="checkbox" />
                  En vedette
                </label>
              </div>
              <div class="form-actions">
                <button type="button" @click="editingTestimonial = null" class="btn-secondary">
                  Annuler
                </button>
                <button type="submit" class="btn-primary">
                  Enregistrer
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <div v-if="deletingTestimonial" class="modal-overlay" @click="deletingTestimonial = null">
        <div class="modal-content modal-small" @click.stop>
          <div class="modal-header">
            <h2>Confirmer la suppression</h2>
            <button @click="deletingTestimonial = null" class="close-btn">×</button>
          </div>
          <div class="modal-body">
            <p>Êtes-vous sûr de vouloir supprimer ce témoignage de <strong>{{ deletingTestimonial.name }}</strong> ?</p>
            <p class="warning-text">Cette action est irréversible.</p>
            <div class="form-actions">
              <button @click="deletingTestimonial = null" class="btn-secondary">
                Annuler
              </button>
              <button @click="deleteTestimonial" class="btn-danger">
                Supprimer
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../../services/api'

const testimonials = ref([])
const loading = ref(false)
const filter = ref('all')
const editingTestimonial = ref(null)
const deletingTestimonial = ref(null)

const editForm = ref({
  name: '',
  content: '',
  rating: 5,
  event_type: 'wedding',
  approved: false,
  featured: false
})

const eventTypeLabels = {
  proposal: '💍 Demande en mariage',
  wedding: '💐 Mariage',
  birthday: '🎂 Anniversaire',
  baby_shower: '👶 Baby Shower',
  corporate: '🏢 Corporatif',
  other: '✨ Autre'
}

const stats = computed(() => {
  const total = testimonials.value.length
  const approved = testimonials.value.filter(t => t.approved).length
  const pending = testimonials.value.filter(t => !t.approved).length
  const featured = testimonials.value.filter(t => t.featured).length
  
  return { total, approved, pending, featured }
})

const filteredTestimonials = computed(() => {
  if (filter.value === 'all') return testimonials.value
  if (filter.value === 'pending') return testimonials.value.filter(t => !t.approved)
  if (filter.value === 'approved') return testimonials.value.filter(t => t.approved)
  return testimonials.value
})

async function fetchTestimonials() {
  loading.value = true
  try {
    const response = await api.get('/admin/testimonials')
    testimonials.value = response.data
  } catch (error) {
    console.error('Failed to fetch testimonials:', error)
  } finally {
    loading.value = false
  }
}

function getEventTypeLabel(type) {
  return eventTypeLabels[type] || type
}

async function approveTestimonial(testimonial) {
  try {
    await api.put(`/admin/testimonials/${testimonial.id}`, {
      ...testimonial,
      approved: true
    })
    testimonial.approved = true
  } catch (error) {
    console.error('Failed to approve testimonial:', error)
    alert('Erreur lors de l\'approbation')
  }
}

async function toggleFeatured(testimonial) {
  try {
    await api.put(`/admin/testimonials/${testimonial.id}`, {
      ...testimonial,
      featured: !testimonial.featured
    })
    testimonial.featured = !testimonial.featured
  } catch (error) {
    console.error('Failed to toggle featured:', error)
    alert('Erreur lors de la modification')
  }
}

function editTestimonial(testimonial) {
  editingTestimonial.value = testimonial
  editForm.value = {
    name: testimonial.name,
    content: testimonial.content,
    rating: testimonial.rating,
    event_type: testimonial.event_type,
    approved: testimonial.approved,
    featured: testimonial.featured
  }
}

async function saveTestimonial() {
  try {
    await api.put(`/admin/testimonials/${editingTestimonial.value.id}`, editForm.value)
    
    // Update local data
    const index = testimonials.value.findIndex(t => t.id === editingTestimonial.value.id)
    if (index !== -1) {
      testimonials.value[index] = { ...testimonials.value[index], ...editForm.value }
    }
    
    editingTestimonial.value = null
  } catch (error) {
    console.error('Failed to update testimonial:', error)
    alert('Erreur lors de la mise à jour')
  }
}

function confirmDelete(testimonial) {
  deletingTestimonial.value = testimonial
}

async function deleteTestimonial() {
  try {
    await api.delete(`/admin/testimonials/${deletingTestimonial.value.id}`)
    
    // Remove from local data
    testimonials.value = testimonials.value.filter(t => t.id !== deletingTestimonial.value.id)
    
    deletingTestimonial.value = null
  } catch (error) {
    console.error('Failed to delete testimonial:', error)
    alert('Erreur lors de la suppression')
  }
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('fr-FR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

onMounted(() => {
  fetchTestimonials()
})
</script>

<style scoped>
.testimonial-manager {
  padding: var(--spacing-3xl);
  max-width: 1400px;
  margin: 0 auto;
}

.manager-header {
  margin-bottom: var(--spacing-3xl);
}

.manager-header h1 {
  color: var(--color-gold);
  font-size: var(--font-size-4xl);
  margin-bottom: var(--spacing-sm);
}

.subtitle {
  color: var(--color-gray);
  font-size: var(--font-size-lg);
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-xl);
  margin-bottom: var(--spacing-3xl);
}

.stat-card {
  background: white;
  padding: var(--spacing-xl);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  transition: transform var(--transition-base);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.stat-icon {
  font-size: var(--font-size-4xl);
}

.stat-content h3 {
  font-size: var(--font-size-3xl);
  color: var(--color-gold);
  margin-bottom: var(--spacing-xs);
}

.stat-content p {
  color: var(--color-gray);
  font-size: var(--font-size-sm);
}

/* Controls */
.controls {
  margin-bottom: var(--spacing-2xl);
}

.filter-buttons {
  display: flex;
  gap: var(--spacing-md);
}

.filter-btn {
  padding: var(--spacing-md) var(--spacing-xl);
  border: 2px solid var(--color-border);
  background: white;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-base);
  font-weight: var(--font-weight-semibold);
  text-transform: uppercase;
  font-size: var(--font-size-sm);
  letter-spacing: 0.05em;
}

.filter-btn:hover {
  border-color: var(--color-gold);
}

.filter-btn.active {
  background: var(--color-gold);
  border-color: var(--color-gold);
  color: white;
}

/* Testimonials Grid */
.testimonials-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: var(--spacing-xl);
}

.testimonial-card {
  background: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  overflow: hidden;
  transition: transform var(--transition-base);
}

.testimonial-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.card-header {
  padding: var(--spacing-lg);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.author-info strong {
  display: block;
  color: var(--color-text);
  margin-bottom: var(--spacing-xs);
}

.rating {
  display: flex;
  gap: 2px;
}

.star {
  font-size: var(--font-size-sm);
  opacity: 0.3;
}

.star.filled {
  opacity: 1;
}

.badges {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  align-items: flex-end;
}

.badge {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  text-transform: uppercase;
}

.badge-success {
  background: #2D7A4F;
  color: white;
}

.badge-warning {
  background: #FFA500;
  color: white;
}

.badge-featured {
  background: var(--color-gold);
  color: white;
}

.card-body {
  padding: var(--spacing-lg);
}

.content {
  color: var(--color-text);
  line-height: var(--line-height-relaxed);
  margin-bottom: var(--spacing-md);
}

.meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: var(--font-size-sm);
  color: var(--color-gray);
}

.card-actions {
  padding: var(--spacing-md);
  background: var(--color-gray-lighter);
  display: flex;
  gap: var(--spacing-sm);
  flex-wrap: wrap;
}

.btn-action {
  padding: var(--spacing-xs) var(--spacing-md);
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  transition: all var(--transition-base);
}

.btn-success {
  background: #2D7A4F;
  color: white;
}

.btn-success:hover {
  background: #246339;
}

.btn-secondary {
  background: var(--color-gray-light);
  color: var(--color-text);
}

.btn-secondary:hover {
  background: var(--color-gray);
}

.btn-warning {
  background: var(--color-gold);
  color: white;
}

.btn-warning:hover {
  background: #C5A028;
}

.btn-danger {
  background: #C1292E;
  color: white;
}

.btn-danger:hover {
  background: #A01F23;
}

/* Loading & Empty States */
.loading-state,
.empty-state {
  text-align: center;
  padding: var(--spacing-5xl);
  color: var(--color-gray);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--color-border);
  border-top-color: var(--color-gold);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto var(--spacing-md);
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Modals */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal);
  padding: var(--spacing-xl);
}

.modal-content {
  background: white;
  border-radius: var(--radius-xl);
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: var(--shadow-xl);
}

.modal-small {
  max-width: 400px;
}

.modal-header {
  padding: var(--spacing-xl);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h2 {
  color: var(--color-gold);
  font-size: var(--font-size-2xl);
}

.close-btn {
  background: none;
  border: none;
  font-size: var(--font-size-3xl);
  cursor: pointer;
  color: var(--color-gray);
  line-height: 1;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  transition: background var(--transition-base);
}

.close-btn:hover {
  background: var(--color-gray-lighter);
}

.modal-body {
  padding: var(--spacing-xl);
}

/* Forms */
.edit-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-lg);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.form-group label {
  font-weight: var(--font-weight-semibold);
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
  transition: border-color var(--transition-base);
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--color-gold);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: auto;
  cursor: pointer;
}

.form-actions {
  display: flex;
  gap: var(--spacing-md);
  justify-content: flex-end;
  margin-top: var(--spacing-lg);
}

.btn-primary {
  padding: var(--spacing-md) var(--spacing-xl);
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-semibold);
  cursor: pointer;
  transition: all var(--transition-base);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-size: var(--font-size-sm);
  background: var(--color-gold);
  color: white;
}

.btn-primary:hover {
  background: #C5A028;
}

.warning-text {
  color: #C1292E;
  font-size: var(--font-size-sm);
  margin-top: var(--spacing-sm);
}

@media (max-width: 768px) {
  .testimonial-manager {
    padding: var(--spacing-xl);
  }

  .testimonials-grid {
    grid-template-columns: 1fr;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .filter-buttons {
    flex-wrap: wrap;
  }
}
</style>
