//go:generate gin-swagger error
package globals

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morlay/gin-swagger/http_error_code/httplib"
)

const (
	ServiceB = 2
)

const (
	// 未定义
	HTTP_ERROR_UNKNOWN httplib.HttpErrorCode = iota + ServiceB*1e3 + http.StatusBadRequest*1e6
	// @errTalk Summary
	HTTP_ERROR__TEST
	_
	_
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
