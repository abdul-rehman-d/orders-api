package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/abdul-rehman-d/orders-api/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

var ErrNotFound = errors.New("order does not exist")

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := orderIDKey(order.OrderID)

	pipe := r.Client.TxPipeline()

	res := pipe.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		pipe.Discard()
		return fmt.Errorf("failed to set %w", err)
	}

	if err := pipe.SAdd(ctx, "orders", key).Err(); err != nil {
		pipe.Discard()
		return fmt.Errorf("failed to set %w", err)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec %w", err)
	}

	return nil
}

type Page struct {
	Cursor uint64
	Count  int64
}

type GetAllResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepo) GetAll(ctx context.Context, page Page) (GetAllResult, error) {
	fmt.Println(page.Count)
	keys, cursor, err := r.Client.SScan(ctx, "orders", page.Cursor, "*", page.Count).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return GetAllResult{}, ErrNotFound
		} else {
			return GetAllResult{}, fmt.Errorf("failed to get %w", err)
		}
	}

	if len(keys) == 0 {
		return GetAllResult{
			Orders: []model.Order{},
			Cursor: 0,
		}, nil
	}

	all, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return GetAllResult{}, fmt.Errorf("failed to get %w", err)
	}

	orders := make([]model.Order, len(all))

	for idx, data := range all {
		data := data.(string)
		err = json.Unmarshal([]byte(data), &orders[idx])
		if err != nil {
			return GetAllResult{}, fmt.Errorf("failed to decode order %w", err)
		}
	}

	return GetAllResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}

func (r *RedisRepo) Get(ctx context.Context, id uint64) (model.Order, error) {
	key := orderIDKey(id)

	res, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Order{}, ErrNotFound
		} else {
			return model.Order{}, fmt.Errorf("failed to get %w", err)
		}
	}

	o := model.Order{}

	err = json.Unmarshal([]byte(res), &o)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode order %w", err)
	}

	return o, nil
}

func (r *RedisRepo) Update(ctx context.Context, id uint64, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := orderIDKey(id)

	err = r.Client.SetXX(ctx, key, data, 0).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		} else {
			return fmt.Errorf("failed to get %w", err)
		}
	}

	return nil
}

func (r *RedisRepo) Delete(ctx context.Context, id uint64) error {
	key := orderIDKey(id)

	pipe := r.Client.TxPipeline()

	err := pipe.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		pipe.Discard()
		return ErrNotFound
	} else if err != nil {
		pipe.Discard()
		return fmt.Errorf("failed to delete %w", err)
	}

	if err := pipe.SRem(ctx, "orders", key).Err(); err != nil {
		pipe.Discard()
		return fmt.Errorf("failed to remove from orders %w", err)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec %w", err)
	}

	return nil
}
