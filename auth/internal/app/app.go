package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/viacheslavek/grpcauth/auth/internal/app/grpcapp"
	"github.com/viacheslavek/grpcauth/auth/internal/config"
)

type App struct {
	GRPCServer *grpcapp.App
	log        *slog.Logger
}

func New(log *slog.Logger, grpcPort int, database config.StorageConfig, tokenTTL time.Duration) *App {
	// TODO: создаю storage

	// TODO: сервисный слой

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
		log:        log,
	}
}

func (a *App) GracefulStop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	a.GRPCServer.Stop()

	// TODO: сделать Stop() для БД

	a.log.Info("Gracefully stopped")
}
