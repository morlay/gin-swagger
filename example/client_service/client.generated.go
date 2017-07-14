package client_service

import (
	"mime/multipart"

	"github.com/morlay/gin-swagger/example/test2"
	"github.com/morlay/gin-swagger/swagger_to_client/client"
)

type ClientService struct {
	client.Client
	Host string `default:"service-service" docker:"@external_link=$$/$$:$$"`
}

type GetUserRequest struct {
	//
	Id string `json:"id" in:"query"`
	//
	Age string `json:"age" in:"query"`
}

type GetUserResponse struct {
	//
	Body GetUser `json:"body"`
}

func (c ClientService) GetUser(req GetUserRequest) (resp GetUserResponse, err error) {
	err = c.DoRequest("ClientService.GetUser", "GET", "/auto", req, &resp)
	return
}
func (c *ClientService) InjectGetUser(resp GetUserResponse, err error) {
	c.Inject("ClientService.GetUser", resp, err)
}
func (c *ClientService) ResetGetUser() {
	c.Reset("ClientService.GetUser")
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
	File multipart.FileHeader `json:"file" in:"formData"`
}

type TestResponse struct {
	//
	Body SomeTest `json:"body"`
}

// @httpError(40000200,HTTP_ERROR_UNKNOWN,"未定义","",false);
// @httpError(400002002,HTTP_ERROR__TEST2,"Test2","Description",true);
// 正常返回
func (c ClientService) Test(req TestRequest) (resp TestResponse, err error) {
	err = c.DoRequest("ClientService.Test", "POST", "/", req, &resp)
	return
}
func (c *ClientService) InjectTest(resp TestResponse, err error) {
	c.Inject("ClientService.Test", resp, err)
}
func (c *ClientService) ResetTest() {
	c.Reset("ClientService.Test")
}

type Test2Request struct {
	//
	Authorization string `json:"authorization" in:"header"`
	// 分页大小
	Size int8 `json:"size" in:"query" default:"10" validate:"@int8[-1,20)"`
	// 分页偏移
	Offset int8 `json:"offset" in:"query" default:"0" validate:"@int8[-1,100]"`
	// ids
	Ids []int8 `json:"ids" in:"query"`
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

// @httpError(400002000,HTTP_ERROR_UNKNOWN,"未定义","",false);
// @httpError(400002001,HTTP_ERROR__TEST,"Summary","",true);
// 正常返回
func (c ClientService) Test2(req Test2Request) (resp Test2Response, err error) {
	err = c.DoRequest("ClientService.Test2", "GET", "/user/test/:name/0", req, &resp)
	return
}
func (c *ClientService) InjectTest2(resp Test2Response, err error) {
	c.Inject("ClientService.Test2", resp, err)
}
func (c *ClientService) ResetTest2() {
	c.Reset("ClientService.Test2")
}

type Test3Request struct {
	//
	Authorization string `json:"authorization" in:"header"`
}

type Test3Response struct {
	//
	Body Some `json:"body"`
}

// @httpError(400002000,HTTP_ERROR_UNKNOWN,"未定义","",false);
func (c ClientService) Test3(req Test3Request) (resp Test3Response, err error) {
	err = c.DoRequest("ClientService.Test3", "GET", "/test", req, &resp)
	return
}
func (c *ClientService) InjectTest3(resp Test3Response, err error) {
	c.Inject("ClientService.Test3", resp, err)
}
func (c *ClientService) ResetTest3() {
	c.Reset("ClientService.Test3")
}
