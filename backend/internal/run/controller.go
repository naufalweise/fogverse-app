// handlers/deploy.go
// This package handles the business logic for run container API endpoints.

package run

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/run/job"
	"gitlab.cs.ui.ac.id/fogverse/app/internal/run/service"
	"k8s.io/client-go/kubernetes"
)

type Controller struct {
	kubeClient kubernetes.Interface
	serviceSvc service.Svc
	jobSvc     job.Svc
}

func NewController(kubeClient kubernetes.Interface) *Controller {
	return &Controller{kubeClient: kubeClient, serviceSvc: service.NewSvc(kubeClient), jobSvc: job.NewSvc(kubeClient)}
}

func (controller *Controller) CreateService(c *gin.Context) {
	var request service.Service

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
		return
	}
	err := controller.serviceSvc.Create(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create resource: " + err.Error(), "status": "error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Service '" + request.Name + "' created successfully!", "status": "success"})
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

func (controller *Controller) ListJob(c *gin.Context) {
	ctx := c.Request.Context()
	jobs, err := controller.jobSvc.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get jobs: " + err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": jobs})
}

func (controller *Controller) CreateJob(c *gin.Context) {
	var request job.Job
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "error"})
		return
	}

	ctx := c.Request.Context()
	err := controller.jobSvc.Create(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create job: " + err.Error(), "status": "error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Job '" + request.Name + "' created successfully!", "status": "success"})
}
