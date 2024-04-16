package usecase

import (
	"github/user_service_evrone_microservces/internal/entity"
	"github/user_service_evrone_microservces/internal/infrastructure/repository"
	"github/user_service_evrone_microservces/internal/pkg/otlp"
	"context"
	"time"
)





const (
	serviceNameUsers = "contentService"
	spanNameUsers    = "contentUsecase"
)

type Users interface {
	Create(ctx context.Context, users *entity.Users) (string, error)
	Get(ctx context.Context, params map[string]string) (*entity.Users, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Users, error)
	Delete(ctx context.Context, guid string) error
}

type usersService struct {
	BaseUseCase
	repo       repository.Users
	ctxTimeout time.Duration
}

func NewUsersService(ctxTimeout time.Duration, repo repository.Users) newsService {
	return newsService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (u usersService) Create(ctx context.Context, Users *entity.Users) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameUsers, spanNameUsers+"Create")
	defer span.End()

	u.beforeRequest(&Users.GUID, &Users.CreatedAt, &Users.UpdatedAt)

	return Users.GUID, u.repo.Create(ctx, Users)
}
func (u usersService) Get(ctx context.Context, params map[string]string) (*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameUsers, spanNameUsers+"Get")
	defer span.End()

	return u.repo.Get(ctx, params)
}
func (u usersService) List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameUsers, spanNameUsers+"List")
	defer span.End()

	return u.repo.List(ctx, limit, offset, filter)
}

func (u usersService) Delete(ctx context.Context, guid string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameUsers, spanNameUsers+"Delete")
	defer span.End()

	return u.repo.Delete(ctx, guid)
}