package main

import (
	"fmt"

	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/hexennacht/ping-pong-pubsub/configuration"
	"github.com/hexennacht/ping-pong-pubsub/server"
)

func main() {
	conf := configuration.ReadConfiguration()

	logger := zap.Must(zap.NewProduction())

	if conf.ApplicationEnvironment == "development" {
		logger = zap.Must(zap.NewDevelopment())
	}

	defer logger.Sync()

	redisAddress := fmt.Sprintf("%s:%s", conf.RedisAddress, conf.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})

	server.StartServer(
		server.WithConfiguration(conf),
		server.WithLogger(logger),
		server.WithRedis(rdb),
	)
}
