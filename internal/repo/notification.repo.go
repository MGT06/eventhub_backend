package repo

import (
	"context"

	"github.com/MGT06/EventHub_Backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepo struct {
	db *pgxpool.Pool
}

func NewNotificationRepo(db *pgxpool.Pool) *NotificationRepo {
	return &NotificationRepo{
		db: db,
	}
}

func (n *NotificationRepo) GetMyNotification(ctx context.Context, userId int) ([]model.NotificationDetail, error) {
	query := "SELECT n.id, a.name, n.title, n.message, n.type, n.read_at FROM notifications n JOIN accounts a ON n.account_id = a.id WHERE a.id = $1 ORDER BY read_at DESC"
	args := []any{userId}

	res, err := n.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var notifications []model.NotificationDetail

	for res.Next() {
		var notification model.NotificationDetail
		if err := res.Scan(&notification.NotifId, &notification.Name, &notification.Message, &notification.Type, &notification.Read_at); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
