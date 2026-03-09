package structs

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

// User represents the users table for authentication
type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Hidden from JSON
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Profile represents the profiles table with user information
type Profile struct {
	ID                int       `json:"id" db:"id"`
	UserID            int       `json:"user_id" db:"user_id"`
	FirstName         string    `json:"first_name,omitempty" db:"first_name"`
	LastName          string    `json:"last_name,omitempty" db:"last_name"`
	Username          string    `json:"username,omitempty" db:"username"`
	AboutMe           string    `json:"about_me,omitempty" db:"about_me"`
	ProfilePictureURL string    `json:"profile_picture_url,omitempty" db:"profile_picture_url"`
	Location          string    `json:"location" db:"location"`
	Latitude          *float64  `json:"latitude,omitempty" db:"latitude"`
	Longitude         *float64  `json:"longitude,omitempty" db:"longitude"`
	IsComplete        bool      `json:"is_complete" db:"is_complete"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// IsValid checks if profile has minimum required fields
func (p *Profile) IsValid() bool {
	hasName := (p.FirstName != "" && p.LastName != "") || p.Username != ""
	return hasName && p.Location != ""
}

// GetDisplayName returns the appropriate name to display
func (p *Profile) GetDisplayName() string {
	if p.Username != "" {
		return p.Username
	}
	return p.FirstName + " " + p.LastName
}

// BioData represents the bio_data table with biographical information for matching
type BioData struct {
	ID               int            `json:"id" db:"id"`
	UserID           int            `json:"user_id" db:"user_id"`
	Interests        pq.StringArray `json:"interests" db:"interests"`
	Hobbies          pq.StringArray `json:"hobbies" db:"hobbies"`
	MusicPreferences pq.StringArray `json:"music_preferences" db:"music_preferences"`
	FoodPreferences  pq.StringArray `json:"food_preferences" db:"food_preferences"`
	LookingFor       string         `json:"looking_for,omitempty" db:"looking_for"`
	AgeRangeMin      *int           `json:"age_range_min,omitempty" db:"age_range_min"`
	AgeRangeMax      *int           `json:"age_range_max,omitempty" db:"age_range_max"`
	MaxDistanceKm    *int           `json:"max_distance_km,omitempty" db:"max_distance_km"`
	AdditionalData   JSONB          `json:"additional_data,omitempty" db:"additional_data"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

// IsComplete checks if bio has at least 5 data points
func (b *BioData) IsComplete() bool {
	count := 0

	// Each field needs at least one non-empty value
	if len(b.Interests) > 0 && b.Interests[0] != "" {
		count++
	}
	if len(b.Hobbies) > 0 && b.Hobbies[0] != "" {
		count++
	}
	if len(b.MusicPreferences) > 0 && b.MusicPreferences[0] != "" {
		count++
	}
	if len(b.FoodPreferences) > 0 && b.FoodPreferences[0] != "" {
		count++
	}
	if b.LookingFor != "" && len(strings.TrimSpace(b.LookingFor)) > 0 {
		count++
	}

	// Could also count age ranges as one data point:
	if b.AgeRangeMin != nil && b.AgeRangeMax != nil {
		count++
	}

	return count >= 5
}

// JSONB is a custom type for PostgreSQL JSONB fields
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(src interface{}) error {
	if src == nil {
		*j = make(JSONB)
		return nil
	}

	var source []byte
	switch v := src.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	default:
		return fmt.Errorf("incompatible type for JSONB")
	}

	return json.Unmarshal(source, j)
}

// ConnectionRequest represents the connection_requests table
type ConnectionRequest struct {
	ID         int       `json:"id" db:"id"`
	FromUserID int       `json:"from_user_id" db:"from_user_id"`
	ToUserID   int       `json:"to_user_id" db:"to_user_id"`
	Status     string    `json:"status" db:"status"` // "pending", "accepted", "rejected"
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Connection represents the connections table for accepted connections
type Connection struct {
	ID        int       `json:"id" db:"id"`
	User1ID   int       `json:"user_1_id" db:"user_1_id"`
	User2ID   int       `json:"user_2_id" db:"user_2_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// DismissedRecommendation represents the dismissed_recommendations table
type DismissedRecommendation struct {
	ID              int       `json:"id" db:"id"`
	UserID          int       `json:"user_id" db:"user_id"`
	DismissedUserID int       `json:"dismissed_user_id" db:"dismissed_user_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// Chat represents the chats table (one chat per user pair)
type Chat struct {
	ID        int       `json:"id" db:"id"`
	User1ID   int       `json:"user_1_id" db:"user_1_id"`
	User2ID   int       `json:"user_2_id" db:"user_2_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Message represents the messages table
type Message struct {
	ID         int       `json:"id" db:"id"`
	ChatID     int       `json:"chat_id" db:"chat_id"`
	SenderID   int       `json:"sender_id" db:"sender_id"`
	ReceiverID int       `json:"receiver_id" db:"-"` // Computed field, not in messages table
	Content    string    `json:"content" db:"message_text"`
	IsRead     bool      `json:"is_read" db:"is_read"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type ChatPreview struct {
	UserID            int       `json:"user_id"`
	DisplayName       string    `json:"display_name"`
	ProfilePictureURL *string   `json:"profile_picture_url"`
	LastMessage       string    `json:"last_message"`
	LastMessageTime   time.Time `json:"last_message_time"`
	UnreadCount       int       `json:"unread_count"`
}

// UnreadMessage represents the unread_messages tracking table
type UnreadMessage struct {
	ID          int `json:"id" db:"id"`
	ChatID      int `json:"chat_id" db:"chat_id"`
	UserID      int `json:"user_id" db:"user_id"`
	UnreadCount int `json:"unread_count" db:"unread_count"`
}

// ViewedRecommendation represents the viewed_recommendations table
type ViewedRecommendation struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	ViewedUserID int       `json:"viewed_user_id" db:"viewed_user_id"`
	ViewedAt     time.Time `json:"viewed_at" db:"viewed_at"`
}

// ConnectionRequestWithUsers includes user information with connection request
type ConnectionRequestWithUsers struct {
	ConnectionRequest
	FromUser    *User    `json:"from_user,omitempty"`
	ToUser      *User    `json:"to_user,omitempty"`
	FromProfile *Profile `json:"from_profile,omitempty"`
	ToProfile   *Profile `json:"to_profile,omitempty"`
}

// ConnectionWithUsers includes user information with connection
type ConnectionWithUsers struct {
	Connection
	User1        *User    `json:"user_1,omitempty"`
	User2        *User    `json:"user_2,omitempty"`
	User1Profile *Profile `json:"user_1_profile,omitempty"`
	User2Profile *Profile `json:"user_2_profile,omitempty"`
	UnreadCount  int      `json:"unread_count"`
}

// ChatWithUsers includes user information with chat
type ChatWithUsers struct {
	Chat
	User1        *User    `json:"user_1,omitempty"`
	User2        *User    `json:"user_2,omitempty"`
	User1Profile *Profile `json:"user_1_profile,omitempty"`
	User2Profile *Profile `json:"user_2_profile,omitempty"`
	LastMessage  *Message `json:"last_message,omitempty"`
	UnreadCount  int      `json:"unread_count,omitempty"`
}

// MessageWithSender includes sender information with message
type MessageWithSender struct {
	Message
	Sender        *User    `json:"sender,omitempty"`
	SenderProfile *Profile `json:"sender_profile,omitempty"`
}

// UserWithProfile represents a combined user and profile data
type UserWithProfile struct {
	User    User     `json:"user"`
	Profile Profile  `json:"profile"`
	BioData *BioData `json:"bio_data,omitempty"`
}

// Recommendation represents a user recommendation with match score
type Recommendation struct {
	UserID      int     `json:"user_id"`
	Profile     Profile `json:"profile"`
	BioData     BioData `json:"bio_data"`
	MatchScore  float64 `json:"match_score"`
	IsViewed    bool    `json:"is_viewed"`
	IsDismissed bool    `json:"is_dismissed"`
}

// Constants for connection request statuses
const (
	ConnectionStatusPending  = "pending"
	ConnectionStatusAccepted = "accepted"
	ConnectionStatusRejected = "rejected"
)

// UserBasicInfo is for /users/{id} endpoint - just name and picture
type UserBasicInfo struct {
	ID                int    `json:"id"`
	DisplayName       string `json:"name"` // For frontend compatibility
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Username          string `json:"username"`
	Email             string `json:"email,omitempty"`
	ProfilePictureURL string `json:"profile_picture_url"`
	ProfileLink       string `json:"profile_link"`
}

// UserProfileInfo is for /users/{id}/profile endpoint - the "about me" section
type UserProfileInfo struct {
	ID        int      `json:"id"`
	AboutMe   string   `json:"about_me"`
	Location  string   `json:"location"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// UserBioInfo is for /users/{id}/bio endpoint - biographical data for matching
type UserBioInfo struct {
	ID               int            `json:"id"`
	Interests        pq.StringArray `json:"interests"`
	Hobbies          pq.StringArray `json:"hobbies"`
	MusicPreferences pq.StringArray `json:"music_preferences"`
	FoodPreferences  pq.StringArray `json:"food_preferences"`
	LookingFor       string         `json:"looking_for"`
	AgeRangeMin      *int           `json:"age_range_min,omitempty"`
	AgeRangeMax      *int           `json:"age_range_max,omitempty"`
	MaxDistanceKm    *int           `json:"max_distance_km,omitempty"`
}

// RecommendationResponse is for /recommendations endpoint
type RecommendationResponse struct {
	Recommendations []RecommendationItem `json:"recommendations"`
}

type RecommendationItem struct {
	UserID     int     `json:"user_id"`
	MatchScore float64 `json:"match_score"`
}

// ConnectionsResponse is for /connections endpoint - just IDs
type ConnectionsResponse struct {
	Connections []int `json:"connections"`
}
