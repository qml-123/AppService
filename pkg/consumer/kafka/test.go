package kafka

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/qml-123/app_log/logger"
)

type testHandle struct{}

func (*testHandle) HandleMessage(ctx context.Context, msg *kafka.Message) error {
	var s string
	err := json.Unmarshal(msg.Value, &s)
	if err != nil {
		return err
	}
	logger.Info(ctx, "get test message: %v", s)
	return nil
}
