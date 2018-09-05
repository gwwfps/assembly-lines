package db

import "github.com/go-redis/redis"

type DB struct {
	rc *redis.Client
}

func NewDB(redisUrl string, redisPassword string) *DB {
	return &DB{
		rc: redis.NewClient(&redis.Options{
			Addr:     redisUrl,
			Password: redisPassword,
			DB:       0,
		}),
	}
}
