package repository

import (
	"github/user_service_evrone_microservces/internal/entity"
	"context"
)

type Users interface {
	Create(ctx context.Context, kyc *entity.Users) error
	Get(ctx context.Context, params map[string]string) (*entity.Users, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Users, error)
	Delete(ctx context.Context, guid string) error
}