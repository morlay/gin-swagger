package client_service

import (
	"github.com/morlay/gin-swagger/swagger_to_client/client"
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
