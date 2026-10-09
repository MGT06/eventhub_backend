package service

import (
	"context"

	"github.com/MGT06/EventHub_Backend.git/internal/dto"
	"github.com/MGT06/EventHub_Backend.git/internal/repo"
)

type AdminService struct {
	ar *repo.AdminRepo
}

func NewAdminService(ar *repo.AdminRepo) *AdminService {
	return &AdminService{
		ar: ar,
	}
}

func (a *AdminService) GetDataDashboard(ctx context.Context) (dto.DashboardAdmin, error) {
	res, err := a.ar.GetDataDashboard(ctx)
	if err != nil {
		return dto.DashboardAdmin{}, err
	}

	data := dto.DashboardAdmin{
		TotalUsers:       res.TotalUsers,
		TotalEvents:      res.TotalEvents,
		TotalCommunities: res.TotalCommunities,
	}

	return data, nil
}

func (a *AdminService) GetUsers(ctx context.Context, search string) ([]dto.AllUser, error) {
	res, err := a.ar.GetUsers(ctx, search)
	if err != nil {
		return nil, err
	}

	data := make([]dto.AllUser, 0, len(res))
	for _, v := range res {
		data = append(data, dto.AllUser{
			Name:   v.Name,
			Role:   v.Role,
			Status: v.Status,
			Joined: v.Created_at,
		})
	}

	return data, nil
}

func (a *AdminService) GetEvents(ctx context.Context) ([]dto.Event, error) {
	res, err := a.ar.GetEvents(ctx)
	if err != nil {
		return nil, err
	}

	data := make([]dto.Event, 0, len(res))
	for _, v := range res {
		data = append(data, dto.Event{
			Id: v.Id_event,
			Title: v.Title,
			Image_url: v.Image_event_url,
			Start_at: v.Start_at,
			Location: v.City,
			Capacity: v.Capacity,
			Attendees: v.Attendees,
		})
	}

	return data, nil
}

func (a *AdminService) GetCommunities(ctx context.Context) ([]dto.Community, error) {
	res, err := a.ar.GetCommunities(ctx)
	if err != nil {
		return nil, err
	}

	data := make([]dto.Community, 0, len(res))
	for _, v := range res {
		data = append(data, dto.Community{
			Id: v.Id_community,
			Community_name: v.Community_name,
			Image_url: v.Image_community_url,
			IsActive: v.Is_active,
			TotalMembers: v.TotalMembers,
			UpcomingEvent: v.UpcomingEvent,
		})
	}

	return data, nil
}