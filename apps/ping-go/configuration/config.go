package configuration

import "github.com/kelseyhightower/envconfig"

type Configuration struct {
	ApplicationName        string `envconfig:"APPLICATION_NAME" default:"ping-go"`
	ApplicationHost        string `envconfig:"APPLICATION_HOST" default:"127.0.0.1"`
	ApplicationEnvironment string `envconfig:"APPLICATION_ENVIRONMENT" default:"development"`
	ApplicationPort        string `envconfig:"APPLICATION_PORT" default:"3000"`
	RedisAddress           string `envconfig:"REDIS_ADDRESS" default:"127.0.0.1"`
	RedisPort              string `envconfig:"REDIS_PORT" default:"6379"`
}

func ReadConfiguration() *Configuration {
	var config Configuration
	envconfig.MustProcess("", &config)
	return &config
}
