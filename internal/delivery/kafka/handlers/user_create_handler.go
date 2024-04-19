package handlers

import (
	"context"
	"encoding/json"
	"github/user_service_evrone_microservces/internal/entity"
	"github/user_service_evrone_microservces/internal/infrastructure/kafka"
	"github/user_service_evrone_microservces/internal/pkg/config"
	"github/user_service_evrone_microservces/internal/usecase"
	"github/user_service_evrone_microservces/internal/usecase/event"

	"go.uber.org/zap"
)

type userCreateHandler struct {
	config         *config.Config
	brokerConsumer event.BrokerConsumer
	logger         *zap.Logger
	userUsecase    usecase.Users
}

func NewUserCreateHandler(config *config.Config,
	brokerConsumer event.BrokerConsumer,
	logger *zap.Logger,
	userUsecase usecase.Users) *userCreateHandler {
	return &userCreateHandler{
		config:         config,
		brokerConsumer: brokerConsumer,
		logger:         logger,
		userUsecase:    userUsecase,
	}
}

func (h *userCreateHandler) HandlerEvents() error {
	consumerConfig := kafka.NewConsumerConfig(
		h.config.Kafka.Address,
		"api.user.create",
		"1",
		func(ctx context.Context, key, value []byte) error {
			var user *entity.Users

			if err := json.Unmarshal(value, &user); err != nil {
				return err
			}

			if _, err := h.userUsecase.Create(ctx, user); err != nil {
				return err
			}

			return nil
		},
	)

	h.brokerConsumer.RegisterConsumer(consumerConfig)
	h.brokerConsumer.Run()

	return nil

}
