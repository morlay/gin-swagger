package client_service

import (
	"github.com/morlay/gin-swagger/example/service/test2"
)

type Common struct {
	// 总数
	Total int8 `json:"total"`
}

type ErrorMap map[string]map[string]int64

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

type ReqBody struct {
	//
	Name string `json:"name"`
	//
	UserName string `json:"username"`
}

type Some struct {
	// Test
	State test2.State `json:"state" validate:"@string{,TWO}"`
	//
	Data []ItemData `json:"data"`
	//
	Name uint64 `json:"name,string"`
	//
	StartTime test2.Date `json:"startTime,string"`
}

type SomeTest struct {
	Common
	//
	State test2.State `json:"state" validate:"@string{TWO}"`
	//
	ErrorMap ErrorMap `json:"errorMap"`
}
