package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mazong/angel_event/internal/database"
	"github.com/mazong/angel_event/internal/models"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found")
	}

	// Set database path to the main database
	os.Setenv("DATABASE_PATH", "../../angel_event.db")

	// Initialize database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	storagePath := "../../../storage"

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
			relPath, _ := filepath.Rel(storagePath, path)
			publicURL := fmt.Sprintf("/storage/%s", relPath)

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
				fmt.Printf("Imported: %s (category: %s)\n", info.Name(), category)
			} else {
				skippedCount++
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal("Failed to scan storage folder:", err)
	}

	fmt.Printf("\n✅ Scan completed!\n")
	fmt.Printf("Imported: %d\n", importedCount)
	fmt.Printf("Skipped: %d\n", skippedCount)
}
