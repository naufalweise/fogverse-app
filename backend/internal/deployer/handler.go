// handlers/deploy.go
// This package handles the business logic for run container API endpoints.

package deployer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/deployer/service"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/kubeutils"
)

// DeployHandler processes the API request to create a new deployment or job.
func DeployHandler(c *gin.Context) {
	var request DeployRequest

	// Bind the incoming JSON payload to the DeployRequest struct.
	// If there's an error (e.g., missing fields), return a 400 Bad Request.
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
		return
	}

	// Get a Kubernetes clientset. This function is smart and will work
	// whether the app is running inside the cluster or locally.
	clientset, err := kubeutils.GetKubeClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to connect to Kubernetes cluster: " + err.Error(), "status": "error"})
		return
	}

	d := Deployer{KubeInterface: clientset, DeployRequest: request}
	creationErr := d.Deploy()

	if creationErr != nil {
		// If the resource already exists or another error occurs, return an error message.
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create resource: " + creationErr.Error(), "status": "error"})
		return
	}

	// If everything is successful, return a 201 Created status and a success message.
	c.JSON(http.StatusCreated, gin.H{"message": "Resource '" + request.Name + "' of type '" + request.Type + "' created successfully!", "status": "success"})
}

func ListServicesHandler(c *gin.Context) {

	clientset, err := kubeutils.GetKubeClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to connect to Kubernetes cluster: " + err.Error(), "status": "error"})
		return
	}

	services, err := service.List(clientset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get services: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": services})
}
