package ownerCtl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
	"github.com/viacheslavek/grpcauth/auth/internal/lib/jwt"
	"github.com/viacheslavek/grpcauth/auth/internal/storage"
)

func (oc OwnerCtl) CreateOwner(ctx context.Context, owner models.Owner) error {
	const op = "ownerCtl.CreateOwner"

	log := oc.log.With(
		slog.String("op", op),
		slog.String("login", owner.Login()),
	)

	log.Info("create owner")

	passwordHash, errGPH := getPasswordHash(owner.Password())
	if errGPH != nil {
		return errGPH
	}
	owner.SetPassHash(passwordHash)

	if err := oc.ownerSaver.SaveOwner(ctx, owner); err != nil {
		if errors.Is(err, storage.ErrOwnerExists) {
			return fmt.Errorf("%s: %w", op, storage.ErrOwnerExists)
		}
		return fmt.Errorf("failed to save owner %w", err)
	}

	log.Info("owner created")

	return nil
}

func (oc OwnerCtl) UpdateOwner(ctx context.Context, owner models.Owner) error {
	const op = "ownerCtl.UpdateOwner"

	log := oc.log.With(
		slog.String("op", op),
		slog.Int("id", int(owner.Id())),
	)

	log.Info("update owner")

	if len(owner.Password()) != 0 {
		passwordHash, errGPH := getPasswordHash(owner.Password())
		if errGPH != nil {
			return errGPH
		}
		owner.SetPassHash(passwordHash)
	}

	if err := oc.ownerProvider.UpdateOwner(ctx, owner); err != nil {
		if errors.Is(err, storage.ErrOwnerNotFound) {
			return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		return fmt.Errorf("failed to update owner %w", err)
	}

	log.Info("owner updated")

	return nil
}

func (oc OwnerCtl) DeleteOwner(ctx context.Context, owner models.Owner) error {
	const op = "ownerCtl.DeleteOwner"

	log := oc.log.With(
		slog.String("op", op),
		slog.String("login", owner.Login()),
		slog.Int("id", int(owner.Id())),
	)

	log.Info("delete owner")

	ownerKey := models.OwnerKey{Id: owner.Id(), Login: owner.Login()}
	if err := oc.ownerProvider.DeleteOwner(ctx, ownerKey); err != nil {
		if errors.Is(err, storage.ErrOwnerNotFound) {
			return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		return fmt.Errorf("failed to delete owner %w", err)
	}

	log.Info("owner deleted")

	return nil
}

func (oc OwnerCtl) GetOwner(ctx context.Context, owner models.Owner) (models.Owner, error) {
	const op = "ownerCtl.GetOwner"

	log := oc.log.With(
		slog.String("op", op),
		slog.String("login", owner.Login()),
		slog.Int("id", int(owner.Id())),
	)

	log.Info("get owner")

	ownerKey := models.OwnerKey{Id: owner.Id(), Login: owner.Login()}
	newOwner, errGO := oc.ownerProvider.GetOwner(ctx, ownerKey)
	if errGO != nil {
		if errors.Is(errGO, storage.ErrOwnerNotFound) {
			return models.Owner{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		return models.Owner{}, fmt.Errorf("failed to get owner %w", errGO)
	}

	log.Info("owner got")

	return newOwner, nil
}

func (oc OwnerCtl) LoginOwner(ctx context.Context, owner models.Owner, appId int) (token string, err error) {
	const op = "ownerCtl.LoginOwner"

	// TODO: разобраться с нужности ручки loginOwner
	fmt.Println("app id", appId)

	log := oc.log.With(
		slog.String("op", op),
		slog.String("login", owner.Login()),
	)

	log.Info("login owner")

	ownerKey := models.OwnerKey{Id: owner.Id(), Login: owner.Login()}
	dbOwner, errGO := oc.ownerProvider.GetOwner(ctx, ownerKey)
	if errGO != nil {
		if errors.Is(errGO, storage.ErrOwnerNotFound) {
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		return "", fmt.Errorf("%s: failed get owner %w", op, errGO)
	}

	if err = bcrypt.CompareHashAndPassword(dbOwner.PassHash(), []byte(owner.Password())); err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("owner logged in successfully")

	token, err = jwt.NewToken(owner, oc.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("%s: failed to generate token %w", op, err)
	}

	return token, nil
}

func getPasswordHash(password string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to get password hash %w", err)
	}
	return passwordHash, nil
}
