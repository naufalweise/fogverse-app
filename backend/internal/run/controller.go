// handlers/deploy.go
// This package handles the business logic for run container API endpoints.

package run

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/run/service"
	"k8s.io/client-go/kubernetes"
)

type Controller struct {
	kubeClient kubernetes.Interface
	serviceSvc service.Svc
}

func NewController(kubeClient kubernetes.Interface) *Controller {
	return &Controller{kubeClient: kubeClient, serviceSvc: service.NewSvc(kubeClient)}
}

// DeployHandler processes the API request to create a new deployment or job.
func (controller *Controller) DeployHandler(c *gin.Context) {
	var request DeployRequest

	// Bind the incoming JSON payload to the DeployRequest struct.
	// If there's an error (e.g., missing fields), return a 400 Bad Request.
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
		return
	}
	d := Deployer{KubeInterface: controller.kubeClient, DeployRequest: request}
	creationErr := d.Deploy()

	if creationErr != nil {
		// If the resource already exists or another error occurs, return an error message.
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create resource: " + creationErr.Error(), "status": "error"})
		return
	}

	// If everything is successful, return a 201 Created status and a success message.
	c.JSON(http.StatusCreated, gin.H{"message": "Resource '" + request.Name + "' of type '" + request.Type + "' created successfully!", "status": "success"})
}

func (controller *Controller) ListServices(c *gin.Context) {
	ctx := c.Request.Context()
	services, err := controller.serviceSvc.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get services: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": services})
}
