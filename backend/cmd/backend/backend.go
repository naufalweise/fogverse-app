// main.go
// This is the main entry point for our application.
// It sets up the Gin web server and defines the API endpoint.

package main

import (
	"log"

	"gitlab.cs.ui.ac.id/fogverse/app/internal/backend/handlers"

	"github.com/gin-gonic/gin" // Import the Gin framework
)

func main() {
	// Initialize a default Gin router
	router := gin.Default()

	// Set a lower memory limit for multipart forms (defaults to 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Define our API endpoint.
	// When a POST request is made to /deploy, it will be handled by DeployHandler.
	router.POST("/deploy", handlers.DeployHandler)

	log.Println("Starting server on port 8081...")
	// Start the HTTP server on port 8080.
	// The server will run indefinitely until the program is terminated.
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
	log.Println("Server stopped")
}
