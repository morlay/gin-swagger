package globals

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/morlay/gin-swagger/http_error_code/httplib"
)

const (
	ServiceB = 2
)

//go:generate gin-swagger error
type HttpErrorCode int32

const (
	// 未定义
	HTTP_ERROR_UNKNOWN HttpErrorCode = iota + ServiceB*1e3 + http.StatusBadRequest*1e6
	// @errTalk Summary
	HTTP_ERROR__TEST
	// @errTalk Test2
	// Description
	HTTP_ERROR__TEST2
)

func getError() *httplib.GeneralError {
	return HTTP_ERROR__TEST2.ToError()
}

func WriteErr(c *gin.Context) {
	c.JSON(http.StatusBadRequest, getError())
}
