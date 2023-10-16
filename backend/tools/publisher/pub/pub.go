package pub

import (
	"fmt"
	"log/slog"

	"github.com/nats-io/stan.go"
	"github.com/taaanechka/order-service/pkg/client/nats"
)

type Publisher struct {
	sconn stan.Conn
	lg    *slog.Logger
}

func NewPublisher(lg *slog.Logger, cfg nats.Config) (*Publisher, error) {
	lg.Info("Pub: new publisher")
	sconn, err := nats.NewClient(cfg, true)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	return &Publisher{
		sconn: sconn,
		lg:    lg,
	}, nil
}

func (p *Publisher) Publish(byteData []byte) error {
	p.lg.Info("Pub: publish data")
	if err := p.sconn.Publish("orders", byteData); err != nil {
		p.lg.Error("Pub: failed to publish data", "err", err)
		return err
	}
	return nil
}
