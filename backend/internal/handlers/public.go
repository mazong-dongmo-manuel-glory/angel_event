package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mazong/angel_event/internal/database"
	"github.com/mazong/angel_event/internal/models"
)

// ContactFormRequest represents a contact form submission
type ContactFormRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// SubmitContactForm handles contact form submissions
func SubmitContactForm(c *fiber.Ctx) error {
	var req ContactFormRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Send email to admin
	if err := emailService.SendContactFormEmail(req.Name, req.Email, req.Phone, req.Message); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send message",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Message sent successfully",
	})
}

// GetPublicAvailabilities returns availability calendar for public
func GetPublicAvailabilities(c *fiber.Ctx) error {
	var availabilities []models.Availability

	// Get month parameter
	month := c.Query("month", time.Now().Format("2006-01"))

	startDate := fmt.Sprintf("%s-01", month)
	endDate := fmt.Sprintf("%s-31", month)

	// Fetch availabilities
	if err := database.DB.Where("date >= ? AND date <= ?", startDate, endDate).
		Order("date ASC").
		Find(&availabilities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch availabilities",
		})
	}

	// For public view, we might want to filter or transform data if needed
	// But since Availability model is simple, we can return as is for now
	// If notes contained sensitive info, we would create a response struct

	return c.JSON(availabilities)
}

// GetRandomGalleryImages returns random gallery images for a category
func GetRandomGalleryImages(c *fiber.Ctx) error {
	category := c.Query("category")
	limit := c.QueryInt("limit", 3)
	excludeIDs := c.Query("exclude", "")

	var images []models.GalleryImage
	query := database.DB

	// Filter by category if provided
	if category != "" && category != "all" {
		query = query.Where("category = ?", category)
	}

	// Exclude specific IDs if provided
	if excludeIDs != "" {
		query = query.Where("id NOT IN (?)", excludeIDs)
	}

	// Get random images using ORDER BY RANDOM()
	// Note: For SQLite use RANDOM(), for PostgreSQL use RANDOM(), for MySQL use RAND()
	if err := query.Order("RANDOM()").Limit(limit).Find(&images).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch gallery images",
		})
	}

	return c.JSON(images)
}
