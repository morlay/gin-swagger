package globals

import (
	"github.com/morlay/gin-swagger/http_error_code/httplib"
)

func init() {
	httplib.RegisterError("HTTP_ERROR_UNKNOWN", HTTP_ERROR_UNKNOWN, "未定义", "", false)
	httplib.RegisterError("HTTP_ERROR__TEST", HTTP_ERROR__TEST, "Summary", "", true)
	httplib.RegisterError("HTTP_ERROR__TEST2", HTTP_ERROR__TEST2, "Test2", "Description", true)
}
