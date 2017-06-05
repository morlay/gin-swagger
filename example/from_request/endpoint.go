package from_request

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

// Get User
type GetUser struct {
	Id  string `json:"id" in:"query"`
	Age string `json:"age" in:"query"`
}

func (req GetUser) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, req)
}
