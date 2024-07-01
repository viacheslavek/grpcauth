package main

import (
	"context"

	"github.com/viacheslavek/grpcauth/auth/internal/app"
	"github.com/viacheslavek/grpcauth/auth/internal/config"
	"github.com/viacheslavek/grpcauth/auth/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()
	lg := logger.SetupLogger(cfg.Env)
	ctx := context.Background()

	application := app.New(ctx, lg, cfg.GRPC.Port, cfg.DB, cfg.TokenTTL)

	go func() {
		application.GRPCServer.MustRun()
	}()

	application.GracefulStop()
}
