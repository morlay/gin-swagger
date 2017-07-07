package client_service

import (
	"github.com/morlay/gin-swagger/example/test2"
)

type Common struct {
	// 总数
	Total int8 `json:"total"`
}

type ErrorFieldModel struct {
	// 输入中出错的字段信息,这个信息为一个json字符串,方便客户端进行定位错误原因
	// 例如输入中{"name":{"alias" : "test"}}中的alias出错,则返回"name.alias"
	// 如果 alias 是数组, 且第2个元素的a字段错误,则返回"name.alias[2].a"
	Field string `json:"field"`
	// 错误字段位置, body, query, header, path, formData
	In string `json:"in"`
	// 错误信息
	Msg string `json:"msg"`
}

type ErrorMap map[string]map[string]int64

type GeneralError struct {
	// 是否能作为错误话术
	CanBeErrorTalk bool `json:"canBeTalkError"`
	// 详细描述
	Code int32 `json:"code"`
	// 错误代码
	Desc string `json:"desc"`
	// 出错字段
	ErrorFields []ErrorFieldModel `json:"errorFields"`
	// Request Id
	Id string `json:"id"`
	// 详细描述
	Key string `json:"key"`
	// 错误信息
	Msg string `json:"msg"`
	// 错误溯源
	Source []string `json:"source"`
}

type GetUser struct {
	//
	Age string `json:"age"`
	//
	Id string `json:"id"`
}

type ItemData struct {
	//
	Id string `json:"id"`
	//
	Name string `json:"name" validate:"@string[0,)"`
	//
	StartTime test2.Date `json:"startTime,string"`
	//
	State test2.State `json:"state"`
}

type Some struct {
	//
	Data []ItemData `json:"data"`
	//
	Name uint64 `json:"name,string"`
	//
	StartTime test2.Date `json:"startTime,string"`
	// Test
	State test2.State `json:"state" validate:"@string{,TWO}"`
}

type SomeTest struct {
	Common
	//
	ErrorMap ErrorMap `json:"errorMap"`
	//
	State test2.State `json:"state" validate:"@string{TWO}"`
}
