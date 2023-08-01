package gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	rdb "github.com/ayo-ajayi/chitchat/redis"
	"github.com/ayo-ajayi/chitchat/types/chat"
	"github.com/ayo-ajayi/chitchat/types/models"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
)

const apiURL string = "https://api.openai.com/v1"

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
			BaseURL: apiURL,
			APIKey:  APIKey,
		},
	}, nil
}

//reduces call to models API and also keeps the models list up-to-date
func listModels() ([]string, error) {
	c := rdb.DefaultClient()
	defer c.C.Close()
	res, err := c.GetList(rdb.MODELS_KEY)
	if err != nil {
		return nil, err
	}
	var ModelsList []string
	exists := c.Exists(rdb.MODELS_KEY)
	if res == nil || exists == 0 {
		ModelsList, err = LS()
		if err != nil {
			return nil, err
		}
		err = c.SetList(rdb.MODELS_KEY, ModelsList)
		if err != nil {
			return nil, err
		}
		return ModelsList, nil
	}
	return res, nil
}

var ListOfModels []string = (func() []string {
	res, err := listModels()
	if err != nil {
		ls, _ := LS()
		return ls
	}
	return res
})()
var defaultModelFunc = func() chat.ChatGPTModel {
	models := ListOfModels
	lastindex := len(models) - 1
	return chat.ChatGPTModel(models[lastindex])
}

func DefaultModel() chat.ChatGPTModel {
	return defaultModelFunc()
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
func (c *Client) exec(message string, model *chat.ChatGPTModel) (*http.Response, error) {
	x := chat.ChatGPTModel("")
	if model != nil {
		x = *model
	} else {
		x = DefaultModel()
	}
	reqStruct := &chat.RequestStruct{
		Model: x,
		Messages: []chat.ChatMessage{
			{
				Role:    chat.ChatGPTModelRoleUser,
				Content: message,
			},
		},
	}
	endpoint := "/chat/completions"
	return c.do("POST", endpoint, reqStruct)
}
func Chat(input string, model *string) string {
	c, err := NewClient()
	if err != nil {
		panic(err)
	}
	x := chat.ChatGPTModel(*model)
	res, err := c.exec(input, &x)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	var resStruct chat.ResponseStruct
	if err = json.NewDecoder(res.Body).Decode(&resStruct); err != nil {
		panic(err)
	}
	return resStruct.Choices[0].Message.Content
}
func (c *Client) listAvailableModels() (*http.Response, error) {
	endpoint := "/models"
	return c.do("GET", endpoint, nil)
}
func LS() ([]string, error) {
	c, err := NewClient()
	if err != nil {
		return nil, err
	}
	res, err := c.listAvailableModels()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var resStruct models.AvailableGPTModels
	if err = json.NewDecoder(res.Body).Decode(&resStruct); err != nil {
		return nil, err
	}
	var arr []string
	for _, model := range resStruct.Data {
		prefix := "gpt-"
		if strings.HasPrefix(model.Root, prefix) {
			arr = append(arr, model.Root)
		}
	}
	return arr, nil
}
