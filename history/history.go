package history

import (
	"github.com/ayo-ajayi/chitchat/util"
	"github.com/ayo-ajayi/chitchat/types"
)
type DBClient interface {
	AddStream(key string, value string) error
	ReadStream(limit int64) ([]types.Message, error)
}

type History struct {
	dbclient DBClient
}

func NewHistory(C DBClient) *History {
	return &History{
		dbclient: C,
	}
}
func (h *History)SaveChat(input, response string) error {
	if err := h.dbclient.AddStream(input, response); err != nil {
		return err
	}
	return nil
}
func (h *History)GetChat(limit int64) ([]types.Message, error) {
	res, err := h.dbclient.ReadStream(limit)
	if err != nil {
		return nil, err
	}
	return res, err
}
func (h *History)GetDate(id string) (string, error) {
	date, err := util.GetDateFromStreamID(id)
	if err != nil {
		return "", err
	}
	return date, nil
}
