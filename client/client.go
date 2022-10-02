package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Config struct {
	// BaseURL is the base url for the example API
	// This field is mandatory.
	BaseURL string
}

// Client is the client to interact with the example HTTP API.
type Client struct {
	resty *resty.Client
}

// New creates a new client to interact with the example HTTP API.
func New(config Config) (*Client, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base url is required")
	}

	c := resty.New().SetHostURL(config.BaseURL)

	c.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		switch resp.StatusCode() {
		case http.StatusNotFound:
			return ErrNotFound
		}

		if resp.IsError() {
			if resp.Header().Get("Content-Type") == "application/json" {
				var msgErr ErrorMessage
				err := json.Unmarshal(resp.Body(), &msgErr)
				if err != nil {
					return err
				}
				return &msgErr
			}
			return fmt.Errorf("status code %d", resp.StatusCode())
		}

		return nil
	})

	return &Client{resty: c}, nil
}

// Ping check the availability of the API.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.resty.R().SetContext(ctx).Get("/api/v1/ping")
	return err
}

type MsgResponse struct {
	Message string `json:"message"`
}

// Hello say hello to the world
func (c *Client) Hello(ctx context.Context) (*MsgResponse, error) {
	var resp MsgResponse

	_, err := c.resty.R().SetContext(ctx).SetResult(&resp).Get("/api/v1/hello")
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type IssueResponse struct {
	Project string `json:"project" bson:"project,omitempty"`
	// Key unique key format : VTM-84
	Key      string `json:"key" bson:"key,omitempty"`
	Type     string `json:"type" bson:"type,omitempty"`
	Desc     string `json:"desc" bson:"desc,omitempty"`
	Assigned string `json:"assigned" bson:"assigned,omitempty"`
	//Date format : yyyy-n°week = 2021-18
	Date string `json:"date" bson:"date,omitempty"`
}
type IssuesResponse struct {
	Issues []IssueResponse `json:"issues" bson:"issues,omitempty"`
}
