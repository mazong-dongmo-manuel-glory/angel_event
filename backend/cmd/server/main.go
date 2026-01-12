package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/mazong/angel_event/internal/database"
	"github.com/mazong/angel_event/internal/handlers"
	"github.com/mazong/angel_event/internal/middleware"
	"github.com/mazong/angel_event/internal/models"
	"github.com/mazong/angel_event/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed default data
	if err := database.SeedDefaultData(); err != nil {
		log.Fatal("Failed to seed default data:", err)
	}

	// Scan storage for new images
	services.ScanStorage()

	// Hash default admin password if needed
	// detailed synchronization of admin password
	var user models.User
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPass := os.Getenv("ADMIN_PASSWORD")

	if adminEmail != "" && adminPass != "" {
		if err := database.DB.Where("email = ?", adminEmail).First(&user).Error; err == nil {
			// Update password to ensure it matches .env
			hashedPassword, _ := handlers.HashPassword(adminPass)
			user.Password = hashedPassword
			database.DB.Save(&user)
			log.Println("Admin password synced with environment variable")
		}
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Angel Event API",
		ServerHeader: "Angel Event",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// CORS configuration
	allowOrigins := os.Getenv("FRONTEND_URL")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:5173,http://127.0.0.1:5173"
	} else {
		allowOrigins += ",http://localhost:5173,http://127.0.0.1:5173"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "Angel Event API",
		})
	})

	// Public Uploads and Storage
	app.Static("/uploads", "./uploads")
	app.Static("/storage", "../storage")

	// API routes
	api := app.Group("/api")

	// Public routes
	public := api.Group("/public")
	public.Post("/contact", handlers.SubmitContactForm)
	public.Post("/bookings", handlers.CreateBooking)
	public.Post("/bookings", handlers.CreateBooking)
	public.Get("/bookings/availability", handlers.CheckAvailability)
	public.Get("/availabilities", handlers.GetPublicAvailabilities)
	public.Get("/testimonials", handlers.GetTestimonials)
	public.Post("/testimonials", handlers.CreateTestimonial)
	public.Post("/newsletter/subscribe", handlers.SubscribeNewsletter)
	public.Get("/gallery", handlers.GetGalleryImages)
	public.Get("/gallery/random", handlers.GetRandomGalleryImages)
	public.Get("/rentals", handlers.GetRentalItems)
	public.Get("/categories", handlers.GetCategories)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Get("/me", middleware.AuthMiddleware(), handlers.GetCurrentUser)
	auth.Post("/change-password", middleware.AuthMiddleware(), handlers.ChangePassword)

	// Admin routes (protected)
	admin := api.Group("/admin", middleware.AuthMiddleware(), middleware.AdminMiddleware())

	// Dashboard
	admin.Get("/dashboard/stats", handlers.GetDashboardStats)

	// Bookings
	admin.Get("/bookings", handlers.GetBookings)
	admin.Get("/bookings/:id", handlers.GetBooking)
	admin.Put("/bookings/:id/status", handlers.UpdateBookingStatus)
	admin.Get("/availabilities", handlers.GetAvailabilities)
	admin.Post("/availabilities", handlers.UpdateAvailability)

	// Clients
	admin.Get("/clients", handlers.GetClients)
	admin.Get("/clients/:id", handlers.GetClient)
	admin.Put("/clients/:id", handlers.UpdateClient)
	admin.Delete("/clients/:id", handlers.DeleteClient)
	admin.Post("/clients/:id/email", handlers.SendClientEmail)

	// Testimonials
	admin.Get("/testimonials", handlers.GetTestimonials)
	admin.Put("/testimonials/:id", handlers.UpdateTestimonial)
	admin.Delete("/testimonials/:id", handlers.DeleteTestimonial)

	// Newsletter
	admin.Get("/newsletter/subscribers", handlers.GetNewsletterSubscribers)
	admin.Put("/newsletter/subscribers/:id", handlers.UpdateNewsletterSubscriber)
	admin.Delete("/newsletter/subscribers/:id", handlers.DeleteNewsletterSubscriber)
	admin.Post("/newsletter/send", handlers.SendNewsletter)

	// Gallery
	admin.Post("/gallery", handlers.CreateGalleryImage)
	admin.Put("/gallery/:id", handlers.UpdateGalleryImage)
	admin.Delete("/gallery/:id", handlers.DeleteGalleryImage)
	admin.Post("/gallery/scan", handlers.ScanStorageFolder)
	admin.Get("/gallery/categories", handlers.GetGalleryCategories)

	// Site Content
	admin.Get("/content", handlers.GetSiteContent)

	// Rentals
	admin.Get("/rentals", handlers.GetRentalItems)
	admin.Post("/rentals", handlers.CreateRentalItem)
	admin.Put("/rentals/:id", handlers.UpdateRentalItem)
	admin.Delete("/rentals/:id", handlers.DeleteRentalItem)

	// Categories
	admin.Get("/categories", handlers.GetCategories)
	admin.Post("/categories", handlers.CreateCategory)
	admin.Put("/categories/:id", handlers.UpdateCategory)
	admin.Delete("/categories/:id", handlers.DeleteCategory)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
