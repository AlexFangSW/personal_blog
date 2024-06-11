package interfaces

import (
	"blog/entities"
	"context"
	"database/sql"
)

// Concrete implementations are at db/models/<db name>/
type UsersModel interface {
	Get(ctx context.Context, db *sql.DB) (*entities.User, error)
	UpdateJWT(ctx context.Context, tx *sql.Tx, jwt string) error
	ClearJWT(ctx context.Context, tx *sql.Tx) error

	Create(ctx context.Context, tx *sql.Tx, user entities.InUser) (*entities.User, error)
	Update(ctx context.Context, tx *sql.Tx, user entities.InUser) (*entities.User, error)
	Delete(ctx context.Context, tx *sql.Tx) error
}
