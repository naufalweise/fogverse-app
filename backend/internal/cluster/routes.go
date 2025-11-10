package cluster

import "github.com/gin-gonic/gin"

func (controller *Controller) SetupRoutes(r gin.IRoutes) {
	r.GET("/topics", controller.ListTopics)
	r.POST("/topics", controller.CreateTopic)
	r.PUT("/topics", controller.UpdateTopic)
	r.GET("/topics/:topicName", controller.GetTopic)
	r.DELETE("/topics/:topicName", controller.DeleteTopic)
}
