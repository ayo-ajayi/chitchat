package gpt

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ayo-ajayi/chitchat/config"
	"github.com/ayo-ajayi/chitchat/redis"
	"github.com/ayo-ajayi/chitchat/types/chat"
	"github.com/ayo-ajayi/chitchat/types/models"
)



type DB interface {
	GetList(key string) ([]string, error)
	Exists(key string) int64
	SetList(key string, ls []string) error
	Close() error
}
type DBClient struct {
	db DB
}

func NewDBClient(db DB) *DBClient {
	return &DBClient{db: db}
}

type GPT struct {
	httpclient *Client
	dbclient   DBClient
}

func NewGPT() (*GPT, error) {
	c, err := NewClient()
	if err != nil {
		return nil, err
	}
	return &GPT{
		httpclient: c,
		dbclient:   *NewDBClient(redis.DefaultClient()),
	}, nil
}


func (g *GPT) exec(message string, model *chat.ChatGPTModel) (*http.Response, error) {
	x := chat.ChatGPTModel("")
	if model != nil {
		x = *model
	} else {
		x = g.DefaultModel()
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
	return g.httpclient.do("POST", endpoint, reqStruct)
}

func (g *GPT) Chat(input string, model *string) string {
	x := chat.ChatGPTModel(*model)
	res, err := g.exec(input, &x)
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

func (g *GPT) LS() ([]string, error) {
	res, err := g.httpclient.listAvailableModels()
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

//reduces call to models API and also keeps the models list up-to-date
func (g *GPT) listModels(modelskey string) ([]string, error) {
	res, err := g.dbclient.db.GetList(modelskey)
	if err != nil {
		return nil, err
	}
	var ModelsList []string
	exists := g.dbclient.db.Exists(modelskey)
	if res == nil || exists == 0 {
		ModelsList, err = g.LS()
		if err != nil {
			return nil, err
		}
		err = g.dbclient.db.SetList(modelskey, ModelsList)
		if err != nil {
			return nil, err
		}
		return ModelsList, nil
	}
	return res, nil
}

func (g *GPT) ListOfModels() []string {
	res, err := g.listModels(config.MODELS_KEY)
	if err != nil {
		ls, _ := g.LS()
		return ls
	}
	return res
}

func (g *GPT) defaultModelFunc() chat.ChatGPTModel {
	models := g.ListOfModels()
	lastindex := len(models) - 1
	return chat.ChatGPTModel(models[lastindex])
}

func (g *GPT) DefaultModel() chat.ChatGPTModel {
	return g.defaultModelFunc()
}
