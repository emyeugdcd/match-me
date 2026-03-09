package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	structs "server/src/structs"
	"strings"
)

// ProfileService handles business logic for user profiles
type ProfileService struct {
	db *sql.DB
}

type ProfileUpdateRequest struct {
	FirstName         *string `json:"first_name"`
	LastName          *string `json:"last_name"`
	Username          *string `json:"username"`
	AboutMe           *string `json:"about_me"`
	ProfilePictureURL *string `json:"profile_picture_url"`
	Location          *string `json:"location"`
}

// NewProfileService creates a new profile service
func NewProfileService(db *sql.DB) *ProfileService {
	return &ProfileService{db: db}
}

// CreateOrUpdateProfile creates a new profile or updates existing one
func (s *ProfileService) CreateOrUpdateProfile(userID int, profile structs.Profile) (*structs.Profile, error) {
	// Validate user exists
	var userExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&userExists)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}
	if !userExists {
		return nil, errors.New("user does not exist")
	}

	// Validate profile data
	if !profile.IsValid() {
		return nil, errors.New("profile must have either first/last name or username, and location")
	}

	// Check if bio data exists and is complete
	bioData, err := s.GetBioData(userID)
	if err != nil {
		return nil, err
	}

	isComplete := profile.IsValid() && bioData != nil && bioData.IsComplete()

	// Upsert profile (INSERT or UPDATE if exists)
	query := `
        INSERT INTO profiles (
            user_id, first_name, last_name, username, about_me, 
            profile_picture_url, location, latitude, longitude, is_complete
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (user_id) 
        DO UPDATE SET
            first_name = EXCLUDED.first_name,
            last_name = EXCLUDED.last_name,
            username = EXCLUDED.username,
            about_me = EXCLUDED.about_me,
            profile_picture_url = EXCLUDED.profile_picture_url,
            location = EXCLUDED.location,
            latitude = EXCLUDED.latitude,
            longitude = EXCLUDED.longitude,
            is_complete = EXCLUDED.is_complete,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id, user_id, first_name, last_name, username, about_me,
                  profile_picture_url, location, latitude, longitude,
                  is_complete, created_at, updated_at
    `

	var newProfile structs.Profile
	err = s.db.QueryRow(
		query,
		userID,
		profile.FirstName,
		profile.LastName,
		profile.Username,
		profile.AboutMe,
		profile.ProfilePictureURL,
		profile.Location,
		profile.Latitude,
		profile.Longitude,
		isComplete,
	).Scan(
		&newProfile.ID,
		&newProfile.UserID,
		&newProfile.FirstName,
		&newProfile.LastName,
		&newProfile.Username,
		&newProfile.AboutMe,
		&newProfile.ProfilePictureURL,
		&newProfile.Location,
		&newProfile.Latitude,
		&newProfile.Longitude,
		&newProfile.IsComplete,
		&newProfile.CreatedAt,
		&newProfile.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error saving profile: %w", err)
	}

	log.Printf("Profile saved successfully for user %d", userID)
	return &newProfile, nil
}

// GetProfile retrieves a user's profile by user ID
func (s *ProfileService) GetProfile(userID int) (*structs.Profile, error) {
	query := `
		SELECT id, user_id, first_name, last_name, username, about_me, 
		       profile_picture_url, location, latitude, longitude,
		       is_complete, created_at, updated_at
		FROM profiles 
		WHERE user_id = $1
	`

	var profile structs.Profile
	err := s.db.QueryRow(query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Username,
		&profile.AboutMe,
		&profile.ProfilePictureURL,
		&profile.Location,
		&profile.Latitude,
		&profile.Longitude,
		&profile.IsComplete,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("profile not found")
		}
		return nil, fmt.Errorf("error getting profile: %w", err)
	}

	return &profile, nil
}

// DeleteProfile deletes a user's profile
func (s *ProfileService) DeleteProfile(userID int) error {
	result, err := s.db.Exec("DELETE FROM profiles WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("error deleting profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("profile not found")
	}

	log.Printf("Profile deleted successfully for user %d", userID)
	return nil
}

// CreateBioData creates or updates bio data for a user
func (s *ProfileService) CreateBioData(userID int, bioData structs.BioData) (*structs.BioData, error) {
	// Check if user exists
	var userExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&userExists)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}
	if !userExists {
		return nil, errors.New("user does not exist")
	}

	// Upsert bio data (INSERT or UPDATE if exists)
	query := `
		INSERT INTO bio_data (
			user_id, interests, hobbies, music_preferences, food_preferences,
			looking_for, age_range_min, age_range_max, max_distance_km, additional_data
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (user_id) 
		DO UPDATE SET
			interests = EXCLUDED.interests,
			hobbies = EXCLUDED.hobbies,
			music_preferences = EXCLUDED.music_preferences,
			food_preferences = EXCLUDED.food_preferences,
			looking_for = EXCLUDED.looking_for,
			age_range_min = EXCLUDED.age_range_min,
			age_range_max = EXCLUDED.age_range_max,
			max_distance_km = EXCLUDED.max_distance_km,
			additional_data = EXCLUDED.additional_data,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, user_id, interests, hobbies, music_preferences,
				  food_preferences, looking_for, age_range_min, age_range_max,
				  max_distance_km, additional_data, created_at, updated_at
	`

	var newBioData structs.BioData
	err = s.db.QueryRow(
		query,
		userID,
		bioData.Interests,
		bioData.Hobbies,
		bioData.MusicPreferences,
		bioData.FoodPreferences,
		bioData.LookingFor,
		bioData.AgeRangeMin,
		bioData.AgeRangeMax,
		bioData.MaxDistanceKm,
		bioData.AdditionalData,
	).Scan(
		&newBioData.ID,
		&newBioData.UserID,
		&newBioData.Interests,
		&newBioData.Hobbies,
		&newBioData.MusicPreferences,
		&newBioData.FoodPreferences,
		&newBioData.LookingFor,
		&newBioData.AgeRangeMin,
		&newBioData.AgeRangeMax,
		&newBioData.MaxDistanceKm,
		&newBioData.AdditionalData,
		&newBioData.CreatedAt,
		&newBioData.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error saving bio data: %w", err)
	}

	// Update profile completeness after saving bio data
	// Profile is complete only if it's valid and bio data has 5+ data points
	_, err = s.db.Exec(`
		UPDATE profiles 
		SET is_complete = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE user_id = $2
	`, newBioData.IsComplete(), userID)

	if err != nil {
		log.Printf("Warning: Failed to update profile completeness for user %d: %v", userID, err)
	}

	log.Printf("Bio data saved successfully for user %d (complete: %v)", userID, newBioData.IsComplete())
	return &newBioData, nil
}

// GetBioData retrieves bio data for a user
func (s *ProfileService) GetBioData(userID int) (*structs.BioData, error) {
	query := `
		SELECT id, user_id, interests, hobbies, music_preferences, food_preferences,
		       looking_for, age_range_min, age_range_max, max_distance_km,
		       additional_data, created_at, updated_at
		FROM bio_data 
		WHERE user_id = $1
	`

	var bioData structs.BioData
	err := s.db.QueryRow(query, userID).Scan(
		&bioData.ID,
		&bioData.UserID,
		&bioData.Interests,
		&bioData.Hobbies,
		&bioData.MusicPreferences,
		&bioData.FoodPreferences,
		&bioData.LookingFor,
		&bioData.AgeRangeMin,
		&bioData.AgeRangeMax,
		&bioData.MaxDistanceKm,
		&bioData.AdditionalData,
		&bioData.CreatedAt,
		&bioData.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil instead of error if not found
		}
		return nil, fmt.Errorf("error getting bio data: %w", err)
	}

	return &bioData, nil
}

// GetUserBasicInfo gets some basic user info for display
func (s *ProfileService) GetUserBasicInfo(userID int, includeEmail bool) (*structs.UserBasicInfo, error) {
	profile, err := s.GetProfile(userID)
	if err != nil {
		return nil, err
	}

	var email string
	if includeEmail {
		// Get email from users table
		err = s.db.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&email)
		if err != nil {
			// Log error but don't fail - email is optional/extra info
			log.Printf("Error getting user email: %v", err)
		}
	}

	return &structs.UserBasicInfo{
		ID:                userID,
		DisplayName:       profile.GetDisplayName(),
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		Username:          profile.Username,
		Email:             email,
		ProfilePictureURL: profile.ProfilePictureURL,
	}, nil
}

// IsProfileComplete checks if a user's profile is complete enough for matching
func (s *ProfileService) IsProfileComplete(userID int) (bool, error) {
	profile, err := s.GetProfile(userID)
	if err != nil {
		return false, err
	}

	bioData, err := s.GetBioData(userID)
	if err != nil {
		return false, err
	}

	// Profile must be complete AND bio data must exist and be complete
	return profile.IsComplete && bioData != nil && bioData.IsComplete(), nil
}

// CanViewProfile checks if requesterID can view targetUserID's profile
func (s *ProfileService) CanViewProfile(requesterID, targetUserID int) (bool, error) {
	// 1. If viewing own profile
	if requesterID == targetUserID {
		return true, nil
	}

	// 2. Check if connected
	var connected bool
	query := `
        SELECT EXISTS(
            SELECT 1 FROM connections 
            WHERE (user_1_id = $1 AND user_2_id = $2) 
               OR (user_1_id = $2 AND user_2_id = $1)
        )
    `
	err := s.db.QueryRow(query, requesterID, targetUserID).Scan(&connected)
	if err != nil {
		return false, err
	}
	if connected {
		return true, nil
	}

	// 3. Check if there's a pending connection request (either direction)
	query = `
        SELECT EXISTS(
            SELECT 1 FROM connection_requests 
            WHERE status = 'pending' 
            AND ((from_user_id = $1 AND to_user_id = $2) 
                 OR (from_user_id = $2 AND to_user_id = $1))
        )
    `
	var hasPendingRequest bool
	err = s.db.QueryRow(query, requesterID, targetUserID).Scan(&hasPendingRequest)
	if err != nil {
		return false, err
	}
	if hasPendingRequest {
		return true, nil
	}

	// 4. Check if target is in requester's viewed recommendations
	query = `
        SELECT EXISTS(
            SELECT 1 FROM viewed_recommendations 
            WHERE user_id = $1 AND viewed_user_id = $2
        )
    `
	var isRecommended bool
	err = s.db.QueryRow(query, requesterID, targetUserID).Scan(&isRecommended)
	if err != nil {
		return false, err
	}

	return isRecommended, nil
}

// UpdateProfileCompleteness recalculates and updates the is_complete flag
func (s *ProfileService) UpdateProfileCompleteness(userID int) error {
	profile, err := s.GetProfile(userID)
	if err != nil {
		return err
	}

	bioData, err := s.GetBioData(userID)
	if err != nil {
		return err
	}

	isComplete := profile.IsValid() && bioData != nil && bioData.IsComplete()

	_, err = s.db.Exec(`
        UPDATE profiles 
        SET is_complete = $1, updated_at = CURRENT_TIMESTAMP 
        WHERE user_id = $2
    `, isComplete, userID)

	return err
}

// ValidateProfilePicture checks if a file is a valid image
func (s *ProfileService) ValidateProfilePicture(fileHeader *multipart.FileHeader) error {
	// Check file size (max 5MB)
	const maxFileSize = 5 * 1024 * 1024
	if fileHeader.Size > maxFileSize {
		return errors.New("file size exceeds 5MB limit")
	}

	// Check MIME type
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		return errors.New("file is not an image")
	}

	return nil
}

// UpdateProfilePicture updates just the profile picture URL
func (s *ProfileService) UpdateProfilePicture(userID int, pictureURL string) error {
	_, err := s.db.Exec(`
        UPDATE profiles 
        SET profile_picture_url = $1, updated_at = CURRENT_TIMESTAMP 
        WHERE user_id = $2
    `, pictureURL, userID)

	return err
}

// // GetUserWithProfile gets complete user data with profile
// func (s *ProfileService) GetUserWithProfile(userID int) (*structs.UserWithProfile, error) {
// 	// Get user
// 	var user structs.User
// 	err := s.db.QueryRow(
// 		"SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = $1",
// 		userID,
// 	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, errors.New("user not found")
// 		}
// 		return nil, fmt.Errorf("error getting user: %w", err)
// 	}

// 	// Get profile
// 	profile, err := s.GetProfile(userID)
// 	if err != nil {
// 		// User exists but profile doesn't - return user with empty profile
// 		emptyProfile := structs.Profile{UserID: userID}
// 		return &structs.UserWithProfile{
// 			User:    user,
// 			Profile: emptyProfile,
// 		}, nil
// 	}

// 	// Get bio data (optional)
// 	bioData, _ := s.GetBioData(userID)

// 	return &structs.UserWithProfile{
// 		User:    user,
// 		Profile: *profile,
// 		BioData: bioData,
// 	}, nil
// }

// // SearchProfiles searches profiles by location and interests
// func (s *ProfileService) SearchProfiles(location string, interests []string, limit int) ([]structs.Profile, error) {
// 	query := `
// 		SELECT p.id, p.user_id, p.first_name, p.last_name, p.username, p.about_me,
// 		       p.profile_picture_url, p.location, p.is_complete, p.created_at, p.updated_at
// 		FROM profiles p
// 		LEFT JOIN bio_data b ON p.user_id = b.user_id
// 		WHERE p.location ILIKE $1 AND p.is_complete = true
// 	`

// 	args := []interface{}{"%" + location + "%"}
// 	argCount := 1

// 	// Add interests filter if provided
// 	if len(interests) > 0 {
// 		query += fmt.Sprintf(" AND b.interests && $%d", argCount+1)
// 		args = append(args, interests)
// 		argCount++
// 	}

// 	query += " ORDER BY p.updated_at DESC LIMIT $2"
// 	args = append(args, limit)

// 	rows, err := s.db.Query(query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("error searching profiles: %w", err)
// 	}
// 	defer rows.Close()

// 	var profiles []structs.Profile
// 	for rows.Next() {
// 		var profile structs.Profile
// 		err := rows.Scan(
// 			&profile.ID,
// 			&profile.UserID,
// 			&profile.FirstName,
// 			&profile.LastName,
// 			&profile.Username,
// 			&profile.AboutMe,
// 			&profile.ProfilePictureURL,
// 			&profile.Location,
// 			&profile.IsComplete,
// 			&profile.CreatedAt,
// 			&profile.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("error scanning profile: %w", err)
// 		}
// 		profiles = append(profiles, profile)
// 	}

// 	return profiles, nil
// }
