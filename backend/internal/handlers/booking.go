package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mazong/angel_event/internal/database"
	"github.com/mazong/angel_event/internal/models"
	"github.com/mazong/angel_event/internal/services"
)

var emailService = services.NewEmailService()

// CreateBookingRequest represents a booking creation request
type CreateBookingRequest struct {
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Phone           string  `json:"phone"`
	EventDate       string  `json:"event_date"`
	EventType       string  `json:"event_type"`
	EventLocation   string  `json:"event_location"`
	GuestCount      int     `json:"guest_count"`
	Budget          float64 `json:"budget"`
	Message         string  `json:"message"`
	SpecialRequests string  `json:"special_requests"`
	Language        string  `json:"language"`
	RentalItemIDs   []uint  `json:"rental_item_ids"`
}

// CreateBooking creates a new booking
func CreateBooking(c *fiber.Ctx) error {
	var req CreateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Parse event date
	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format",
		})
	}

	// Check availability
	var availability models.Availability
	result := database.DB.Where("date = ?", eventDate.Format("2006-01-02")).First(&availability)

	if result.Error == nil && !availability.Available {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Cette date n'est pas disponible",
		})
	}

	// Find or create client
	var client models.Client
	result = database.DB.Where("email = ?", req.Email).First(&client)

	if result.Error != nil {
		// Create new client
		client = models.Client{
			Name:  req.Name,
			Email: req.Email,
			Phone: req.Phone,
		}
		if err := database.DB.Create(&client).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create client",
			})
		}
	}

	// Calculate deposit (30% of budget)
	depositAmount := req.Budget * 0.30

	// Create booking
	booking := models.Booking{
		ClientID:        client.ID,
		EventDate:       eventDate,
		EventType:       models.EventType(req.EventType),
		EventLocation:   req.EventLocation,
		GuestCount:      req.GuestCount,
		Budget:          req.Budget,
		Message:         req.Message,
		SpecialRequests: req.SpecialRequests,
		Status:          models.BookingStatusPending,
		TotalAmount:     req.Budget,
		DepositAmount:   depositAmount,
	}

	// Associate rental items if selected
	if len(req.RentalItemIDs) > 0 {
		var rentalItems []models.RentalItem
		if err := database.DB.Where("id IN ?", req.RentalItemIDs).Find(&rentalItems).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch rental items",
			})
		}
		booking.RentalItems = rentalItems
	}

	if err := database.DB.Create(&booking).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create booking",
		})
	}

	// Send confirmation email to client
	go emailService.SendBookingConfirmation(
		client.Email,
		client.Name,
		string(booking.EventType),
		eventDate.Format("2 January 2006"),
		req.Language,
	)

	// Send notification email to admin
	go emailService.SendAdminBookingNotification(
		client.Name,
		client.Email,
		client.Phone,
		string(booking.EventType),
		eventDate.Format("2 January 2006"),
		req.EventLocation,
		req.GuestCount,
		req.Budget,
		req.Message,
	)

	// Log email
	database.DB.Create(&models.EmailLog{
		To:       client.Email,
		Subject:  "Confirmation de réservation",
		Type:     "booking_confirmation",
		Status:   "sent",
		ClientID: &client.ID,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"booking": booking,
		"message": "Booking created successfully",
	})
}

// GetBookings returns all bookings (admin only)
func GetBookings(c *fiber.Ctx) error {
	var bookings []models.Booking

	query := database.DB.Preload("Client").Preload("RentalItems")

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("event_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("event_date <= ?", endDate)
	}

	if err := query.Order("event_date DESC").Find(&bookings).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch bookings",
		})
	}

	return c.JSON(bookings)
}

// GetBooking returns a single booking
func GetBooking(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID",
		})
	}

	var booking models.Booking
	if err := database.DB.Preload("Client").Preload("RentalItems").First(&booking, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Booking not found",
		})
	}

	return c.JSON(booking)
}

// UpdateBookingStatus updates a booking's status
func UpdateBookingStatus(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking ID",
		})
	}

	type UpdateStatusRequest struct {
		Status     string `json:"status"`
		AdminNotes string `json:"admin_notes"`
	}

	var req UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var booking models.Booking
	if err := database.DB.First(&booking, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Booking not found",
		})
	}

	booking.Status = models.BookingStatus(req.Status)
	booking.AdminNotes = req.AdminNotes

	if err := database.DB.Save(&booking).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update booking",
		})
	}

	return c.JSON(booking)
}

// CheckAvailability checks if a date is available
func CheckAvailability(c *fiber.Ctx) error {
	dateStr := c.Query("date")
	if dateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Date parameter is required",
		})
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format",
		})
	}

	var availability models.Availability
	result := database.DB.Where("date = ?", date.Format("2006-01-02")).First(&availability)

	if result.Error != nil {
		// No specific availability set, assume available
		return c.JSON(fiber.Map{
			"available": true,
			"date":      dateStr,
		})
	}

	// Check if max events reached
	var bookingCount int64
	database.DB.Model(&models.Booking{}).
		Where("event_date = ? AND status != ?", date.Format("2006-01-02"), models.BookingStatusCancelled).
		Count(&bookingCount)

	available := availability.Available && int(bookingCount) < availability.MaxEvents

	return c.JSON(fiber.Map{
		"available":  available,
		"date":       dateStr,
		"max_events": availability.MaxEvents,
		"booked":     bookingCount,
	})
}

// GetAvailabilities returns availability calendar (admin)
func GetAvailabilities(c *fiber.Ctx) error {
	var availabilities []models.Availability

	// Get month parameter
	month := c.Query("month", time.Now().Format("2006-01"))

	startDate := fmt.Sprintf("%s-01", month)
	endDate := fmt.Sprintf("%s-31", month)

	if err := database.DB.Where("date >= ? AND date <= ?", startDate, endDate).
		Order("date ASC").
		Find(&availabilities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch availabilities",
		})
	}

	return c.JSON(availabilities)
}

// UpdateAvailability updates or creates availability for a date
func UpdateAvailability(c *fiber.Ctx) error {
	type AvailabilityRequest struct {
		Date      string `json:"date"`
		Available bool   `json:"available"`
		MaxEvents int    `json:"max_events"`
		Notes     string `json:"notes"`
	}

	var req AvailabilityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var availability models.Availability
	result := database.DB.Where("date = ?", req.Date).First(&availability)

	if result.Error != nil {
		// Create new
		date, _ := time.Parse("2006-01-02", req.Date)
		availability = models.Availability{
			Date:      date,
			Available: req.Available,
			MaxEvents: req.MaxEvents,
			Notes:     req.Notes,
		}
		if err := database.DB.Create(&availability).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create availability",
			})
		}
	} else {
		// Update existing
		availability.Available = req.Available
		availability.MaxEvents = req.MaxEvents
		availability.Notes = req.Notes
		if err := database.DB.Save(&availability).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update availability",
			})
		}
	}

	return c.JSON(availability)
}
