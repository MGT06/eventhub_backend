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

func (n *NotificationRepo) GetMyNotification(ctx context.Context, userId int) ([]model.Notification, error) {
	query := "SELECT n.id, n.title, n.message, n.type, n.read_at, n.created_at FROM notifications n JOIN accounts a ON n.account_id = a.id WHERE a.id = $1 ORDER BY read_at DESC"
	args := []any{userId}

	res, err := n.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var notifications []model.Notification

	for res.Next() {
		var notification model.Notification
		if err := res.Scan(&notification.NotifId, &notification.Title, &notification.Message, &notification.Type, &notification.Read_at, &notification.Created_at); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (n *NotificationRepo) ListAfterID(ctx context.Context, userId, afterId int) ([]model.Notification, error) {
	query := `SELECT id, account_id, type, title, message, read_at, created_at
		FROM notifications
		WHERE user_id = $1 AND id > $2
		ORDER BY id ASC`
	args := []any{userId, afterId}

	res, err := n.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var data []model.Notification
	for res.Next() {
		var notify model.Notification
		if err := res.Scan(&notify.NotifId, &notify.Account_id, &notify.Type, &notify.Title,
			&notify.Message, &notify.Read_at, &notify.Created_at); err != nil {
			return nil, err
		}
		data = append(data, notify)
	}

	return data, res.Err()
}

func (n *NotificationRepo) CreateNotif(ctx context.Context, notify *model.Notification) error {
	query := `INSERT INTO notifications (account_id, type, title, message, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, read_at, created_at`
	args := []any{notify.Account_id, notify.Type, notify.Title, notify.Message}

	err := n.db.QueryRow(ctx, query, args...).Scan(&notify.NotifId, &notify.Read_at, &notify.Created_at)
	if err != nil {
		return err
	}

	return nil
}
