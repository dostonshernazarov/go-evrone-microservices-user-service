package app

import (
	"fmt"
	"github/user_service_evrone_microservces/internal/delivery/kafka/handlers"
	"github/user_service_evrone_microservces/internal/infrastructure/kafka"
	"github/user_service_evrone_microservces/internal/infrastructure/repository/postgresql"
	"github/user_service_evrone_microservces/internal/pkg/config"
	logPkg "github/user_service_evrone_microservces/internal/pkg/logger"
	"github/user_service_evrone_microservces/internal/pkg/postgres"
	"github/user_service_evrone_microservces/internal/usecase"
	"github/user_service_evrone_microservces/internal/usecase/event"

	"go.uber.org/zap"
)

type UserCreateConsumerCLI struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	BrokerConsumer event.BrokerConsumer
}

func NewUserConsumer(config *config.Config) (*UserCreateConsumerCLI, error) {
	logger, err := logPkg.New(config.LogLevel, config.Environment, config.APP+"_cli"+".log")
	if err != nil {
		return nil, err
	}

	consumer := kafka.NewConsumer(logger)

	db, err := postgres.New(config)
	if err != nil {
		return nil, err
	}

	return &UserCreateConsumerCLI{
		Config:         config,
		Logger:         logger,
		DB:             db,
		BrokerConsumer: consumer,
	}, nil
}

func (c *UserCreateConsumerCLI) Run() error {
	fmt.Print("consume is running ....")
	// repo init
	userRepo := postgresql.NewUsersRepo(c.DB)

	// usecase init
	userUsecase := usecase.NewUsersService(c.DB.Config().ConnConfig.ConnectTimeout, userRepo)

	eventHandler := handlers.NewUserCreateHandler(c.Config, c.BrokerConsumer, c.Logger, userUsecase)

	return eventHandler.HandlerEvents()
}

func (c *UserCreateConsumerCLI) Close() {
	c.BrokerConsumer.Close()

	c.Logger.Sync()
}
