package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"log"
	"os"
)

type Config struct {
	Addr     string
	Password string
	DB       int
	DBurl    string
}

func New(config *Config) *redsync.Redsync {
	var red *redis.Client

	url, err := redis.ParseURL(config.Addr)
	if err != nil {
		red = redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
		})
	} else {
		red = redis.NewClient(url)
	}

	pong, err := red.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println(pong)

	pool := goredis.NewPool(red)
	redisSync := redsync.New(pool)

	return redisSync
}
