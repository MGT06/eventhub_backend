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
	query := "SELECT id, name, email, role, status, created_at FROM accounts WHERE LOWER(name) LIKE LOWER($1)"
	args := []any{"%" + search + "%"}

	res, err := a.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var users []model.Account
	for res.Next() {
		var user model.Account
		if err := res.Scan(&user.UserId, &user.Name, &user.Email, &user.Role, &user.Status, &user.Created_at); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (a *AdminRepo) GetEvents(ctx context.Context) ([]model.EventDetail, error) {
	query := `SELECT e.id, e.title, e.image_event_url, e.start_at, le.city, e.capacity,
       COUNT(DISTINCT je.account_id) AS attendees
	FROM events e
	JOIN accounts a ON e.organizer_id = a.id
	LEFT JOIN join_event je ON je.event_id = e.id
	JOIN location_event le ON e.location_event_id = le.id
	GROUP BY e.id, a.id, le.city;`

	res, err := a.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var events []model.EventDetail
	for res.Next() {
		var event model.EventDetail
		if err := res.Scan(&event.Id_event, &event.Title, &event.Image_event_url, &event.Start_at, &event.City, &event.Capacity, &event.Attendees); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (a *AdminRepo) GetCommunities(ctx context.Context) ([]model.CommunityDetail, error) {
	query := `SELECT 
    c.id, 
    c.community_name, 
    c.image_community_url, 
    c.is_active,
    COUNT(DISTINCT cm.account_id) AS total_members,
    COUNT(DISTINCT e.id) FILTER (WHERE e.start_at > NOW()) AS upcoming_events_count
FROM communities c 
LEFT JOIN community_members cm ON cm.community_id = c.id
LEFT JOIN events e ON e.community_id = c.id
GROUP BY c.id, c.community_name, c.description, c.image_community_url;`

	res, err := a.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var communities []model.CommunityDetail
	for res.Next() {
		var community model.CommunityDetail
		if err := res.Scan(&community.Id_community, &community.Community_name, &community.Image_community_url, &community.Is_active, &community.TotalMembers, &community.UpcomingEvent); err != nil {
			return nil, err
		}
		communities = append(communities, community)
	}

	return communities, nil
}
