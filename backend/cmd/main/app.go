package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"log/slog"

	"github.com/julienschmidt/httprouter"
	v1 "github.com/taaanechka/order-service/internal/api-server/api/http/v1"
	"github.com/taaanechka/order-service/internal/api-server/repositories/order/postgresql"
	"github.com/taaanechka/order-service/internal/api-server/services/orderservice"
	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
	"github.com/taaanechka/order-service/internal/config"
)

func main() {
	lg := slog.New(slog.NewTextHandler(os.Stdout, nil))
	lg.Info("create router")
	router := httprouter.New()

	lg.Info("get application configuration")
	cfg, err := config.GetConfig(lg)
	if err != nil {
		return
	} 

	rep, err := postgresql.NewRepository(lg, ordersrepository.Config(cfg.Repository))
	if err != nil {
		lg.Error("failed to init storage", "err", err)
		return
	}

	service := orderservice.NewService(lg, rep)

	lg.Info("register order handler")
	handler := v1.NewHandler(lg, service)
	handler.Register(router)

	start(lg, router, cfg)
}

func start(lg *slog.Logger, router *httprouter.Router, cfg *config.Config) {
	lg.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		lg.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			lg.Error("Cant detect app path", "err", err)
			return
		}
		lg.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		lg.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		lg.Info("server is listening unix socket", "socket_path", socketPath)
	} else {
		lg.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		lg.Info("server is listening port", "bind_ip", cfg.Listen.BindIP, "port", cfg.Listen.Port)
	}

	if listenErr != nil {
		lg.Error("failed to listen source", "err", listenErr)
		return
	}

	server := &http.Server{
		Handler: router,
	}

	if err := server.Serve(listener); err != nil {
		lg.Error("failed to serve API", "err", err)
		return
	}
}
