package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
)

func (a Auth) CreateOwner(ctx context.Context, owner models.Owner) error {
	const op = "auth.CreateOwner"

	log := a.log.With(
		slog.String("op", op),
		slog.String("login", owner.Login),
	)

	log.Info("create owner")

	passwordHash, errGPH := bcrypt.GenerateFromPassword([]byte(owner.Password), bcrypt.DefaultCost)
	if errGPH != nil {
		return fmt.Errorf("failed to get password hash %w", errGPH)
	}
	owner.PassHash = passwordHash

	if err := a.ownerSaver.SaveOwner(ctx, owner); err != nil {
		return fmt.Errorf("failed to save owner %w", err)
	}

	log.Info("owner created")

	return nil
}

func (a Auth) UpdateOwner(ctx context.Context, owner models.Owner) error {
	const op = "auth.UpdateOwner"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("id", int(owner.Id)),
	)

	log.Info("update owner")

	if err := a.ownerProvider.UpdateOwner(ctx, owner); err != nil {
		return fmt.Errorf("failed to update owner %w", err)
	}

	log.Info("owner updated")

	return nil
}

func (a Auth) DeleteOwner(ctx context.Context, owner models.Owner) error {
	const op = "auth.DeleteOwner"

	log := a.log.With(
		slog.String("op", op),
		slog.String("login", owner.Login),
		slog.Int("id", int(owner.Id)),
	)

	log.Info("delete owner")

	ownerKey := models.OwnerKey{Id: owner.Id, Login: owner.Login}
	if err := a.ownerProvider.DeleteOwner(ctx, ownerKey); err != nil {
		return fmt.Errorf("failed to delete owner %w", err)
	}

	log.Info("owner deleted")

	return nil
}

func (a Auth) GetOwner(ctx context.Context, owner models.Owner) (models.Owner, error) {
	const op = "auth.GetOwner"

	log := a.log.With(
		slog.String("op", op),
		slog.String("login", owner.Login),
		slog.Int("id", int(owner.Id)),
	)

	log.Info("get owner")

	ownerKey := models.OwnerKey{Id: owner.Id, Login: owner.Login}
	newOwner, errGO := a.ownerProvider.GetOwner(ctx, ownerKey)
	if errGO != nil {
		return models.Owner{}, fmt.Errorf("failed to get owner %w", errGO)
	}

	log.Info("owner got")

	return newOwner, nil
}

func (a Auth) LoginOwner(ctx context.Context, owner models.Owner, appId int) (token string, err error) {
	const op = "auth.LoginOwner"

	log := a.log.With(
		slog.String("op", op),
		slog.String("login", owner.Login),
	)

	log.Info("login owner")

	// TODO: делаю дальше

	return "nil", nil
}
