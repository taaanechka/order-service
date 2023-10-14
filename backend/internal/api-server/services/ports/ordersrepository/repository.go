package ordersrepository

import "context"

type Config struct {
	Host       string   `yaml:"host"`
	Port       string   `yaml:"port"`
	Database   string   `yaml:"database"`
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
}

type Repository interface {
	// Create(ctx context.Context, order Order) (string, error)
	FindOne(ctx context.Context, id string) (Order, error)
}
