package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
)

type Auth struct {
	log           *slog.Logger
	ownerSaver    OwnerSaver
	ownerProvider OwnerProvider
	appProvider   AppProvider
	tokenTTL      time.Duration
}

type OwnerSaver interface {
	SaveOwner(ctx context.Context, owner models.Owner) error
}

type OwnerProvider interface {
	GetOwner(ctx context.Context, key models.OwnerKey) (models.Owner, error)
	UpdateOwner(ctx context.Context, owner models.Owner) error
	DeleteOwner(ctx context.Context, key models.OwnerKey) error
}

type AppProvider interface {
	GetApp(ctx context.Context, appId int) models.App
}

func New(
	log *slog.Logger,
	ownerSaver OwnerSaver,
	ownerProvider OwnerProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:           log,
		ownerSaver:    ownerSaver,
		ownerProvider: ownerProvider,
		appProvider:   appProvider,
		tokenTTL:      tokenTTL,
	}
}
