package orderservice

import (
	"context"
	"fmt"
	"log"

	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
)

type Service struct {
	repository ordersrepository.Repository
	lg         *log.Logger
}

func NewService(lg *log.Logger, rep ordersrepository.Repository) *Service {
	return &Service{
		repository: rep,
		lg:         lg,
	}
}

func (s *Service) GetByUUID(ctx context.Context, id string) (ordersrepository.Order, error) {
	o, err := s.repository.FindOne(ctx, id)
	if err != nil {
		return ordersrepository.Order{}, fmt.Errorf("failed to get order: %w", err)
	}
	return o, nil
}