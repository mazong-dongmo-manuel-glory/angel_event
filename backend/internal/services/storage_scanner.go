package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mazong/angel_event/internal/database"
	"github.com/mazong/angel_event/internal/models"
)

// ScanStorage scans the storage directory and populates the database
func ScanStorage() {
	log.Println("Starting storage scan...")
	scanRentals()
	scanGallery()
	log.Println("Storage scan completed")
}

func scanRentals() {
	cwd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", cwd)

	// Try absolute path if we can guess it, or just print what we are trying
	storagePath := "../storage/location"
	fmt.Printf("Scanning rentals in relative path: %s\n", storagePath)

	absPath, _ := filepath.Abs(storagePath)
	fmt.Printf("Absolute path: %s\n", absPath)

	files, err := os.ReadDir(storagePath)
	if err != nil {
		fmt.Printf("Could not read rental storage directory: %v\n", err)
		return
	}
	fmt.Printf("Found %d files in rental storage\n", len(files))

	// Pre-fetch categories
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		log.Printf("Failed to fetch categories: %v", err)
	}
	log.Printf("Loaded %d categories for mapping", len(categories))
	catMap := make(map[string]uint)
	for _, c := range categories {
		catMap[c.Slug] = c.ID
	}

	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}

		// Check if exists
		var count int64
		// ImageURL format: /storage/location/filename
		imageURL := fmt.Sprintf("/storage/location/%s", file.Name())
		database.DB.Model(&models.RentalItem{}).Where("image_url = ?", imageURL).Count(&count)
		if count > 0 {
			continue
		}

		// Determine category based on filename keywords
		lowerName := strings.ToLower(file.Name())
		catID := catMap["other"] // default

		if strings.Contains(lowerName, "fleur") || strings.Contains(lowerName, "bouquet") {
			if id, ok := catMap["flower"]; ok {
				catID = id
			}
		} else if strings.Contains(lowerName, "table") {
			if id, ok := catMap["centerpiece"]; ok {
				catID = id
			}
		} else if strings.Contains(lowerName, "arche") || strings.Contains(lowerName, "mur") || strings.Contains(lowerName, "backdrop") {
			if id, ok := catMap["backdrop"]; ok {
				catID = id
			}
		} else if strings.Contains(lowerName, "animation") {
			if id, ok := catMap["animation"]; ok {
				catID = id
			}
		}

		// Create Title from filename
		title := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		title = strings.ReplaceAll(title, "_", " ")
		title = strings.ReplaceAll(title, "-", " ")
		// Capitalize first letter (simple)
		if len(title) > 0 {
			title = strings.ToUpper(title[:1]) + title[1:]
		}

		item := models.RentalItem{
			Title:        title,
			Description:  "Importé automatiquement depuis le stockage",
			Price:        0.00, // Default price
			CategoryID:   catID,
			CategoryEnum: models.RentalCategoryOther, // Satisfy legacy constraint
			ImageURL:     imageURL,
			Featured:     false,
			Available:    true,
		}

		if err := database.DB.Create(&item).Error; err != nil {
			log.Printf("Failed to import rental item %s: %v", file.Name(), err)
		} else {
			log.Printf("Imported rental item: %s", title)
		}
	}
}

func scanGallery() {
	storageRoot := "../storage"
	dirs, err := os.ReadDir(storageRoot)
	if err != nil {
		log.Printf("Could not read storage root: %v", err)
		return
	}

	for _, dir := range dirs {
		if !dir.IsDir() || strings.HasPrefix(dir.Name(), ".") || dir.Name() == "location" {
			continue
		}

		// Map directory name to category enum/slug
		// The directory names might not match GalleryCategory enum exactly, may need normalization
		// But for now, let's assume valid categories or map common ones.
		// Actually, GalleryImage uses Category field which is GalleryCategory (string).
		// Let's use the directory name as the category, ensuring it's lowercase/slugified.

		category := strings.ToLower(dir.Name())
		// Fix common mismatches if necessary, e.g. "weeding" -> "wedding"
		if category == "weeding" {
			category = "wedding"
		}
		if category == "baby shower" {
			category = "baby_shower"
		}

		scanGalleryDir(filepath.Join(storageRoot, dir.Name()), category)
	}
}

func scanGalleryDir(path string, category string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Could not read gallery dir %s: %v", path, err)
		return
	}

	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}

		// Relative path for URL: /storage/category/filename (roughly)
		// But we need to serve it. The static route is /storage -> ../storage
		// So URL is /storage/dirName/fileName
		dirName := filepath.Base(path)
		imageURL := fmt.Sprintf("/storage/%s/%s", dirName, file.Name())

		var count int64
		database.DB.Model(&models.GalleryImage{}).Where("image_url = ?", imageURL).Count(&count)
		if count > 0 {
			continue
		}

		// Create Title
		title := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		title = strings.ReplaceAll(title, "_", " ")
		title = strings.ReplaceAll(title, "-", " ")

		img := models.GalleryImage{
			Title:         title,
			Description:   "Importé depuis le stockage",
			ImageURL:      imageURL,
			Category:      models.GalleryCategory(category), // Cast to enum
			FileName:      file.Name(),
			IsFromStorage: true,
			Featured:      false,
			SortOrder:     0,
		}

		if err := database.DB.Create(&img).Error; err != nil {
			log.Printf("Failed to import gallery image %s: %v", file.Name(), err)
		} else {
			log.Printf("Imported gallery image: %s (%s)", title, category)
		}
	}
}
