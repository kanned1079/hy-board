package cmd

import (
	"fmt"
	"log"

	"hy-board-backend/config"
	"hy-board-backend/database"
	"hy-board-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	devMode  bool
	prodMode bool
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the HTTP API server",
	Long:  `Run the Gin HTTP server to handle Client, Admin, and UniProxy API endpoints.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		cfg := config.LoadConfig(configPath)

		// Override server mode if dev/prod flags are set
		if devMode {
			cfg.Server.Mode = "debug"
		} else if prodMode {
			cfg.Server.Mode = "release"
		}

		gin.SetMode(cfg.Server.Mode)

		// Init Database & Migrate
		database.InitDB()
		database.Migrate()

		// Start real-time Ticket WebSocket Hub
		go routes.StartHub()

		r := routes.SetupRouter()

		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Printf("Starting server in [%s] mode on %s...", cfg.Server.Mode, addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("Server failed to run: %v", err)
		}
	},
}

func init() {
	serverCmd.Flags().BoolVar(&devMode, "dev", false, "Run in development mode (Gin debug)")
	serverCmd.Flags().BoolVar(&prodMode, "prod", false, "Run in production mode (Gin release)")
	rootCmd.AddCommand(serverCmd)
}
