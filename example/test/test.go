package test

import (
	"fmt"
	"github.com/morlay/gin-swagger/example/test2"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

// ErrorMap
type ErrorMap map[string]map[string]*int

// SomeTest
type SomeTest struct {
	test2.Common
	ErrorMap ErrorMap `json:"errorMap"`
}

// ReqBody
type ReqBody struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
}

// SomeReq
type SomeReq struct {
	// Body
	test2.Pager
	StartTime test2.Date `in:"query" json:"startTime"`
	Body      ReqBody
}

func Test(c *gin.Context) {
	req := SomeReq{}

	fmt.Println(req)

	var res = SomeTest{}

	// 正常返回
	c.JSON(http.StatusOK, res)
}
