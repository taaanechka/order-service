package main

import (
	"log/slog"
	"os"

	"github.com/taaanechka/order-service/internal/config"
	"github.com/taaanechka/order-service/tools/publisher/gen"
	"github.com/taaanechka/order-service/tools/publisher/pub"
)

func main() {
	lg := slog.New(slog.NewTextHandler(os.Stdout, nil))
	lg.Info("Pub: get nats configuration")
	cfg, err := config.GetConfig(lg)
	if err != nil {
		return
	}

	publisher, err := pub.NewPublisher(lg, cfg.Nats)
	if err != nil {
		lg.Error("Pub: failed to create publisher", "err", err)
		return
	}

	// Generate data
	bPosData, err := gen.GeneratePositiveData()
	if err != nil {
		lg.Error("Pub: failed to generate positive data", "err", err)
		return
	}
	bNegData, err := gen.GenerateNegativeData()
	if err != nil {
		lg.Error("Pub: failed to generate negative data", "err", err)
		return
	}

	// Publish data
	if err := publisher.Publish(bPosData); err != nil {
		lg.Error("Pub: failed to publish PosData", "err", err)
		return
	}
	if err := publisher.Publish(bNegData); err != nil {
		lg.Error("Pub: failed to publish NegData", "err", err)
		return
	}
}
