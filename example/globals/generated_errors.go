package globals

import (
	"fmt"
	"strconv"
	"github.com/morlay/gin-swagger/http_error_code/httplib"
)

func (httpErrorCode HttpErrorCode) Code() int32 {
	return int32(httpErrorCode)
}

func (httpErrorCode HttpErrorCode) Status() int {
	status, _ := strconv.Atoi(fmt.Sprintln(httpErrorCode)[:3])
	return status
}

func (httpErrorCode HttpErrorCode) ToError() *httplib.GeneralError {
	return &httplib.GeneralError{
		Code:           httpErrorCode.Code(),
		Msg:            httpErrorCode.Msg(),
		Desc:           httpErrorCode.Desc(),
		CanBeErrorTalk: httpErrorCode.CanBeErrTalk(),
	}
}

func (httpErrorCode HttpErrorCode) ToResp() (int, *httplib.GeneralError) {
	return httpErrorCode.Status(), httpErrorCode.ToError()
}

func (httpErrorCode HttpErrorCode) Msg() string {
	switch httpErrorCode {
	case HTTP_ERROR_UNKNOWN:
		return "未定义"
	case HTTP_ERROR__TEST2:
		return "Test2"
	case HTTP_ERROR__TEST:
		return "Summary"
	}
	return ""
}

func (httpErrorCode HttpErrorCode) Desc() string {
	switch httpErrorCode {
	case HTTP_ERROR_UNKNOWN:
		return ""
	case HTTP_ERROR__TEST2:
		return "Description"
	case HTTP_ERROR__TEST:
		return ""
	}
	return ""
}

func (httpErrorCode HttpErrorCode) CanBeErrTalk() bool {
	switch httpErrorCode {
	case HTTP_ERROR_UNKNOWN:
		return false

	case HTTP_ERROR__TEST2:
		return true

	case HTTP_ERROR__TEST:
		return true

	}
	return false
}
