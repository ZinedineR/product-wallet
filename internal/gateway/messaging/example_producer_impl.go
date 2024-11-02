package messaging

import (
	"boiler-plate-clean/internal/model"
	kafkaserver "boiler-plate-clean/pkg/broker/kafkaservice"
)

type ExampleProducerImpl struct {
	ProducerKafka[*model.ExampleMessage]
}

func NewExampleKafkaProducerImpl(producer *kafkaserver.KafkaService, topic string) ExampleProducer {
	return &ExampleProducerImpl{
		ProducerKafka: ProducerKafka[*model.ExampleMessage]{
			Topic:         topic,
			KafkaProducer: producer,
		},
	}
}
