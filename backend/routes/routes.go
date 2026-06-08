package routes

import (
	"hy-board-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Simple CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Token, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api/v1")
	{
		// 1. UniProxy API Group (Used by XrayR Backend)
		server := api.Group("/server/UniProxy")
		server.Use(middleware.UniProxyAuth())
		{
			server.GET("/config", GetNodeConfig)
			server.GET("/user", GetNodeUsers)
			server.POST("/push", PushNodeData)
		}

		// 2. Client API Group (User Portal)
		client := api.Group("/client")
		{
			client.GET("/subscribe", Subscribe) // Subscription link doesn't use JWT, uses token query param
		}

		// 4. GraphQL API Endpoint
		api.POST("/graphql", middleware.GraphQLAuthContext(), GraphQLHandler)

		// 5. WebSocket Real-time Ticket Endpoint
		api.GET("/ws/tickets", HandleTicketsWS)

		// 6. Tickets REST Endpoint Group
		tickets := api.Group("/tickets")
		tickets.Use(middleware.UserAuth(false))
		{
			tickets.GET("", GetTicketsList)
			tickets.GET("/:id", GetTicketDetails)
			tickets.POST("", CreateTicket)
		}
	}

	return r
}
