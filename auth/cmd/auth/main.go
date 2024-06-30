package main

import (
	"fmt"

	"github.com/viacheslavek/grpcauth/auth/internal/config"
	"github.com/viacheslavek/grpcauth/auth/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	fmt.Printf("%+v\n", cfg)

	lg := logger.SetupLogger(cfg.Env)
	lg.Info("lalal")
	lg.Warn("pupu")

	// TODO: приложение
	// TODO: запуск
}
