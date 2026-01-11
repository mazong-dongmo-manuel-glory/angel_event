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
		// Try root .env
		if err := godotenv.Load("../../../.env"); err != nil {
			log.Println("No .env file found")
		}
	}

	// Initialize database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	galleryPath := "backend/uploads/gallery"

	// Map folder names to GalleryCategory
	categoryMap := map[string]models.GalleryCategory{
		"marryme":     models.CategoryMarryMe,
		"weeding":     models.CategoryWedding,
		"wedding":     models.CategoryWedding,
		"birthday":    models.CategoryBirthday,
		"baby shower": models.CategoryBabyShower,
		"baby_shower": models.CategoryBabyShower,
		"bapteme":     models.CategoryBapteme,
		"congrats":    models.CategoryCongrats,
		"loveroom":    models.CategoryLoveroom,
	}

	err := filepath.Walk(galleryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Skip hidden files
			if strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			// Determine category from parent folder
			dir := filepath.Base(filepath.Dir(path))

			category, exists := categoryMap[dir]
			if !exists {
				// Fallback to birthday
				category = models.CategoryBirthday
			}

			// Construct public URL
			// Accessing via http://localhost:8080/uploads/gallery/...
			// Path relative to backend/
			relPath, _ := filepath.Rel("backend", path)
			publicURL := fmt.Sprintf("/%s", relPath)

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
				database.DB.Create(&image)
				fmt.Printf("Inserted: %s\n", info.Name())
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Gallery seeding completed!")
}
