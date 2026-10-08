package repo

import (
	"context"
	"errors"
	"fmt"

	errorTemplate "github.com/MGT06/EventHub_Backend.git/internal/error"
	"github.com/MGT06/EventHub_Backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (a *AuthRepo) Register(ctx context.Context, body model.Account) error {
	query := "INSERT INTO accounts (name, email, password) VALUES ($1, $2, $3)"
	args := []any{body.Name, body.Email, body.Password}

	cmt, err := a.db.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "accounts_email_key" {
			return errorTemplate.ErrEmailExist
		}
		return err
	}

	if cmt.RowsAffected() == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (a *AuthRepo) FindAccount(ctx context.Context, email string) (model.Account, error) {
	query := "SELECT id, role, password from accounts WHERE email=$1"
	args := []any{email}

	var data model.Account
	if err := a.db.QueryRow(ctx, query, args...).Scan(&data.UserId, &data.Role, &data.Password); err != nil {
		return model.Account{}, err
	}

	return data, nil
}

func (a *AuthRepo) ResetPassword(ctx context.Context, newPassword string, userId int) error {
	query := "UPDATE accounts SET password = $1 WHERE id = $2"
	args := []any{newPassword, userId}

	cmt, err := a.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmt.RowsAffected() == 0 {
		return fmt.Errorf("no row affected")
	}
	return nil
}
