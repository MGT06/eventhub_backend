package repo

import (
	"context"

	"github.com/MGT06/EventHub_Backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepo struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{
		db: db,
	}
}

func (a *AdminRepo) GetDataDashboard(ctx context.Context) (model.DashboardAdmin, error) {
	query := `SELECT
    (SELECT COUNT(*)
     FROM accounts) AS "total_users",

    (SELECT COUNT(*)
     FROM events) AS "total_events",

    (SELECT COUNT(*)
     FROM communities) AS "total_communities"`

	res := a.db.QueryRow(ctx, query)

	var data model.DashboardAdmin
	if err := res.Scan(&data.TotalUsers, &data.TotalEvents, &data.TotalCommunities); err != nil {
		return model.DashboardAdmin{}, err
	}

	return data, nil
}

func (a *AdminRepo) GetUsers(ctx context.Context, search string) ([]model.Account, error) {
	query := "SELECT name, role, status, created_at FROM accounts WHERE LOWER(name) LIKE LOWER($1)"
	args := []any{"%" + search + "%"}

	res, err := a.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var users []model.Account
	for res.Next() {
		var user model.Account
		if err := res.Scan(&user.Name, &user.Role, &user.Status, &user.Created_at); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}