package credstor

import (
	"time"
)

const defaultTimeout = 10 * time.Second

type Client struct{}

func NewClient() (*Client, error) {
	return &Client{}, nil
}
