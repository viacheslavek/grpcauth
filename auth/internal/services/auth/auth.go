package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
)

// TODO: потом реализовать интерфейсы для Storage

type Auth struct {
	log *slog.Logger
	//ownerSaver    OwnerSaver
	//ownerProvider OwnerProvider
	//appProvider   AppProvider
	tokenTTL time.Duration
}

//// OwnerSaver TODO: допилить
//type OwnerSaver interface {
//	SaveOwner()
//}
//
//// OwnerProvider TODO: допилить
//type OwnerProvider interface {
//	GetOwner()
//	UpdateOwner()
//	DeleteOwner()
//}
//
//// AppProvider TODO: допилить
//type AppProvider interface {
//	GetApp()
//}

func New(
	log *slog.Logger,
	//ownerSaver OwnerSaver,
	//ownerProvider OwnerProvider,
	//appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log: log,
		//ownerSaver:    ownerSaver,
		//ownerProvider: ownerProvider,
		//appProvider:   appProvider,
		tokenTTL: tokenTTL,
	}
}

func (a Auth) CreateOwner(ctx context.Context, owner models.Owner) error {
	//TODO implement me
	panic("implement me")
}

func (a Auth) UpdateOwner(ctx context.Context, owner models.Owner) error {
	//TODO implement me
	panic("implement me")
}

func (a Auth) DeleteOwner(ctx context.Context, owner models.Owner) error {
	//TODO implement me
	panic("implement me")
}

func (a Auth) GetOwner(ctx context.Context, owner models.Owner) (models.Owner, error) {
	//TODO implement me
	panic("implement me")
}

func (a Auth) LoginOwner(ctx context.Context, owner models.Owner, appId int) (token string, err error) {
	//TODO implement me
	panic("implement me")
}
