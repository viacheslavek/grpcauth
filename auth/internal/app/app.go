package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/viacheslavek/grpcauth/auth/internal/app/grpcapp"
	"github.com/viacheslavek/grpcauth/auth/internal/config"
	"github.com/viacheslavek/grpcauth/auth/internal/services/auth"
	"github.com/viacheslavek/grpcauth/auth/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
	log        *slog.Logger
}

func New(
	ctx context.Context, log *slog.Logger,
	grpcPort int, database config.StorageConfig, tokenTTL time.Duration,
) *App {
	db, errN := postgres.New(ctx, log, database)
	if errN != nil {
		log.Error("failed to init database")
		panic(errN)
	}

	if err := db.Ping(); err != nil {
		log.Error("failed to ping database")
		panic(err)
	}

	authService := auth.New(log, db, db, db, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

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
