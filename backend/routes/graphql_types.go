package routes

import (
	"time"

	"hy-board-backend/database"
	"hy-board-backend/models"

	"github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":              &graphql.Field{Type: graphql.Int},
		"email":           &graphql.Field{Type: graphql.String},
		"username":        &graphql.Field{Type: graphql.String},
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
		"group_id":        &graphql.Field{Type: graphql.Int},
		"register_ip":     &graphql.Field{Type: graphql.String},
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
		"group_id":     &graphql.Field{Type: graphql.Int},
		"group_ids":    &graphql.Field{Type: graphql.String},
		"online": &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				node, ok := p.Source.(models.Node)
				if !ok {
					return false, nil
				}
				// Node is online if it has checked in within the last 45 seconds
				return time.Since(node.UpdatedAt) < 45*time.Second, nil
			},
		},
		"status": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				node, ok := p.Source.(models.Node)
				if !ok {
					return "offline", nil
				}

				// 1. Red: Offline / Cannot connect
				if time.Since(node.UpdatedAt) >= 45*time.Second {
					return "offline", nil
				}

				// 2. Green: Active / Someone is connected
				var count int64
				oneMinuteAgo := time.Now().Add(-1 * time.Minute)
				database.DB.Model(&models.TrafficLog{}).
					Where("node_id = ? AND created_at > ? AND (up > 0 OR down > 0)", node.ID, oneMinuteAgo).
					Count(&count)

				if count > 0 {
					return "active", nil
				}

				// 3. Yellow: Idle / Online but no one connected
				return "idle", nil
			},
		},
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

var groupType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Group",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.Int},
		"name":        &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"created_at":  &graphql.Field{Type: graphql.String},
	},
})

var planType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Plan",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.Int},
		"name":         &graphql.Field{Type: graphql.String},
		"description":  &graphql.Field{Type: graphql.String},
		"price":        &graphql.Field{Type: graphql.Float},
		"traffic":      &graphql.Field{Type: graphql.Int},
		"speed_limit":  &graphql.Field{Type: graphql.Int},
		"device_limit": &graphql.Field{Type: graphql.Int},
		"expiry_days":  &graphql.Field{Type: graphql.Int},
		"group_id":     &graphql.Field{Type: graphql.Int},
		"show":         &graphql.Field{Type: graphql.Boolean},
		"created_at":   &graphql.Field{Type: graphql.String},
	},
})

var trafficLogType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TrafficLog",
	Fields: graphql.Fields{
		"id":      &graphql.Field{Type: graphql.Int},
		"user_id": &graphql.Field{Type: graphql.Int},
		"node_id": &graphql.Field{Type: graphql.Int},
		"node_name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				log, ok := p.Source.(models.TrafficLog)
				if ok {
					return log.Node.Name, nil
				}
				return "", nil
			},
		},
		"node_rate": &graphql.Field{
			Type: graphql.Float,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				log, ok := p.Source.(models.TrafficLog)
				if ok {
					return log.Node.TrafficRate, nil
				}
				return 1.0, nil
			},
		},
		"up":         &graphql.Field{Type: graphql.Float},
		"down":       &graphql.Field{Type: graphql.Float},
		"created_at": &graphql.Field{Type: graphql.String},
	},
})

var systemSettingsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SystemSettings",
	Fields: graphql.Fields{
		"site_name":                 &graphql.Field{Type: graphql.String},
		"site_description":          &graphql.Field{Type: graphql.String},
		"site_url":                  &graphql.Field{Type: graphql.String},
		"tos_url":                   &graphql.Field{Type: graphql.String},
		"stop_register":             &graphql.Field{Type: graphql.Boolean},
		"currency_unit":             &graphql.Field{Type: graphql.String},
		"currency_symbol":           &graphql.Field{Type: graphql.String},
		"email_verify":              &graphql.Field{Type: graphql.Boolean},
		"ban_gmail_alias":           &graphql.Field{Type: graphql.Boolean},
		"ip_register_limit":         &graphql.Field{Type: graphql.Boolean},
		"ip_register_limit_count":   &graphql.Field{Type: graphql.Int},
		"ip_register_limit_penalty": &graphql.Field{Type: graphql.Int},
		"theme_color":               &graphql.Field{Type: graphql.String},
		"home_background":           &graphql.Field{Type: graphql.String},
		"uniproxy_token":            &graphql.Field{Type: graphql.String},
		"node_pull_interval":        &graphql.Field{Type: graphql.Int},
		"node_push_interval":        &graphql.Field{Type: graphql.Int},
		"smtp_host":                 &graphql.Field{Type: graphql.String},
		"smtp_port":                 &graphql.Field{Type: graphql.Int},
		"smtp_encryption":           &graphql.Field{Type: graphql.String},
		"smtp_username":             &graphql.Field{Type: graphql.String},
		"smtp_password":             &graphql.Field{Type: graphql.String},
		"smtp_from":                 &graphql.Field{Type: graphql.String},
		"app_win":                   &graphql.Field{Type: graphql.String},
		"app_macos":                 &graphql.Field{Type: graphql.String},
		"app_linux":                 &graphql.Field{Type: graphql.String},
		"app_android":               &graphql.Field{Type: graphql.String},
		"app_ios":                   &graphql.Field{Type: graphql.String},
	},
})

var publicSettingsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PublicSettings",
	Fields: graphql.Fields{
		"site_name":        &graphql.Field{Type: graphql.String},
		"site_description": &graphql.Field{Type: graphql.String},
		"site_url":         &graphql.Field{Type: graphql.String},
		"tos_url":          &graphql.Field{Type: graphql.String},
		"stop_register":    &graphql.Field{Type: graphql.Boolean},
		"currency_unit":    &graphql.Field{Type: graphql.String},
		"currency_symbol":  &graphql.Field{Type: graphql.String},
		"theme_color":      &graphql.Field{Type: graphql.String},
		"home_background":  &graphql.Field{Type: graphql.String},
		"app_win":          &graphql.Field{Type: graphql.String},
		"app_macos":        &graphql.Field{Type: graphql.String},
		"app_linux":        &graphql.Field{Type: graphql.String},
		"app_android":      &graphql.Field{Type: graphql.String},
		"app_ios":          &graphql.Field{Type: graphql.String},
	},
})
