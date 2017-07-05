package client

import (
	"time"
)

type Client struct {
	BaseURL string
	Timeout time.Duration
}

func (c Client) DoRequest(operationID string, method, uri string, req interface{}, res interface{}) error {
	// Write your own HTTP call
	return nil
}

func (c Client) Inject(operationID string, res interface{}, err error) {
	// write yor own mock method
}

func (c Client) Reset(operationID string) {
	// write yor own mock clear
}
