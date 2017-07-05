package test

import (
	"github.com/morlay/gin-swagger/example/test2"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(userRouter *gin.RouterGroup) {
	userRouteWith := userRouter.Group("/test")
	userRouteWith2 := userRouteWith.Group("")
	{
		userRouteWith2.GET("/:name/:action", test2.Test2)
	}
}
