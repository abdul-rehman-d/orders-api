package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/abdul-rehman-d/orders-api/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	_, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	return nil
}
