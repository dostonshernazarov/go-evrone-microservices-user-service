package kafka


import (
	// "github/user_service_evrone_microservces/internal/pkg/otlp"
	"github/user_service_evrone_microservces/internal/usecase/event"
	"context"
	// "errors"
	// "fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	MinBytes = 10e3 // 10KB
	MaxBytes = 10e6 // 10MB
)

type HandlerFunc func(ctx context.Context, key, value []byte) error

type consumer struct {
	logger          *zap.Logger
	consumerConfigs []event.ConsumerConfig
	readers         []*kafka.Reader
}

func NewConsumer(logger *zap.Logger) *consumer {
	return &consumer{
		logger: logger,
	}
}

func (c *consumer) RegisterConsumer(consumerConfig event.ConsumerConfig) {
	c.consumerConfigs = append(c.consumerConfigs, consumerConfig)
}

func (c *consumer) Run() {
	for _, consumerConfig := range c.consumerConfigs {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  consumerConfig.GetBrokers(),
			Topic:    consumerConfig.GetTopic(),
			GroupID:  consumerConfig.GetGroupID(),
			MinBytes: MinBytes,
			MaxBytes: MaxBytes,
		})
		c.readers = append(c.readers, r)
		go runReader(r, consumerConfig, c.logger)
	}
}

func (c *consumer) Close() {
	for _, reader := range c.readers {
		if err := reader.Close(); err != nil {
			c.logger.Error("consumer reader close", zap.Error(err))
		}
	}
}


func runReader(r *kafka.Reader, consumerConfig event.ConsumerConfig, logger *zap.Logger) {
	var (
		topic    = consumerConfig.GetTopic()
		handler  = consumerConfig.GetHandler()
	)
	for {
		ctx := context.Background()
		m, err := r.FetchMessage(ctx)
		if err != nil {
			logger.Error("consumer failed to fetch message:", zap.String("topic", topic), zap.Error(err))
			break
		}

	

		if err := handler(ctx, m.Key, m.Value); err != nil {
			logger.Error("consumer failed to handler message:", zap.ByteString("value", m.Value), zap.String("topic", topic), zap.Error(err))
			continue
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			logger.Error("consumer failed to commit messages:", zap.String("topic", topic), zap.Error(err))
		}

	}
}

type ConsumerConfig struct {
	brokers []string
	topic   string
	groupID string
	handler HandlerFunc
}

func NewConsumerConfig(
	brokers []string,
	topic string,
	groupID string,
	handler HandlerFunc,
) *ConsumerConfig {
	return &ConsumerConfig{
		brokers: brokers,
		topic:   topic,
		groupID: groupID,
		handler: handler,
	}
}

func (c *ConsumerConfig) GetBrokers() []string {
	return c.brokers
}

func (c *ConsumerConfig) GetTopic() string {
	return c.topic
}

func (c *ConsumerConfig) GetGroupID() string {
	return c.groupID
}

func (c *ConsumerConfig) GetHandler() func(ctx context.Context, key, value []byte) error {
	return c.handler
}
