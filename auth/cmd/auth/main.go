package main

import (
	"github.com/viacheslavek/grpcauth/auth/internal/app"
	"github.com/viacheslavek/grpcauth/auth/internal/config"
	"github.com/viacheslavek/grpcauth/auth/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()
	lg := logger.SetupLogger(cfg.Env)

	// TODO: после реализации бизнес логики буду передавать конфиг бд
	application := app.New(lg, cfg.GRPC.Port, cfg.DB.Type, cfg.TokenTTL)

	application.GRPCServer.MustRun()
}
