package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hy-board-backend",
	Short: "HY-Board is a high-aesthetic XrayR control panel backend",
	Long:  `A fast and lightweight management system backend for XrayR subscription and node management, written in Go.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add global flags here if needed (e.g. config path)
	rootCmd.PersistentFlags().StringP("config", "c", "config.yaml", "config file path")
}
