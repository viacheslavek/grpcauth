package app

import (
	"log/slog"
	"time"

	"github.com/viacheslavek/grpcauth/auth/internal/app/grpcapp"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, database string, tokenTTL time.Duration) *App {
	// TODO: создаю storage

	// TODO: сервисный слой

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
