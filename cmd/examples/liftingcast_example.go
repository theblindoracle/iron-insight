package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"iron-insight/internal/liftingcast"

	"github.com/joho/godotenv"
)

func main() {
	// Create WebSocket client

	godotenv.Load()

	apiKey := os.Getenv("LIFTINGCAST_APIKEY")
	client := liftingcast.New("wss://backup.liftingcast.com/websocket", "mbho66s0ddh9", "test", apiKey)

	// Start the client
	if err := client.Start(); err != nil {
		log.Fatalf("Failed to start WebSocket client: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages and errors
	go func() {
		for {
			select {
			case message := <-client.Messages():
				log.Printf("Received message: %s", string(message))

			case err := <-client.Errors():
				log.Printf("WebSocket error: %v", err)
			}
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down...")

	// Stop the client
	if err := client.Stop(); err != nil {
		log.Printf("Error stopping client: %v", err)
	}

	log.Println("Shutdown complete")
}

