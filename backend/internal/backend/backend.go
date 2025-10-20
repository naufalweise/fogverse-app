package backend

import (
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/kubeutils"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/run"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/run/job"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/run/keda"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/run/service"
)

func StartServer() {
	// Initialize a default Gin router
	router := gin.Default()

	// Set a lower memory limit for multipart forms (defaults to 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	kubeConfig, err := kubeutils.GetKubeConfig()
	if err != nil {
		log.Fatal(err)
	}

	kubeClient, err := kubeutils.GetKubeClient(kubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	kedaClient, err := keda.GetKedaClient(kubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	kedaSvc := keda.NewSvc(kubeClient, kedaClient)
	serviceSvc := service.NewSvc(kubeClient, kedaClient)
	jobSvc := job.NewSvc(kubeClient)

	runRoute := router.Group("/run")
	runController := run.NewController(serviceSvc, jobSvc, kedaSvc)
	runController.SetupRoutes(runRoute)

	log.Println("Starting server on port 8081...")
	// The server will run indefinitely until the program is terminated.
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
	log.Println("Server stopped")
}
