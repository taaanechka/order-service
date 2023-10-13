package config

import (
	"log"
	"os"
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

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

		infoLog.Println("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			infoLog.Println(help)
			errorLog.Fatal(err)
		}
	})

	return instance
}
