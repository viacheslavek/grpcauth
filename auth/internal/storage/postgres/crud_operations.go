package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

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

// TODO: делал достаточно второпях, надо будет найти узкие места и их улучшить
// Например, можно ли дублировать меньше кода в update и лучше обрабатывать ошибки -> в техдолг

func (s *Storage) GetOwner(ctx context.Context, key models.OwnerKey) (models.Owner, error) {
	if key.Id != 0 {
		return s.getOwnerById(ctx, key.Id)
	} else if key.Login != "" {
		return s.getOwnerByLogin(ctx, key.Login)
	}
	return models.Owner{}, fmt.Errorf("unattainable error: either id or login must be provided")
}

func (s *Storage) getOwnerById(ctx context.Context, id int64) (models.Owner, error) {
	var owner models.Owner
	query := `
		SELECT id, email, login, password_hash 
		FROM owners
		WHERE id=$1
	`
	err := s.conn.QueryRow(ctx, query, id).Scan(&owner.Id, &owner.Email, &owner.Login, &owner.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Owner{}, fmt.Errorf("owner with id %d not found", id)
		}
		return models.Owner{}, fmt.Errorf("failed to get owner by id: %w", err)
	}

	s.log.Info("Owner retrieved successfully by id",
		slog.Int64("id", owner.Id),
		slog.String("email", owner.Email),
		slog.String("login", owner.Login),
	)

	return owner, nil
}

func (s *Storage) getOwnerByLogin(ctx context.Context, login string) (models.Owner, error) {
	var owner models.Owner
	query := `
		SELECT id, email, login, password_hash
		FROM owners WHERE
		login=$1
	`
	err := s.conn.QueryRow(ctx, query, login).Scan(&owner.Id, &owner.Email, &owner.Login, &owner.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Owner{}, fmt.Errorf("owner with login %s not found", login)
		}
		return models.Owner{}, fmt.Errorf("failed to get owner by login: %w", err)
	}

	s.log.Info("Owner retrieved successfully by login",
		slog.Int64("id", owner.Id),
		slog.String("email", owner.Email),
		slog.String("login", owner.Login),
	)

	return owner, nil
}

func (s *Storage) UpdateOwner(ctx context.Context, owner models.Owner) error {
	setClauses := make([]string, 0)

	if owner.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email=$%s", owner.Email))
	}
	if owner.Login != "" {
		setClauses = append(setClauses, fmt.Sprintf("login=$%s", owner.Login))
	}
	if len(owner.PassHash) > 0 {
		setClauses = append(setClauses, fmt.Sprintf("password_hash=$%s", owner.PassHash))
	}

	query := fmt.Sprintf(`
		UPDATE owners 
		SET %s
		WHERE id=$%d
	`, strings.Join(setClauses, ", "), owner.Id)

	_, err := s.conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to update owner: %w", err)
	}

	s.log.Info("Owner updated successfully", slog.Int64("id", owner.Id))

	return nil
}

func (s *Storage) DeleteOwner(ctx context.Context, key models.OwnerKey) error {
	if key.Id != 0 {
		return s.deleteOwnerById(ctx, key.Id)
	} else if key.Login != "" {
		return s.deleteOwnerByLogin(ctx, key.Login)
	}
	return fmt.Errorf("either id or login must be provided")
}

func (s *Storage) deleteOwnerById(ctx context.Context, id int64) error {
	query := `DELETE FROM owners WHERE id=$1`
	commandTag, err := s.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete owner by id: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("owner with id %d not found", id)
	}

	s.log.Info("Owner deleted successfully by id", slog.Int64("id", id))

	return nil
}

func (s *Storage) deleteOwnerByLogin(ctx context.Context, login string) error {
	query := `DELETE FROM owners WHERE login=$1`
	commandTag, err := s.conn.Exec(ctx, query, login)
	if err != nil {
		return fmt.Errorf("failed to delete owner by login: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("owner with login %s not found", login)
	}

	s.log.Info("Owner deleted successfully by login", slog.String("login", login))

	return nil
}
