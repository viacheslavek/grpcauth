package postgres

import (
	"context"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
)

func (s *Storage) SaveOwner(ctx context.Context, owner models.Owner) error {
	//TODO implement me
	panic("implement me")
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
