package nats

import (
	"fmt"

	"github.com/nats-io/stan.go"
)

type Config struct {
	ClusterID   string `yaml:"cluster_id"`
	ClientIDSub string `yaml:"client_id_sub"`
	ClientIDPub string `yaml:"client_id_pub"`
}

func NewClient(cfg Config, isPub bool) (stan.Conn, error) {
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
