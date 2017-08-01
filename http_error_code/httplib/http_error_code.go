package httplib

import (
	"fmt"
	"strconv"
)

var HttpErrorCodes = map[HttpErrorCode]*GeneralError{}

type HttpErrorCode int32

func RegisterError(key string, code HttpErrorCode, msg string, desc string, canBeErrTalk bool) {
	HttpErrorCodes[code] = &GeneralError{
		Key:            key,
		Code:           int32(code),
		Msg:            msg,
		Desc:           desc,
		CanBeErrorTalk: canBeErrTalk,
	}
}

func (c HttpErrorCode) ToError() *GeneralError {
	generalErr, ok := HttpErrorCodes[c]

	if ok {
		return generalErr
	}

	return nil
}

func (c HttpErrorCode) Status() int {
	status, _ := strconv.Atoi(fmt.Sprintln(c)[:3])
	return status
}
