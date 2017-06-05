package from_request

import (
	"gopkg.in/gin-gonic/gin.v1"
	"reflect"
)

type Request interface {
	Handle(c *gin.Context)
}

func TypeOfRequest(req Request) reflect.Type {
	reqInterface := reflect.Indirect(reflect.ValueOf(req)).Interface()
	return reflect.Indirect(reflect.ValueOf(reqInterface)).Type()
}

func FromRequest(req Request) func(c *gin.Context) {
	return func(c *gin.Context) {
		reqType := TypeOfRequest(req)
		finalReq := reflect.New(reqType)

		reflectValue := reflect.Indirect(finalReq)
		reflectType := reflectValue.Type()

		for i := 0; i < reflectValue.NumField(); i++ {
			field := reflectType.Field(i)
			fieldValue := reflectValue.Field(i)
			location := field.Tag.Get("in")
			name := field.Tag.Get("json")

			if !fieldValue.CanSet() || !fieldValue.IsValid() {
				panic("un setable")
			}

			switch location {
			case "query":
				ConvertFromStr(c.Query(name), fieldValue)
			}
		}

		reflectValue.Interface().(Request).Handle(c)
	}
}
