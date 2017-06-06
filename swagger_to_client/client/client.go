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
