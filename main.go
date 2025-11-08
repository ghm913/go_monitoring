package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"

	"ghm913/go_monitoring/api"
	"ghm913/go_monitoring/config"
	"ghm913/go_monitoring/services"
)

func main() {
	// Setup logging
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("\n%s", cfg)

	// Create monitor service
	mon := services.NewMonitor(cfg)

	// Configure routes
	handler := api.NewHandler(mon)
	r := gin.New()
	r.Use(gin.Recovery())
	handler.SetupRoutes(r)

	// Create cancellable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start monitoring in background
	wg.Add(1)
	go func() {
		defer wg.Done()
		mon.Start(ctx)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		log.Printf("Starting server on :8080, monitoring %s", cfg.TargetURL)
		if err := r.Run(":8080"); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutdown signal received")
	cancel()

	// Wait for monitor to finish and save logs
	wg.Wait()
	log.Println("Shutdown complete")
}
