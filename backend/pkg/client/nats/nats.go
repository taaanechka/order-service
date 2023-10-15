package nats

import (
	"fmt"

	"github.com/nats-io/stan.go"
	"github.com/taaanechka/order-service/internal/config"
)

func NewClient(cfg config.NatsConfig, isPub bool) (stan.Conn, error) {
	var clientID string
	if isPub {
		clientID = cfg.ClientIDPub
	} else {
		clientID = cfg.ClientIDSub
	}

	sc, err := stan.Connect(cfg.ClusterID, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats-streaming")
	}
	return sc, nil
}
