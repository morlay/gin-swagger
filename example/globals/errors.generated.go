package globals

import (
	"fmt"
	"strconv"

	"github.com/morlay/gin-swagger/http_error_code/httplib"
)

func (c HttpErrorCode) Code() int32 {
	return int32(c)
}

func (c HttpErrorCode) Status() int {
	status, _ := strconv.Atoi(fmt.Sprintln(c)[:3])
	return status
}

func (c HttpErrorCode) ToError() *httplib.GeneralError {
	return &httplib.GeneralError{
		Key:            c.Key(),
		Code:           c.Code(),
		Msg:            c.Msg(),
		Desc:           c.Desc(),
		CanBeErrorTalk: c.CanBeErrTalk(),
	}
}

func (c HttpErrorCode) Key() string {
	switch c {
	case HTTP_ERROR_UNKNOWN:
		return "HTTP_ERROR_UNKNOWN"
	case HTTP_ERROR__TEST:
		return "HTTP_ERROR__TEST"
	case HTTP_ERROR__TEST2:
		return "HTTP_ERROR__TEST2"
	}
	return ""
}

func (c HttpErrorCode) Msg() string {
	switch c {
	case HTTP_ERROR_UNKNOWN:
		return "未定义"
	case HTTP_ERROR__TEST:
		return "Summary"
	case HTTP_ERROR__TEST2:
		return "Test2"
	}
	return ""
}

func (c HttpErrorCode) Desc() string {
	switch c {
	case HTTP_ERROR_UNKNOWN:
		return ""
	case HTTP_ERROR__TEST:
		return ""
	case HTTP_ERROR__TEST2:
		return "Description"
	}
	return ""
}

func (c HttpErrorCode) CanBeErrTalk() bool {
	switch c {
	case HTTP_ERROR_UNKNOWN:
		return false
	case HTTP_ERROR__TEST:
		return true
	case HTTP_ERROR__TEST2:
		return true
	}
	return false
}
