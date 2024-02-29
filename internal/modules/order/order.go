package order

import (
	"context"

	wrapper "notification-service/internal/pkg/helpers"
)

type MongodbRepositoryQuery interface {
	FindOrderById(ctx context.Context, id string) <-chan wrapper.Result
}
