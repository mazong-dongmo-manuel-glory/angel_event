package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mazong/angel_event/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

const (
	defaultAdminEmail    = "admin@angelevent.com"
	defaultAdminPassword = "ChangeThisPassword123!"
)

// Connect initializes the database connection
func Connect() error {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./angel_event.db"
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// Migrate runs database migrations
func Migrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.Booking{},
		&models.Availability{},
		&models.Testimonial{},
		&models.Newsletter{},
		&models.GalleryImage{},
		&models.SiteContent{},
		&models.EmailLog{},
		&models.Category{},
		&models.RentalItem{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	if err := migrateSiteContentSchema(); err != nil {
		return err
	}

	log.Println("Database migrated successfully")
	return nil
}

// SeedDefaultData creates initial data if needed
func SeedDefaultData() error {
	// Check if admin user exists
	var count int64
	DB.Model(&models.User{}).Count(&count)

	if count == 0 {
		// Create default admin user
		adminEmail := getEnvOrDefault("ADMIN_EMAIL", defaultAdminEmail)
		adminPassword := getEnvOrDefault("ADMIN_PASSWORD", defaultAdminPassword)
		hashedPassword, err := hashPassword(adminPassword)
		if err != nil {
			return err
		}

		admin := models.User{
			Email:    adminEmail,
			Password: hashedPassword,
			Name:     "Administrator",
			Role:     "admin",
		}

		if err := DB.Create(&admin).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		log.Printf("Default admin user created: %s", adminEmail)
	}

	// Seed default categories
	var catCount int64
	DB.Model(&models.Category{}).Count(&catCount)
	if catCount == 0 {
		categories := []models.Category{
			{Name: "Centre de table", Slug: "centerpiece", Type: "rental"},
			{Name: "Backdrop / Panneau", Slug: "backdrop", Type: "rental"},
			{Name: "Fleur", Slug: "flower", Type: "rental"},
			{Name: "Autre article", Slug: "other", Type: "rental"},
			{Name: "Animation", Slug: "animation", Type: "rental", Description: "Articles d'animation"},
		}
		if err := DB.Create(&categories).Error; err != nil {
			return fmt.Errorf("failed to seed default categories: %w", err)
		}
		log.Println("Default categories seeded")
	}

	if err := seedMissingSiteContentEntries(defaultSiteContentEntries()); err != nil {
		return err
	}

	// Seed default rental items
	var rentalCount int64
	DB.Model(&models.RentalItem{}).Count(&rentalCount)

	if rentalCount == 0 {
		// Get categories first
		var cats []models.Category
		DB.Find(&cats)
		catMap := make(map[string]uint)
		for _, c := range cats {
			catMap[c.Slug] = c.ID
		}

		defaultRentals := []models.RentalItem{
			{
				Title:       "Arche Ronde Dorée",
				Description: "Une magnifique arche ronde dorée pour sublimer votre décor de cérémonie ou photobooth.",
				Price:       50.00,
				CategoryID:  catMap["backdrop"],
				ImageURL:    "/uploads/rentals/default-rental.png",
				Featured:    true,
				Available:   true,
				Stock:       2,
			},
			{
				Title:       "Centre de Table Floral",
				Description: "Composition florale élégante dans des tons blanc et crème, idéale pour les mariages.",
				Price:       25.00,
				CategoryID:  catMap["centerpiece"],
				ImageURL:    "/uploads/rentals/default-rental.png",
				Featured:    true,
				Available:   true,
				Stock:       2,
			},
			{
				Title:       "Mur de Fleurs Blanc",
				Description: "Mur de fleurs artificielles haute qualité, dimensions 2m x 2m. Impact visuel garanti.",
				Price:       120.00,
				CategoryID:  catMap["backdrop"],
				ImageURL:    "/uploads/rentals/default-rental.png",
				Featured:    false,
				Available:   true,
				Stock:       2,
			},
			{
				Title:       "Vase Haut Cylindrique",
				Description: "Vase en verre transparent, hauteur 60cm. Parfait pour les compositions modernes.",
				Price:       15.00,
				CategoryID:  catMap["other"],
				ImageURL:    "/uploads/rentals/default-rental.png",
				Featured:    false,
				Available:   true,
				Stock:       2,
			},
			{
				Title:       "Bouquet de Pivoines",
				Description: "Bouquet de pivoines artificielles réalistes.",
				Price:       10.00,
				CategoryID:  catMap["flower"],
				ImageURL:    "/uploads/rentals/default-rental.png",
				Featured:    false,
				Available:   true,
				Stock:       2,
			},
		}

		if err := DB.Create(&defaultRentals).Error; err != nil {
			return fmt.Errorf("failed to seed default rentals: %w", err)
		}

		log.Println("Default rental items seeded")
	}

	return nil
}

// SyncAdminUserCredentials keeps the configured admin password usable and
// repairs legacy databases where the default admin password was stored in plain text.
func SyncAdminUserCredentials() error {
	adminEmail, emailSet := os.LookupEnv("ADMIN_EMAIL")
	adminPassword, passwordSet := os.LookupEnv("ADMIN_PASSWORD")

	if emailSet && passwordSet && adminEmail != "" && adminPassword != "" {
		return syncAdminPassword(adminEmail, adminPassword, false)
	}

	return syncAdminPassword(defaultAdminEmail, defaultAdminPassword, true)
}

func syncAdminPassword(email, rawPassword string, onlyIfPlaintextMatch bool) error {
	var user models.User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("failed to find admin user: %w", err)
	}

	if onlyIfPlaintextMatch {
		if user.Password != rawPassword {
			return nil
		}
	} else if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawPassword)); err == nil {
		return nil
	}

	hashedPassword, err := hashPassword(rawPassword)
	if err != nil {
		return err
	}

	if err := DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
		return fmt.Errorf("failed to update admin password: %w", err)
	}

	if onlyIfPlaintextMatch {
		log.Println("Legacy default admin password repaired")
	} else {
		log.Println("Admin password synced with environment variable")
	}

	return nil
}

func getEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

func migrateSiteContentSchema() error {
	if err := DB.Exec("DROP INDEX IF EXISTS idx_site_contents_key").Error; err != nil {
		return fmt.Errorf("failed to drop legacy site content index: %w", err)
	}

	if err := DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_site_content_key_language ON site_contents(key, language)").Error; err != nil {
		return fmt.Errorf("failed to create site content composite index: %w", err)
	}

	if err := DB.Model(&models.SiteContent{}).
		Where("type IS NULL OR type = ''").
		Update("type", "text").Error; err != nil {
		return fmt.Errorf("failed to normalize site content types: %w", err)
	}

	return nil
}

func seedMissingSiteContentEntries(entries []models.SiteContent) error {
	for _, entry := range entries {
		if entry.Type == "" {
			entry.Type = "text"
		}

		var existing models.SiteContent
		err := DB.Where("key = ? AND language = ?", entry.Key, entry.Language).First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := DB.Create(&entry).Error; err != nil {
				return fmt.Errorf("failed to seed site content %s: %w", entry.Key, err)
			}
			continue
		}

		if err != nil {
			return fmt.Errorf("failed to query site content %s: %w", entry.Key, err)
		}

		updates := map[string]interface{}{}
		if existing.Section != entry.Section {
			updates["section"] = entry.Section
		}
		if existing.Type == "" {
			updates["type"] = entry.Type
		}

		if len(updates) > 0 {
			if err := DB.Model(&existing).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update seeded site content %s: %w", entry.Key, err)
			}
		}
	}

	log.Println("Default site content ensured")
	return nil
}

func defaultSiteContentEntries() []models.SiteContent {
	return []models.SiteContent{
		newSiteContent("global", "site_logo_url", "/logo.jpeg", "image"),
		newSiteContent("global", "site_brand_name", "ANGEL EVENT", "text"),
		newSiteContent("global", "site_brand_tagline", "PLANIFICATION D'ÉVÉNEMENTS", "text"),
		newSiteContent("global", "site_contact_email", "contact@angelevent.com", "text"),
		newSiteContent("global", "site_contact_phone", "+1 (819) 244-4702", "text"),
		newSiteContent("global", "site_contact_location", "Trois-Rivieres, QC", "text"),
		newSiteContent("global", "site_instagram_url", "https://www.instagram.com/angel_eventt/", "text"),
		newSiteContent("global", "site_tiktok_url", "https://www.tiktok.com/@angel_eventt", "text"),
		newSiteContent("global", "footer_description", "Créer l'instant parfait avec élégance et passion.", "text"),

		newSiteContent("home", "home_hero_pretitle", "PLANIFICATION D'ÉVÉNEMENTS DE LUXE", "text"),
		newSiteContent("home", "home_hero_title", "Créer des Moments Inoubliables", "text"),
		newSiteContent("home", "home_hero_subtitle", "Nous transformons vos rêves en réalité avec élégance et précision.", "text"),
		newSiteContent("home", "home_hero_cta_label", "DÉCOUVRIR NOS SERVICES", "text"),
		newSiteContent("home", "home_hero_image", "https://images.unsplash.com/photo-1519741497674-611481863552?w=1920", "image"),
		newSiteContent("home", "home_welcome_title", "Bienvenue chez Angel Event", "text"),
		newSiteContent("home", "home_welcome_text", "Nous sommes spécialisés dans la création d'événements privés et corporatifs élégants, conçus sur mesure pour refléter votre style, vos objectifs et votre personnalité.", "text"),
		newSiteContent("home", "home_welcome_signature", "L'équipe Angel Event", "text"),
		newSiteContent("home", "home_service_1_title", "Mariages", "text"),
		newSiteContent("home", "home_service_1_description", "Planification complète pour votre grand jour.", "text"),
		newSiteContent("home", "home_service_2_title", "Corporatif", "text"),
		newSiteContent("home", "home_service_2_description", "Événements professionnels qui marquent les esprits.", "text"),
		newSiteContent("home", "home_service_3_title", "Célébrations", "text"),
		newSiteContent("home", "home_service_3_description", "Anniversaires, fiançailles et plus encore.", "text"),
		newSiteContent("home", "home_services_more", "EN SAVOIR PLUS", "text"),
		newSiteContent("home", "home_quote_text", "Chaque détail compte pour créer la perfection.", "text"),
		newSiteContent("home", "home_quote_author", "Fondatrice", "text"),
		newSiteContent("home", "home_gallery_image_1", "https://images.unsplash.com/photo-1511795409834-ef04bbd61622?w=600", "image"),
		newSiteContent("home", "home_gallery_image_2", "https://images.unsplash.com/photo-1519225421980-715cb0215aed?w=600", "image"),
		newSiteContent("home", "home_gallery_image_3", "https://images.unsplash.com/photo-1469334031218-e382a71b716b?w=600", "image"),
		newSiteContent("home", "home_gallery_image_4", "https://images.unsplash.com/photo-1519741497674-611481863552?w=600", "image"),

		newSiteContent("about", "about_hero_title", "Notre Histoire", "text"),
		newSiteContent("about", "about_hero_subtitle", "Passionnés par l'art de la célébration", "text"),
		newSiteContent("about", "about_story_title", "Qui sommes-nous?", "text"),
		newSiteContent("about", "about_story_p1", "Angel Event est né d'une passion pour l'art de recevoir et de célébrer les moments importants de la vie.", "text"),
		newSiteContent("about", "about_story_p2", "Nous sommes spécialisés dans la planification d'événements privés et corporatifs, en créant des expériences raffinées, mémorables et parfaitement alignées à votre vision.", "text"),
		newSiteContent("about", "about_story_p3", "Votre satisfaction est notre priorité absolue, et nous mettons tout en œuvre pour dépasser vos attentes avec créativité, rigueur et élégance.", "text"),
		newSiteContent("about", "about_story_image", "https://images.unsplash.com/photo-1511285560929-80b456fea0bc?w=800", "image"),
		newSiteContent("about", "about_values_title", "Nos Valeurs", "text"),
		newSiteContent("about", "about_values_excellence_title", "Excellence", "text"),
		newSiteContent("about", "about_values_excellence_desc", "Nous visons la perfection dans chaque détail.", "text"),
		newSiteContent("about", "about_values_passion_title", "Passion", "text"),
		newSiteContent("about", "about_values_passion_desc", "Nous aimons ce que nous faisons et cela se voit.", "text"),
		newSiteContent("about", "about_values_custom_title", "Sur Mesure", "text"),
		newSiteContent("about", "about_values_custom_desc", "Chaque événement est unique, tout comme vous.", "text"),
		newSiteContent("about", "about_values_trust_title", "Confiance", "text"),
		newSiteContent("about", "about_values_trust_desc", "Nous sommes votre partenaire de confiance.", "text"),
		newSiteContent("about", "about_why_title", "Pourquoi nous choisir?", "text"),
		newSiteContent("about", "about_why_1_title", "Expérience", "text"),
		newSiteContent("about", "about_why_1_desc", "Des années de savoir-faire à votre service.", "text"),
		newSiteContent("about", "about_why_2_title", "Clé en main", "text"),
		newSiteContent("about", "about_why_2_desc", "On s'occupe de tout, pour votre tranquillité d'esprit.", "text"),
		newSiteContent("about", "about_why_3_title", "Souci du détail", "text"),
		newSiteContent("about", "about_why_3_desc", "Rien n'est laissé au hasard.", "text"),
		newSiteContent("about", "about_why_4_title", "Réseau", "text"),
		newSiteContent("about", "about_why_4_desc", "Accès aux meilleurs prestataires de l'industrie.", "text"),
		newSiteContent("about", "about_cta_title", "Prêt à commencer?", "text"),
		newSiteContent("about", "about_cta_subtitle", "Contactez-nous dès aujourd'hui pour planifier votre événement.", "text"),
		newSiteContent("about", "about_cta_book", "RÉSERVER", "text"),
		newSiteContent("about", "about_cta_contact", "CONTACTEZ-NOUS", "text"),

		newSiteContent("services", "services_hero_title", "Nos Services", "text"),
		newSiteContent("services", "services_hero_subtitle", "Mariage, demande en mariage, baptême, anniversaire et baby shower pensés avec élégance", "text"),
		newSiteContent("services", "services_hero_description", "Mariage, demande en mariage, baptême, anniversaire et baby shower pensés avec élégance", "text"),
		newSiteContent("services", "services_list_cta", "RÉSERVER CE SERVICE", "text"),
		newSiteContent("services", "services_wedding_title", "Mariage", "text"),
		newSiteContent("services", "services_wedding_desc", "Un accompagnement mariage avec touche Event Planner pour orchestrer une journée raffinée, fluide et mémorable, de la cérémonie jusqu'à la réception.", "text"),
		newSiteContent("services", "services_wedding_feature_1", "Concept et ambiance sur mesure", "text"),
		newSiteContent("services", "services_wedding_feature_2", "Décoration de cérémonie et réception", "text"),
		newSiteContent("services", "services_wedding_feature_3", "Accompagnement Event Planner", "text"),
		newSiteContent("services", "services_wedding_feature_4", "Coordination des prestataires", "text"),
		newSiteContent("services", "services_wedding_feature_5", "Mise en place élégante", "text"),
		newSiteContent("services", "services_wedding_feature_6", "Présence et soutien le jour J", "text"),
		newSiteContent("services", "services_wedding_image", "https://images.unsplash.com/photo-1519741497674-611481863552?w=800", "image"),
		newSiteContent("services", "services_proposal_title", "Demande en Mariage", "text"),
		newSiteContent("services", "services_proposal_desc", "Une mise en scène romantique et personnalisée pour transformer votre demande en un moment inoubliable.", "text"),
		newSiteContent("services", "services_proposal_feature_1", "Repérage d'un lieu romantique", "text"),
		newSiteContent("services", "services_proposal_feature_2", "Décor signature et fleurs", "text"),
		newSiteContent("services", "services_proposal_feature_3", "Mise en lumière de l'instant", "text"),
		newSiteContent("services", "services_proposal_feature_4", "Photographe ou vidéaste", "text"),
		newSiteContent("services", "services_proposal_feature_5", "Installation discrète et coordination", "text"),
		newSiteContent("services", "services_proposal_feature_6", "Ambiance célébration après le grand oui", "text"),
		newSiteContent("services", "services_proposal_image", "https://images.unsplash.com/photo-1520854221256-17451cc331bf?w=800", "image"),
		newSiteContent("services", "services_baptism_title", "Baptêmes", "text"),
		newSiteContent("services", "services_baptism_desc", "Une célébration douce et harmonieuse pour marquer ce moment précieux entouré de votre famille et de vos proches.", "text"),
		newSiteContent("services", "services_baptism_feature_1", "Palette délicate et élégante", "text"),
		newSiteContent("services", "services_baptism_feature_2", "Décoration de table et coin douceur", "text"),
		newSiteContent("services", "services_baptism_feature_3", "Ballons, fleurs et signalétique", "text"),
		newSiteContent("services", "services_baptism_feature_4", "Mise en scène pour photos souvenirs", "text"),
		newSiteContent("services", "services_baptism_feature_5", "Ambiance familiale chaleureuse", "text"),
		newSiteContent("services", "services_baptism_feature_6", "Coordination du décor et du déroulement", "text"),
		newSiteContent("services", "services_baptism_image", "https://images.unsplash.com/photo-1513151233558-d860c5398176?w=800", "image"),
		newSiteContent("services", "services_birthday_title", "Anniversaire", "text"),
		newSiteContent("services", "services_birthday_desc", "Des anniversaires sur mesure, chic, festifs ou intimistes, conçus pour faire vivre une expérience à votre image.", "text"),
		newSiteContent("services", "services_birthday_feature_1", "Thématique personnalisée", "text"),
		newSiteContent("services", "services_birthday_feature_2", "Décor de salle ou de domicile", "text"),
		newSiteContent("services", "services_birthday_feature_3", "Table dessert et mise en scène", "text"),
		newSiteContent("services", "services_birthday_feature_4", "Coins photo et animations visuelles", "text"),
		newSiteContent("services", "services_birthday_feature_5", "Ambiance adulte ou enfant", "text"),
		newSiteContent("services", "services_birthday_feature_6", "Installation clé en main", "text"),
		newSiteContent("services", "services_birthday_image", "https://images.unsplash.com/photo-1530103862676-de8c9debad1d?w=800", "image"),
		newSiteContent("services", "services_baby_shower_title", "Baby Shower", "text"),
		newSiteContent("services", "services_baby_shower_desc", "Une baby shower tendre et esthétique pour célébrer cette belle attente dans une atmosphère conviviale et raffinée.", "text"),
		newSiteContent("services", "services_baby_shower_feature_1", "Univers doux et personnalisé", "text"),
		newSiteContent("services", "services_baby_shower_feature_2", "Arche de ballons et fleurs", "text"),
		newSiteContent("services", "services_baby_shower_feature_3", "Table dessert harmonisée", "text"),
		newSiteContent("services", "services_baby_shower_feature_4", "Coin cadeaux et souvenirs", "text"),
		newSiteContent("services", "services_baby_shower_feature_5", "Scénographie adaptée à votre espace", "text"),
		newSiteContent("services", "services_baby_shower_feature_6", "Installation et finition soignées", "text"),
		newSiteContent("services", "services_baby_shower_image", "https://images.unsplash.com/photo-1542042161784-26ab9e041e89?w=800", "image"),
		newSiteContent("services", "services_cta_title", "Un projet en tête?", "text"),
		newSiteContent("services", "services_cta_subtitle", "Parlons-en et rendons-le réel!", "text"),
		newSiteContent("services", "services_cta_quote", "OBTENIR UN DEVIS", "text"),
		newSiteContent("services", "services_cta_contact", "CONTACT", "text"),
	}
}

func newSiteContent(section, key, value, contentType string) models.SiteContent {
	return models.SiteContent{
		Key:      key,
		Value:    value,
		Language: "fr",
		Section:  section,
		Type:     contentType,
	}
}
