package producer

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	Topic    string
	producer *kafka.Producer
}

func (p *Producer) Send(value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: bytes,
	}
	p.producer.ProduceChannel() <- msg

	e := <-p.producer.Events()
	msg = e.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		return msg.TopicPartition.Error
	}
	return nil
}

func newProcuder(address, topic string) (*Producer, error) {
	configMap := &kafka.ConfigMap{"bootstrap.servers": address}
	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		return nil, err
	}
	return &Producer{
		Topic:    topic,
		producer: producer,
	}, nil
}