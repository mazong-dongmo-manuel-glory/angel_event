package handlers

import (
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

// GetRentalItems returns all rental items, optionally filtered
func GetRentalItems(c *fiber.Ctx) error {
	var items []models.RentalItem
	query := database.DB

	// Filter by category
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Featured only
	if featured := c.Query("featured"); featured == "true" {
		query = query.Where("featured = ?", true)
	}

	if err := query.Order("created_at DESC").Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch rental items",
		})
	}

	return c.JSON(items)
}

// CreateRentalItem creates a new rental item with image upload
func CreateRentalItem(c *fiber.Ctx) error {
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
	uploadDir := "./uploads/rentals"
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

	// Parse price
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)

	// Create database record
	item := models.RentalItem{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Category:    models.RentalCategory(c.FormValue("category")),
		Price:       price,
		ImageURL:    "/uploads/rentals/" + filename,
		Featured:    c.FormValue("featured") == "true",
		Available:   true,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create rental item",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

// UpdateRentalItem updates a rental item
func UpdateRentalItem(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid item ID",
		})
	}

	var item models.RentalItem
	if err := database.DB.First(&item, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Item not found",
		})
	}

	// Helper function for updates
	updateData := make(map[string]interface{})
	if title := c.FormValue("title"); title != "" {
		updateData["title"] = title
	}
	if desc := c.FormValue("description"); desc != "" {
		updateData["description"] = desc
	}
	if cat := c.FormValue("category"); cat != "" {
		updateData["category"] = cat
	}
	if priceStr := c.FormValue("price"); priceStr != "" {
		price, _ := strconv.ParseFloat(priceStr, 64)
		updateData["price"] = price
	}
	if featured := c.FormValue("featured"); featured != "" {
		updateData["featured"] = featured == "true"
	}
	if available := c.FormValue("available"); available != "" {
		updateData["available"] = available == "true"
	}

	// Handle new image upload if present
	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["image"]
		if len(files) > 0 {
			file := files[0]
			ext := strings.ToLower(filepath.Ext(file.Filename))

			uploadDir := "./uploads/rentals"
			os.MkdirAll(uploadDir, 0755)

			filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
			filePath := filepath.Join(uploadDir, filename)

			if err := c.SaveFile(file, filePath); err == nil {
				updateData["image_url"] = "/uploads/rentals/" + filename
				// Don't delete old image just in case, or implement cleanup logic
			}
		}
	}

	if err := database.DB.Model(&item).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update item",
		})
	}

	return c.JSON(item)
}

// DeleteRentalItem deletes a rental item
func DeleteRentalItem(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid item ID",
		})
	}

	if err := database.DB.Delete(&models.RentalItem{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete item",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Item deleted successfully",
	})
}
