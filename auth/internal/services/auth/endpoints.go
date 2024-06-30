package auth

import (
	"context"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
)

// TODO: реализую бизнес логику

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
