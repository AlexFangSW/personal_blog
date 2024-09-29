package repositories

import (
	"blog/config"
	"blog/db/models/interfaces"
	"blog/entities"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UsersRepoModels struct {
	users interfaces.UsersModel
}

func NewUsersRepoModels(
	users interfaces.UsersModel,
) *UsersRepoModels {

	return &UsersRepoModels{
		users: users,
	}
}

type Users struct {
	db     *sql.DB
	config config.DBSetting
	models UsersRepoModels
}

func NewUsers(db *sql.DB, config config.DBSetting, models UsersRepoModels) *Users {
	return &Users{
		db:     db,
		config: config,
		models: models,
	}
}

func (t *Users) Get(ctx context.Context) (*entities.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	user, err := t.models.users.Get(ctxTimeout, t.db)
	if err != nil {
		return &entities.User{}, fmt.Errorf("Get: model get user failed: %w", err)
	}

	return user, nil
}

func (t *Users) UpdateJWT(ctx context.Context, jwt string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tx, err := t.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("UpdateJWT: begin transaction failed: %w", err)
	}

	if err := t.models.users.UpdateJWT(ctxTimeout, tx, jwt); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("UpdateJWT: model update jwt rollback failed: %w", err)
		}
		return fmt.Errorf("UpdateJWT: model update jwt failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("UpdateJWT: commit failed: %w", err)
	}

	return nil
}

func (t *Users) ClearJWT(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tx, err := t.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("ClearJWT: begin transaction failed: %w", err)
	}

	if err := t.models.users.ClearJWT(ctxTimeout, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("ClearJWT: model clear jwt rollback failed: %w", err)
		}
		return fmt.Errorf("ClearJWT: model clear jwt failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ClearJWT: commit failed: %w", err)
	}

	return nil
}

func (t *Users) Create(ctx context.Context, user entities.InUser) (*entities.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tx, err := t.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return &entities.User{}, fmt.Errorf("Create: begin transaction failed: %w", err)
	}

	newUser, err := t.models.users.Create(ctxTimeout, tx, user)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.User{}, fmt.Errorf("Create: model create user rollback failed: %w", err)
		}
		return &entities.User{}, fmt.Errorf("Create: model create user failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return &entities.User{}, fmt.Errorf("Create: commit failed: %w", err)
	}

	return newUser, nil
}

func (t *Users) Update(ctx context.Context, user entities.InUser) (*entities.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tx, err := t.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return &entities.User{}, fmt.Errorf("Update: begin transaction failed: %w", err)
	}

	newUser, err := t.models.users.Update(ctxTimeout, tx, user)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.User{}, fmt.Errorf("Update: model update user rollback failed: %w", err)
		}
		return &entities.User{}, fmt.Errorf("Update: model update user failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return &entities.User{}, fmt.Errorf("Update: commit failed: %w", err)
	}

	return newUser, nil
}

func (t *Users) Delete(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tx, err := t.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("Delete: begin transaction failed: %w", err)
	}

	if err := t.models.users.Delete(ctxTimeout, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("Delete: model delete user rollback failed: %w", err)
		}
		return fmt.Errorf("Delete: model delete user failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Delete: commit failed: %w", err)
	}

	return nil
}
