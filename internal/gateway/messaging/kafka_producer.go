package messaging

import (
	"context"
	"encoding/json"
	"log/slog"
	kafkaserver "product-wallet/pkg/broker/kafkaservice"

	"github.com/segmentio/kafka-go"
)

type ProducerKafka[T any] struct {
	Topic         string
	KafkaProducer *kafkaserver.KafkaService
}

func (p *ProducerKafka[T]) GetTopic() string {
	return p.Topic
}

func (p *ProducerKafka[T]) Send(ctx context.Context, data ...T) error {
	var payloads []kafka.Message
	for _, d := range data {
		payload, err := json.Marshal(d)
		if err != nil {
			slog.Error("error when marshal payload", slog.String("error", err.Error()))
			return err
		}
		payloads = append(payloads, kafka.Message{Value: payload})
	}

	writer := p.KafkaProducer.NewWriter(p.Topic)
	defer writer.Close()

	if err := writer.WriteMessages(ctx, payloads...); err != nil {
		slog.Error("failed to write messages", slog.String("error", err.Error()))
		return err
	}

	slog.Info("Success publish transaction", "topic", p.Topic)

	return nil
}
