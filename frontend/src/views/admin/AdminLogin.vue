<template>
  <div class="admin-login">
    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <h1 class="font-script text-gold">Angel Event</h1>
          <p>Administration</p>
        </div>

        <form v-if="!loading" @submit.prevent="handleLogin" class="login-form">
          <div class="form-group">
            <label for="email">Email</label>
            <input
              id="email"
              v-model="credentials.email"
              type="email"
              required
              placeholder="admin@angelevent.com"
              autocomplete="username"
            />
          </div>

          <div class="form-group">
            <label for="password">Mot de passe</label>
            <input
              id="password"
              v-model="credentials.password"
              type="password"
              required
              placeholder="••••••••"
              autocomplete="current-password"
            />
          </div>

          <div v-if="error" class="error-message">
            {{ error }}
          </div>

          <Button type="submit" size="lg" block :loading="loading">
            Se connecter
          </Button>
        </form>

        <div v-else class="loading-state">
          <div class="spinner"></div>
          <p>Connexion en cours...</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import Button from '../../components/ui/Button.vue'

const router = useRouter()
const authStore = useAuthStore()

const credentials = ref({
  email: '',
  password: ''
})

const loading = ref(false)
const error = ref(null)

async function handleLogin() {
  loading.value = true
  error.value = null

  const success = await authStore.login(credentials.value.email, credentials.value.password)

  if (success) {
    const redirect = router.currentRoute.value.query.redirect || '/admin'
    router.push(redirect)
  } else {
    error.value = authStore.error || 'Identifiants invalides'
  }

  loading.value = false
}
</script>

<style scoped>
.admin-login {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--color-ivory) 0%, var(--color-champagne) 100%);
  padding: var(--spacing-xl);
}

.login-container {
  width: 100%;
  max-width: 450px;
}

.login-card {
  background: var(--color-white);
  border-radius: var(--radius-xl);
  padding: var(--spacing-3xl);
  box-shadow: var(--shadow-2xl);
}

.login-header {
  text-align: center;
  margin-bottom: var(--spacing-3xl);
}

.login-header h1 {
  font-size: var(--font-size-5xl);
  margin-bottom: var(--spacing-sm);
}

.login-header p {
  color: var(--color-gray);
  font-size: var(--font-size-lg);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.form-group label {
  font-weight: var(--font-weight-medium);
  color: var(--color-text);
}

.form-group input {
  padding: var(--spacing-md);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-sans);
  font-size: var(--font-size-base);
  transition: border-color var(--transition-base);
}

.form-group input:focus {
  outline: none;
  border-color: var(--color-gold);
}

.error-message {
  background: rgba(193, 41, 46, 0.1);
  color: var(--color-error);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  text-align: center;
  font-size: var(--font-size-sm);
}

.loading-state {
  text-align: center;
  padding: var(--spacing-2xl);
  color: var(--color-gray);
}
</style>
