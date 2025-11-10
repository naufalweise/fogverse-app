package cluster

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	svc Svc
}

func NewController(clusterSvc Svc) *Controller {
	return &Controller{svc: clusterSvc}
}

func (controller *Controller) CreateTopic(c *gin.Context) {
	var topic Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	err := controller.svc.CreateTopic(c, topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "topic created"})

}

func (controller *Controller) DeleteTopic(c *gin.Context) {
	topicName := c.Param("topicName")
	if topicName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "topic name is required"})
		return
	}
	err := controller.svc.DeleteTopic(c, topicName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "topic deleted"})
}

func (controller *Controller) UpdateTopic(c *gin.Context) {
	var topic Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	err := controller.svc.UpdateTopic(c, topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "topic updated"})
}

func (controller *Controller) GetTopic(c *gin.Context) {
	topicName := c.Param("topicName")
	if topicName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "topic name is required"})
		return
	}
	topic, err := controller.svc.GetTopic(c, topicName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": topic})
}

func (controller *Controller) ListTopics(c *gin.Context) {
	topics, err := controller.svc.ListTopics(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": topics})
}
