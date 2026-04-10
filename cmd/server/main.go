package main

import (
	"log"

	"opentunnel/server/config"
	"opentunnel/server/controller"
	"opentunnel/server/database"
	"opentunnel/server/endpoints"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env")
	cfg := config.NewConfig()

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer database.Close()

	// Start tunnel controller first (non-blocking goroutine)
	go func() {
		ctrl := controller.NewTunnelController(cfg)
		ctrl.Run()
	}()

	// HTTP server blocks here - must be last
	router := endpoints.SetupRouter(cfg)
	log.Printf("HTTP server starting on :%s", cfg.HTTPPort)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
