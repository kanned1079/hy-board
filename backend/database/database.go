package database

import (
	"fmt"
	"log"

	"hy-board-backend/config"
	"hy-board-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	var err error
	cfg := config.GlobalConfig.Database

	switch cfg.Type {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(cfg.Sqlite.File), &gorm.Config{})
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.Mysql.User,
			cfg.Mysql.Password,
			cfg.Mysql.Host,
			cfg.Mysql.Port,
			cfg.Mysql.DBName,
			cfg.Mysql.Charset,
		)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
			cfg.Postgres.Host,
			cfg.Postgres.User,
			cfg.Postgres.Password,
			cfg.Postgres.DBName,
			cfg.Postgres.Port,
			cfg.Postgres.SSLMode,
		)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		log.Fatalf("Unsupported database type: %s", cfg.Type)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully.")
	return DB
}

func Migrate() {
	if DB == nil {
		log.Fatal("Database not initialized. Cannot migrate.")
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.Node{},
		&models.TrafficLog{},
		&models.Announcement{},
		&models.Knowledge{},
		&models.Ticket{},
		&models.TicketMessage{},
		&models.Group{},
		&models.Plan{},
		&models.SystemSetting{},
	)
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Println("Database schema auto-migrated successfully.")

	// Auto-seed default settings if none exist
	var settingsCount int64
	DB.Model(&models.SystemSetting{}).Count(&settingsCount)
	if settingsCount == 0 {
		token := "secret-uniproxy-token"
		if config.GlobalConfig != nil && config.GlobalConfig.Auth.UniProxyToken != "" {
			token = config.GlobalConfig.Auth.UniProxyToken
		}
		defaultSettings := map[string]string{
			"site_name":                 "HY-Board",
			"site_description":          "High-Aesthetic XrayR Management Panel",
			"site_url":                  "http://localhost:3000",
			"tos_url":                   "",
			"stop_register":             "false",
			"currency_unit":             "CNY",
			"currency_symbol":           "¥",
			"email_verify":              "false",
			"ban_gmail_alias":           "false",
			"ip_register_limit":         "false",
			"ip_register_limit_count":   "5",
			"ip_register_limit_penalty": "60",
			"theme_color":               "green",
			"home_background":           "",
			"uniproxy_token":            token,
			"node_pull_interval":        "60",
			"node_push_interval":        "60",
			"smtp_host":                 "",
			"smtp_port":                 "465",
			"smtp_encryption":           "SSL",
			"smtp_username":             "",
			"smtp_password":             "",
			"smtp_from":                 "",
			"app_win":                   "",
			"app_macos":                 "",
			"app_linux":                 "",
			"app_android":               "",
			"app_ios":                   "",
		}
		for k, v := range defaultSettings {
			DB.Create(&models.SystemSetting{Key: k, Value: v})
		}
		log.Println("Default system settings seeded successfully.")
	}

	// Auto-seed default permission groups if none exist
	var count int64
	DB.Model(&models.Group{}).Count(&count)
	if count == 0 {
		DB.Create(&models.Group{ID: 1, Name: "S1 亞太標準訂閱", Description: "提供亞太地區標準節點，限速 150Mbps"})
		DB.Create(&models.Group{ID: 2, Name: "S2 全球至尊訂閱", Description: "提供全球（亞太+美洲）節點，限速 300Mbps"})
		DB.Create(&models.Group{ID: 99, Name: "Admin 系統管理組", Description: "系統管理與內部測試用權限組"})
		log.Println("Default permission groups seeded successfully.")
	}

	// Auto-seed default plans if none exist
	var planCount int64
	DB.Model(&models.Plan{}).Count(&planCount)
	if planCount == 0 {
		DB.Create(&models.Plan{
			ID:          1,
			Name:        "S1 亞太標準版",
			Description: "適合日常網頁瀏覽、社交媒體與普通通訊",
			Price:       5.99,
			Traffic:     150,
			SpeedLimit:  150,
			DeviceLimit: 3,
			ExpiryDays:  30,
			GroupID:     1,
			Show:        true,
		})
		DB.Create(&models.Plan{
			ID:          2,
			Name:        "S2 全球至尊版",
			Description: "適合影音串流、跨國遊戲與極速傳輸需求",
			Price:       11.99,
			Traffic:     400,
			SpeedLimit:  300,
			DeviceLimit: 5,
			ExpiryDays:  30,
			GroupID:     2,
			Show:        true,
		})
		log.Println("Default subscription plans seeded successfully.")
	}
}
