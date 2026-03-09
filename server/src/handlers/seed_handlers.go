package handlers

import (
	"net/http"

	"server/src/services"

	"github.com/gin-gonic/gin"
)

// SeedHandlers exposes endpoints for seeding fictitious users into the system.
type SeedHandlers struct {
	seedUsers *services.SeedUsers
}

// NewSeedHandlers constructs a new SeedHandlers instance.
func NewSeedHandlers(seedUsers *services.SeedUsers) *SeedHandlers {
	return &SeedHandlers{
		seedUsers: seedUsers,
	}
}

// SeedRequest represents the optional payload for the seeding endpoint.
type SeedRequest struct {
	Count int `json:"count"` // optional; default is 100 if missing or invalid
}

// SeedUsers handles POST /dev/seed-users
// It creates fictitious users until there are at least Count (or 100) users.
func (h *SeedHandlers) SeedUsers(c *gin.Context) {
	var req SeedRequest

	// Try to bind JSON; if body is empty or invalid, use the default.
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Count = 100
	}
	if req.Count <= 0 {
		req.Count = 100
	}

	created, err := h.seedUsers.SeedUsers(req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Fictitious users seeded successfully",
		"requested": req.Count,
		"created":   created,
	})
}
