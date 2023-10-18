package orderservice

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/taaanechka/order-service/internal/api-server/services/ports/cacherepository"
	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
	"github.com/taaanechka/order-service/internal/apperror"
)

type Service struct {
	rep      ordersrepository.Repository
	cache    cacherepository.Repository
	validate *validator.Validate
	lg       *slog.Logger
}

func NewService(lg *slog.Logger, rep ordersrepository.Repository,
	cache cacherepository.Repository, validate *validator.Validate) *Service {
	lg.Info("Service: new order service")
	return &Service{
		rep:      rep,
		cache:    cache,
		validate: validate,
		lg:       lg,
	}
}

func (s *Service) Init(ctx context.Context) {
	s.lg.Info("Init order service")
	orders, err := s.rep.FindAll(ctx)
	if err != nil {
		s.lg.Error("INIT: failed to get all orders", "err", err)
		return
	}

	for _, order := range orders {
		if err := s.validate.Struct(order); err != nil {
			s.lg.Error("INIT: failed to validate orders", "err", err)
			return
		}
	}

	if _, err = s.cache.CreateMany(ctx, orders); err != nil {
		s.lg.Error("INIT: failed to write orders from db to cache", "err", err)
	}
}

func (s *Service) Create(ctx context.Context, order ordersrepository.Order) (string, error) {
	if err := s.validate.Struct(order); err != nil {
		s.lg.Error("Service: failed to validate orders", "err", err)
		return "", apperror.ErrValidate
	}

	uid, err := s.rep.Create(ctx, order)
	if err != nil {
		s.lg.Error("Service: failed to create order", "err", err)
		return "", apperror.ErrCreate
	}

	go func() {
		if _, err := s.cache.CreateOne(ctx, order); err != nil {
			s.lg.Warn("Service: failed to write order to cache", "err", err)
		}
	}()
	return uid, nil
}

func (s *Service) GetByUUID(ctx context.Context, id string) (ordersrepository.Order, error) {
	order, err := s.cache.FindOne(ctx, id)
	if err == nil {
		if err := s.validate.Struct(order); err != nil {
			s.lg.Error("Service: failed to validate order", "err", err)
			return ordersrepository.Order{}, apperror.ErrValidate
		}
		return order, nil
	}

	order, err = s.rep.FindOne(ctx, id)
	if err != nil {
		if err := s.validate.Struct(order); err != nil {
			s.lg.Error("Service: failed to validate order", "err", err)
			return ordersrepository.Order{}, apperror.ErrValidate
		}
		return ordersrepository.Order{}, fmt.Errorf("Service: failed to get order: %w", err)
	}

	go func() {
		if _, err := s.cache.CreateOne(ctx, order); err != nil {
			s.lg.Warn("Service: failed to write order to cache", "err", err)
		}
	}()
	return order, nil
}
