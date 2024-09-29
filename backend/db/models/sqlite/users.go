package sqlite

import (
	"blog/entities"
	"blog/util"
	"context"
	"database/sql"
	"fmt"
)

type Users struct{}

func NewUsers() *Users {
	return &Users{}
}

func (t *Users) Get(ctx context.Context, db *sql.DB) (*entities.User, error) {
	stmt := `SELECT * FROM users WHERE id = 0;`
	util.LogQuery(ctx, "GetUser:", stmt)

	row := db.QueryRowContext(ctx, stmt)
	if err := row.Err(); err != nil {
		return &entities.User{}, fmt.Errorf("Get: query failed: %w", err)
	}

	user := entities.User{}
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.JWT,
	)
	if err != nil {
		return &entities.User{}, fmt.Errorf("Get: row scan failed: %w", err)
	}

	return &user, nil
}

func (t *Users) UpdateJWT(ctx context.Context, tx *sql.Tx, jwt string) error {
	stmt := `
	UPDATE users 
	SET
		jwt = ?
	WHERE id = 0;
	`
	util.LogQuery(ctx, "UpdateJWT:", stmt)

	_, err := tx.ExecContext(
		ctx,
		stmt,
		jwt,
	)
	if err != nil {
		return fmt.Errorf("UpdateJWT: update query failed: %w", err)
	}

	return nil
}

func (t *Users) ClearJWT(ctx context.Context, tx *sql.Tx) error {
	stmt := `
	UPDATE users SET jwt = "" WHERE id = 0;
	`
	util.LogQuery(ctx, "ClearJWT:", stmt)

	_, err := tx.ExecContext(ctx, stmt)
	if err != nil {
		return fmt.Errorf("ClearJWT: delete error: %w", err)
	}

	return nil
}

func (t *Users) Create(ctx context.Context, tx *sql.Tx, user entities.InUser) (*entities.User, error) {
	stmt := `
	INSERT INTO users
	(
		id,
		name,
		password,
		jwt
	)
	VALUES (0,?,?,'')
	RETURNING name, password;
	`
	util.LogQuery(ctx, "CreateUser:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		user.Name,
		user.Password,
	)
	if err := row.Err(); err != nil {
		return &entities.User{}, fmt.Errorf("Create: create error: %w", err)
	}

	newUser, err := scanUser(row)
	if err != nil {
		return &entities.User{}, fmt.Errorf("Create: scan error: %w", err)
	}

	return newUser, nil
}

func (t *Users) Update(ctx context.Context, tx *sql.Tx, user entities.InUser) (*entities.User, error) {
	stmt := `
	UPDATE users
	SET
		name = ?,
		password = ?,
		jwt = ''
	WHERE id = 0
	RETURNING name, password;
	`
	util.LogQuery(ctx, "UpdateUser:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		user.Name,
		user.Password,
	)
	if err := row.Err(); err != nil {
		return &entities.User{}, fmt.Errorf("Update: update error: %w", err)
	}

	newUser, err := scanUser(row)
	if err != nil {
		return &entities.User{}, fmt.Errorf("Update: scan error: %w", err)
	}

	return newUser, nil
}

func (t *Users) Delete(ctx context.Context, tx *sql.Tx) error {
	stmt := `
	DELETE FROM users WHERE id = 0;
	`
	util.LogQuery(ctx, "DeleteUser:", stmt)

	_, err := tx.ExecContext(
		ctx,
		stmt,
	)
	if err != nil {
		return fmt.Errorf("Delete: delete user error: %w", err)
	}

	return nil
}

func scanUser(row *sql.Row) (*entities.User, error) {
	newUser := entities.User{}
	err := row.Scan(
		&newUser.Name,
		&newUser.Password,
	)
	if err != nil {
		return &entities.User{}, fmt.Errorf("scanUser: scan user failed: %w", err)
	}
	return &newUser, nil
}
