package test

import (
	"fmt"
	"github.com/morlay/gin-swagger/example/service/test2"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

// ErrorMap
type ErrorMap map[string]map[string]*int

// SomeTest
type SomeTest struct {
	test2.Common
	State    test2.State `json:"state" validate:"@string{TWO}"`
	ErrorMap ErrorMap    `json:"errorMap"`
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
	StartTime test2.Date  `in:"query" json:"startTime"`
	State     test2.State `in:"query" json:"state" validate:"@string{TWO}"`
	Body      ReqBody
}

func Test(c *gin.Context) {
	req := SomeReq{}

	fmt.Println(req)

	var res = SomeTest{
		State: test2.STATE__ONE,
	}

	// 正常返回
	c.JSON(http.StatusOK, res)
}
