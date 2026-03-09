package handlers

import (
	"fmt"
	"net/http"
	"server/src/services"
	"server/src/structs"

	"github.com/gin-gonic/gin"
)

type RecommendationHandlers struct {
	RecommendationService *services.RecommendationService
}

func NewRecommendationHandlers(service *services.RecommendationService) *RecommendationHandlers {
	return &RecommendationHandlers{RecommendationService: service}
}

func (h *RecommendationHandlers) GetRecommendations(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	recommendations, err := h.RecommendationService.GetRecommendations(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommendations"})
		return
	}

	// Extract IDs and scores
	var items []structs.RecommendationItem
	for _, rec := range recommendations {
		items = append(items, structs.RecommendationItem{
			UserID:     rec.UserID,
			MatchScore: rec.MatchScore,
		})
	}

	c.JSON(http.StatusOK, gin.H{"recommendations": items})
}

func (h *RecommendationHandlers) DismissRecommendation(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	dismissedUserIDStr := c.Param("userId")
	// Convert string to int (omitted error check for brevity, assuming valid int)
	// In real app, use strconv.Atoi
	var dismissedUserID int
	// simple scan
	fmt.Sscanf(dismissedUserIDStr, "%d", &dismissedUserID)

	if dismissedUserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err := h.RecommendationService.DismissRecommendation(userID.(int), dismissedUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dismiss recommendation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recommendation dismissed"})
}
