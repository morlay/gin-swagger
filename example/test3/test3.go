package test3

import (
	"github.com/morlay/gin-swagger/example/test2"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

// Summary
//
// Others
// heheheh
func Test3(c *gin.Context) {
	c.JSON(http.StatusOK, test2.Some{
		State: test2.STATE__ONE,
	})
}
