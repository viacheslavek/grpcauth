package ownerCtl

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
)

type OwnerCtl struct {
	log           *slog.Logger
	ownerSaver    OwnerSaver
	ownerProvider OwnerProvider
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

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func New(
	log *slog.Logger,
	ownerSaver OwnerSaver,
	ownerProvider OwnerProvider,
	tokenTTL time.Duration,
) *OwnerCtl {
	return &OwnerCtl{
		log:           log,
		ownerSaver:    ownerSaver,
		ownerProvider: ownerProvider,
		tokenTTL:      tokenTTL,
	}
}
