package database

import (
	"fmt"
	"log"
	"os"

	"github.com/mazong/angel_event/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

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
		&models.RentalItem{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
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
		adminEmail := os.Getenv("ADMIN_EMAIL")
		if adminEmail == "" {
			adminEmail = "admin@angelevent.com"
		}

		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if adminPassword == "" {
			adminPassword = "ChangeThisPassword123!"
		}

		// Hash password (we'll implement this in auth handler)
		admin := models.User{
			Email:    adminEmail,
			Password: adminPassword, // Will be hashed by auth service
			Name:     "Administrator",
			Role:     "admin",
		}

		if err := DB.Create(&admin).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		log.Printf("Default admin user created: %s", adminEmail)
	}

	// Seed default site content
	var contentCount int64
	DB.Model(&models.SiteContent{}).Count(&contentCount)

	if contentCount == 0 {
		defaultContent := []models.SiteContent{
			{
				Key:      "hero_title",
				Value:    "Chaque événement mérite une mise en scène inoubliable",
				Language: "fr",
				Section:  "home",
			},
			{
				Key:      "hero_subtitle",
				Value:    "Créer l'instant parfait",
				Language: "fr",
				Section:  "home",
			},
			{
				Key:      "about_title",
				Value:    "À propos d'Angel Event",
				Language: "fr",
				Section:  "about",
			},
			{
				Key:      "about_description",
				Value:    "Nous transformons vos rêves en réalité avec passion, excellence et un sens du détail incomparable.",
				Language: "fr",
				Section:  "about",
			},
		}

		if err := DB.Create(&defaultContent).Error; err != nil {
			return fmt.Errorf("failed to seed default content: %w", err)
		}

		log.Println("Default site content seeded")
	}

	// Seed default rental items
	var rentalCount int64
	DB.Model(&models.RentalItem{}).Count(&rentalCount)

	if rentalCount == 0 {
		defaultRentals := []models.RentalItem{
			{
				Title:       "Arche Ronde Dorée",
				Description: "Une magnifique arche ronde dorée pour sublimer votre décor de cérémonie ou photobooth.",
				Price:       50.00,
				Category:    models.RentalCategoryBackdrop,
				ImageURL:    "/uploads/rentals/default-rental.png",
				Featured:    true,
				Available:   true,
			},
			{
				Title:       "Centre de Table Floral",
				Description: "Composition florale élégante dans des tons blanc et crème, idéale pour les mariages.",
				Price:       25.00,
				Category:    models.RentalCategoryCenterpiece,
				ImageURL:    "/uploads/rentals/default-rental.png",
				Featured:    true,
				Available:   true,
			},
			{
				Title:       "Mur de Fleurs Blanc",
				Description: "Mur de fleurs artificielles haute qualité, dimensions 2m x 2m. Impact visuel garanti.",
				Price:       120.00,
				Category:    models.RentalCategoryBackdrop,
				ImageURL:    "/uploads/rentals/default-rental.png",
				Featured:    false,
				Available:   true,
			},
			{
				Title:       "Vase Haut Cylindrique",
				Description: "Vase en verre transparent, hauteur 60cm. Parfait pour les compositions modernes.",
				Price:       15.00,
				Category:    models.RentalCategoryOther,
				ImageURL:    "/uploads/rentals/default-rental.png",
				Featured:    false,
				Available:   true,
			},
			{
				Title:       "Bouquet de Pivoines",
				Description: "Bouquet de pivoines artificielles réalistes.",
				Price:       10.00,
				Category:    models.RentalCategoryFlower,
				ImageURL:    "/uploads/rentals/default-rental.png",
				Featured:    false,
				Available:   true,
			},
		}

		if err := DB.Create(&defaultRentals).Error; err != nil {
			return fmt.Errorf("failed to seed default rentals: %w", err)
		}

		log.Println("Default rental items seeded")
	}

	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
