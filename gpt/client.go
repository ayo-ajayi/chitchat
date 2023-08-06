package gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ayo-ajayi/chitchat/config"
	"github.com/spf13/viper"
)


type Config struct {
	BaseURL string
	APIKey  string
	OrgID   string
}

type Client struct {
	client *http.Client
	config *Config
}

func NewClient() (*Client, error) {
	APIKey := viper.GetViper().GetString("openapi_key")
	if APIKey == "" {
		return nil, errors.New("api key not provided")
	}
	return &Client{
		client: &http.Client{},
		config: &Config{
			BaseURL: config.APIURL,
			APIKey:  APIKey,
		},
	}, nil
}

func (c *Client) do(method, endpoint string, reqBody any) (*http.Response, error) {
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	var x io.Reader
	if reqBody != nil {
		x = bytes.NewBuffer(reqBytes)
	} else {
		x = nil
	}
	r, err := http.NewRequest(method, c.config.BaseURL+endpoint, x)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	w, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (c *Client) listAvailableModels() (*http.Response, error) {
	endpoint := "/models"
	return c.do("GET", endpoint, nil)
}