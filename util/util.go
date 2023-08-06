package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ayo-ajayi/chitchat/types"
)

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


func ConvertToXM(redisMessages []redis.XMessage) []types.Message {
	customMessages := make([]types.Message, len(redisMessages))
	for i, redisMsg := range redisMessages {
		customMessages[i] = types.Message{
			ID:     redisMsg.ID,
			Values: redisMsg.Values,
		}
	}
	return customMessages
}
