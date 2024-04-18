package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/hexennacht/ping-pong-pubsub/configuration"
	"github.com/hexennacht/ping-pong-pubsub/handler/api"
	"github.com/hexennacht/ping-pong-pubsub/module/repository"
	"github.com/hexennacht/ping-pong-pubsub/module/service"
	"github.com/hexennacht/ping-pong-pubsub/server/middleware"
)

type server struct {
	conf   *configuration.Configuration
	redis  *redis.Client
	logger *zap.Logger
	router *chi.Mux
	svc    service.PingService
}

type ServerOption func(*server)

func StartServer(opts ...ServerOption) {
	s := &server{
		router: chi.NewRouter(),
	}

	for _, opt := range opts {
		opt(s)
	}

	s.registerEndpoints()

	err := make(chan error)

	go func() {
		err <- s.serve()
	}()

	go func() {
		for {
			s.svc.ReceiveMessage()
		}
	}()

	<-err
}

func WithLogger(logger *zap.Logger) ServerOption {
	return func(s *server) {
		s.logger = logger
	}
}

func WithRedis(redis *redis.Client) ServerOption {
	return func(s *server) {
		s.redis = redis
	}
}

func WithConfiguration(conf *configuration.Configuration) ServerOption {
	return func(s *server) {
		s.conf = conf
	}
}

func (s *server) registerEndpoints() {
	r := s.router

	r.Use(middleware.NewZapLoggerMiddleware(s.logger))
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	pingRepo := repository.NewPingRepository(s.redis)
	pingSvc := service.NewPingService(pingRepo, s.logger)

	s.svc = pingSvc

	api.RegisterPingHandler(s.router, pingSvc)
}

func (s *server) serve() error {
	address := fmt.Sprintf("%s:%s", s.conf.ApplicationHost, s.conf.ApplicationPort)

	return http.ListenAndServe(address, s.router)
}
