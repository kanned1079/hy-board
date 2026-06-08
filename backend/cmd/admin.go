package cmd

import (
	"log"
	"time"

	"hy-board-backend/config"
	"hy-board-backend/database"
	"hy-board-backend/models"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var (
	adminEmail    string
	adminPassword string
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Manage administrative commands",
}

var initAdminCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize default administrator user",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		config.LoadConfig(configPath)

		database.InitDB()
		database.Migrate()

		if adminEmail == "" || adminPassword == "" {
			log.Fatal("Email and password are required. Use flags --email and --password.")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		// Generate random UUID for VMess/VLESS client subscription
		newUUID := uuid.New().String()

		admin := models.User{
			Email:          adminEmail,
			Password:       string(hashedPassword),
			V2rayUUID:      newUUID,
			TrojanPassword: uuid.New().String()[:12],
			TotalTraffic:   1024 * 1024 * 1024 * 1024, // 1TB default
			ExpiredAt:      time.Now().AddDate(10, 0, 0), // 10 years expiration
			Status:         1,
			IsAdmin:        true,
		}

		// Check if user already exists
		var existingUser models.User
		if err := database.DB.Where("email = ?", adminEmail).First(&existingUser).Error; err == nil {
			log.Fatalf("User with email %s already exists.", adminEmail)
		}

		if err := database.DB.Create(&admin).Error; err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}

		log.Printf("Administrator user created successfully!")
		log.Printf("Email: %s", adminEmail)
		log.Printf("V2ray UUID: %s", newUUID)
	},
}

func init() {
	initAdminCmd.Flags().StringVar(&adminEmail, "email", "admin@example.com", "Administrator email address")
	initAdminCmd.Flags().StringVar(&adminPassword, "password", "admin123456", "Administrator password")
	adminCmd.AddCommand(initAdminCmd)
	rootCmd.AddCommand(adminCmd)
}
