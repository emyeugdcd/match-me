package handlers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"server/src/services"
	"server/src/structs"

	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	profileService    *services.ProfileService
	connectionService *services.ConnectionService
}

func NewUserHandlers(profileService *services.ProfileService, connectionService *services.ConnectionService) *UserHandlers {
	return &UserHandlers{
		profileService:    profileService,
		connectionService: connectionService,
	}
}

// getAuthenticatedUserID extracts user ID from Gin context
func (h *UserHandlers) getAuthenticatedUserID(c *gin.Context) (int, error) {
	userID, exists := GetCurrentUserID(c)
	if !exists {
		return 0, fmt.Errorf("user not authenticated")
	}
	return userID, nil
}

// GetUserBasicInfo handles GET /users/{id}
func (h *UserHandlers) GetUserBasicInfo(c *gin.Context) {
	// Get authenticated user
	authUserID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get target user ID from URL
	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user can view this profile
	canView, err := h.profileService.CanViewProfile(authUserID, targetUserID)
	if err != nil || !canView {
		c.JSON(404, gin.H{"error": "User not found or no permission to view"})
		return
	}

	// Get user basic info
	userInfo, err := h.profileService.GetUserBasicInfo(targetUserID, false)
	if err != nil {
		log.Printf("Error getting user basic info: %v", err)
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Add profile link
	userInfo.ProfileLink = fmt.Sprintf("/profile/%d", targetUserID)

	c.JSON(200, userInfo)
}

// GetUserProfileInfo handles GET /users/{id}/profile
func (h *UserHandlers) GetUserProfileInfo(c *gin.Context) {
	// Get authenticated user
	authUserID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get target user ID from URL
	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user can view this profile
	canView, err := h.profileService.CanViewProfile(authUserID, targetUserID)
	if err != nil || !canView {
		c.JSON(404, gin.H{"error": "User not found or no permission to view"})
		return
	}

	// Get profile
	profile, err := h.profileService.GetProfile(targetUserID)
	if err != nil {
		log.Printf("Error getting profile: %v", err)
		c.JSON(404, gin.H{"error": "Profile not found"})
		return
	}

	// Convert to UserProfileInfo response
	profileInfo := structs.UserProfileInfo{
		ID:        targetUserID,
		AboutMe:   profile.AboutMe,
		Location:  profile.Location,
		Latitude:  profile.Latitude,
		Longitude: profile.Longitude,
	}

	c.JSON(200, profileInfo)
}

// GetUserBioInfo handles GET /users/{id}/bio
func (h *UserHandlers) GetUserBioInfo(c *gin.Context) {
	// Get authenticated user
	authUserID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get target user ID from URL
	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user can view this profile
	canView, err := h.profileService.CanViewProfile(authUserID, targetUserID)
	if err != nil || !canView {
		c.JSON(404, gin.H{"error": "User not found or no permission to view"})
		return
	}

	// Get bio data
	bioData, err := h.profileService.GetBioData(targetUserID)
	if err != nil {
		log.Printf("Error getting bio data: %v", err)
		c.JSON(404, gin.H{"error": "Bio data not found"})
		return
	}

	if bioData == nil {
		c.JSON(404, gin.H{"error": "Bio data not found"})
		return
	}

	// Convert to UserBioInfo response
	bioInfo := structs.UserBioInfo{
		ID:               targetUserID,
		Interests:        bioData.Interests,
		Hobbies:          bioData.Hobbies,
		MusicPreferences: bioData.MusicPreferences,
		FoodPreferences:  bioData.FoodPreferences,
		LookingFor:       bioData.LookingFor,
		AgeRangeMin:      bioData.AgeRangeMin,
		AgeRangeMax:      bioData.AgeRangeMax,
		MaxDistanceKm:    bioData.MaxDistanceKm,
	}

	c.JSON(200, bioInfo)
}

// GetCurrentUserBasicInfo handles GET /me
func (h *UserHandlers) GetCurrentUserBasicInfo(c *gin.Context) {
	// Get authenticated user
	authUserID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get user basic info
	userInfo, err := h.profileService.GetUserBasicInfo(authUserID, true)
	if err != nil {
		log.Printf("Error getting current user basic info: %v", err)
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Add profile link
	userInfo.ProfileLink = fmt.Sprintf("/profile/%d", authUserID)

	c.JSON(200, userInfo)
}

// GetCurrentUserProfileInfo handles GET /me/profile
func (h *UserHandlers) GetCurrentUserProfileInfo(c *gin.Context) {
	// Get authenticated user
	authUserID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get profile
	profile, err := h.profileService.GetProfile(authUserID)
	if err != nil {
		log.Printf("Error getting current user profile: %v", err)
		c.JSON(404, gin.H{"error": "Profile not found"})
		return
	}

	// Convert to UserProfileInfo response
	profileInfo := structs.UserProfileInfo{
		ID:        authUserID,
		AboutMe:   profile.AboutMe,
		Location:  profile.Location,
		Latitude:  profile.Latitude,
		Longitude: profile.Longitude,
	}

	c.JSON(200, profileInfo)
}

// GetCurrentUserBioInfo handles GET /me/bio
func (h *UserHandlers) GetCurrentUserBioInfo(c *gin.Context) {
	// Get authenticated user
	authUserID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get bio data
	bioData, err := h.profileService.GetBioData(authUserID)
	if err != nil {
		log.Printf("Error getting current user bio data: %v", err)
		c.JSON(404, gin.H{"error": "Bio data not found"})
		return
	}

	if bioData == nil {
		c.JSON(404, gin.H{"error": "Bio data not found"})
		return
	}

	// Convert to UserBioInfo response
	bioInfo := structs.UserBioInfo{
		ID:               authUserID,
		Interests:        bioData.Interests,
		Hobbies:          bioData.Hobbies,
		MusicPreferences: bioData.MusicPreferences,
		FoodPreferences:  bioData.FoodPreferences,
		LookingFor:       bioData.LookingFor,
		AgeRangeMin:      bioData.AgeRangeMin,
		AgeRangeMax:      bioData.AgeRangeMax,
		MaxDistanceKm:    bioData.MaxDistanceKm,
	}

	c.JSON(200, bioInfo)
}

// GetUserConnections handles GET /connections
func (h *UserHandlers) GetUserConnections(c *gin.Context) {
	// Get authenticated user
	authUserID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get user connections
	connections, err := h.connectionService.GetUserConnections(authUserID)
	if err != nil {
		log.Printf("Error getting user connections: %v", err)
		c.JSON(500, gin.H{"error": "Failed to get connections"})
		return
	}

	// Extract just the connection IDs
	connectionIDs := make([]int, 0, len(connections))
	for _, conn := range connections {
		// Determine which user is the other user in the connection
		otherUserID := conn.User2ID
		if conn.User2ID == authUserID {
			otherUserID = conn.User1ID
		}
		connectionIDs = append(connectionIDs, otherUserID)
	}

	response := structs.ConnectionsResponse{
		Connections: connectionIDs,
	}

	c.JSON(200, response)
}

// POST /me/profile-picture
func (h *UserHandlers) UploadProfilePicture(c *gin.Context) {
	authUserID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Parse multipart form
	file, err := c.FormFile("profile_picture")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate
	if err := h.profileService.ValidateProfilePicture(file); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Save file to disk
	filename := fmt.Sprintf("%d_%d.jpg", authUserID, time.Now().Unix())
	filepath := "./uploads/profiles/" + filename

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	// Update database
	pictureURL := "/uploads/profiles/" + filename
	if err := h.profileService.UpdateProfilePicture(authUserID, pictureURL); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(200, gin.H{"profile_picture_url": pictureURL})
}

// DELETE /me/profile-picture
func (h *UserHandlers) DeleteProfilePicture(c *gin.Context) {
	authUserID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get current picture to delete file
	profile, err := h.profileService.GetProfile(authUserID)
	if err == nil && profile.ProfilePictureURL != "" {
		// Delete file from disk
		os.Remove("." + profile.ProfilePictureURL)
	}

	// Clear from database
	if err := h.profileService.UpdateProfilePicture(authUserID, ""); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete picture"})
		return
	}

	c.JSON(200, gin.H{"message": "Profile picture deleted"})
}

// POST /me/profile - Create or update profile
func (h *UserHandlers) UpdateProfile(c *gin.Context) {
	authUserID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	var req structs.Profile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	profile, err := h.profileService.CreateOrUpdateProfile(authUserID, req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, profile)
}

// POST /me/bio - Create or update bio data
func (h *UserHandlers) UpdateBioData(c *gin.Context) {
	authUserID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	var req structs.BioData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	bioData, err := h.profileService.CreateBioData(authUserID, req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, bioData)
}
