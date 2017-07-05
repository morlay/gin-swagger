package test3

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/morlay/gin-swagger/example/test2"
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
