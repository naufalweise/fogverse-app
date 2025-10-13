package run

import "github.com/gin-gonic/gin"

func (controller *Controller) SetupRoutes(r gin.IRoutes) {
	r.GET("/jobs", controller.ListJob)
	r.POST("/jobs", controller.CreateJob)

	r.GET("/services", controller.ListServices)
	r.POST("/services", controller.CreateService)
}
