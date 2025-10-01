package backend

import (
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/deployer"
)

func StartServer() {
	// Initialize a default Gin router
	router := gin.Default()

	// Set a lower memory limit for multipart forms (defaults to 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Define our API endpoint.
	// When a POST request is made to /deploy, it will be handled by DeployHandler.
	router.POST("/deploy", deployer.DeployHandler)
	runRoute := router.Group("/run")

	runRoute.GET("/services", deployer.ListServicesHandler)
	runRoute.POST("/services", deployer.DeployHandler)

	log.Println("Starting server on port 8081...")
	// Start the HTTP server on port 8080.
	// The server will run indefinitely until the program is terminated.
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
	log.Println("Server stopped")
}
