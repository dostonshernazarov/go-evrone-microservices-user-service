package services

import (
	"context"
	userproto "github/user_service_evrone_microservces/genproto/user_proto"
	"github/user_service_evrone_microservces/internal/entity"
	"github/user_service_evrone_microservces/internal/pkg/otlp"
	"github/user_service_evrone_microservces/internal/usecase"
	"github/user_service_evrone_microservces/internal/usecase/event"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type userRPC struct {
	logger         *zap.Logger
	userUsecase    usecase.Users
	brokerProducer event.BrokerProducer
}

func NewRPC(logger *zap.Logger, userUsecase usecase.Users) userproto.UserServiceServer {
	return &userRPC{
		logger:         logger,
		userUsecase:    userUsecase,
		// brokerProducer: brokerProducer,
	}
}

func (s userRPC) Create(ctx context.Context, in *userproto.User) (*userproto.UserRes, error) {
	//AddAtribute in jeager
	ctx, span := otlp.Start(ctx, "user_grpc_delivery", "Create")
	span.SetAttributes(
		attribute.Key("id").String(in.Id),
	)
	defer span.End()

	guid, err := s.userUsecase.Create(ctx, &entity.Users{
		GUID:      in.Id,
		FullName: in.FullName,
		Username:  in.Username,
		Email:     in.Email,
		Password:  in.Password,
		Bio:       in.Bio,
		Website:   in.Website,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &userproto.UserRes{
		Id:                   guid,
		FullName:             in.FullName,
		Username:  in.Username,
		Email:     in.Email,
		Password:  in.Password,
		Bio:       in.Bio,
		Website:   in.Website,
		Role:                 "user",
	}, nil
}

// func (s userRPC) Update(ctx context.Context, in *userproto.User) (*userproto.User, error) {
// 	err := s.userUsecase.Update(ctx, &entity.Users{
// 		GUID:      in.Id,
// 		FirstName: in.FirstName,
// 		LastName:  in.LastName,
// 		Username:  in.Username,
// 		Email:     in.Email,
// 		Password:  in.Password,
// 		Bio:       in.Password,
// 		Website:   in.Website,
// 		UpdatedAt: time.Now().UTC(),
// 	})
// 	if err != nil {
// 		s.logger.Error(err.Error())
// 		return nil, err
// 	}
// 	return &userproto.User{
// 		Id:        in.Id,
// 		FirstName: in.FirstName,
// 		LastName:  in.LastName,
// 		Username:  in.Username,
// 		Email:     in.Email,
// 		Password:  in.Password,
// 		Bio:       in.Bio,
// 		Website:   in.Website,
// 		CreatedAt: in.CreatedAt,
// 		UpdatedAt: in.UpdatedAt,
// 	}, nil
// }

func (s userRPC) DeleteUserByID(ctx context.Context, in *userproto.IdRequest) (*userproto.DeleteUserByIDRespons, error) {
	ctx, span := otlp.Start(ctx, "user_grpc_delivery", "Delete")
	span.SetAttributes(
		attribute.Key("id").String(in.Id),
	)
	defer span.End()

	if err := s.userUsecase.Delete(ctx, in.Id); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &userproto.DeleteUserByIDRespons{}, nil
}

func (s userRPC) GetUserByID(ctx context.Context, in *userproto.IdRequest) (*userproto.UserRes, error) {
	ctx, span := otlp.Start(ctx, "user_grpc_delivery", "Get")
	span.SetAttributes(
		attribute.Key("id").String(in.Id),
	)
	defer span.End()

	user, err := s.userUsecase.Get(ctx, map[string]string{
		"id": in.Id,
	})

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &userproto.UserRes{
		Id:        user.GUID,
		FullName: user.FullName,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Bio:       user.Bio,
		Website:   user.Website,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s userRPC) GetAllUsers(ctx context.Context, in *userproto.GetAllUsersRequest) (*userproto.GetAllUsersRespons, error) {
	ctx, span := otlp.Start(ctx, "user_grpc_delivery", "List")
	
	defer span.End()

	offset := in.Limit * (in.Page - 1)
	users, err := s.userUsecase.List(ctx, uint64(in.Limit), uint64(offset), map[string]string{})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	var resp userproto.GetAllUsersRespons
	for _, u := range users {
		resp.User = append(resp.User, &userproto.UserRes{
			Id:        u.GUID,
			FullName: u.FullName,
			Username:  u.Username,
			Email:     u.Email,
			Password:  u.Password,
			Bio:       u.Bio,
			Website:   u.Website,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
			UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &resp, nil
}
