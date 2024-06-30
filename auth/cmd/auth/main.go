package main

import (
	"fmt"

	"github.com/viacheslavek/grpcauth/auth/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Printf("%+v", cfg)
	// TODO: логгер
	// TODO: приложение
	// TODO: запуск
}
