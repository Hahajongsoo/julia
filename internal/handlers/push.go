package handlers

import (
	"log"
	"net/http"
	"os"

	"julia/internal/models"
	"julia/internal/services"

	"github.com/gin-gonic/gin"
)

type PushHandler struct {
	pushService services.PushService
}

func NewPushHandler(pushService services.PushService) *PushHandler {
	return &PushHandler{pushService: pushService}
}

func (h *PushHandler) GetVapidPublic(c *gin.Context) {
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Cache-Control", "no-store")
	c.String(200, os.Getenv("VAPID_PUBLIC"))
}

func (h *PushHandler) CreateSubscription(c *gin.Context) {
	var createRequest models.SubscriptionRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		log.Printf("Invalid request body for subscription creation: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if createRequest.UserID == "" || createRequest.Endpoint == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id and endpoint are required",
		})
		return
	}

	subscription := &models.Subscription{
		UserID:   createRequest.UserID,
		Endpoint: createRequest.Endpoint,
		P256dh:   createRequest.Keys.P256dh,
		Auth:     createRequest.Keys.Auth,
	}

	err := h.pushService.CreateSubscription(subscription)
	if err != nil {
		log.Printf("Failed to create subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create subscription",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Subscription created successfully",
		"data":    subscription,
	})
}

func (h *PushHandler) DeleteSubscription(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "userID parameter is required",
		})
		return
	}

	var deleteRequest models.DeleteSubscriptionRequest
	if err := c.ShouldBindJSON(&deleteRequest); err != nil {
		log.Printf("Invalid request body for subscription deletion: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if deleteRequest.Endpoint == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "endpoint is required in request body",
		})
		return
	}

	log.Printf("Deleting subscription for user %s with endpoint: %s", userID, deleteRequest.Endpoint)

	err := h.pushService.DeleteSubscription(userID, deleteRequest.Endpoint)
	if err != nil {
		log.Printf("Failed to delete subscription for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete subscription",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscription deleted successfully",
	})
}

func (h *PushHandler) SendNotification(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "userID parameter is required",
		})
		return
	}

	err := h.pushService.SendNotification(userID)
	if err != nil {
		log.Printf("Failed to send notification to user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send notification",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification sent successfully",
	})
}

func (h *PushHandler) GetSubscriptions(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "userID parameter is required",
		})
		return
	}

	subscriptions, err := h.pushService.GetSubscription(userID)
	if err != nil {
		log.Printf("Failed to get subscriptions for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get subscriptions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": subscriptions,
	})
}
