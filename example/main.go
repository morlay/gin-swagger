// swagger:meta
//
package main

import (
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/morlay/gin-swagger/example/test"
	"github.com/morlay/gin-swagger/example/test2"
	"github.com/morlay/gin-swagger/example/test3"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/", test.Test)
	router.GET("/test", test3.Test3)

	userRoute := router.Group("/user", test.AuthMiddleware)
	userRouteWith := userRoute.Group("/test")
	{
		userRouteWith.GET("/:name/:action", test2.Test2)
	}

	router.Run(":8080")
}
