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
	query := "SELECT bio, user_location, position FROM accounts where id = $1"
	args := []any{userId}

	res := u.db.QueryRow(ctx, query, args...)

	var profile model.Account
	if err := res.Scan(&profile.Bio, &profile.User_location, &profile.Position); err != nil {
		return model.Account{}, nil
	}

	return profile, nil
}

func (u *UserRepo) EditProfileUser(ctx context.Context, body model.Account, userId int) (model.Account, error) {
	query := `UPDATE accounts 
    SET 
		name = COALESCE(NULLIF($1, ''), name),
        bio = COALESCE($2, bio), 
        user_location = COALESCE($3, user_location),
        position = COALESCE($4, position), 
        avatar_url = COALESCE($5, avatar_url) 
    WHERE id = $6 RETURNING name, bio, user_location, position, avatar_url`
	args := []any{body.Name, body.Bio, body.User_location, body.Position, body.Avatar_url, userId}

	res := u.db.QueryRow(ctx, query, args...)

	var profile model.Account
	if err := res.Scan(&profile.Name, &profile.Bio, &profile.User_location, &profile.Position, &profile.Avatar_url); err != nil {
		return model.Account{}, err
	}

	return profile, nil
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
	query := `SELECT name, email, avatar_url, role, COALESCE(ARRAY_AGG(DISTINCT je.event_id) FILTER (WHERE je.event_id IS NOT NULL), '{}'::int[]) AS "event_join", COALESCE(ARRAY_AGG(DISTINCT cm.community_id) FILTER (WHERE cm.community_id IS NOT NULL), '{}'::int[]) AS "community_member", COALESCE(ARRAY_AGG(DISTINCT se.event_id) FILTER (WHERE se.event_id IS NOT NULL), '{}'::int[]) AS "event_save"
		FROM accounts a
		LEFT JOIN join_event je ON je.account_id = a.id
		LEFT JOIN community_members cm ON cm.account_id = a.id
        LEFT JOIN saved_events se ON se.account_id = a.id
		WHERE a.id = $1
		GROUP BY a.id;`
	args := []any{userId}

	res := u.db.QueryRow(ctx, query, args...)

	var userInformation model.AccountDetail
	if err := res.Scan(&userInformation.Name, &userInformation.Email, &userInformation.Avatar_url, &userInformation.Role, &userInformation.EventJoin, &userInformation.CommunityMember, &userInformation.EventSave); err != nil {
		return model.AccountDetail{}, err
	}

	return userInformation, nil
}
