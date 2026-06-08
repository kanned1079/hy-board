package routes

import (
	"errors"
	"fmt"
	"time"

	"hy-board-backend/database"
	"hy-board-backend/middleware"
	"hy-board-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"
)

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
				"group_id":     &graphql.ArgumentConfig{Type: graphql.Int},
				"group_ids":    &graphql.ArgumentConfig{Type: graphql.String},
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

				groupID := 1
				if val, ok := p.Args["group_id"]; ok {
					groupID = val.(int)
				}

				groupIDsStr := fmt.Sprintf("%d", groupID)
				if val, ok := p.Args["group_ids"]; ok {
					groupIDsStr = val.(string)
				}

				node := models.Node{
					Name:        name,
					Type:        nodeProto,
					Address:     address,
					Port:        uint16(port),
					TrafficRate: float32(trafficRate),
					Settings:    settings,
					Show:        true,
					GroupID:     uint(groupID),
					GroupIDs:    groupIDsStr,
				}

				if err := database.DB.Create(&node).Error; err != nil {
					return nil, errors.New("failed to create node")
				}

				return &node, nil
			},
		},
		// Update node mutation (Admin only)
		"updateNode": &graphql.Field{
			Type: nodeType,
			Args: graphql.FieldConfigArgument{
				"id":           &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":         &graphql.ArgumentConfig{Type: graphql.String},
				"type":         &graphql.ArgumentConfig{Type: graphql.String},
				"address":      &graphql.ArgumentConfig{Type: graphql.String},
				"port":         &graphql.ArgumentConfig{Type: graphql.Int},
				"traffic_rate": &graphql.ArgumentConfig{Type: graphql.Float},
				"settings":     &graphql.ArgumentConfig{Type: graphql.String},
				"group_id":     &graphql.ArgumentConfig{Type: graphql.Int},
				"group_ids":    &graphql.ArgumentConfig{Type: graphql.String},
				"show":         &graphql.ArgumentConfig{Type: graphql.Boolean},
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
				var node models.Node
				if err := database.DB.First(&node, id).Error; err != nil {
					return nil, errors.New("node not found")
				}

				if val, ok := p.Args["name"]; ok {
					node.Name = val.(string)
				}
				if val, ok := p.Args["type"]; ok {
					node.Type = val.(string)
				}
				if val, ok := p.Args["address"]; ok {
					node.Address = val.(string)
				}
				if val, ok := p.Args["port"]; ok {
					node.Port = uint16(val.(int))
				}
				if val, ok := p.Args["traffic_rate"]; ok {
					node.TrafficRate = float32(val.(float64))
				}
				if val, ok := p.Args["settings"]; ok {
					node.Settings = val.(string)
				}
				if val, ok := p.Args["group_id"]; ok {
					node.GroupID = uint(val.(int))
				}
				if val, ok := p.Args["group_ids"]; ok {
					node.GroupIDs = val.(string)
				}
				if val, ok := p.Args["show"]; ok {
					node.Show = val.(bool)
				}

				if err := database.DB.Save(&node).Error; err != nil {
					return nil, errors.New("failed to update node")
				}

				return &node, nil
			},
		},
		// Delete node mutation (Admin only)
		"deleteNode": &graphql.Field{
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
				if err := database.DB.Delete(&models.Node{}, id).Error; err != nil {
					return nil, errors.New("failed to delete node")
				}

				return true, nil
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
				"group_id":      &graphql.ArgumentConfig{Type: graphql.Int},
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
							t, err = time.Parse("2006-01-02", expiredStr)
						}
						if err != nil {
							t, err = time.Parse("2006-01-02 15:04:05", expiredStr)
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
				if val, ok := p.Args["group_id"]; ok {
					user.GroupID = uint(val.(int))
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
		// Create group mutation (Admin only)
		"createGroup": &graphql.Field{
			Type: groupType,
			Args: graphql.FieldConfigArgument{
				"name":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"description": &graphql.ArgumentConfig{Type: graphql.String},
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
				desc, _ := p.Args["description"].(string)

				group := models.Group{
					Name:        name,
					Description: desc,
				}

				if err := database.DB.Create(&group).Error; err != nil {
					return nil, errors.New("failed to create group")
				}
				return &group, nil
			},
		},
		// Update group mutation (Admin only)
		"updateGroup": &graphql.Field{
			Type: groupType,
			Args: graphql.FieldConfigArgument{
				"id":          &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"description": &graphql.ArgumentConfig{Type: graphql.String},
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
				name := p.Args["name"].(string)
				desc, _ := p.Args["description"].(string)

				var group models.Group
				if err := database.DB.First(&group, id).Error; err != nil {
					return nil, errors.New("group not found")
				}

				group.Name = name
				group.Description = desc

				if err := database.DB.Save(&group).Error; err != nil {
					return nil, errors.New("failed to update group")
				}
				return &group, nil
			},
		},
		// Delete group mutation (Admin only)
		"deleteGroup": &graphql.Field{
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
				if err := database.DB.Delete(&models.Group{}, id).Error; err != nil {
					return nil, errors.New("failed to delete group")
				}
				return true, nil
			},
		},
		// Create subscription plan (Admin only)
		"createPlan": &graphql.Field{
			Type: planType,
			Args: graphql.FieldConfigArgument{
				"name":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"description":  &graphql.ArgumentConfig{Type: graphql.String},
				"price":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				"traffic":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"speed_limit":  &graphql.ArgumentConfig{Type: graphql.Int},
				"device_limit": &graphql.ArgumentConfig{Type: graphql.Int},
				"expiry_days":  &graphql.ArgumentConfig{Type: graphql.Int},
				"group_id":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"show":         &graphql.ArgumentConfig{Type: graphql.Boolean},
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
				description, _ := p.Args["description"].(string)
				price := p.Args["price"].(float64)
				traffic := p.Args["traffic"].(int)

				speedLimit := 0
				if val, ok := p.Args["speed_limit"]; ok {
					speedLimit = val.(int)
				}
				deviceLimit := 0
				if val, ok := p.Args["device_limit"]; ok {
					deviceLimit = val.(int)
				}
				expiryDays := 30
				if val, ok := p.Args["expiry_days"]; ok {
					expiryDays = val.(int)
				}
				groupID := p.Args["group_id"].(int)
				show := true
				if val, ok := p.Args["show"]; ok {
					show = val.(bool)
				}

				plan := models.Plan{
					Name:        name,
					Description: description,
					Price:       price,
					Traffic:     uint64(traffic),
					SpeedLimit:  uint32(speedLimit),
					DeviceLimit: uint32(deviceLimit),
					ExpiryDays:  uint32(expiryDays),
					GroupID:     uint(groupID),
					Show:        show,
				}

				if err := database.DB.Create(&plan).Error; err != nil {
					return nil, errors.New("failed to create plan")
				}

				return &plan, nil
			},
		},
		// Update subscription plan (Admin only)
		"updatePlan": &graphql.Field{
			Type: planType,
			Args: graphql.FieldConfigArgument{
				"id":           &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":         &graphql.ArgumentConfig{Type: graphql.String},
				"description":  &graphql.ArgumentConfig{Type: graphql.String},
				"price":        &graphql.ArgumentConfig{Type: graphql.Float},
				"traffic":      &graphql.ArgumentConfig{Type: graphql.Int},
				"speed_limit":  &graphql.ArgumentConfig{Type: graphql.Int},
				"device_limit": &graphql.ArgumentConfig{Type: graphql.Int},
				"expiry_days":  &graphql.ArgumentConfig{Type: graphql.Int},
				"group_id":     &graphql.ArgumentConfig{Type: graphql.Int},
				"show":         &graphql.ArgumentConfig{Type: graphql.Boolean},
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
				var plan models.Plan
				if err := database.DB.First(&plan, id).Error; err != nil {
					return nil, errors.New("plan not found")
				}

				if val, ok := p.Args["name"]; ok {
					plan.Name = val.(string)
				}
				if val, ok := p.Args["description"]; ok {
					plan.Description = val.(string)
				}
				if val, ok := p.Args["price"]; ok {
					plan.Price = val.(float64)
				}
				if val, ok := p.Args["traffic"]; ok {
					plan.Traffic = uint64(val.(int))
				}
				if val, ok := p.Args["speed_limit"]; ok {
					plan.SpeedLimit = uint32(val.(int))
				}
				if val, ok := p.Args["device_limit"]; ok {
					plan.DeviceLimit = uint32(val.(int))
				}
				if val, ok := p.Args["expiry_days"]; ok {
					plan.ExpiryDays = uint32(val.(int))
				}
				if val, ok := p.Args["group_id"]; ok {
					plan.GroupID = uint(val.(int))
				}
				if val, ok := p.Args["show"]; ok {
					plan.Show = val.(bool)
				}

				if err := database.DB.Save(&plan).Error; err != nil {
					return nil, errors.New("failed to update plan")
				}

				return &plan, nil
			},
		},
		// Delete subscription plan (Admin only)
		"deletePlan": &graphql.Field{
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
				if err := database.DB.Delete(&models.Plan{}, id).Error; err != nil {
					return nil, errors.New("failed to delete plan")
				}

				return true, nil
			},
		},
		// Purchase plan (Users)
		"purchasePlan": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"plan_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
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

				planID := p.Args["plan_id"].(int)
				var plan models.Plan
				if err := database.DB.First(&plan, planID).Error; err != nil {
					return nil, errors.New("subscription plan not found")
				}

				if !plan.Show && !currentUser.IsAdmin {
					return nil, errors.New("plan is not available for purchase")
				}

				// Find user in database to ensure we have fresh data
				var dbUser models.User
				if err := database.DB.First(&dbUser, currentUser.ID).Error; err != nil {
					return nil, errors.New("user not found")
				}

				if dbUser.Balance < plan.Price {
					return nil, errors.New("insufficient balance, please top up first")
				}

				// Start transaction
				tx := database.DB.Begin()

				// Deduct balance
				dbUser.Balance -= plan.Price

				// Calculate expiry
				var newExpiry time.Time
				if dbUser.ExpiredAt.After(time.Now()) && dbUser.GroupID == plan.GroupID {
					newExpiry = dbUser.ExpiredAt.Add(time.Duration(plan.ExpiryDays) * 24 * time.Hour)
				} else {
					newExpiry = time.Now().Add(time.Duration(plan.ExpiryDays) * 24 * time.Hour)
				}
				dbUser.ExpiredAt = newExpiry

				// Reset and set traffic
				dbUser.UsedTraffic = 0
				dbUser.TotalTraffic = plan.Traffic * 1024 * 1024 * 1024 // GB to Bytes

				// Set limits and group
				dbUser.SpeedLimit = plan.SpeedLimit
				dbUser.DeviceLimit = plan.DeviceLimit
				dbUser.GroupID = plan.GroupID

				if err := tx.Save(&dbUser).Error; err != nil {
					tx.Rollback()
					return nil, errors.New("failed to purchase plan")
				}

				tx.Commit()

				// Update user in context as well
				c.Set("user", &dbUser)

				return &dbUser, nil
			},
		},
	},
})
