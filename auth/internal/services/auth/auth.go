package auth

import (
	"log/slog"
	"time"
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
