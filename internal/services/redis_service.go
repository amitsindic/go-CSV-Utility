package services

import (
	"context"
	"encoding/json"

	"ctlSolution.com/internal/models"
	"github.com/go-redis/redis/v8"
)

const (
	redisTimeout = 5 * 60000 * 10000 * 100
)

type RedisService interface {
	CacheRecord(record *models.Record) error
	GetData() ([]models.Record, error)
	GetDataByKey(key string) ([]models.Record, error)
	DeleteRecordByFirstName(firstName string) error
}

type redisService struct {
	redis *redis.Client
}

func NewRedisService(redis *redis.Client) RedisService {
	return &redisService{redis: redis}
}

func (s *redisService) CacheRecord(record *models.Record) error {
	
	jsonRecord, err := json.Marshal(record)
	if err != nil {
		return err
	}

	ctx := context.Background()

	err = s.redis.Set(ctx, record.FirstName, jsonRecord, redisTimeout).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *redisService) GetDataByKey(key string) ([]models.Record, error) {
	ctx := context.Background()
	keys, err := s.redis.Keys(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var records []models.Record
	for _, key := range keys {
		data, err := s.redis.Get(ctx, key).Bytes()
		if err != nil {
			return nil, err
		}
		var record models.Record
		if err := json.Unmarshal(data, &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (s *redisService) GetData() ([]models.Record, error) {
	ctx := context.Background()
	keys, err := s.redis.Keys(ctx, "*").Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var records []models.Record
	for _, key := range keys {
		data, err := s.redis.Get(ctx, key).Bytes()
		if err != nil {
			return nil, err
		}
		var record models.Record
		if err := json.Unmarshal(data, &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (s *redisService) DeleteRecordByFirstName(firstName string) error {
	ctx := context.Background()
	if err := s.redis.Del(ctx, firstName).Err(); err != nil {
		return err
	}
	return nil
}
