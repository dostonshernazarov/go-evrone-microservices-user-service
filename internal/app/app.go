package app

import (
	"fmt"
	pb "github/user_service_evrone_microservces/genproto/user_proto"
	grpc_server "github/user_service_evrone_microservces/internal/delivery/grpc/server"
	clean_grpc "github/user_service_evrone_microservces/internal/delivery/grpc/services"
	"github/user_service_evrone_microservces/internal/infrastructure/grpc_service_clients"
	"github/user_service_evrone_microservces/internal/infrastructure/kafka"
	"github/user_service_evrone_microservces/internal/pkg/config"
	"github/user_service_evrone_microservces/internal/pkg/logger"
	"github/user_service_evrone_microservces/internal/pkg/otlp"
	"github/user_service_evrone_microservces/internal/pkg/postgres"
	"github/user_service_evrone_microservces/internal/usecase"
	"github/user_service_evrone_microservces/internal/usecase/event"
	"time"

	repo "github/user_service_evrone_microservces/internal/infrastructure/repository/postgresql"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	GrpcServer     *grpc.Server
	ShutdownOTLP   func() error
	ServiceClients grpc_service_clients.ServiceClients
	BrokerProducer event.BrokerProducer
	BrokerConsumer event.BrokerConsumer
}

func NewApp(cfg *config.Config) (*App, error) {
	// init logger
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	kafkaProducer := kafka.NewProducer(cfg, logger)
	kafkaConsumer := kafka.NewConsumer(logger)

	// otlp collector initialization
	shutdownOTLP, err := otlp.InitOTLPProvider(cfg)
	if err != nil {
		return nil, err
	}

	// init db
	db, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}

	// grpc server init
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_server.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_recovery.UnaryServerInterceptor(),
			),
			grpc_server.UnaryInterceptorData(logger),
		)),
	)

	return &App{
		Config:     cfg,
		Logger:     logger,
		DB:         db,
		GrpcServer: grpcServer,
		ShutdownOTLP:   shutdownOTLP,
		BrokerProducer: kafkaProducer,
		BrokerConsumer: kafkaConsumer,
	}, nil
}
func (a *App) Run() error {
	var (
		contextTimeout time.Duration
	)

	// context timeout initialization
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error during parse duration for context timeout : %w", err)
	}
	// Initialize Service Clients
	serviceClients, err := grpc_service_clients.New(a.Config)
	if err != nil {
		return fmt.Errorf("error during initialize service clients: %w", err)
	}
	a.ServiceClients = serviceClients

	// repositories initialization
	articleRepo := repo.NewUsersRepo(a.DB)

	// usecase initialization
	articleUsecase := usecase.NewUsersService(contextTimeout, articleRepo)

	pb.RegisterUserServiceServer(a.GrpcServer, clean_grpc.NewRPC(a.Logger, articleUsecase, a.BrokerProducer))
	a.Logger.Info("gRPC Server Listening", zap.String("url", a.Config.RPCPort))
	if err := grpc_server.Run(a.Config, a.GrpcServer); err != nil {
		return fmt.Errorf("gRPC fatal to serve grpc server over %s %w", a.Config.RPCPort, err)
	}

	return nil
}

func (a *App) Stop() {
	// close broker producer
	a.BrokerProducer.Close()
	// closing client service connections
	a.ServiceClients.Close()
	// stop gRPC server
	a.GrpcServer.Stop()

	// database connection
	a.DB.Close()

	// shutdown otlp collector
	if err := a.ShutdownOTLP(); err != nil {
		a.Logger.Error("shutdown otlp collector", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}
