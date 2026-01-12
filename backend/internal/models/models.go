package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents an admin user
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Name      string         `gorm:"not null" json:"name"`
	Role      string         `gorm:"default:'admin'" json:"role"`
}

// Client represents a customer
type Client struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Phone     string         `json:"phone"`
	Notes     string         `gorm:"type:text" json:"notes"`
	Bookings  []Booking      `json:"bookings,omitempty"`
}

// EventType represents the type of event
type EventType string

const (
	EventTypeProposal   EventType = "proposal"
	EventTypeWedding    EventType = "wedding"
	EventTypeBirthday   EventType = "birthday"
	EventTypeBabyShower EventType = "baby_shower"
	EventTypeCorporate  EventType = "corporate"
	EventTypeOther      EventType = "other"
)

// BookingStatus represents the status of a booking
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusPaid      BookingStatus = "paid"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

// Booking represents a client reservation
type Booking struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	ClientID        uint           `gorm:"not null" json:"client_id"`
	Client          *Client        `json:"client,omitempty"`
	EventDate       time.Time      `gorm:"not null" json:"event_date"`
	EventType       EventType      `gorm:"not null" json:"event_type"`
	EventLocation   string         `json:"event_location"`
	GuestCount      int            `json:"guest_count"`
	Budget          float64        `json:"budget"`
	Status          BookingStatus  `gorm:"default:'pending'" json:"status"`
	Message         string         `gorm:"type:text" json:"message"`
	SpecialRequests string         `gorm:"type:text" json:"special_requests"`
	TotalAmount     float64        `json:"total_amount"`
	DepositAmount   float64        `json:"deposit_amount"`

	AdminNotes  string       `gorm:"type:text" json:"admin_notes"`
	RentalItems []RentalItem `gorm:"many2many:booking_rental_items;" json:"rental_items"`
}

// Availability represents available dates for bookings
type Availability struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Date      time.Time      `gorm:"uniqueIndex;not null" json:"date"`
	Available bool           `gorm:"default:true" json:"available"`
	MaxEvents int            `gorm:"default:1" json:"max_events"`
	Notes     string         `json:"notes"`
}

// Testimonial represents client feedback
type Testimonial struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	ClientID  *uint          `json:"client_id,omitempty"`
	Client    *Client        `json:"client,omitempty"`
	Name      string         `gorm:"not null" json:"name"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Rating    int            `gorm:"default:5" json:"rating"`
	EventType EventType      `json:"event_type"`
	Approved  bool           `gorm:"default:false" json:"approved"`
	Featured  bool           `gorm:"default:false" json:"featured"`
}

// Newsletter represents newsletter subscribers
type Newsletter struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Email      string         `gorm:"uniqueIndex;not null" json:"email"`
	Name       string         `json:"name"`
	Active     bool           `gorm:"default:true" json:"active"`
	Language   string         `gorm:"default:'fr'" json:"language"`
	UnsubToken string         `gorm:"uniqueIndex" json:"-"`
}

// GalleryCategory represents gallery image categories
type GalleryCategory string

const (
	CategoryWedding    GalleryCategory = "wedding"
	CategoryMarryMe    GalleryCategory = "marryme"
	CategoryBirthday   GalleryCategory = "birthday"
	CategoryBabyShower GalleryCategory = "baby_shower"
	CategoryBapteme    GalleryCategory = "bapteme"
	CategoryLoveroom   GalleryCategory = "loveroom"
	CategoryCongrats   GalleryCategory = "congrats"
)

// GalleryImage represents images in the gallery
type GalleryImage struct {
	ID            uint            `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
	Title         string          `json:"title"`
	Description   string          `gorm:"type:text" json:"description"`
	ImageURL      string          `gorm:"not null" json:"image_url"`
	ThumbnailURL  string          `json:"thumbnail_url"`
	Category      GalleryCategory `json:"category"` // baby_shower, bapteme, birthday, etc.
	FileName      string          `json:"file_name"`
	IsFromStorage bool            `gorm:"default:false" json:"is_from_storage"`
	Featured      bool            `gorm:"default:false" json:"featured"`
	SortOrder     int             `gorm:"default:0" json:"sort_order"`
}

// SiteContent represents editable site content
type SiteContent struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Key       string         `gorm:"uniqueIndex;not null" json:"key"`
	Value     string         `gorm:"type:text;not null" json:"value"`
	Language  string         `gorm:"default:'fr'" json:"language"`
	Section   string         `json:"section"` // home, services, about, etc.
}

// EmailLog represents sent emails for tracking
type EmailLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	To        string         `gorm:"not null" json:"to"`
	Subject   string         `gorm:"not null" json:"subject"`
	Type      string         `json:"type"` // booking_confirmation, newsletter, custom
	Status    string         `gorm:"default:'sent'" json:"status"`
	Error     string         `gorm:"type:text" json:"error,omitempty"`
	ClientID  *uint          `json:"client_id,omitempty"`
}

// Category represents a rental or gallery category
type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null" json:"name"`
	Slug        string         `gorm:"uniqueIndex;not null" json:"slug"`
	Description string         `json:"description"`
	Type        string         `gorm:"default:'rental'" json:"type"` // 'rental', 'gallery'
}

// RentalCategory represents rental item categories
type RentalCategory string

const (
	RentalCategoryCenterpiece RentalCategory = "centerpiece"
	RentalCategoryBackdrop    RentalCategory = "backdrop"
	RentalCategoryFlower      RentalCategory = "flower"
	RentalCategoryOther       RentalCategory = "other"
)

// RentalItem represents an item available for rent
type RentalItem struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	CategoryID  uint           `json:"category_id"`
	Category    Category       `json:"category"`
	// Deprecated: Use CategoryID instead
	CategoryEnum RentalCategory `gorm:"column:category" json:"category_enum"`
	ImageURL     string         `gorm:"not null" json:"image_url"`
	Featured     bool           `gorm:"default:false" json:"featured"`
	Available    bool           `gorm:"default:true" json:"available"`
}
