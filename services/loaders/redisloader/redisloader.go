package redisloader

import (
	"context"
	"log"

	"github.com/csv-loader/types"
	redis "github.com/go-redis/redis/v8"
)

type Service struct {
	client        *redis.Client
	failedRecords []*types.KeyValPair
}

func New(client *redis.Client) *Service {
	s := &Service{
		client:        client,
		failedRecords: []*types.KeyValPair{},
	}

	go s.retryFailedRecords()

	return s
}

func (s *Service) Store(pairs []*types.KeyValPair) error {
	p := []interface{}{}

	for i := range pairs {
		p = append(p, pairs[i].Key, pairs[i].Value)
	}

	if err := s.client.MSet(context.Background(), p...).Err(); err != nil {
		log.Printf("Failed to set the key/val pairs: %v", pairs)
		s.failedRecords = pairs

		return err
	}

	return nil
}

func (s *Service) retryFailedRecords() {
	for {
		if len(s.failedRecords) != 0 {
			p := []interface{}{}

			for i := range s.failedRecords {
				p = append(p, s.failedRecords[i].Key, s.failedRecords[i].Value)
			}

			if err := s.client.MSet(context.Background(), p...).Err(); err != nil {
				log.Printf("Failed to set the key/val pairs on retry: %v, error: %v", s.failedRecords, err)
			} else {
				s.failedRecords = []*types.KeyValPair{}
			}
		}
	}
}
