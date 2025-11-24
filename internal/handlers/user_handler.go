package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/pokemon-chatbot-api/internal/services"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type RegisterRequest struct {
	TelegramID string `json:"telegram_id" binding:"required"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name"`
	Username   string `json:"username"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "telegram_id and first_name are required",
		})
		return
	}

	response, err := h.service.Register(
		req.TelegramID,
		req.FirstName,
		req.LastName,
		req.Username,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	telegramID := c.Param("telegramId")

	user, err := h.service.GetUserByTelegramID(telegramID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"exists": false,
			"user":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exists": true,
		"user":   user,
	})
}

func (h *UserHandler) CheckRegistration(c *gin.Context) {
	telegramID := c.Param("telegramId")

	registered, err := h.service.IsUserRegistered(telegramID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"registered": registered,
	})
}
