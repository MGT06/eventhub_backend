package repo

import (
	"context"
	"fmt"

	"github.com/MGT06/EventHub_Backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) GetProfileUser(ctx context.Context, userId int) (model.Account, error) {
	query := "SELECT name, email, bio, user_location, position, avatar_url, role FROM accounts where id = $1"
	args := []any{userId}

	res := u.db.QueryRow(ctx, query, args...)

	var profile model.Account
	if err := res.Scan(&profile.Name, &profile.Email, &profile.Bio, &profile.User_location, &profile.Position, &profile.Avatar_url, &profile.Role); err != nil {
		return model.Account{}, nil
	}

	fmt.Println(profile.Name)

	return profile, nil
}

func (u *UserRepo) EditProfileUser(ctx context.Context, body model.Account, userId int) error {
	query := `UPDATE accounts 
    SET 
		name = COALESCE($1, name),
        bio = $2, 
        user_location = $3,
        position = $4, 
        avatar_url = $5 
    WHERE id = $6`
	args := []any{body.Name, body.Bio, body.User_location, body.Position, body.Avatar_url, userId}

	cmt, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmt.RowsAffected() == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (u *UserRepo) GetCurrPassword(ctx context.Context, userId int) (string, error) {
	query := "SELECT password FROM accounts WHERE id = $1"
	args := []any{userId}

	res := u.db.QueryRow(ctx, query, args...)

	var currPassword string
	if err := res.Scan(&currPassword); err != nil {
		return "", err
	}

	return currPassword, nil
}

func (u *UserRepo) ChangePassword(ctx context.Context, userId int, newPassword string) error {
	query := "UPDATE accounts SET password = $1 WHERE id = $2"
	args := []any{newPassword, userId}

	cmt, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmt.RowsAffected() == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (u *UserRepo) GetUserHeaderInformation(ctx context.Context, userId int) (model.AccountDetail, error) {
	query := `SELECT name, email, avatar_url, role, COALESCE(json_agg(DISTINCT je.event_id) FILTER (WHERE je.event_id IS NOT NULL), '[]'::json) AS "event_join", COALESCE(json_agg(DISTINCT cm.community_id) FILTER (WHERE cm.community_id IS NOT NULL), '[]'::json) AS "community_member"
		FROM accounts a
		LEFT JOIN join_event je ON je.account_id = a.id
		LEFT JOIN community_members cm ON cm.account_id = a.id
		WHERE a.id = $1
		GROUP BY a.id`
	args := []any{userId}

	res := u.db.QueryRow(ctx, query, args...)

	var userInformation model.AccountDetail
	if err := res.Scan(&userInformation.Name, &userInformation.Email, &userInformation.Avatar_url, &userInformation.Role, &userInformation.EventJoin, &userInformation.CommunityMember); err != nil {
		return model.AccountDetail{}, err
	}

	return userInformation, nil
}
