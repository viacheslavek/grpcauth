package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/viacheslavek/grpcauth/auth/internal/config"
)

type Storage struct {
	pool *pgxpool.Pool
	ctx  context.Context
	log  *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, cfg config.StorageConfig) (*Storage, error) {
	const op = "storage.postgresql.new"

	postgresURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:   cfg.DBName,
	}

	log.Info("current postgres url", slog.String("url", postgresURL.String()))

	pool, err := pgxpool.New(ctx, postgresURL.String())
	if err != nil {
		return &Storage{}, fmt.Errorf("%s: failed to connect db %w", op, err)
	}

	log.Info("Postgres conn init")

	return &Storage{
		pool: pool,
		ctx:  ctx,
		log:  log,
	}, nil
}

func (s *Storage) Ping() error {
	if err := s.pool.Ping(s.ctx); err != nil {
		return fmt.Errorf("Ping is failed: %w\n", err)
	}
	s.log.Info("Postgres ping success")

	var greeting string
	err := s.pool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		return fmt.Errorf("QueryRow failed: %w\n", err)
	}

	s.log.Info("QueryRow success")

	return nil
}
