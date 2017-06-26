package globals

import (
	"fmt"
	"github.com/morlay/gin-swagger/http_error_code/httplib"
	"strconv"
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
		Code:           c.Code(),
		Msg:            c.Msg(),
		Desc:           c.Desc(),
		CanBeErrorTalk: c.CanBeErrTalk(),
	}
}

func (c HttpErrorCode) ToResp() (int, *httplib.GeneralError) {
	return c.Status(), c.ToError()
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
