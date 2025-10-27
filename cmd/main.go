package main

import (
	"github.com/Asylann/Auth/internal/config"
	"github.com/Asylann/Auth/internal/route"
	logger2 "github.com/Asylann/Auth/lib/logger"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"os/signal"
	"syscall"
)

func main() {
	logger := logger2.NewLogger()
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("No .env variables are loaded: %s", err.Error())
		return
	}

	cfg := config.LoadConfig()
	server := route.NewRoute(cfg, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		server.Run()
	}()

	<-ctx.Done()
	stop()

	logger.Warn("Server is shutting down...")
	server.GracefullyShutDown()
}
