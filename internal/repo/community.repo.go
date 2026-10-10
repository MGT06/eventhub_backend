package repo

import (
	"context"
	"fmt"

	"github.com/MGT06/EventHub_Backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunityRepo struct {
	db *pgxpool.Pool
}

func NewCommunityRepo(db *pgxpool.Pool) *CommunityRepo {
	return &CommunityRepo{
		db: db,
	}
}

func (c *CommunityRepo) GetCommunity(ctx context.Context, communityId int) (model.CommunityDetail, error) {
	query := `SELECT 
    c.id, 
    c.community_name, 
    c.description, 
    c.image_community_url, 
    COUNT(DISTINCT cm.account_id) AS total_members,
    COUNT(DISTINCT e.id) FILTER (WHERE e.start_at > NOW()) AS upcoming_events_count,
    STRING_AGG(DISTINCT cg.category_name, ', ') AS category
FROM communities c 
LEFT JOIN community_categories cc ON c.id = cc.community_id
LEFT JOIN categories cg ON cg.id = cc.category_id
LEFT JOIN community_members cm ON cm.community_id = c.id
LEFT JOIN events e ON e.community_id = c.id
WHERE c.id = $1
GROUP BY c.id, c.community_name, c.description, c.image_community_url;`

	args := []any{communityId}
	res := c.db.QueryRow(ctx, query, args...)


	var community model.CommunityDetail
	if err := res.Scan(&community.Id_community, &community.Community_name, &community.Description, &community.Image_community_url, &community.TotalMembers, &community.UpcomingEvent, &community.Category_name); err != nil {
		return model.CommunityDetail{}, err
	}

	return community, nil
}

func (c *CommunityRepo) GetUpComingEventCommunity(ctx context.Context, communityId int) ([]model.EventDetail, error) {
	query := `SELECT e.id, e.title, e.description, e.image_event_url, e.start_at, le.city, e.capacity,
       COUNT(DISTINCT je.account_id) AS attendees,
       STRING_AGG(DISTINCT c.category_name, ', ') AS category
	FROM events e
	JOIN accounts a ON e.organizer_id = a.id
	LEFT JOIN join_event je ON je.event_id = e.id
	JOIN location_event le ON e.location_event_id = le.id
	LEFT JOIN event_categories ec ON ec.event_id = e.id
	JOIN categories c ON ec.category_id = c.id
	WHERE e.community_id = $1 AND e.start_at > NOW()
	GROUP BY e.id, a.id, le.city;`

	args := []any{communityId}

	res, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var events []model.EventDetail

	for res.Next() {
		var event model.EventDetail
		if err := res.Scan(&event.Id_event, &event.Title, &event.Description, &event.Image_event_url, &event.Start_at, &event.City, &event.Capacity, &event.Attendees, &event.Category_name); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (c *CommunityRepo) GetCommunityBySearchFilter(ctx context.Context, search string, filter string) ([]model.CommunityDetail, error) {
	query := `SELECT 
    c.id, 
    c.community_name, 
    c.description, 
    c.image_community_url, 
    COUNT(DISTINCT cm.account_id) AS total_members,
    COUNT(DISTINCT e.id) FILTER (WHERE e.start_at > NOW()) AS upcoming_events_count,
    STRING_AGG(DISTINCT cg.category_name, ', ') AS category
FROM communities c 
LEFT JOIN community_categories cc ON c.id = cc.community_id
LEFT JOIN categories cg ON cg.id = cc.category_id
LEFT JOIN community_members cm ON cm.community_id = c.id
LEFT JOIN events e ON e.community_id = c.id
WHERE c.community_name ILIKE $1
GROUP BY c.id, c.community_name, c.description, c.image_community_url
HAVING STRING_AGG(cg.category_name, ', ') ILIKE $2`

	args := []any{"%" + search + "%", "%" + filter + "%"}

	res, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var communities []model.CommunityDetail

	for res.Next() {
		var community model.CommunityDetail
		if err := res.Scan(&community.Id_community, &community.Community_name, &community.Description, &community.Image_community_url, &community.TotalMembers, &community.UpcomingEvent, &community.Category_name); err != nil {
			return nil, err
		}
		communities = append(communities, community)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return communities, nil
}

func (e *CommunityRepo) JoinCommunity(ctx context.Context, idUser int, idCommunity int) error {
	query := "INSERT INTO community_members (account_id, community_id) VALUES ($1, $2)"
	args := []any{idUser, idCommunity}

	cmt, err := e.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmt.RowsAffected() == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (e *CommunityRepo) LeaveCommunity(ctx context.Context, idUser int, idCommunity int) error {
	query := "DELETE FROM community_members WHERE account_id = $1 AND community_id = $2"
	args := []any{idUser, idCommunity}

	cmt, err := e.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmt.RowsAffected() == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (e *CommunityRepo) IsJoin(ctx context.Context, idUser int, idCommunity int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM community_members WHERE account_id = $1 AND community_id = $2)"
	args := []any{idUser, idCommunity}

	res := e.db.QueryRow(ctx, query, args...)

	var isJoin bool
	if err := res.Scan(&isJoin); err != nil {
		return false, err
	}

	return isJoin, nil
}

func (c *CommunityRepo) GetCommunityMembers(ctx context.Context, communityId int) ([]model.CommunityDetail, error) {
	query := "SELECT a.name FROM community_members cm JOIN accounts a ON a.id = cm.account_id WHERE cm.community_id = $1"
	args := []any{communityId}

	res, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var communityMembers []model.CommunityDetail
	for res.Next() {
		var member model.CommunityDetail
		if err := res.Scan(&member.Name); err != nil {
			return nil, err
		}

		communityMembers = append(communityMembers, member)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return communityMembers, nil
}

func (c *CommunityRepo) GetPopularCommunity(ctx context.Context) ([]model.CommunityDetail, error) {
	query := `SELECT c.id, c.community_name, c.description, c.image_community_url, STRING_AGG(cg.category_name, ', ') AS "category"
	FROM communities c
	JOIN community_categories cc ON cc.community_id = c.id
	JOIN categories cg ON cg.id = cc.category_id
	JOIN community_members cm ON cm.community_id = c.id
	GROUP BY  c.id, c.community_name, c.description, c.image_community_url
	ORDER BY COUNT(cm.account_id) DESC
	LIMIT 4;`

	res, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var communities []model.CommunityDetail

	for res.Next() {
		var community model.CommunityDetail
		if err := res.Scan(&community.Id_community, &community.Community_name, &community.Description, &community.Image_community_url, &community.Category_name); err != nil {
			return nil, err
		}
		communities = append(communities, community)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return communities, nil
}
