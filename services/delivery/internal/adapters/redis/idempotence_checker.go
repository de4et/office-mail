package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/de4et/office-mail/services/delivery/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	timeout = time.Second
	ttl     = time.Hour * 24
)

type RedisIdempotenceChecker struct {
	client *redis.Client
}

func MustGetRedisIdempotenceChecker(config Config) *RedisIdempotenceChecker {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	return &RedisIdempotenceChecker{
		client: rdb,
	}
}

func (ic *RedisIdempotenceChecker) CheckTaskForIdempotence(ctx context.Context, mail domain.OutboxTask) (bool, error) {
	key := strconv.Itoa(mail.ID)
	return ic.client.SetNX(ctx, key, true, ttl).Result()
}
