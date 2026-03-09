package handlers

import (
	"encoding/json"
	"net/http"
	"server/src/services"
	"server/src/websocket"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatHandlers struct {
	chatService *services.ChatService // Changed to unexported
	hub         *websocket.Hub        // Added field
}

func NewChatHandlers(chatService *services.ChatService, hub *websocket.Hub) *ChatHandlers { // Modified signature
	return &ChatHandlers{
		chatService: chatService, // Changed to unexported
		hub:         hub,         // Initialized new field
	}
}

func (h *ChatHandlers) SendMessage(c *gin.Context) {
	senderID, exists := GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		ReceiverID int    `json:"receiver_id"`
		Content    string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.chatService.SendMessage(senderID, req.ReceiverID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Broadcast via WebSocket
	wsMsg := map[string]interface{}{
		"type":    "message",
		"payload": msg,
	}

	jsonMsg, _ := json.Marshal(wsMsg)

	// Send to receiver
	h.hub.BroadcastToUser(req.ReceiverID, jsonMsg)

	// Send to sender (for other tabs/devices)
	h.hub.BroadcastToUser(senderID, jsonMsg)

	c.JSON(http.StatusOK, msg)
}

func (h *ChatHandlers) GetMessages(c *gin.Context) {
	senderID, exists := GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	receiverIDStr := c.Param("userId")
	receiverID, err := strconv.Atoi(receiverIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	messages, err := h.chatService.GetMessages(senderID, receiverID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *ChatHandlers) GetChats(c *gin.Context) {
	userID, exists := GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	chats, err := h.chatService.GetChats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}
