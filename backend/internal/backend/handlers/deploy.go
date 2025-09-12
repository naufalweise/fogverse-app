// handlers/deploy.go
// This package handles the business logic for run container API endpoints.

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/deployer"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/kubeutils"
)

// DeployHandler processes the API request to create a new deployment or job.
func DeployHandler(c *gin.Context) {
	var request deployer.DeployRequest

	// Bind the incoming JSON payload to the DeployRequest struct.
	// If there's an error (e.g., missing fields), return a 400 Bad Request.
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get a Kubernetes clientset. This function is smart and will work
	// whether the app is running inside the cluster or locally.
	clientset, err := kubeutils.GetKubeClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Kubernetes cluster: " + err.Error()})
		return
	}

	d := deployer.Deployer{KubeInterface: clientset, DeployRequest: request}
	creationErr := d.Deploy()

	if creationErr != nil {
		// If the resource already exists or another error occurs, return an error message.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resource: " + creationErr.Error()})
		return
	}

	// If everything is successful, return a 201 Created status and a success message.
	c.JSON(http.StatusCreated, gin.H{"message": "Resource '" + request.Name + "' of type '" + request.Type + "' created successfully!"})
}
