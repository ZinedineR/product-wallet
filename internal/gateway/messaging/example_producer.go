package messaging

import (
	"context"
	"product-wallet/internal/model"
)

type ExampleProducer interface {
	GetTopic() string
	Send(ctx context.Context, order ...*model.ExampleMessage) error
}
