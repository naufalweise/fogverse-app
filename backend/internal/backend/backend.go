package backend

import (
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/kubeutils"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/run"
)

func StartServer() {
	// Initialize a default Gin router
	router := gin.Default()

	// Set a lower memory limit for multipart forms (defaults to 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	kubeClient, err := kubeutils.GetKubeClient()
	if err != nil {
		log.Fatal(err)
	}

	runRoute := router.Group("/run")
	runController := run.NewController(kubeClient)
	runController.SetupRoutes(runRoute)

	log.Println("Starting server on port 8081...")
	// The server will run indefinitely until the program is terminated.
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
	log.Println("Server stopped")
}
