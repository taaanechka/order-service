package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/taaanechka/order-service/internal/config"
	"github.com/taaanechka/order-service/tools/publisher/gen"
	"github.com/taaanechka/order-service/tools/publisher/pub"
	"github.com/xlab/closer"
)

type GenFunc func() ([]byte, error)

func PublishData(lg *slog.Logger, done chan bool, ticker *time.Ticker,
	gFunc GenFunc, pub *pub.Publisher) {
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			lg.Info("Ticker", "time", t)

			bPosData, err := gFunc()
			if err != nil {
				lg.Error("Pub: failed to generate data", "err", err)
				return
			}

			if err := pub.Publish(bPosData); err != nil {
				lg.Error("Pub: failed to publish data", "err", err)
				return
			}
		}
	}
}

func main() {
	lg := slog.New(slog.NewTextHandler(os.Stdout, nil))
	lg.Info("Pub: get nats configuration")
	cfg, err := config.GetConfig(lg)
	if err != nil {
		return
	}

	defer closer.Close()

	publisher, err := pub.NewPublisher(lg, cfg.Nats)
	if err != nil {
		lg.Error("Pub: failed to create publisher", "err", err)
		return
	}
	closer.Bind(publisher.Drain)

	ticker := time.NewTicker(2 * time.Second)
	done := make(chan bool)

	// Publish positive data
	go PublishData(lg, done, ticker, gen.GeneratePositiveData, publisher)

	// Publish negative data
	go PublishData(lg, done, ticker, gen.GenerateNegativeData, publisher)

	closer.Hold()

	time.Sleep(21 * time.Second)
	ticker.Stop()
	done <- true
	lg.Info("Ticker stopped")
}
