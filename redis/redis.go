package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

const ADDR string = ""
const DEFAULT_DB int = 4
const HISTORY_KEY string = "chitchat:history"
const MODELS_KEY string = "chitchat:model"

type Client struct {
	C *redis.Client
}
type options struct {
	opt        *redis.Options
	historyKey string
	modelsKey  string
}

func NewClient(opt *options) *Client {
	client := Client{
		C: redis.NewClient(opt.opt),
	}
	return &client
}
func DefaultClient() *Client {

	var historyKey string = HISTORY_KEY
	var defaultDB int = DEFAULT_DB
	var addr string = ADDR
	var modelsKey string = MODELS_KEY

	const redis_yaml = "redis"
	const models_yaml = "models_key"
	const history_yaml = "history_key"
	const db_yaml = "default_db"
	const addr_yaml = "addr"

	if viper.IsSet(redis_yaml) {
		redis := viper.GetViper().GetStringMap("redis")
		if redis[history_yaml] != nil {
			historyKey = redis[history_yaml].(string)
		}
		if redis[db_yaml] != nil {
			defaultDB = redis[db_yaml].(int)
		}
		if redis[addr_yaml] != nil {
			addr = redis[addr_yaml].(string)
		}
		if redis[models_yaml] != nil {
			modelsKey = redis[models_yaml].(string)
		}
	}

	client := NewClient(&options{
		opt: &redis.Options{
			Addr: addr,
			DB:   defaultDB,
		},
		historyKey: historyKey,
		modelsKey:  modelsKey,
	})
	return client
}
func (c *Client) DelKey(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.C.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
func (c *Client) Exists(key string) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.C.Exists(ctx, key).Val()
}
func (c *Client) SetList(key string, ls []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	list := make([]interface{}, len(ls))
	for i, l := range ls {
		list[i] = l
	}
	if exists := c.C.Exists(ctx, key).Val(); exists == 1 {
		if err := c.C.Del(ctx, key).Err(); err != nil {
			return err
		}
	}
	if err := c.C.RPush(ctx, key, list...).Err(); err != nil {
		return err
	}
	return nil
}
func (c *Client) GetList(key string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	arr, err := c.C.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return arr, nil
}
func (c *Client) AddStream(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.C.XAdd(ctx, &redis.XAddArgs{
		Stream: HISTORY_KEY,
		Values: map[string]any{key: value},
	}).Err(); err != nil {
		return err
	}
	return nil
}
func (c *Client) ReadStream(limit int64) ([]redis.XMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := c.C.XRevRangeN(ctx, HISTORY_KEY, "+", "-", limit).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}
func GetDateFromStreamID(streamID string) (string, error) {
	parts := strings.Split(streamID, "-")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid stream ID format")
	}

	timestampMs, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse timestamp: %v", err)
	}
	timestampSec := timestampMs / 1000
	t := time.Unix(timestampSec, 0).UTC()
	return t.Format("2006-01-02 15:04:05"), nil
}
