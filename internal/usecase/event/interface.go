package event

import (
	"github/user_service_evrone_microservces/internal/entity"
	"context"
)

type ConsumerConfig interface {
	GetBrokers() []string
	GetTopic() string
	GetGroupID() string
	GetHandler() func(ctx context.Context, key, value []byte) error
}

type BrokerConsumer interface {
	Run()
	RegisterConsumer(config ConsumerConfig)
	Close()
}

type BrokerProducer interface {
	ProduceContent(ctx context.Context, key string, value *entity.Users) error
	Close()
}
