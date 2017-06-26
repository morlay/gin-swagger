package test

import (
	"github.com/morlay/gin-swagger/example/test2"
	"gopkg.in/gin-gonic/gin.v1"
)

func SetupUserRoutes(userRouter *gin.RouterGroup) {
	userRouteWith := userRouter.Group("/test")
	userRouteWith2 := userRouteWith.Group("")
	{
		userRouteWith2.GET("/:name/:action", test2.Test2)
	}
}
