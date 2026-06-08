package routes

import (
	"net/http"

	"hy-board-backend/database"
	"hy-board-backend/models"

	"github.com/gin-gonic/gin"
)

// GetTicketsList returns the list of tickets. Admins see all, users see only their own.
func GetTicketsList(c *gin.Context) {
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := u.(*models.User)

	var list []models.Ticket
	if user.IsAdmin {
		if err := database.DB.Preload("User").Order("updated_at desc, id desc").Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
			return
		}
	} else {
		if err := database.DB.Where("user_id = ?", user.ID).Preload("User").Order("updated_at desc, id desc").Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
			return
		}
	}

	c.JSON(http.StatusOK, list)
}

// GetTicketDetails returns a single ticket with all messages. Security checked.
func GetTicketDetails(c *gin.Context) {
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := u.(*models.User)

	id := c.Param("id")

	var ticket models.Ticket
	if err := database.DB.Preload("User").Preload("Messages.User").First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Security: Non-admin can only access their own tickets
	if !user.IsAdmin && ticket.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

type CreateTicketInput struct {
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// CreateTicket handles HTTP POST /api/v1/tickets to create a new ticket
func CreateTicket(c *gin.Context) {
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := u.(*models.User)

	var input CreateTicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	ticket := models.Ticket{
		UserID: user.ID,
		Title:  input.Title,
		Status: "open",
	}
	if err := database.DB.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	msg := models.TicketMessage{
		TicketID: ticket.ID,
		UserID:   user.ID,
		Message:  input.Message,
		IsAdmin:  user.IsAdmin,
	}
	if err := database.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save initial message"})
		return
	}

	// Fetch full details and broadcast the ticket_created event to connected admins
	var fullTicket models.Ticket
	database.DB.Preload("User").Preload("Messages.User").First(&fullTicket, ticket.ID)
	broadcastToRoleAndOwner("ticket_created", fullTicket, user.ID)

	c.JSON(http.StatusCreated, fullTicket)
}
