package grpc

import (
	"github/user_service_evrone_microservces/internal/entity"
	"context"
	"errors"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errNotFound   *entity.ErrNotFound
	errConflict   *entity.ErrConflict
	errValidation *entity.ErrValidation
)

func ErrorStatus(ctx context.Context, err error) *status.Status {
	var (
		st *status.Status
	)
	switch {
	// error not found
	case errors.As(err, &errNotFound):
		st = status.New(codes.NotFound, err.Error())
	// error conflict
	case errors.As(err, &errConflict):
		st = status.New(codes.AlreadyExists, err.Error())
	// error validation errors
	case errors.As(err, &errValidation):
		st = status.New(codes.InvalidArgument, codes.InvalidArgument.String())
		br := &epb.BadRequest{}
		for field, des := range errValidation.Errors {
			br.FieldViolations = append(br.FieldViolations, &epb.BadRequest_FieldViolation{
				Field:       field,
				Description: des,
			})
		}
		st, _ = st.WithDetails(br)
	// error internal
	default:
		st = status.New(codes.Internal, codes.Internal.String())
		errInfo := &epb.ErrorInfo{
			Reason: err.Error(),
		}
		st, _ = st.WithDetails(errInfo)
	}
	return st
}

func Error(ctx context.Context, err error) error {
	return ErrorStatus(ctx, err).Err()
}
