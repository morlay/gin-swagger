package from_request

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get User
type GetUser struct {
	Id  string `json:"id" in:"query"`
	Age string `json:"age" in:"query"`
}

func (req GetUser) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, req)
}
