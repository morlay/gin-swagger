package test2

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"time"
)

// swagger:enum State
type (
	State int
)

const (
	STATE_UNKNOWN = iota
	STATE__ONE    // one
	STATE__TWO    // two
)

// swagger:strfmt date-time
type Date time.Time

// Common
type Common struct {
	// 总数
	Total int8 `json:"total"`
}

// 分页类
type Pager struct {
	// 分页大小
	Size int8 `json:"size" in:"query" default:"10" validate:"@int8[-1,20)"`
	// 分页偏移
	Offset int8 `json:"offset" in:"query" default:"0" validate:"@int8[-1,100]"`
}

type SomeReq struct {
	Pager
	// ids
	Ids       []int8 `json:"ids" in:"query" validate:"@int8[-1,100]"`
	Id        int8   `json:"id" in:"query"`
	Name      string `json:"name" in:"path"`
	State     State  `json:"state" in:"query" validate:"@string{ONE}"`
	StartTime Date   `json:"startTime" in:"query"`
}

type ItemData struct {
	Name      string `json:"name" validate:"@string[0,)"`
	Id        string `json:"id"`
	State     State  `json:"state"`
	StartTime Date   `json:"startTime"`
}

// Some
type (
	// test
	Some struct {
		State     State      `json:"state" validate:"@string{TWO}"`
		Name      uint64     `json:"name,string"`
		Data      []ItemData `json:"data"`
		StartTime Date       `json:"startTime"`
	}
)

// Summary
//
// Others
// heheheh
func Test2(c *gin.Context) {
	var req = SomeReq{}

	var res = Some{}

	res.Name = uint64(req.Size)

	c.JSON(http.StatusOK, res) // 正常返回
}
