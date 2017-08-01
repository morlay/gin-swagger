package globals

import (
	"github.com/morlay/gin-swagger/http_error_code/httplib"
)

func init() {
	httplib.RegisterError("HTTP_ERROR_UNKNOWN", 400002000, "未定义", "", false)
	httplib.RegisterError("HTTP_ERROR__TEST", 400002001, "Summary", "", true)
	httplib.RegisterError("HTTP_ERROR__TEST2", 400002004, "Test2", "Description", true)
}
