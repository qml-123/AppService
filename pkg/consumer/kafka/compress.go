package kafka

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type CompressHandler struct{}

func (*CompressHandler) HandleMessage(ctx context.Context, msg *kafka.Message) error {
	var file_keys []string
	if err := json.Unmarshal(msg.Value, &file_keys); err != nil {
		return err
	}
	//go file.Compress(ctx, file_keys)
	return nil
}
