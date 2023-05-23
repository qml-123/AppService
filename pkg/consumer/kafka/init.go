package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/qml-123/AppService/pkg/id"
	"github.com/qml-123/app_log/common"
	"github.com/qml-123/app_log/logger"
)

type Consumer interface {
	HandleMessage(context.Context, *kafka.Message) error
}

var (
	handles = map[string]Consumer{
		"compress_defer": &CompressHandler{},
		"test":           &testHandle{},
	}
)

func InitConsumer(conf *common.Conf) error {
	for _, consumerInfo := range conf.KafkaConsumer {
		configMap := &kafka.ConfigMap{
			"bootstrap.servers": conf.KafkaAddress[0],
			"group.id":          consumerInfo.ConsumerName,
			"auto.offset.reset": "earliest",
		}

		c, err := kafka.NewConsumer(configMap)
		if err != nil {
			return err
		}

		topics := []string{consumerInfo.Topic}
		if err = c.SubscribeTopics(topics, nil); err != nil {
			return err
		}
		go register(c, consumerInfo.Topic)
	}
	return nil
}

func register(c *kafka.Consumer, topic string) {
	run := true
	for run == true {
		ctx := id.NewContext()
		ev := c.Poll(50000)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			if err := handleMsg(e, topic); err != nil {
				logger.Error(ctx, "handle message failed, err: %v", err)
			}
		case kafka.Error:
			logger.Error(ctx, "get kafka Error: %v: %v", e.Code(), e)
			if e.Code() == kafka.ErrAllBrokersDown {
				run = false
			}
		default:
			logger.Info(ctx, "Ignored %v", e)
		}
	}
}

func handleMsg(e *kafka.Message, topic string) error {
	if e == nil || e.TopicPartition.Topic == nil {
		return fmt.Errorf("message is nil")
	}
	ctx := id.NewContext()
	logger.Info(ctx, "get message from kafka, topic: %v, value: %v", *e.TopicPartition.Topic, e.Value)
	if *e.TopicPartition.Topic != topic {
		return nil
	}
	consumer, ok := handles[*e.TopicPartition.Topic]
	if !ok {
		return fmt.Errorf("not found handler, topic: %v", *e.TopicPartition.Topic)
	}

	return consumer.HandleMessage(ctx, e)
}
