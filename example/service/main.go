// swagger:meta
//
package main

import (
	"github.com/morlay/gin-swagger/example/service/test"
	"github.com/morlay/gin-swagger/example/service/test2"
	"github.com/morlay/gin-swagger/example/service/test3"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/", test.Test)
	router.GET("/test", test3.Test3)

	userRoute := router.Group("/user")
	userRouteWith := userRoute.Group("/test")
	{
		userRouteWith.GET("/:name/:action", test2.Test2)
	}

	router.Run(":8080")
}
