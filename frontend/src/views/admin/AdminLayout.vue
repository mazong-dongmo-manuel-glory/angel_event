<template>
  <div class="admin-layout">
    <aside class="sidebar">
      <div class="sidebar-header">
        <h2 class="font-script">Angel Event</h2>
        <p>Administration</p>
      </div>

      <nav class="sidebar-nav">
        <RouterLink
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          active-class="active"
        >
          <span class="nav-icon">{{ item.icon }}</span>
          <span class="nav-label">{{ item.label }}</span>
        </RouterLink>
      </nav>

      <div class="sidebar-footer">
        <button class="logout-btn" @click="handleLogout">
          <span>🚪</span> Déconnexion
        </button>
      </div>
    </aside>

    <main class="main-content">
      <header class="top-bar">
        <h1>{{ currentPageTitle }}</h1>
        <div class="user-info">
          <span>{{ authStore.user?.name || 'Admin' }}</span>
        </div>
      </header>

      <div class="content-area">
        <RouterView />
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { RouterLink, RouterView, useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const navItems = [
  { path: '/admin', icon: '📊', label: 'Dashboard' },
  { path: '/admin/bookings', icon: '📅', label: 'Réservations' },
  { path: '/admin/clients', icon: '👥', label: 'Clients' },
  { path: '/admin/gallery', icon: '🖼️', label: 'Galerie' },
  { path: '/admin/testimonials', icon: '💬', label: 'Témoignages' },
  { path: '/admin/newsletter', icon: '📧', label: 'Newsletter' },
  { path: '/admin/content', icon: '✏️', label: 'Contenu' },
]

const currentPageTitle = computed(() => {
  const item = navItems.find(item => item.path === route.path)
  return item?.label || 'Administration'
})

async function handleLogout() {
  await authStore.logout()
  router.push('/admin/login')
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: var(--color-gray-lighter);
}

/* Sidebar */
.sidebar {
  width: 260px;
  background: var(--color-black-soft);
  color: var(--color-white);
  display: flex;
  flex-direction: column;
  position: fixed;
  height: 100vh;
  left: 0;
  top: 0;
}

.sidebar-header {
  padding: var(--spacing-2xl);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  text-align: center;
}

.sidebar-header h2 {
  color: var(--color-gold);
  font-size: var(--font-size-3xl);
  margin-bottom: var(--spacing-xs);
}

.sidebar-header p {
  color: var(--color-champagne);
  font-size: var(--font-size-sm);
}

.sidebar-nav {
  flex: 1;
  padding: var(--spacing-lg) 0;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md) var(--spacing-xl);
  color: var(--color-white);
  text-decoration: none;
  transition: all var(--transition-base);
  border-left: 3px solid transparent;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  border-left-color: var(--color-gold);
}

.nav-item.active {
  background: rgba(212, 175, 55, 0.1);
  border-left-color: var(--color-gold);
  color: var(--color-gold);
}

.nav-icon {
  font-size: var(--font-size-xl);
}

.nav-label {
  font-weight: var(--font-weight-medium);
}

.sidebar-footer {
  padding: var(--spacing-xl);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.logout-btn {
  width: 100%;
  padding: var(--spacing-md);
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: var(--color-white);
  border-radius: var(--radius-md);
  cursor: pointer;
  font-family: var(--font-sans);
  font-size: var(--font-size-base);
  transition: all var(--transition-base);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

/* Main Content */
.main-content {
  flex: 1;
  margin-left: 260px;
  display: flex;
  flex-direction: column;
}

.top-bar {
  background: var(--color-white);
  padding: var(--spacing-xl) var(--spacing-2xl);
  box-shadow: var(--shadow-sm);
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: var(--z-sticky);
}

.top-bar h1 {
  font-size: var(--font-size-2xl);
  color: var(--color-black-soft);
}

.user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  color: var(--color-gray);
}

.content-area {
  flex: 1;
  padding: var(--spacing-2xl);
}

@media (max-width: 768px) {
  .sidebar {
    width: 70px;
  }
  
  .sidebar-header p,
  .nav-label {
    display: none;
  }
  
  .main-content {
    margin-left: 70px;
  }
  
  .nav-item {
    justify-content: center;
  }
}
</style>
