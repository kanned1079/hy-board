package routes

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"hy-board-backend/database"
	"hy-board-backend/middleware"
	"hy-board-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":              &graphql.Field{Type: graphql.Int},
		"email":           &graphql.Field{Type: graphql.String},
		"v2ray_uuid":      &graphql.Field{Type: graphql.String},
		"trojan_password": &graphql.Field{Type: graphql.String},
		"speed_limit":     &graphql.Field{Type: graphql.Int},
		"device_limit":    &graphql.Field{Type: graphql.Int},
		"total_traffic":   &graphql.Field{Type: graphql.Float}, // using float64 to prevent int32 overflow in JS for large numbers
		"used_traffic":    &graphql.Field{Type: graphql.Float},
		"expired_at":      &graphql.Field{Type: graphql.String},
		"status":          &graphql.Field{Type: graphql.Int},
		"is_admin":        &graphql.Field{Type: graphql.Boolean},
		"balance":         &graphql.Field{Type: graphql.Float},
	},
})

var nodeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Node",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.Int},
		"name":         &graphql.Field{Type: graphql.String},
		"type":         &graphql.Field{Type: graphql.String},
		"address":      &graphql.Field{Type: graphql.String},
		"port":         &graphql.Field{Type: graphql.Int},
		"traffic_rate": &graphql.Field{Type: graphql.Float},
		"settings":     &graphql.Field{Type: graphql.String},
		"show":         &graphql.Field{Type: graphql.Boolean},
	},
})

var announcementType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Announcement",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.Int},
		"title":      &graphql.Field{Type: graphql.String},
		"content":    &graphql.Field{Type: graphql.String},
		"show":       &graphql.Field{Type: graphql.Boolean},
		"created_at": &graphql.Field{Type: graphql.String},
	},
})

var knowledgeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Knowledge",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.Int},
		"title":      &graphql.Field{Type: graphql.String},
		"content":    &graphql.Field{Type: graphql.String},
		"show":       &graphql.Field{Type: graphql.Boolean},
		"sort":       &graphql.Field{Type: graphql.Int},
		"created_at": &graphql.Field{Type: graphql.String},
		"updated_at": &graphql.Field{Type: graphql.String},
	},
})

var loginPayloadType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LoginPayload",
	Fields: graphql.Fields{
		"token": &graphql.Field{Type: graphql.String},
		"user":  &graphql.Field{Type: userType},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		// Query current user info
		"userInfo": &graphql.Field{
			Type: userType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				return u.(*models.User), nil
			},
		},
		// Query node list (available to all logged in users)
		"nodes": &graphql.Field{
			Type: graphql.NewList(nodeType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				_, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				var nodes []models.Node
				if err := database.DB.Where("show = ?", true).Find(&nodes).Error; err != nil {
					return nil, err
				}
				return nodes, nil
			},
		},
		// Query all users (Admin only)
		"adminUsers": &graphql.Field{
			Type: graphql.NewList(userType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}
				var users []models.User
				if err := database.DB.Find(&users).Error; err != nil {
					return nil, err
				}
				return users, nil
			},
		},
		// Query active announcements
		"announcements": &graphql.Field{
			Type: graphql.NewList(announcementType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var items []models.Announcement
				if err := database.DB.Where("show = ?", true).Order("created_at desc").Find(&items).Error; err != nil {
					return nil, err
				}
				return items, nil
			},
		},
		// Query all announcements (Admin only)
		"adminAnnouncements": &graphql.Field{
			Type: graphql.NewList(announcementType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}
				var items []models.Announcement
				if err := database.DB.Order("created_at desc").Find(&items).Error; err != nil {
					return nil, err
				}
				return items, nil
			},
		},
		// Query active knowledges (Visible to all logged in users)
		"knowledges": &graphql.Field{
			Type: graphql.NewList(knowledgeType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				_, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				var items []models.Knowledge
				if err := database.DB.Where("show = ?", true).Order("sort desc, created_at desc").Find(&items).Error; err != nil {
					return nil, err
				}
				return items, nil
			},
		},
		// Query all knowledges (Admin only)
		"adminKnowledges": &graphql.Field{
			Type: graphql.NewList(knowledgeType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}
				var items []models.Knowledge
				if err := database.DB.Order("sort desc, created_at desc").Find(&items).Error; err != nil {
					return nil, err
				}
				return items, nil
			},
		},
	},
})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		// Login mutation
		"login": &graphql.Field{
			Type: loginPayloadType,
			Args: graphql.FieldConfigArgument{
				"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				email := p.Args["email"].(string)
				password := p.Args["password"].(string)

				var user models.User
				if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
					return nil, errors.New("invalid email or password")
				}

				if user.Status != 1 {
					return nil, errors.New("account disabled")
				}

				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
				if err != nil {
					return nil, errors.New("invalid email or password")
				}

				token, err := middleware.GenerateToken(user.ID)
				if err != nil {
					return nil, errors.New("failed to generate token")
				}

				return map[string]interface{}{
					"token": token,
					"user":  &user,
				}, nil
			},
		},
		// Register mutation
		"register": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				email := p.Args["email"].(string)
				password := p.Args["password"].(string)

				var existingUser models.User
				if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
					return nil, errors.New("email address already registered")
				}

				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if err != nil {
					return nil, errors.New("failed to hash password")
				}

				newUser := models.User{
					Email:          email,
					Password:       string(hashedPassword),
					V2rayUUID:      uuid.New().String(),
					TrojanPassword: uuid.New().String()[:12],
					TotalTraffic:   50 * 1024 * 1024 * 1024,
					UsedTraffic:    0,
					ExpiredAt:      time.Now().AddDate(0, 1, 0),
					Status:         1,
					IsAdmin:        false,
				}

				if err := database.DB.Create(&newUser).Error; err != nil {
					return nil, errors.New("failed to create user")
				}

				return "Registration successful", nil
			},
		},
		// Create node mutation (Admin only)
		"createNode": &graphql.Field{
			Type: nodeType,
			Args: graphql.FieldConfigArgument{
				"name":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"type":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"address":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"port":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"traffic_rate": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				"settings":     &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}

				name := p.Args["name"].(string)
				nodeProto := p.Args["type"].(string)
				address := p.Args["address"].(string)
				port := p.Args["port"].(int)
				trafficRate := p.Args["traffic_rate"].(float64)
				settings, _ := p.Args["settings"].(string)

				node := models.Node{
					Name:        name,
					Type:        nodeProto,
					Address:     address,
					Port:        uint16(port),
					TrafficRate: float32(trafficRate),
					Settings:    settings,
					Show:        true,
				}

				if err := database.DB.Create(&node).Error; err != nil {
					return nil, errors.New("failed to create node")
				}

				return &node, nil
			},
		},
		// Create announcement mutation (Admin only)
		"createAnnouncement": &graphql.Field{
			Type: announcementType,
			Args: graphql.FieldConfigArgument{
				"title":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"content": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}

				title := p.Args["title"].(string)
				content := p.Args["content"].(string)

				item := models.Announcement{
					Title:   title,
					Content: content,
					Show:    true,
				}

				if err := database.DB.Create(&item).Error; err != nil {
					return nil, errors.New("failed to create announcement")
				}

				return &item, nil
			},
		},
		// Toggle announcement visibility mutation (Admin only)
		"toggleAnnouncement": &graphql.Field{
			Type: announcementType,
			Args: graphql.FieldConfigArgument{
				"id":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"show": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Boolean)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}

				id := p.Args["id"].(int)
				show := p.Args["show"].(bool)

				var item models.Announcement
				if err := database.DB.First(&item, id).Error; err != nil {
					return nil, errors.New("announcement not found")
				}

				item.Show = show
				if err := database.DB.Save(&item).Error; err != nil {
					return nil, errors.New("failed to update announcement")
				}

				return &item, nil
			},
		},
		// Update announcement mutation (Admin only)
		"updateAnnouncement": &graphql.Field{
			Type: announcementType,
			Args: graphql.FieldConfigArgument{
				"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"title":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"content": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"show":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Boolean)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}

				id := p.Args["id"].(int)
				title := p.Args["title"].(string)
				content := p.Args["content"].(string)
				show := p.Args["show"].(bool)

				var item models.Announcement
				if err := database.DB.First(&item, id).Error; err != nil {
					return nil, errors.New("announcement not found")
				}

				item.Title = title
				item.Content = content
				item.Show = show

				if err := database.DB.Save(&item).Error; err != nil {
					return nil, errors.New("failed to update announcement")
				}

				return &item, nil
			},
		},
		// Delete announcement mutation (Admin only)
		"deleteAnnouncement": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}

				id := p.Args["id"].(int)

				if err := database.DB.Delete(&models.Announcement{}, id).Error; err != nil {
					return nil, errors.New("failed to delete announcement")
				}

				return true, nil
			},
		},
		// Create knowledge mutation (Admin only)
		"createKnowledge": &graphql.Field{
			Type: knowledgeType,
			Args: graphql.FieldConfigArgument{
				"title":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"content": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"sort":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}

				title := p.Args["title"].(string)
				content := p.Args["content"].(string)
				sort := p.Args["sort"].(int)

				item := models.Knowledge{
					Title:   title,
					Content: content,
					Show:    true,
					Sort:    sort,
				}

				if err := database.DB.Create(&item).Error; err != nil {
					return nil, errors.New("failed to create knowledge")
				}

				return &item, nil
			},
		},
		// Update knowledge mutation (Admin only)
		"updateKnowledge": &graphql.Field{
			Type: knowledgeType,
			Args: graphql.FieldConfigArgument{
				"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"title":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"content": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"show":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Boolean)},
				"sort":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}

				id := p.Args["id"].(int)
				title := p.Args["title"].(string)
				content := p.Args["content"].(string)
				show := p.Args["show"].(bool)
				sort := p.Args["sort"].(int)

				var item models.Knowledge
				if err := database.DB.First(&item, id).Error; err != nil {
					return nil, errors.New("knowledge article not found")
				}

				item.Title = title
				item.Content = content
				item.Show = show
				item.Sort = sort

				if err := database.DB.Save(&item).Error; err != nil {
					return nil, errors.New("failed to update knowledge article")
				}

				return &item, nil
			},
		},
		// Delete knowledge mutation (Admin only)
		"deleteKnowledge": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}

				id := p.Args["id"].(int)

				if err := database.DB.Delete(&models.Knowledge{}, id).Error; err != nil {
					return nil, errors.New("failed to delete knowledge article")
				}

				return true, nil
			},
		},
		// Update user mutation (Admin only)
		"updateUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id":            &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"email":         &graphql.ArgumentConfig{Type: graphql.String},
				"password":      &graphql.ArgumentConfig{Type: graphql.String},
				"balance":       &graphql.ArgumentConfig{Type: graphql.Float},
				"speed_limit":   &graphql.ArgumentConfig{Type: graphql.Int},
				"device_limit":  &graphql.ArgumentConfig{Type: graphql.Int},
				"total_traffic": &graphql.ArgumentConfig{Type: graphql.Float},
				"used_traffic":  &graphql.ArgumentConfig{Type: graphql.Float},
				"expired_at":    &graphql.ArgumentConfig{Type: graphql.String},
				"status":        &graphql.ArgumentConfig{Type: graphql.Int},
				"is_admin":      &graphql.ArgumentConfig{Type: graphql.Boolean},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}

				id := p.Args["id"].(int)
				var user models.User
				if err := database.DB.First(&user, id).Error; err != nil {
					return nil, errors.New("user not found")
				}

				if val, ok := p.Args["email"]; ok {
					email := val.(string)
					// Check if email already exists
					var existing models.User
					if err := database.DB.Where("email = ? AND id != ?", email, id).First(&existing).Error; err == nil {
						return nil, errors.New("email already in use")
					}
					user.Email = email
				}

				if val, ok := p.Args["password"]; ok {
					password := val.(string)
					if password != "" {
						hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
						if err != nil {
							return nil, errors.New("failed to hash password")
						}
						user.Password = string(hashed)
					}
				}

				if val, ok := p.Args["balance"]; ok {
					user.Balance = val.(float64)
				}

				if val, ok := p.Args["speed_limit"]; ok {
					user.SpeedLimit = uint32(val.(int))
				}

				if val, ok := p.Args["device_limit"]; ok {
					user.DeviceLimit = uint32(val.(int))
				}

				if val, ok := p.Args["total_traffic"]; ok {
					user.TotalTraffic = uint64(val.(float64))
				}

				if val, ok := p.Args["used_traffic"]; ok {
					user.UsedTraffic = uint64(val.(float64))
				}

				if val, ok := p.Args["expired_at"]; ok {
					expiredStr := val.(string)
					if expiredStr == "" {
						user.ExpiredAt = time.Time{}
					} else {
						t, err := time.Parse(time.RFC3339, expiredStr)
						if err != nil {
							// Try simple date format YYYY-MM-DD
							t, err = time.Parse("2006-01-02", expiredStr)
							if err != nil {
								// Try standard date time
								t, err = time.Parse("2006-01-02 15:04:05", expiredStr)
							}
						}
						if err == nil {
							user.ExpiredAt = t
						} else {
							return nil, fmt.Errorf("invalid date format: %v", err)
						}
					}
				}

				if val, ok := p.Args["status"]; ok {
					user.Status = int8(val.(int))
				}

				if val, ok := p.Args["is_admin"]; ok {
					user.IsAdmin = val.(bool)
				}

				if err := database.DB.Save(&user).Error; err != nil {
					return nil, errors.New("failed to update user")
				}

				return &user, nil
			},
		},
		// Delete user mutation (Admin only)
		"deleteUser": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				u, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := u.(*models.User)
				if !currentUser.IsAdmin {
					return nil, errors.New("admin privilege required")
				}

				id := p.Args["id"].(int)
				if currentUser.ID == uint(id) {
					return nil, errors.New("cannot delete yourself")
				}

				var targetUser models.User
				if err := database.DB.First(&targetUser, id).Error; err != nil {
					return nil, errors.New("user not found")
				}
				if targetUser.IsAdmin {
					return nil, errors.New("cannot delete an administrator")
				}

				if err := database.DB.Delete(&targetUser).Error; err != nil {
					return nil, errors.New("failed to delete user")
				}

				return true, nil
			},
		},
	},
})

var graphqlSchema graphql.Schema

func init() {
	var err error
	graphqlSchema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize GraphQL Schema: %v", err))
	}
}

type GraphQLRequestBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func GraphQLHandler(c *gin.Context) {
	var body GraphQLRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:         graphqlSchema,
		RequestString:  body.Query,
		VariableValues: body.Variables,
		Context:        c,
	})

	if len(result.Errors) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"errors": result.Errors,
			"data":   result.Data,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
