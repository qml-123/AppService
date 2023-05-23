package producer

import "github.com/qml-123/app_log/common"

var (
	producerMap = make(map[string]*Producer)
)

func InitProducer(conf *common.Conf) (err error) {
	for _, topic := range conf.KafkaProducerTopic {
		producerMap[topic], err = newProcuder(conf.KafkaAddress[0], topic)
		if err != nil {
			return
		}
	}
	return nil
}

func GetProducer(topic string) *Producer {
	return producerMap[topic]
}
