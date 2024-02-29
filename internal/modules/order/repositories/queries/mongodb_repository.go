package queries

import (
	"context"
	"notification-service/internal/modules/order"
	"notification-service/internal/modules/order/models/entity"
	"notification-service/internal/pkg/databases/mongodb"
	wrapper "notification-service/internal/pkg/helpers"
	"notification-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) order.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindOrderById(ctx context.Context, id string) <-chan wrapper.Result {
	var order entity.Order
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &order,
			CollectionName: "order",
			Filter: bson.M{
				"orderId": id,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
