package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/clinstid/schools_api/handlers"
)

// SetupRouter adds routes to a gin HTTP server
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// /schools routes
	r.GET("/schools", handlers.ListSchools)
	r.POST("/schools", handlers.AddSchool)

	// /schools/{id} routes
	r.GET("/schools/:schoolID", handlers.GetSchool)
	r.PUT("/schools/:schoolID", handlers.UpdateSchool)

	// Documentation routes
	r.Static("/docs/", "./dist/")

	return r
}
