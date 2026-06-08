package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"hy-board-backend/config"
	"hy-board-backend/database"
	"hy-board-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow CORS for websockets
	},
}

// WSClient represents a connected websocket client
type WSClient struct {
	Conn    *websocket.Conn
	UserID  uint
	IsAdmin bool
}

type Hub struct {
	clients    map[*WSClient]bool
	broadcast  chan []byte
	register   chan *WSClient
	unregister chan *WSClient
	mu         sync.Mutex
}

var TicketHub = Hub{
	clients:    make(map[*WSClient]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *WSClient),
	unregister: make(chan *WSClient),
}

func StartHub() {
	for {
		select {
		case client := <-TicketHub.register:
			TicketHub.mu.Lock()
			TicketHub.clients[client] = true
			TicketHub.mu.Unlock()
			log.Printf("WS client connected: UserID=%d, IsAdmin=%t", client.UserID, client.IsAdmin)
		case client := <-TicketHub.unregister:
			TicketHub.mu.Lock()
			if _, ok := TicketHub.clients[client]; ok {
				delete(TicketHub.clients, client)
				client.Conn.Close()
				log.Printf("WS client disconnected: UserID=%d", client.UserID)
			}
			TicketHub.mu.Unlock()
		case message := <-TicketHub.broadcast:
			TicketHub.mu.Lock()
			for client := range TicketHub.clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("WS write error: %v", err)
					client.Conn.Close()
					delete(TicketHub.clients, client)
				}
			}
			TicketHub.mu.Unlock()
		}
	}
}

// WS request payload format
type WSRequest struct {
	Action   string `json:"action"`
	Title    string `json:"title,omitempty"`
	Message  string `json:"message,omitempty"`
	TicketID uint   `json:"ticket_id,omitempty"`
}

func HandleTicketsWS(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		return
	}

	// Validate JWT token using configuration
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.Auth.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}
	userID := uint(userIDFloat)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	wsClient := &WSClient{
		Conn:    conn,
		UserID:  user.ID,
		IsAdmin: user.IsAdmin,
	}

	TicketHub.register <- wsClient

	// Read loop
	go func() {
		defer func() {
			TicketHub.unregister <- wsClient
		}()

		// On connection, send initial ticket list
		sendInitialTickets(wsClient)

		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				break
			}

			var req WSRequest
			if err := json.Unmarshal(msgBytes, &req); err != nil {
				sendError(wsClient, "Invalid request format")
				continue
			}

			handleWSAction(wsClient, req)
		}
	}()
}

func sendError(client *WSClient, msg string) {
	payload, _ := json.Marshal(map[string]interface{}{
		"event":   "error",
		"message": msg,
	})
	_ = client.Conn.WriteMessage(websocket.TextMessage, payload)
}

func sendInitialTickets(client *WSClient) {
	var list []models.Ticket
	if client.IsAdmin {
		// Admin sees all tickets with preloaded users and messages
		database.DB.Preload("User").Preload("Messages.User").Find(&list)
	} else {
		// Normal user sees only their own tickets
		database.DB.Where("user_id = ?", client.UserID).Preload("User").Preload("Messages.User").Find(&list)
	}

	payload, err := json.Marshal(map[string]interface{}{
		"event": "tickets_list",
		"data":  list,
	})
	if err == nil {
		_ = client.Conn.WriteMessage(websocket.TextMessage, payload)
	}
}

func handleWSAction(client *WSClient, req WSRequest) {
	switch req.Action {
	case "get_tickets":
		sendInitialTickets(client)

	case "create_ticket":
		if req.Title == "" || req.Message == "" {
			sendError(client, "Title and message are required")
			return
		}
		ticket := models.Ticket{
			UserID: client.UserID,
			Title:  req.Title,
			Status: "open",
		}
		if err := database.DB.Create(&ticket).Error; err != nil {
			sendError(client, "Failed to create ticket")
			return
		}

		msg := models.TicketMessage{
			TicketID: ticket.ID,
			UserID:   client.UserID,
			Message:  req.Message,
			IsAdmin:  client.IsAdmin,
		}
		database.DB.Create(&msg)

		// Reload full ticket info to broadcast
		var fullTicket models.Ticket
		database.DB.Preload("User").Preload("Messages.User").First(&fullTicket, ticket.ID)

		// Broadcast to all admins, and send to the owner client
		broadcastToRoleAndOwner("ticket_created", fullTicket, client.UserID)

	case "reply_ticket":
		if req.TicketID == 0 || req.Message == "" {
			sendError(client, "Ticket ID and message are required")
			return
		}

		var ticket models.Ticket
		if err := database.DB.First(&ticket, req.TicketID).Error; err != nil {
			sendError(client, "Ticket not found")
			return
		}

		// Security: Non-admin can only reply to their own tickets
		if !client.IsAdmin && ticket.UserID != client.UserID {
			sendError(client, "Unauthorized to reply to this ticket")
			return
		}

		msg := models.TicketMessage{
			TicketID: ticket.ID,
			UserID:   client.UserID,
			Message:  req.Message,
			IsAdmin:  client.IsAdmin,
		}
		if err := database.DB.Create(&msg).Error; err != nil {
			sendError(client, "Failed to send message")
			return
		}

		// Update ticket status
		if client.IsAdmin {
			ticket.Status = "replied"
		} else {
			ticket.Status = "open"
		}
		database.DB.Save(&ticket)

		// Reload ticket
		var fullTicket models.Ticket
		database.DB.Preload("User").Preload("Messages.User").First(&fullTicket, ticket.ID)

		broadcastToRoleAndOwner("ticket_updated", fullTicket, ticket.UserID)

	case "close_ticket":
		if req.TicketID == 0 {
			sendError(client, "Ticket ID is required")
			return
		}

		var ticket models.Ticket
		if err := database.DB.First(&ticket, req.TicketID).Error; err != nil {
			sendError(client, "Ticket not found")
			return
		}

		// Security: Non-admin can only close their own tickets
		if !client.IsAdmin && ticket.UserID != client.UserID {
			sendError(client, "Unauthorized")
			return
		}

		ticket.Status = "closed"
		database.DB.Save(&ticket)

		var fullTicket models.Ticket
		database.DB.Preload("User").Preload("Messages.User").First(&fullTicket, ticket.ID)

		broadcastToRoleAndOwner("ticket_updated", fullTicket, ticket.UserID)
	}
}

func broadcastToRoleAndOwner(event string, data models.Ticket, ownerID uint) {
	payload, err := json.Marshal(map[string]interface{}{
		"event": event,
		"data":  data,
	})
	if err != nil {
		return
	}

	TicketHub.mu.Lock()
	defer TicketHub.mu.Unlock()

	for client := range TicketHub.clients {
		// Send to all admins, or the user who owns the ticket
		if client.IsAdmin || client.UserID == ownerID {
			_ = client.Conn.WriteMessage(websocket.TextMessage, payload)
		}
	}
}
