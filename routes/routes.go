package routes

import (
	"github.com/gin-gonic/gin"

	"gitlab.com/clinstid/schools_api/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/schools", handlers.ListSchools)
	r.POST("/schools", handlers.AddSchool)

	r.GET("/schools/:schoolID", handlers.GetSchool)
	r.PUT("/schools/:schoolID", handlers.UpdateSchool)

	r.Static("/docs/", "./dist/")

	return r
}
