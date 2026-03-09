package handlers

import (
	"log"
	"strconv"

	"server/src/services"
	"server/src/structs"

	"github.com/gin-gonic/gin"
)

type ConnectionHandlers struct {
	connectionService *services.ConnectionService
	profileService    *services.ProfileService
}

func NewConnectionHandlers(connectionService *services.ConnectionService, profileService *services.ProfileService) *ConnectionHandlers {
	return &ConnectionHandlers{
		connectionService: connectionService,
		profileService:    profileService,
	}
}

// SendConnectionRequest handles POST /connections/request/{userId}
func (h *ConnectionHandlers) SendConnectionRequest(c *gin.Context) {
	// Get authenticated user
	authUserID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get target user ID from URL parameter
	targetUserID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if requester can view this profile (must be in recommendations)
	canView, err := h.profileService.CanViewProfile(authUserID, targetUserID)
	if err != nil || !canView {
		c.JSON(403, gin.H{"error": "Cannot send connection request to this user"})
		return
	}

	// Send connection request
	request, err := h.connectionService.SendConnectionRequest(authUserID, targetUserID)
	if err != nil {
		log.Printf("Error sending connection request from %d to %d: %v", authUserID, targetUserID, err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if it was auto-accepted (mutual interest case)
	if request.Status == "auto_accepted" {
		c.JSON(200, gin.H{
			"message":      "Connection automatically established (mutual interest)",
			"status":       "connected",
			"connected_at": request.CreatedAt,
		})
		return
	}

	// Return the connection request
	c.JSON(201, gin.H{
		"message":    "Connection request sent successfully",
		"request_id": request.ID,
		"status":     request.Status,
		"created_at": request.CreatedAt,
	})
}

// GetReceivedConnectionRequests handles GET /connections/requests
func (h *ConnectionHandlers) GetReceivedConnectionRequests(c *gin.Context) {
	// Get authenticated user
	authUserID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get all pending requests received by this user
	requests, err := h.connectionService.GetReceivedRequests(authUserID)
	if err != nil {
		log.Printf("Error getting received requests for user %d: %v", authUserID, err)
		c.JSON(500, gin.H{"error": "Failed to retrieve connection requests"})
		return
	}

	// Transform to simpler response format with just essential info
	type RequestResponse struct {
		ID          int    `json:"id"`
		FromUserID  int    `json:"from_user_id"`
		DisplayName string `json:"display_name"`
		ProfilePic  string `json:"profile_picture_url"`
		CreatedAt   string `json:"created_at"`
	}

	responses := make([]RequestResponse, 0, len(requests))
	for _, req := range requests {
		responses = append(responses, RequestResponse{
			ID:          req.ID,
			FromUserID:  req.FromUserID,
			DisplayName: req.FromProfile.GetDisplayName(),
			ProfilePic:  req.FromProfile.ProfilePictureURL,
			CreatedAt:   req.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	c.JSON(200, gin.H{
		"requests": responses,
		"count":    len(responses),
	})
}

// GetSentConnectionRequests handles GET /connections/requests/sent
// This is an extra endpoint to help users track their pending requests
func (h *ConnectionHandlers) GetSentConnectionRequests(c *gin.Context) {
	// Get authenticated user
	authUserID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get all pending requests sent by this user
	requests, err := h.connectionService.GetSentRequests(authUserID)
	if err != nil {
		log.Printf("Error getting sent requests for user %d: %v", authUserID, err)
		c.JSON(500, gin.H{"error": "Failed to retrieve sent requests"})
		return
	}

	// Transform to simpler response
	type SentRequestResponse struct {
		ID          int    `json:"id"`
		ToUserID    int    `json:"to_user_id"`
		DisplayName string `json:"display_name"`
		ProfilePic  string `json:"profile_picture_url"`
		CreatedAt   string `json:"created_at"`
	}

	responses := make([]SentRequestResponse, 0, len(requests))
	for _, req := range requests {
		responses = append(responses, SentRequestResponse{
			ID:          req.ID,
			ToUserID:    req.ToUserID,
			DisplayName: req.ToProfile.GetDisplayName(),
			ProfilePic:  req.ToProfile.ProfilePictureURL,
			CreatedAt:   req.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	c.JSON(200, gin.H{
		"requests": responses,
		"count":    len(responses),
	})
}

// AcceptConnectionRequest handles POST /connections/requests/{requestId}/accept
func (h *ConnectionHandlers) AcceptConnectionRequest(c *gin.Context) {
	// Get authenticated user
	authUserID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get request ID from URL parameter
	requestID, err := strconv.Atoi(c.Param("requestId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request ID"})
		return
	}

	// Accept the request (service will verify user is the recipient)
	connection, err := h.connectionService.AcceptConnectionRequest(requestID, authUserID)
	if err != nil {
		log.Printf("Error accepting request %d by user %d: %v", requestID, authUserID, err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Determine the connected user ID
	connectedUserID := connection.User1ID
	if connection.User1ID == authUserID {
		connectedUserID = connection.User2ID
	}

	c.JSON(200, gin.H{
		"message":           "Connection request accepted",
		"connection_id":     connection.ID,
		"connected_user_id": connectedUserID,
		"connected_at":      connection.CreatedAt,
	})
}

// DismissConnectionRequest handles POST /connections/requests/{requestId}/dismiss
func (h *ConnectionHandlers) DismissConnectionRequest(c *gin.Context) {
	// Get authenticated user
	authUserID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get request ID from URL parameter
	requestID, err := strconv.Atoi(c.Param("requestId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request ID"})
		return
	}

	// Reject the request (service will verify user is the recipient)
	err = h.connectionService.RejectConnectionRequest(requestID, authUserID)
	if err != nil {
		log.Printf("Error dismissing request %d by user %d: %v", requestID, authUserID, err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Connection request dismissed",
	})
}

// GetConnections handles GET /connections
// Returns just the IDs as per spec
func (h *ConnectionHandlers) GetConnections(c *gin.Context) {
	// Get authenticated user
	authUserID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get user connections
	connections, err := h.connectionService.GetUserConnections(authUserID)
	if err != nil {
		log.Printf("Error getting connections for user %d: %v", authUserID, err)
		c.JSON(500, gin.H{"error": "Failed to get connections"})
		return
	}

	// Extract just the connection IDs (the other user's ID)
	connectionIDs := make([]int, 0, len(connections))
	for _, conn := range connections {
		// Determine which user is the "other" user
		otherUserID := conn.User2ID
		if conn.User2ID == authUserID {
			otherUserID = conn.User1ID
		}
		connectionIDs = append(connectionIDs, otherUserID)
	}

	// Return response per spec format
	response := structs.ConnectionsResponse{
		Connections: connectionIDs,
	}

	c.JSON(200, response)
}

// GetConnectionsDetailed handles GET /connections/detailed
// Extra endpoint that returns full connection info with profiles
func (h *ConnectionHandlers) GetConnectionsDetailed(c *gin.Context) {
	// Get authenticated user
	authUserID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get user connections with profiles
	connections, err := h.connectionService.GetUserConnections(authUserID)
	if err != nil {
		log.Printf("Error getting detailed connections for user %d: %v", authUserID, err)
		c.JSON(500, gin.H{"error": "Failed to get connections"})
		return
	}

	// Transform to detailed response
	type DetailedConnection struct {
		UserID            int    `json:"user_id"`
		DisplayName       string `json:"display_name"`
		ProfilePictureURL string `json:"profile_picture_url"`
		Location          string `json:"location"`
		ConnectedAt       string `json:"connected_at"`
	}

	responses := make([]DetailedConnection, 0, len(connections))
	for _, conn := range connections {
		responses = append(responses, DetailedConnection{
			UserID:            conn.User2Profile.UserID,
			DisplayName:       conn.User2Profile.GetDisplayName(),
			ProfilePictureURL: conn.User2Profile.ProfilePictureURL,
			Location:          conn.User2Profile.Location,
			ConnectedAt:       conn.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	c.JSON(200, gin.H{
		"connections": responses,
		"count":       len(responses),
	})
}

// DisconnectUser handles DELETE /connections/{userId}
func (h *ConnectionHandlers) DisconnectUser(c *gin.Context) {
	// Get authenticated user
	authUserID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get target user ID from URL parameter
	targetUserID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Cannot disconnect from yourself
	if authUserID == targetUserID {
		c.JSON(400, gin.H{"error": "Cannot disconnect from yourself"})
		return
	}

	// Remove the connection
	err = h.connectionService.RemoveConnection(authUserID, targetUserID)
	if err != nil {
		log.Printf("Error disconnecting user %d from user %d: %v", authUserID, targetUserID, err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":              "Successfully disconnected",
		"disconnected_user_id": targetUserID,
	})
}

// CheckConnectionStatus handles GET /connections/status/{userId}
// Extra endpoint to check the relationship status with another user
func (h *ConnectionHandlers) CheckConnectionStatus(c *gin.Context) {
	// Get authenticated user
	authUserID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication required"})
		return
	}

	// Get target user ID
	targetUserID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if connected
	connected, err := h.connectionService.AreUsersConnected(authUserID, targetUserID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to check connection status"})
		return
	}

	if connected {
		c.JSON(200, gin.H{
			"status":    "connected",
			"connected": true,
		})
		return
	}

	// Check for pending request
	pendingRequest, err := h.connectionService.GetPendingRequestBetweenUsers(authUserID, targetUserID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to check request status"})
		return
	}

	if pendingRequest != nil {
		if pendingRequest.FromUserID == authUserID {
			c.JSON(200, gin.H{
				"status":     "request_sent",
				"connected":  false,
				"request_id": pendingRequest.ID,
			})
		} else {
			c.JSON(200, gin.H{
				"status":     "request_received",
				"connected":  false,
				"request_id": pendingRequest.ID,
			})
		}
		return
	}

	// No relationship
	c.JSON(200, gin.H{
		"status":    "none",
		"connected": false,
	})
}

// Helper function to extract user ID from context
func getAuthenticatedUserID(c *gin.Context) (int, error) {
	userID, exists := GetCurrentUserID(c)
	if !exists {
		return 0, gin.Error{Err: gin.Error{}.Err, Type: gin.ErrorTypePrivate}
	}
	return userID, nil
}
