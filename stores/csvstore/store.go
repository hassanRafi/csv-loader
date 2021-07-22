package csvstore

import (
	"context"
	"net/http"
	"strconv"

	"github.com/csv-loader/stores"
	"github.com/csv-loader/types"
	"github.com/go-redis/redis/v8"
)

type Store struct {
	client *redis.Client
}

func New(client *redis.Client) stores.CSVGetter {
	return &Store{
		client: client,
	}
}

func (s *Store) GetByKey(key string) (*types.KeyValPair, error) {
	val, err := s.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, &types.CustomError{
				Code:    http.StatusNotFound,
				Message: http.StatusText(http.StatusNotFound),
			}
		}
		return nil, err
	}

	v, err := strconv.Atoi(val)
	if err != nil {
		return nil, err
	}

	return &types.KeyValPair{Key: key, Value: v}, nil
}
