import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

// Public Views
import HomePage from '../views/public/HomePage.vue'
import ServicesPage from '../views/public/ServicesPage.vue'
import GalleryPage from '../views/public/GalleryPage.vue'
import TestimonialsPage from '../views/public/TestimonialsPage.vue'
import AboutPage from '../views/public/AboutPage.vue'
import ContactPage from '../views/public/ContactPage.vue'
import BookingPage from '../views/public/BookingPage.vue'
import RentalsPage from '../views/public/RentalsPage.vue'

// Admin Views
import AdminLogin from '../views/admin/AdminLogin.vue'
import AdminLayout from '../views/admin/AdminLayout.vue'
import AdminDashboard from '../views/admin/AdminDashboard.vue'
import BookingManager from '../views/admin/BookingManager.vue'
import ClientManager from '../views/admin/ClientManager.vue'
import ContentManager from '../views/admin/ContentManager.vue'
import GalleryManager from '../views/admin/GalleryManager.vue'
import TestimonialManager from '../views/admin/TestimonialManager.vue'
import NewsletterManager from '../views/admin/NewsletterManager.vue'
import RentalManager from '../views/admin/RentalManager.vue'

const routes = [
    // Public Routes
    {
        path: '/',
        name: 'home',
        component: HomePage,
        meta: { title: 'Angel Event - Créer l\'instant parfait' }
    },
    {
        path: '/services',
        name: 'services',
        component: ServicesPage,
        meta: { title: 'Nos Services - Angel Event' }
    },
    {
        path: '/galerie',
        name: 'gallery',
        component: GalleryPage,
        meta: { title: 'Galerie - Angel Event' }
    },
    {
        path: '/location',
        name: 'rentals',
        component: RentalsPage,
        meta: { title: 'Location - Angel Event' }
    },
    {
        path: '/temoignages',
        name: 'testimonials',
        component: TestimonialsPage,
        meta: { title: 'Témoignages - Angel Event' }
    },
    {
        path: '/a-propos',
        name: 'about',
        component: AboutPage,
        meta: { title: 'À Propos - Angel Event' }
    },
    {
        path: '/contact',
        name: 'contact',
        component: ContactPage,
        meta: { title: 'Contact - Angel Event' }
    },
    {
        path: '/reserver',
        name: 'booking',
        component: BookingPage,
        meta: { title: 'Réserver - Angel Event' }
    },

    // Admin Routes
    {
        path: '/admin/login',
        name: 'admin-login',
        component: AdminLogin,
        meta: { title: 'Admin Login - Angel Event' }
    },
    {
        path: '/admin',
        component: AdminLayout,
        meta: { requiresAuth: true },
        children: [
            {
                path: '',
                name: 'admin-dashboard',
                component: AdminDashboard,
                meta: { title: 'Dashboard - Admin' }
            },
            {
                path: 'bookings',
                name: 'admin-bookings',
                component: BookingManager,
                meta: { title: 'Réservations - Admin' }
            },
            {
                path: 'clients',
                name: 'admin-clients',
                component: ClientManager,
                meta: { title: 'Clients - Admin' }
            },
            {
                path: 'content',
                name: 'admin-content',
                component: ContentManager,
                meta: { title: 'Contenu - Admin' }
            },
            {
                path: 'gallery',
                name: 'admin-gallery',
                component: GalleryManager,
                meta: { title: 'Galerie - Admin' }
            },
            {
                path: 'rentals',
                name: 'admin-rentals',
                component: RentalManager,
                meta: { title: 'Location - Admin' }
            },
            {
                path: 'testimonials',
                name: 'admin-testimonials',
                component: TestimonialManager,
                meta: { title: 'Témoignages - Admin' }
            },
            {
                path: 'newsletter',
                name: 'admin-newsletter',
                component: NewsletterManager,
                meta: { title: 'Newsletter - Admin' }
            },
        ]
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
    scrollBehavior(to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        } else {
            return { top: 0, behavior: 'smooth' }
        }
    },
})

// Navigation guard for protected routes
router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStore()

    // Set page title
    document.title = to.meta.title || 'Angel Event'

    // Check if route requires authentication
    if (to.meta.requiresAuth) {
        if (!authStore.isAuthenticated) {
            next({ name: 'admin-login', query: { redirect: to.fullPath } })
        } else {
            // Fetch user if not already loaded
            if (!authStore.user) {
                await authStore.fetchCurrentUser()
            }
            next()
        }
    } else {
        next()
    }
})

export default router
