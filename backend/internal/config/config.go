package config

import (
	"log/slog"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	Repository ordersrepository.Config `yaml:"repository"`
}

func GetConfig(lg *slog.Logger) (*Config, error) {
	return sync.OnceValues(func() (*Config, error) {
		instance := &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			lg.Info(help)
			lg.Error("failed to read config", "err", err)
			return nil, err
		}
		return instance, nil
	})()
}
