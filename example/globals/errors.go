package globals

import (
	"net/http"
)

const (
	ServiceB = 2
)

type HttpErrorCode int32

const (
	// 未定义
	HTTP_ERROR_UNKNOWN HttpErrorCode = iota + ServiceB * 1e3 + http.StatusBadRequest * 1e6
	// @errTalk Summary
	HTTP_ERROR__TEST
	// @errTalk Test2
	// Description
	HTTP_ERROR__TEST2
)

