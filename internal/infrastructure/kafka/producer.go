package kafka

import (
	"context"
	"github/user_service_evrone_microservces/internal/entity"
	"github/user_service_evrone_microservces/internal/pkg/config"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type producer struct {
	logger            *zap.Logger
	investmentCreated *kafka.Writer
}

func NewProducer(config *config.Config, logger *zap.Logger) *producer {
	return &producer{
		logger: logger,
		investmentCreated: &kafka.Writer{
			Addr:                   kafka.TCP(config.Kafka.Address...),
			Topic:                  config.Kafka.Topic.InvestorCreate,
			Balancer:               &kafka.Hash{},
			RequiredAcks:           kafka.RequireAll,
			AllowAutoTopicCreation: true,
			Async:                  true, // make the writer asynchronous
			Completion: func(messages []kafka.Message, err error) {
				if err != nil {
					logger.Error("kafka investmentCreated", zap.Error(err))
				}
				for _, message := range messages {
					logger.Sugar().Info(
						"kafka investmentCreated message",
						zap.Int("partition", message.Partition),
						zap.Int64("offset", message.Offset),
						zap.String("key", string(message.Key)),
						zap.String("value", string(message.Value)),
					)
				}
			},
		},
	}
}

// func (p *producer) buildMessageWithTracing(key string, value []byte, otlpSpan otlp.Span) kafka.Message {
func (p *producer) buildMessageWithTracing(key string, value []byte) kafka.Message {
	return kafka.Message{
		Key:   []byte(key),
		Value: value,
		Headers: []kafka.Header{
			{
				Key: "trace_id",
				// Value: []byte(otlpSpan.SpanContext().TraceID().String()),
			},
			{
				Key: "span_id",
				// Value: []byte(otlpSpan.SpanContext().SpanID().String()),
			},
		},
	}
}

func (p *producer) ProduceContent(ctx context.Context, key string, value *entity.Users) error {
	// tracing
	// ctx, span := otlp.Start(ctx, "kafka producer", "ContentProducer")
	// defer span.End()

	return nil
}

func (p *producer) Close() {
	if err := p.investmentCreated.Close(); err != nil {
		p.logger.Error("error during close writer articleCategoryCreated", zap.Error(err))
	}
}
