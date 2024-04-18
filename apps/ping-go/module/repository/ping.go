package repository

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/hexennacht/ping-pong-pubsub/module/entity"
)

type PingRepository interface {
	Ping(ctx context.Context, req *entity.PingRequest) error
	GetMessage() (*entity.PingRequest, error)
}

type pingRepository struct {
	key    string
	redis  *redis.Client
	pubsub *redis.PubSub
}

func NewPingRepository(redis *redis.Client) PingRepository {
	return &pingRepository{
		key:    "com.github.hexennacht.ping-pong-pubsub.go.ping",
		redis:  redis,
		pubsub: redis.Subscribe("com.github.hexennacht.ping-pong-pubsub.rust.pong"),
	}
}

// Ping implements PingRepository.
func (p *pingRepository) Ping(ctx context.Context, req *entity.PingRequest) error {
	message, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = p.redis.Publish(p.key, message).Err()
	if err != nil {
		return err
	}

	return nil
}

func (p *pingRepository) GetMessage() (*entity.PingRequest, error) {
	message, err := p.pubsub.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	var request entity.PingRequest
	err = json.Unmarshal([]byte(message.Payload), &request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
