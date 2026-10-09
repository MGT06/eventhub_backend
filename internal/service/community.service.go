package service

import (
	"context"
	"strings"

	"github.com/MGT06/EventHub_Backend.git/internal/dto"
	"github.com/MGT06/EventHub_Backend.git/internal/repo"
)

type CommunityService struct {
	cr *repo.CommunityRepo
}

func NewCommunityService(cr *repo.CommunityRepo) *CommunityService {
	return &CommunityService{
		cr: cr,
	}
}

func (c *CommunityService) GetCommunity(ctx context.Context, communityId int) ([]dto.Community, error) {
	res, err := c.cr.GetCommunity(ctx, communityId)
	if err != nil {
		return nil, err
	}

	data := make([]dto.Community, 0, len(res))
	for _, v := range res {
		category := strings.Split(v.Category_name, ", ")
		data = append(data, dto.Community{
			Id:             v.Id_community,
			Community_name: v.Community_name,
			Description:    v.Description,
			Image_url:      v.Image_community_url,
			Category:       category,
		})
	}

	return data, nil
}

func (c *CommunityService) GetCommunityBySearchFilter(ctx context.Context, search string, filter string) ([]dto.Community, error) {
	res, err := c.cr.GetCommunityBySearchFilter(ctx, search, filter)

	if err != nil {
		return nil, err
	}

	data := make([]dto.Community, 0, len(res))
	for _, v := range res {
		category := strings.Split(v.Category_name, ", ")
		data = append(data, dto.Community{
			Id:             v.Id_community,
			Community_name: v.Community_name,
			Description:    v.Description,
			Image_url:      v.Image_community_url,
			TotalMembers:   v.TotalMembers,
			UpcomingEvent:  v.UpcomingEvent,
			Category:      category,
		})
	}

	return data, nil
}

func (e *CommunityService) ToggleJoinCommunity(ctx context.Context, idUser int, idCommunity int) (bool, error) {
	isJoin, err := e.cr.IsJoin(ctx, idUser, idCommunity)
	if err != nil {
		return false, err
	}

	if isJoin {
		if err := e.cr.LeaveCommunity(ctx, idUser, idCommunity); err != nil {
			return false, err
		}
		return isJoin, nil
	}

	if err := e.cr.JoinCommunity(ctx, idUser, idCommunity); err != nil {
		return false, err
	}

	return isJoin, nil
}

func (c *CommunityService) GetCommunityMembers(ctx context.Context, communityId int) ([]dto.CommunityMembers, error) {
	res, err := c.cr.GetCommunityMembers(ctx, communityId)
	if err != nil {
		return nil, err
	}

	data := make([]dto.CommunityMembers, 0, len(res))
	for _, v := range res {
		data = append(data, dto.CommunityMembers{
			MemberName: v.Name,
		})
	}

	return data, nil
}

func (c *CommunityService) GetPopularCommunity(ctx context.Context) ([]dto.Community, error) {
	res, err := c.cr.GetPopularCommunity(ctx)
	if err != nil {
		return nil, err
	}

	data := make([]dto.Community, 0, len(res))
	for _, v := range res {
		category := strings.Split(v.Category_name, ", ")
		data = append(data, dto.Community{
			Id:             v.Id_community,
			Community_name: v.Community_name,
			Description:    v.Description,
			Image_url:      v.Image_community_url,
			Category:       category,
		})
	}

	return data, nil
}
