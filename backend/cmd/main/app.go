package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/taaanechka/order-service/internal/api-server/api/http/v1"
	"github.com/taaanechka/order-service/internal/api-server/repositories/order/postgresql"
	"github.com/taaanechka/order-service/internal/api-server/services/orderservice"
	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
	"github.com/taaanechka/order-service/internal/config"
)

func main() {
	InfoLg := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	InfoLg.Println("create router")
	router := httprouter.New()

	cfg := config.GetConfig()
	errorLg := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	rep, err := postgresql.NewRepository(InfoLg, ordersrepository.Config(cfg.Repository))
	if err != nil {
		errorLg.Printf("failed to init storage: %v", err)
		return
	}

	service := orderservice.NewService(InfoLg, rep)

	InfoLg.Println("register user handler")
	handler := v1.NewHandler(InfoLg, service)
	handler.Register(router)

	start(InfoLg, errorLg, router, cfg)
}

func start(InfoLg *log.Logger, errorLg *log.Logger, router *httprouter.Router, cfg *config.Config) {
	InfoLg.Println("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		InfoLg.Println("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			errorLg.Fatal(err)
		}
		InfoLg.Println("create socket")
		socketPath := path.Join(appDir, "app.sock")

		InfoLg.Println("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		InfoLg.Printf("server is listening unix socket: %s", socketPath)
	} else {
		InfoLg.Println("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		InfoLg.Printf("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		errorLg.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
	}

	errorLg.Fatal(server.Serve(listener))
}
