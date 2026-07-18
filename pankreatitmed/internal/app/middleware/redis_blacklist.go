package middleware

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisBlacklist struct {
	client *redis.Client
	prefix string // например "jwt_blacklist:"
}

func NewRedisBlacklist(addr, password string, db int) *RedisBlacklist {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisBlacklist{client: rdb, prefix: "jwt_blacklist:"}
}

// добавили токен в blacklist до его истечения
func (r *RedisBlacklist) Add(token string, exp time.Time) error {
	ctx := context.Background()
	ttl := time.Until(exp)
	if ttl <= 0 {
		ttl = time.Hour // если не указали на часик забаним
	}
	return r.client.Set(ctx, r.prefix+token, "revoked", ttl).Err()
}

// проверим, есть ли токен в blacklist
func (r *RedisBlacklist) IsBlacklisted(token string) (bool, error) {
	ctx := context.Background()
	exists, err := r.client.Exists(ctx, r.prefix+token).Result()
	return exists > 0, err
}
