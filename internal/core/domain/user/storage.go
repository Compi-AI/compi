package user

import (
	"context"
	"github.com/google/uuid"
)

type Storage interface {
	Save(ctx context.Context, user User) error

	FindOne(ctx context.Context, filter SingleFilter) (User, error)
	FindOneByUsername(ctx context.Context, username string) (User, error)
	FindOneByID(ctx context.Context, id uuid.UUID) (User, error)
	FindOneBySolanaWallet(ctx context.Context, solanaWallet string) (User, error)

	Find(ctx context.Context, filter Filter) ([]User, error)
}
