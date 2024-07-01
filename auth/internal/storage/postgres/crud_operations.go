package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
	"github.com/viacheslavek/grpcauth/auth/internal/storage"
)

func (s *Storage) SaveOwner(ctx context.Context, owner models.Owner) error {
	const op = "postgres.saveOwner"
	var existingEmail, existingLogin sql.NullString
	querySearch := `
		SELECT email, login
		FROM owners 
		WHERE email=$1 OR login=$2
    `
	err := s.conn.QueryRow(ctx, querySearch, owner.Email, owner.Login).Scan(&existingEmail, &existingLogin)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%s: failed to check email and login: %w", op, err)
	}

	if existingEmail.Valid && existingEmail.String == owner.Email ||
		existingLogin.Valid && existingLogin.String == owner.Login {
		return storage.ErrOwnerExists
	}

	queryInsert := `
		INSERT INTO owners (email, login, password_hash)
		VALUES ($1, $2, $3)
    `

	_, err = s.conn.Exec(ctx, queryInsert, owner.Email, owner.Login, owner.PassHash)
	if err != nil {
		return fmt.Errorf("%s: failed to save owner: %w", op, err)
	}

	s.log.Info("Owner created successfully",
		slog.String("email", owner.Email),
		slog.String("login", owner.Login),
	)

	return nil
}

func (s *Storage) GetOwner(ctx context.Context, key models.OwnerKey) (models.Owner, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) UpdateOwner(ctx context.Context, owner models.Owner) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) DeleteOwner(ctx context.Context, key models.OwnerKey) error {
	//TODO implement me
	panic("implement me")
}
