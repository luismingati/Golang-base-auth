package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/luismingati/buymeacoffee/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(ctx context.Context) (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisEndpoint(),
		Password: config.GetRedisPassword(),
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Error("falha ao conectar ao Redis: %v", err.Error(), err)
		return nil, fmt.Errorf("falha ao conectar ao Redis: %w", err)
	}

	return &RedisService{client: client}, nil
}

func (r *RedisService) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		slog.Error("erro ao salvar chave no Redis: %v", err.Error(), err)
		return fmt.Errorf("erro ao salvar chave no Redis: %w", err)
	}
	return nil
}

func (r *RedisService) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("chave não encontrada no Redis")
	} else if err != nil {
		return "", fmt.Errorf("erro ao buscar chave no Redis: %w", err)
	}
	return value, nil
}

func (r *RedisService) Del(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		slog.Error("erro ao remover chave no Redis: %v", err.Error(), err)
		return fmt.Errorf("erro ao remover chave no Redis: %w", err)
	}
	return nil
}

func (r *RedisService) Close() error {
	return r.client.Close()
}
