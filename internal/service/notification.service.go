package service

import (
	"context"
	"fmt"

	"github.com/MGT06/EventHub_Backend.git/internal/dto"
	"github.com/MGT06/EventHub_Backend.git/internal/model"
	"github.com/MGT06/EventHub_Backend.git/internal/repo"
)

type NotificationService struct {
	nr *repo.NotificationRepo
	sh Publisher
}

func NewNotificationService(nr *repo.NotificationRepo, sh Publisher) *NotificationService {
	return &NotificationService{
		nr: nr,
		sh: sh,
	}
}

func (n *NotificationService) GetMyNotification(ctx context.Context, userId int) ([]dto.Notification, error) {
	res, err := n.nr.GetMyNotification(ctx, userId)
	if err != nil {
		return nil, err
	}

	data := make([]dto.Notification, 0, len(res))
	for _, v := range res {
		data = append(data, dto.Notification{
			Id:        v.NotifId,
			Title:     v.Title,
			Message:   v.Message,
			Type:      v.Type,
			Read_at:   v.Read_at,
			CreatedAt: v.Created_at,
		})
	}

	return data, nil
}

func (n *NotificationService) ListAfterID(ctx context.Context, userId, afterId int) ([]dto.Notification, error) {
	res, err := n.nr.ListAfterID(ctx, userId, afterId)
	if err != nil {
		return nil, err
	}

	data := make([]dto.Notification, 0, len(res))
	for _, v := range res {
		data = append(data, dto.Notification{
			Id:        v.NotifId,
			UserId:    v.Account_id,
			Title:     v.Title,
			Message:   v.Message,
			Type:      v.Type,
			Read_at:   v.Read_at,
			CreatedAt: v.Created_at,
		})
	}

	return data, nil
}

func (n *NotificationService) Notify(ctx context.Context, in dto.CreateNotification) (*dto.CreateNotification, error) {
	if in.UserId == 0 || in.Title == "" || in.Message == "" {
		return nil, fmt.Errorf("notify: user_id, title, pesan wajib diisi")
	}

	if err := n.nr.CreateNotif(ctx, &model.Notification{
		Account_id: in.UserId,
		Title:      in.Title,
		Message:    in.Message,
		Type:       in.Type,
	}); err != nil {
		return nil, fmt.Errorf("notify: simpan: %w", err)
	}

	n.sh.Publish(in.UserId, dto.Notification{
		UserId:  in.UserId,
		Title:   in.Title,
		Message: in.Message,
		Type:    in.Type,
	})

	return &in, nil
}
