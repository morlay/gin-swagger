// swagger:meta
//
package main

import (
	"github.com/morlay/gin-swagger/example/test"
	"github.com/morlay/gin-swagger/example/test2"
	"github.com/morlay/gin-swagger/example/test3"
	"gopkg.in/gin-gonic/gin.v1"
	"fmt"
	"github.com/morlay/gin-swagger/example/globals"
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

	fmt.Println(globals.HTTP_ERROR__TEST.Status())

	router.Run(":8080")
}
