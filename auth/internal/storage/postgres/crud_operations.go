package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
	"github.com/viacheslavek/grpcauth/auth/internal/storage"
)

func (s *Storage) SaveOwner(ctx context.Context, owner models.Owner) error {
	const op = "postgres.saveOwner"

	queryInsert := `
		INSERT INTO owners (email, login, password_hash)
		VALUES ($1, $2, $3)
    `

	_, err := s.pool.Exec(ctx, queryInsert, owner.Email, owner.Login, owner.PassHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("%s: failed to save owner: %w", op, storage.ErrOwnerExists)
			}
		}
		return fmt.Errorf("%s: failed to save owner: %w", op, err)
	}

	s.log.Info("Owner created successfully",
		slog.String("email", owner.Email),
		slog.String("login", owner.Login),
	)

	return nil
}

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
	err := s.pool.QueryRow(ctx, query, id).Scan(&owner.Id, &owner.Email, &owner.Login, &owner.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Owner{}, fmt.Errorf("%w with id %d ", storage.ErrOwnerNotFound, id)
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
	err := s.pool.QueryRow(ctx, query, login).Scan(&owner.Id, &owner.Email, &owner.Login, &owner.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Owner{}, fmt.Errorf("%w with login %s", storage.ErrOwnerNotFound, login)
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
	args := make([]interface{}, 0)
	argId := 1

	if owner.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email=$%d", argId))
		args = append(args, owner.Email)
		argId++
	}
	if owner.Login != "" {
		setClauses = append(setClauses, fmt.Sprintf("login=$%d", argId))
		args = append(args, owner.Login)
		argId++
	}
	if len(owner.PassHash) > 0 {
		setClauses = append(setClauses, fmt.Sprintf("password_hash=$%d", argId))
		args = append(args, owner.PassHash)
		argId++
	}

	query := fmt.Sprintf(`
        UPDATE owners 
        SET %s
        WHERE id=$%d
    `, strings.Join(setClauses, ", "), argId)
	args = append(args, owner.Id)

	result, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update owner: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w with id %d", storage.ErrOwnerNotFound, owner.Id)
	}

	s.log.Info("Owner updated successfully", "id", owner.Id)

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
	commandTag, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete owner by id: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%w with id %d", storage.ErrOwnerNotFound, id)
	}

	s.log.Info("Owner deleted successfully by id", slog.Int64("id", id))

	return nil
}

func (s *Storage) deleteOwnerByLogin(ctx context.Context, login string) error {
	query := `DELETE FROM owners WHERE login=$1`
	commandTag, err := s.pool.Exec(ctx, query, login)
	if err != nil {
		return fmt.Errorf("failed to delete owner by login: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%w with login %s", storage.ErrOwnerNotFound, login)
	}

	s.log.Info("Owner deleted successfully by login", slog.String("login", login))

	return nil
}
