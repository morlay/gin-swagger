package from_request

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// Get User
type GetUser struct {
	Id  string `json:"id" in:"query"`
	Age string `json:"age" in:"query"`
}

func (req GetUser) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, req)
}
