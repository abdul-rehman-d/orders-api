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

func (r *RedisRepo) GetAll(ctx context.Context) ([]model.Order, error) {
	return []model.Order{}, nil
}
func (r *RedisRepo) Get(ctx context.Context, id int) (model.Order, error) {
	return model.Order{}, nil
}
func (r *RedisRepo) Update(ctx context.Context, id int, order model.Order) (model.Order, error) {
	return model.Order{}, nil
}
func (r *RedisRepo) Delete(ctx context.Context, id int) (model.Order, error) {
	return model.Order{}, nil
}
