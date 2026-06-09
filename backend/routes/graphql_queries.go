package routes

import (
	"errors"

	"hy-board-backend/database"
	"hy-board-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

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
				uVal, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := uVal.(*models.User)

				var nodes []models.Node
				query := database.DB
				if !currentUser.IsAdmin {
					query = query.Where("show = ?", true)
				}
				if err := query.Find(&nodes).Error; err != nil {
					return nil, err
				}
				return nodes, nil
			},
		},
		// Query all groups
		"groups": &graphql.Field{
			Type: graphql.NewList(groupType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var groups []models.Group
				if err := database.DB.Find(&groups).Error; err != nil {
					return nil, err
				}
				return groups, nil
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
		// Query subscription plans list
		"plans": &graphql.Field{
			Type: graphql.NewList(planType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				uVal, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := uVal.(*models.User)

				var plans []models.Plan
				query := database.DB
				if !currentUser.IsAdmin {
					query = query.Where("show = ?", true)
				}
				if err := query.Find(&plans).Error; err != nil {
					return nil, err
				}
				return plans, nil
			},
		},
		// Query current user's traffic logs
		"trafficLogs": &graphql.Field{
			Type: graphql.NewList(trafficLogType),
			Args: graphql.FieldConfigArgument{
				"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
				"offset": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, ok := p.Context.(*gin.Context)
				if !ok {
					return nil, errors.New("invalid context")
				}
				uVal, exists := c.Get("user")
				if !exists {
					return nil, errors.New("unauthorized")
				}
				currentUser := uVal.(*models.User)

				limit := 50
				if lVal, ok := p.Args["limit"].(int); ok && lVal > 0 {
					limit = lVal
					if limit > 100 {
						limit = 100
					}
				}

				offset := 0
				if oVal, ok := p.Args["offset"].(int); ok && oVal >= 0 {
					offset = oVal
				}

				var logs []models.TrafficLog
				if err := database.DB.Preload("Node").
					Where("user_id = ?", currentUser.ID).
					Order("created_at desc").
					Limit(limit).
					Offset(offset).
					Find(&logs).Error; err != nil {
					return nil, err
				}
				return logs, nil
			},
		},
	},
})
