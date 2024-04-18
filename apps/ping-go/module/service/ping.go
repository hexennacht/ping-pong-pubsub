package service

import (
	"context"
	"fmt"

	"github.com/hexennacht/ping-pong-pubsub/helper"
	"github.com/hexennacht/ping-pong-pubsub/module/entity"
	"github.com/hexennacht/ping-pong-pubsub/module/repository"
	"go.uber.org/zap"
)

type PingService interface {
	Ping(ctx context.Context, req *entity.PingRequest) (*entity.PingResponse, error)
	ReceiveMessage() error
}

type pingService struct {
	repo repository.PingRepository
	log  *zap.Logger
}

func NewPingService(repo repository.PingRepository, log *zap.Logger) PingService {
	service := &pingService{repo: repo, log: log}

	return service
}

// Ping implements PingService.
func (p *pingService) Ping(ctx context.Context, req *entity.PingRequest) (*entity.PingResponse, error) {
	err := p.repo.Ping(ctx, req)
	if err != nil {
		return nil, helper.NewError("Failed to send message to pong service", 500, err)
	}

	return &entity.PingResponse{
		Message: fmt.Sprintf("Success sending message %s to pong service!", req.Message),
	}, nil
}

// ReceiveMessage implements PingService.
func (p *pingService) ReceiveMessage() error {
	defer p.log.Sync()

	message, err := p.repo.GetMessage()
	if err != nil {
		p.log.WithOptions(zap.Fields(zap.String("p.repo.GetMessage", err.Error()))).Error("Failed to receive message from pong service", zap.Error(err))
		return helper.NewError("Failed to receive message from pong service", 500, err)
	}

	p.log.Info("Received message from pong service", zap.String("message", message.Message), zap.Int32("limit", message.Limit))

	if message.Limit <= 0 {
		p.log.Info("Message limit reached", zap.String("message", message.Message), zap.Int32("limit", message.Limit))
		return nil
	}

	message.Limit -= 1

	err = p.repo.Ping(context.Background(), message)
	if err != nil {
		return helper.NewError("Failed to send message to pong service", 500, err)
	}

	return nil
}
