package client_service

import (
	"github.com/morlay/gin-swagger/example/client"
	"github.com/morlay/gin-swagger/example/service/test2"
	"time"
)

func NewClientService(baseURL string, timeout time.Duration) *ClientService {
	return &ClientService{
		Client: client.Client{
			ID:      "ClientService",
			BaseURL: baseURL,
			Timeout: timeout,
		},
	}
}

type ClientService struct {
	client.Client
}

type TestRequest struct {
	// 分页大小
	Size int8 `json:"size" in:"query" default:"10" validate:"@int8[-1,20)"`
	// 分页偏移
	Offset int8 `json:"offset" in:"query" default:"0" validate:"@int8[-1,100]"`
	//
	StartTime test2.Date `json:"startTime,string" in:"query"`
	//
	State test2.State `json:"state,string" in:"query" validate:"@string{TWO}"`
	//
	Body ReqBody `json:"body" in:"body"`
}

type TestResponse struct {
	//
	Body SomeTest `json:"body"`
}

func (c ClientService) Test(req TestRequest) (resp TestResponse, err error) {
	err = c.DoRequest("Test", "POST", "/", req, &resp)
	return
}

type Test2Request struct {
	// 分页大小
	Size int8 `json:"size" in:"query" default:"10" validate:"@int8[-1,20)"`
	// 分页偏移
	Offset int8 `json:"offset" in:"query" default:"0" validate:"@int8[-1,100]"`
	// ids
	Ids string `json:"ids" in:"query"`
	//
	Id int8 `json:"id" in:"query"`
	//
	Name string `json:"name" in:"path"`
	//
	Is test2.Bool `json:"is" in:"path"`
	//
	State test2.State `json:"state,string" in:"query" validate:"@string{ONE}"`
	//
	StartTime test2.Date `json:"startTime,string" in:"query"`
}

type Test2Response struct {
	//
	Body Some `json:"body"`
}

func (c ClientService) Test2(req Test2Request) (resp Test2Response, err error) {
	err = c.DoRequest("Test2", "GET", "/user/test/:name/0", req, &resp)
	return
}

type Test3Response struct {
	//
	Body Some `json:"body"`
}

func (c ClientService) Test3() (resp Test3Response, err error) {
	err = c.DoRequest("Test3", "GET", "/test", nil, &resp)
	return
}
