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
		// Query admin system settings (Admin only)
		"systemSettings": &graphql.Field{
			Type: systemSettingsType,
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

				return map[string]interface{}{
					"site_name":                 database.GetSetting("site_name"),
					"site_description":          database.GetSetting("site_description"),
					"site_url":                  database.GetSetting("site_url"),
					"tos_url":                   database.GetSetting("tos_url"),
					"stop_register":             database.GetSettingBool("stop_register"),
					"currency_unit":             database.GetSetting("currency_unit"),
					"currency_symbol":           database.GetSetting("currency_symbol"),
					"email_verify":              database.GetSettingBool("email_verify"),
					"ban_gmail_alias":           database.GetSettingBool("ban_gmail_alias"),
					"ip_register_limit":         database.GetSettingBool("ip_register_limit"),
					"ip_register_limit_count":   database.GetSettingInt("ip_register_limit_count"),
					"ip_register_limit_penalty": database.GetSettingInt("ip_register_limit_penalty"),
					"theme_color":               database.GetSetting("theme_color"),
					"home_background":           database.GetSetting("home_background"),
					"uniproxy_token":            database.GetSetting("uniproxy_token"),
					"node_pull_interval":        database.GetSettingInt("node_pull_interval"),
					"node_push_interval":        database.GetSettingInt("node_push_interval"),
					"smtp_host":                 database.GetSetting("smtp_host"),
					"smtp_port":                 database.GetSettingInt("smtp_port"),
					"smtp_encryption":           database.GetSetting("smtp_encryption"),
					"smtp_username":             database.GetSetting("smtp_username"),
					"smtp_password":             database.GetSetting("smtp_password"),
					"smtp_from":                 database.GetSetting("smtp_from"),
					"app_win":                   database.GetSetting("app_win"),
					"app_macos":                 database.GetSetting("app_macos"),
					"app_linux":                 database.GetSetting("app_linux"),
					"app_android":               database.GetSetting("app_android"),
					"app_ios":                   database.GetSetting("app_ios"),
				}, nil
			},
		},
		// Query public site settings (accessible to anyone)
		"publicSettings": &graphql.Field{
			Type: publicSettingsType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return map[string]interface{}{
					"site_name":        database.GetSetting("site_name"),
					"site_description": database.GetSetting("site_description"),
					"site_url":         database.GetSetting("site_url"),
					"tos_url":          database.GetSetting("tos_url"),
					"stop_register":    database.GetSettingBool("stop_register"),
					"currency_unit":    database.GetSetting("currency_unit"),
					"currency_symbol":  database.GetSetting("currency_symbol"),
					"theme_color":      database.GetSetting("theme_color"),
					"home_background":  database.GetSetting("home_background"),
					"app_win":          database.GetSetting("app_win"),
					"app_macos":        database.GetSetting("app_macos"),
					"app_linux":        database.GetSetting("app_linux"),
					"app_android":      database.GetSetting("app_android"),
					"app_ios":          database.GetSetting("app_ios"),
				}, nil
			},
		},
	},
})
