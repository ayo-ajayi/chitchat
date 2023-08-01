package history

import (
	rdb "github.com/ayo-ajayi/chitchat/redis"
	"github.com/redis/go-redis/v9"
)

func SaveChat(input, response string) error {
	c := rdb.DefaultClient()
	defer c.C.Close()
	if err := c.AddStream(input, response); err != nil {
		return err
	}
	return nil
}
func GetChat(limit int64) ([]redis.XMessage, error) {
	c := rdb.DefaultClient()
	defer c.C.Close()
	res, err := c.ReadStream(limit)
	if err != nil {
		return nil, err
	}
	return res, err
}
func GetDate(id string) (string, error) {
	date, err := rdb.GetDateFromStreamID(id)
	if err != nil {
		return "", err
	}
	return date, nil
}
