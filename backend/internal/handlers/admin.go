package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mazong/angel_event/internal/database"
	"github.com/mazong/angel_event/internal/models"
)

// GetClients returns all clients
func GetClients(c *fiber.Ctx) error {
	var clients []models.Client

	query := database.DB.Preload("Bookings")

	// Search by name or email
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Order("created_at DESC").Find(&clients).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch clients",
		})
	}

	return c.JSON(clients)
}

// GetClient returns a single client with bookings
func GetClient(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid client ID",
		})
	}

	var client models.Client
	if err := database.DB.Preload("Bookings").First(&client, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Client not found",
		})
	}

	return c.JSON(client)
}

// UpdateClient updates client information
func UpdateClient(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid client ID",
		})
	}

	var client models.Client
	if err := database.DB.First(&client, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Client not found",
		})
	}

	if err := c.BodyParser(&client); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := database.DB.Save(&client).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update client",
		})
	}

	return c.JSON(client)
}

// DeleteClient soft deletes a client
func DeleteClient(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid client ID",
		})
	}

	var client models.Client
	if err := database.DB.First(&client, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Client not found",
		})
	}

	if err := database.DB.Delete(&client).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete client",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Client deleted successfully",
	})
}

// SendClientEmail sends an email to a client
func SendClientEmail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid client ID",
		})
	}

	type EmailRequest struct {
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var req EmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var client models.Client
	if err := database.DB.First(&client, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Client not found",
		})
	}

	// Send email
	if err := emailService.SendEmail(client.Email, req.Subject, req.Message); err != nil {
		// Log failed email
		database.DB.Create(&models.EmailLog{
			To:       client.Email,
			Subject:  req.Subject,
			Type:     "custom",
			Status:   "failed",
			Error:    err.Error(),
			ClientID: &client.ID,
		})

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send email",
		})
	}

	// Log successful email
	database.DB.Create(&models.EmailLog{
		To:       client.Email,
		Subject:  req.Subject,
		Type:     "custom",
		Status:   "sent",
		ClientID: &client.ID,
	})

	return c.JSON(fiber.Map{
		"message": "Email sent successfully",
	})
}

// GetTestimonials returns testimonials
func GetTestimonials(c *fiber.Ctx) error {
	var testimonials []models.Testimonial

	query := database.DB.Preload("Client")

	// Public endpoint only shows approved testimonials
	if c.Locals("user_id") == nil {
		query = query.Where("approved = ?", true)
	}

	if err := query.Order("created_at DESC").Find(&testimonials).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch testimonials",
		})
	}

	return c.JSON(testimonials)
}

// CreateTestimonial creates a new testimonial
func CreateTestimonial(c *fiber.Ctx) error {
	var testimonial models.Testimonial
	if err := c.BodyParser(&testimonial); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Public submissions need approval
	if c.Locals("user_id") == nil {
		testimonial.Approved = false
	}

	if err := database.DB.Create(&testimonial).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create testimonial",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(testimonial)
}

// UpdateTestimonial updates a testimonial
func UpdateTestimonial(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid testimonial ID",
		})
	}

	var testimonial models.Testimonial
	if err := database.DB.First(&testimonial, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Testimonial not found",
		})
	}

	if err := c.BodyParser(&testimonial); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := database.DB.Save(&testimonial).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update testimonial",
		})
	}

	return c.JSON(testimonial)
}

// DeleteTestimonial deletes a testimonial
func DeleteTestimonial(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid testimonial ID",
		})
	}

	if err := database.DB.Delete(&models.Testimonial{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete testimonial",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Testimonial deleted successfully",
	})
}

// Newsletter handlers

// SubscribeNewsletter subscribes an email to newsletter
func SubscribeNewsletter(c *fiber.Ctx) error {
	type SubscribeRequest struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Language string `json:"language"`
	}

	var req SubscribeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Generate unsubscribe token
	token := make([]byte, 16)
	rand.Read(token)
	unsubToken := hex.EncodeToString(token)

	newsletter := models.Newsletter{
		Email:      req.Email,
		Name:       req.Name,
		Language:   req.Language,
		Active:     true,
		UnsubToken: unsubToken,
	}

	if err := database.DB.Create(&newsletter).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to subscribe",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Successfully subscribed to newsletter",
	})
}

// GetNewsletterSubscribers returns all subscribers (admin)
func GetNewsletterSubscribers(c *fiber.Ctx) error {
	var subscribers []models.Newsletter

	query := database.DB

	// Filter by active status
	if active := c.Query("active"); active != "" {
		query = query.Where("active = ?", active == "true")
	}

	if err := query.Order("created_at DESC").Find(&subscribers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch subscribers",
		})
	}

	return c.JSON(subscribers)
}

// SendNewsletter sends newsletter to all active subscribers
func SendNewsletter(c *fiber.Ctx) error {
	type NewsletterRequest struct {
		Subject string `json:"subject"`
		Content string `json:"content"`
	}

	var req NewsletterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var subscribers []models.Newsletter
	if err := database.DB.Where("active = ?", true).Find(&subscribers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch subscribers",
		})
	}

	// Send emails asynchronously
	successCount := 0
	failCount := 0

	for _, sub := range subscribers {
		err := emailService.SendNewsletterEmail(sub.Email, req.Subject, req.Content)
		if err != nil {
			failCount++
		} else {
			successCount++
		}
	}

	return c.JSON(fiber.Map{
		"message": "Newsletter sent",
		"total":   len(subscribers),
		"sent":    successCount,
		"failed":  failCount,
	})
}

// UpdateNewsletterSubscriber updates a newsletter subscriber (admin)
func UpdateNewsletterSubscriber(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subscriber ID",
		})
	}

	var subscriber models.Newsletter
	if err := database.DB.First(&subscriber, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Subscriber not found",
		})
	}

	type UpdateRequest struct {
		Active *bool  `json:"active"`
		Name   string `json:"name"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields if provided
	if req.Active != nil {
		subscriber.Active = *req.Active
	}
	if req.Name != "" {
		subscriber.Name = req.Name
	}

	if err := database.DB.Save(&subscriber).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update subscriber",
		})
	}

	return c.JSON(subscriber)
}

// DeleteNewsletterSubscriber deletes a newsletter subscriber (admin)
func DeleteNewsletterSubscriber(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subscriber ID",
		})
	}

	if err := database.DB.Delete(&models.Newsletter{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete subscriber",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Subscriber deleted successfully",
	})
}

// Gallery handlers

// GetGalleryImages returns gallery images
func GetGalleryImages(c *fiber.Ctx) error {
	var images []models.GalleryImage

	query := database.DB

	// Filter by category
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Filter by event type
	if eventType := c.Query("event_type"); eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	// Featured only
	if featured := c.Query("featured"); featured == "true" {
		query = query.Where("featured = ?", true)
	}

	if err := query.Order("sort_order ASC, created_at DESC").Find(&images).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch gallery images",
		})
	}

	return c.JSON(images)
}

// CreateGalleryImage creates a new gallery image
func CreateGalleryImage(c *fiber.Ctx) error {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form",
		})
	}

	// Get image file
	files := form.File["image"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image file is required",
		})
	}
	file := files[0]

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Allowed: jpg, jpeg, png, gif, webp",
		})
	}

	// Create uploads directory if not exists
	uploadDir := "./uploads/gallery"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Create database record
	title := c.FormValue("title")
	if title == "" {
		title = file.Filename
	}

	image := models.GalleryImage{
		Title:         title,
		Description:   c.FormValue("description"),
		Category:      models.GalleryCategory(c.FormValue("category")),
		Featured:      c.FormValue("featured") == "true",
		ImageURL:      "/uploads/gallery/" + filename,
		FileName:      filename,
		IsFromStorage: false,
	}

	if err := database.DB.Create(&image).Error; err != nil {
		// Cleanup file if db fails
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create gallery image record",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(image)
}

// UpdateGalleryImage updates a gallery image
func UpdateGalleryImage(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid image ID",
		})
	}

	var image models.GalleryImage
	if err := database.DB.First(&image, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Image not found",
		})
	}

	if err := c.BodyParser(&image); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := database.DB.Save(&image).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update image",
		})
	}

	return c.JSON(image)
}

// DeleteGalleryImage deletes a gallery image
func DeleteGalleryImage(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid image ID",
		})
	}

	if err := database.DB.Delete(&models.GalleryImage{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete image",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Image deleted successfully",
	})
}

// ScanStorageFolder scans the storage folder and imports images
func ScanStorageFolder(c *fiber.Ctx) error {
	storagePath := "../storage" // Relative to backend directory

	// Map folder names to GalleryCategory
	categoryMap := map[string]models.GalleryCategory{
		"marryme":     models.CategoryMarryMe,
		"weeding":     models.CategoryWedding,
		"wedding":     models.CategoryWedding,
		"birthday":    models.CategoryBirthday,
		"baby shower": models.CategoryBabyShower,
		"bapteme":     models.CategoryBapteme,
		"congrats":    models.CategoryCongrats,
		"loveroom":    models.CategoryLoveroom,
	}

	importedCount := 0
	skippedCount := 0

	err := filepath.Walk(storagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Skip hidden files
			if strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			// Only process image files
			ext := strings.ToLower(filepath.Ext(info.Name()))
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
				return nil
			}

			// Determine category from parent folder
			dir := filepath.Base(filepath.Dir(path))
			category, exists := categoryMap[dir]
			if !exists {
				category = models.CategoryBirthday // Default fallback
			}

			// Construct public URL - storage folder is served at /storage
			// URL-encode the path to handle spaces and special characters
			relPath, _ := filepath.Rel(storagePath, path)
			// Split path and encode each component
			pathParts := strings.Split(relPath, string(filepath.Separator))
			for i, part := range pathParts {
				pathParts[i] = strings.ReplaceAll(part, " ", "%20")
			}
			publicURL := "/storage/" + strings.Join(pathParts, "/")

			// Check if image already exists
			var count int64
			database.DB.Model(&models.GalleryImage{}).Where("image_url = ?", publicURL).Count(&count)
			if count == 0 {
				image := models.GalleryImage{
					Title:         info.Name(),
					Description:   fmt.Sprintf("Image de %s", dir),
					ImageURL:      publicURL,
					Category:      category,
					FileName:      info.Name(),
					IsFromStorage: true,
				}
				if err := database.DB.Create(&image).Error; err != nil {
					return err
				}
				importedCount++
			} else {
				skippedCount++
			}
		}
		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to scan storage folder: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Storage folder scanned successfully",
		"imported": importedCount,
		"skipped":  skippedCount,
	})
}

// GetGalleryCategories returns all categories with image counts
func GetGalleryCategories(c *fiber.Ctx) error {
	type CategoryStats struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}

	var stats []CategoryStats

	// Get count for each category
	categories := []models.GalleryCategory{
		models.CategoryWedding,
		models.CategoryMarryMe,
		models.CategoryBirthday,
		models.CategoryBabyShower,
		models.CategoryBapteme,
		models.CategoryLoveroom,
		models.CategoryCongrats,
	}

	for _, cat := range categories {
		var count int64
		database.DB.Model(&models.GalleryImage{}).Where("category = ?", cat).Count(&count)
		stats = append(stats, CategoryStats{
			Category: string(cat),
			Count:    count,
		})
	}

	return c.JSON(stats)
}

// Site Content handlers

// GetSiteContent returns site content
func GetSiteContent(c *fiber.Ctx) error {
	var content []models.SiteContent

	query := database.DB

	// Filter by section
	if section := c.Query("section"); section != "" {
		query = query.Where("section = ?", section)
	}

	// Filter by language
	language := c.Query("language", "fr")
	query = query.Where("language = ?", language)

	if err := query.Find(&content).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch content",
		})
	}

	return c.JSON(content)
}

// UpdateSiteContent updates site content
func UpdateSiteContent(c *fiber.Ctx) error {
	type ContentUpdate struct {
		Key      string `json:"key"`
		Value    string `json:"value"`
		Language string `json:"language"`
		Section  string `json:"section"`
	}

	var req ContentUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var content models.SiteContent
	result := database.DB.Where("key = ? AND language = ?", req.Key, req.Language).First(&content)

	if result.Error != nil {
		// Create new
		content = models.SiteContent{
			Key:      req.Key,
			Value:    req.Value,
			Language: req.Language,
			Section:  req.Section,
		}
		if err := database.DB.Create(&content).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create content",
			})
		}
	} else {
		// Update existing
		content.Value = req.Value
		content.Section = req.Section
		if err := database.DB.Save(&content).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update content",
			})
		}
	}

	return c.JSON(content)
}

// GetDashboardStats returns dashboard statistics
func GetDashboardStats(c *fiber.Ctx) error {
	var stats struct {
		TotalBookings     int64   `json:"total_bookings"`
		PendingBookings   int64   `json:"pending_bookings"`
		ConfirmedBookings int64   `json:"confirmed_bookings"`
		TotalClients      int64   `json:"total_clients"`
		TotalRevenue      float64 `json:"total_revenue"`
		MonthRevenue      float64 `json:"month_revenue"`
	}

	database.DB.Model(&models.Booking{}).Count(&stats.TotalBookings)
	database.DB.Model(&models.Booking{}).Where("status = ?", models.BookingStatusPending).Count(&stats.PendingBookings)
	database.DB.Model(&models.Booking{}).Where("status = ?", models.BookingStatusConfirmed).Count(&stats.ConfirmedBookings)
	database.DB.Model(&models.Client{}).Count(&stats.TotalClients)

	database.DB.Model(&models.Booking{}).
		Where("status IN ?", []models.BookingStatus{models.BookingStatusPaid, models.BookingStatusCompleted}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&stats.TotalRevenue)

	return c.JSON(stats)
}
