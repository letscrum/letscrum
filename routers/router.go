package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/letscrum/letscrum/middlewares/jwt"

	"github.com/letscrum/letscrum/routers/apis/v1"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	apisV1 := r.Group("/apis/v1")

	apisV1.Use()
	{
		apisV1.GET("/signin", v1.SignIn)
	}

	apisV1.Use(jwtMiddleware.JWT())
	{
		apisV1.POST("/users", v1.CreateUser)
		apisV1.POST("/projects", v1.CreateProject)
		apisV1.GET("/projects", v1.ListProject)
		apisV1.PUT("/projects/:project_name", v1.UpdateProject)
		apisV1.DELETE("/projects/:project_name", v1.DeleteProject)
		apisV1.GET("/projects/:project_name", v1.GetProject)
		apisV1.GET("/projects/:project_name/members", v1.ListProjectMembers)
	}

	return r
}
