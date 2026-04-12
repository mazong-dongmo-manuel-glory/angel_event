# Angel Event - Site Web Professionnel

Un site web élégant et haut de gamme pour Angel Event, spécialisé dans la planification d'événements, les demandes en mariage et la décoration événementielle.

## 🎨 Caractéristiques

### Frontend (Vue.js 3)
- ✨ Design élégant avec thème blanc/noir/doré
- 🎭 Animations fluides et micro-interactions premium
- 📱 Responsive (mobile, tablette, desktop)
- 🌐 Multilingue (FR/EN prêt)
- ⚡ Performance optimisée avec Vite
- 🎨 Design system complet avec variables CSS
- 💳 Intégration Stripe pour les paiements
- 📧 Système de newsletter
- 📅 Réservation avec calendrier de disponibilités

### Backend (Golang Fiber)
- 🚀 API REST rapide et performante
- 🔐 Authentification JWT sécurisée
- 💾 Base de données SQLite (évolutif vers PostgreSQL)
- 💳 Intégration Stripe complète
- 📧 Service d'envoi d'emails avec templates HTML
- 📊 Panel d'administration complet
- 🔒 Middleware de sécurité (CORS, rate limiting)
- 📝 Gestion de contenu dynamique

## 📁 Structure du Projet

\`\`\`
angel_event/
├── backend/                 # API Golang Fiber
│   ├── cmd/server/         # Point d'entrée
│   ├── internal/
│   │   ├── database/       # Configuration DB
│   │   ├── models/         # Modèles de données
│   │   ├── handlers/       # Gestionnaires de routes
│   │   ├── services/       # Services (email, Stripe)
│   │   └── middleware/     # Middleware (auth, CORS)
│   ├── go.mod
│   └── .env
│
└── frontend/               # Application Vue.js
    ├── src/
    │   ├── assets/styles/  # CSS variables & animations
    │   ├── components/ui/  # Composants réutilisables
    │   ├── views/
    │   │   ├── public/     # Pages publiques
    │   │   └── admin/      # Panel admin
    │   ├── router/         # Configuration routes
    │   ├── stores/         # Pinia stores
    │   ├── services/       # API client
    │   └── main.js
    └── package.json
\`\`\`

## 🚀 Installation

### Prérequis
- Go 1.21+
- Node.js 18+
- npm ou yarn

### Backend

\`\`\`bash
cd backend

# Copier le fichier d'environnement
cp .env.example .env

# Éditer .env avec vos configurations
# - JWT_SECRET
# - STRIPE_SECRET_KEY
# - SMTP credentials

# Installer les dépendances
go mod download

# Lancer le serveur
go run cmd/server/main.go
\`\`\`

Le serveur démarre sur `http://localhost:8080`

### Frontend

\`\`\`bash
cd frontend

# Copier le fichier d'environnement
cp .env.example .env

# Installer les dépendances
npm install

# Lancer le serveur de développement
npm run dev
\`\`\`

L'application démarre sur `http://localhost:5173`

## 🔑 Configuration

### Variables d'Environnement Backend

\`\`\`env
PORT=8080
DATABASE_PATH=./angel_event.db
JWT_SECRET=your-secret-key
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
FRONTEND_URL=http://localhost:5173
ADMIN_EMAIL=admin@angelevent.ca
ADMIN_PASSWORD=ChangeThisPassword123!
\`\`\`

### Variables d'Environnement Frontend

\`\`\`env
VITE_API_URL=http://localhost:8080/api
VITE_STRIPE_PUBLISHABLE_KEY=pk_test_...
\`\`\`

## 👤 Connexion Admin

Par défaut:
- **Email**: admin@angelevent.ca
- **Mot de passe**: ChangeThisPassword123!

⚠️ **IMPORTANT**: Changez ces identifiants après la première connexion!

## 📋 Fonctionnalités

### Pages Publiques
- **Accueil**: Hero section WOW, aperçu services, témoignages
- **Services**: Détails des services offerts
- **Galerie**: Photos d'événements avec filtres
- **Témoignages**: Avis clients
- **À Propos**: Histoire et valeurs
- **Contact**: Formulaire de contact
- **Réservation**: Système de réservation avec paiement Stripe

### Panel d'Administration
- **Dashboard**: Statistiques et aperçu
- **Réservations**: Gestion complète des bookings
- **Clients**: Base de données clients avec envoi d'emails
- **Contenu**: Éditeur de contenu du site
- **Galerie**: Upload et gestion des images
- **Témoignages**: Modération et approbation
- **Newsletter**: Gestion des abonnés et envoi de campagnes

## 🎨 Design System

### Couleurs
- **Ivoire**: #FFFEF7 (fond principal)
- **Champagne**: #F7E7CE (accents doux)
- **Or**: #D4AF37 (primaire)
- **Noir doux**: #1A1A1A (texte)

### Typographies
- **Script**: Great Vibes (logo, titres élégants)
- **Sans-serif**: Inter (texte principal)

### Animations
- Fade in/out
- Slide transitions
- Hover effects premium
- Scroll reveal
- Glassmorphism

## 🛠️ Technologies

### Frontend
- Vue 3 (Composition API)
- Vue Router 4
- Pinia (state management)
- Axios
- Stripe.js
- VueUse

### Backend
- Golang 1.21
- Fiber v2 (framework web)
- GORM (ORM)
- SQLite
- JWT
- Stripe Go SDK
- Gomail (emails)

## 📦 Build Production

### Backend
\`\`\`bash
cd backend
go build -o angel-event-api cmd/server/main.go
./angel-event-api
\`\`\`

### Frontend
\`\`\`bash
cd frontend
npm run build
# Les fichiers sont dans dist/
\`\`\`

## 🚢 Déploiement

### Recommandations
- **Backend**: Railway, Fly.io, DigitalOcean
- **Frontend**: Vercel, Netlify, Cloudflare Pages
- **Base de données**: Migrer vers PostgreSQL en production
- **Fichiers**: S3, Cloudinary pour les images

## 📧 Configuration Email

Pour Gmail:
1. Activer l'authentification à 2 facteurs
2. Générer un mot de passe d'application
3. Utiliser ce mot de passe dans SMTP_PASSWORD

## 💳 Configuration Stripe

1. Créer un compte Stripe
2. Récupérer les clés API (test et production)
3. Configurer les webhooks pour `/api/webhooks/stripe`
4. Ajouter les clés dans les fichiers .env

## 🔒 Sécurité

- JWT pour l'authentification
- Mots de passe hashés avec bcrypt
- CORS configuré
- Validation des entrées
- Protection CSRF
- Rate limiting (à implémenter)

## 📝 API Endpoints

### Public
- `POST /api/public/contact` - Formulaire de contact
- `POST /api/public/bookings` - Créer une réservation
- `GET /api/public/gallery` - Images galerie
- `POST /api/public/newsletter/subscribe` - S'abonner

### Admin (authentifié)
- `GET /api/admin/dashboard/stats` - Statistiques
- `GET /api/admin/bookings` - Liste réservations
- `GET /api/admin/clients` - Liste clients
- `POST /api/admin/newsletter/send` - Envoyer newsletter

## 🤝 Contribution

Ce projet est privé. Pour toute question, contactez l'équipe de développement.

## 📄 Licence

Propriétaire - Angel Event © 2026

---

**Développé avec ❤️ pour Angel Event**
