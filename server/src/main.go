package main

import (
	"log"
	"server/src/db"
	"server/src/handlers"
	"server/src/services"
	"server/src/websocket"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	route := gin.Default()

	// Add CORS middleware for frontend
	route.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Serve uploaded files
	route.Static("/uploads", "./uploads")

	db.ConnectDatabase()
	database := db.GetDB()

	// Initialize services
	profileService := services.NewProfileService(database)
	connectionService := services.NewConnectionService(database)
	chatService := services.NewChatService(database)
	recommendationService := services.NewRecommendationService(database)
	seedUsers := services.NewSeedUsers(database, profileService)

	// Initialize WebSocket Hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Initialize handlers
	// Initialize handlers
	authHandlers := handlers.NewAuthHandlers(database)
	userHandlers := handlers.NewUserHandlers(profileService, connectionService) // Fixed: added connectionService
	connectionHandlers := handlers.NewConnectionHandlers(connectionService, profileService)
	chatHandlers := handlers.NewChatHandlers(chatService, wsHub)
	recommendationHandlers := handlers.NewRecommendationHandlers(recommendationService)
	seedHandlers := handlers.NewSeedHandlers(seedUsers) // Fixed: use seedUsers service

	// Public routes (no authentication required)
	route.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	route.POST("/signup", authHandlers.Signup)
	route.POST("/login", authHandlers.Login)
	route.POST("/logout", authHandlers.Logout)

	// Dev/testing: seed fictitious users
	route.POST("/dev/seed-users", seedHandlers.SeedUsers)

	// WebSocket route
	route.GET("/ws", handlers.AuthMiddleware(), func(c *gin.Context) {
		websocket.ServeWs(wsHub, c)
	})

	// Debug: list online user ids
	route.GET("/debug/online-users", func(c *gin.Context) {
		ids := wsHub.GetOnlineUserIDs()
		c.JSON(200, gin.H{"online_user_ids": ids})
	})

	// Public status endpoint to check if a user is online (for debugging)
	route.GET("/users/:id/status", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid user id"})
			return
		}

		ids := wsHub.GetOnlineUserIDs()
		found := false
		for _, uid := range ids {
			if uid == id {
				found = true
				break
			}
		}
		c.JSON(200, gin.H{"online": found})
	})

	// Protected routes group (require authentication)
	protected := route.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		// User endpoints
		protected.GET("/users/:id", userHandlers.GetUserBasicInfo)
		protected.GET("/users/:id/profile", userHandlers.GetUserProfileInfo)
		protected.GET("/users/:id/bio", userHandlers.GetUserBioInfo)

		// Current user endpoints (shortcuts)
		protected.GET("/me", userHandlers.GetCurrentUserBasicInfo)
		protected.GET("/me/profile", userHandlers.GetCurrentUserProfileInfo)
		protected.GET("/me/bio", userHandlers.GetCurrentUserBioInfo)
		protected.POST("/me/profile", userHandlers.UpdateProfile)
		protected.POST("/me/bio", userHandlers.UpdateBioData)
		protected.POST("/me/profile-picture", userHandlers.UploadProfilePicture)
		protected.DELETE("/me/profile-picture", userHandlers.DeleteProfilePicture)

		// Connection request routes
		protected.POST("/connections/request/:userId", connectionHandlers.SendConnectionRequest)
		protected.GET("/connections/requests", connectionHandlers.GetReceivedConnectionRequests)
		protected.GET("/connections/requests/sent", connectionHandlers.GetSentConnectionRequests)
		protected.POST("/connections/requests/:requestId/accept", connectionHandlers.AcceptConnectionRequest)
		protected.POST("/connections/requests/:requestId/dismiss", connectionHandlers.DismissConnectionRequest)

		// Connection management routes
		protected.GET("/connections", connectionHandlers.GetConnections)
		protected.GET("/connections/detailed", connectionHandlers.GetConnectionsDetailed)
		protected.DELETE("/connections/:userId", connectionHandlers.DisconnectUser)
		protected.GET("/connections/status/:userId", connectionHandlers.CheckConnectionStatus)

		// Chat routes
		protected.POST("/chats/messages", chatHandlers.SendMessage)
		protected.GET("/chats/:userId/messages", chatHandlers.GetMessages)
		protected.GET("/chats", chatHandlers.GetChats)

		// Recommendation routes
		protected.GET("/recommendations", recommendationHandlers.GetRecommendations)
		protected.POST("/recommendations/:userId/dismiss", recommendationHandlers.DismissRecommendation)
	}

	log.Println("Server starting on :8080")

	err := route.Run(":8080")
	if err != nil {
		panic(err)
	}
}
