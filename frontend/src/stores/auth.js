import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'

export const useAuthStore = defineStore('auth', () => {
    const user = ref(null)
    const token = ref(localStorage.getItem('auth_token'))
    const loading = ref(false)
    const error = ref(null)

    const isAuthenticated = computed(() => !!token.value)
    const isAdmin = computed(() => user.value?.role === 'admin')

    async function login(email, password) {
        loading.value = true
        error.value = null
        try {
            const response = await api.post('/auth/login', { email, password })
            token.value = response.data.token
            user.value = response.data.user
            localStorage.setItem('auth_token', response.data.token)
            return true
        } catch (err) {
            error.value = err.response?.data?.error || 'Login failed'
            return false
        } finally {
            loading.value = false
        }
    }

    async function logout() {
        token.value = null
        user.value = null
        localStorage.removeItem('auth_token')
    }

    async function fetchCurrentUser() {
        if (!token.value) return

        try {
            const response = await api.get('/auth/me')
            user.value = response.data
        } catch (err) {
            console.error('Failed to fetch user:', err)
            logout()
        }
    }

    async function changePassword(currentPassword, newPassword) {
        loading.value = true
        error.value = null
        try {
            await api.post('/auth/change-password', {
                current_password: currentPassword,
                new_password: newPassword,
            })
            return true
        } catch (err) {
            error.value = err.response?.data?.error || 'Failed to change password'
            return false
        } finally {
            loading.value = false
        }
    }

    return {
        user,
        token,
        loading,
        error,
        isAuthenticated,
        isAdmin,
        login,
        logout,
        fetchCurrentUser,
        changePassword,
    }
})
