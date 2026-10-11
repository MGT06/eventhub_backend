package service

import (
	"context"
	"fmt"
	"log"

	"github.com/MGT06/EventHub_Backend.git/internal/dto"
	errorTemplate "github.com/MGT06/EventHub_Backend.git/internal/error"
	"github.com/MGT06/EventHub_Backend.git/internal/model"
	"github.com/MGT06/EventHub_Backend.git/internal/repo"
	"github.com/MGT06/EventHub_Backend.git/internal/utils"
	"github.com/MGT06/EventHub_Backend.git/pkg"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	ur *repo.UserRepo
	rc *redis.Client
}

func NewUserService(ur *repo.UserRepo, rc *redis.Client) *UserService {
	return &UserService{
		ur: ur,
		rc: rc,
	}
}

func (u *UserService) GetProfileUser(ctx context.Context, userId int) (dto.UserProfile, error) {
	key := fmt.Sprintf("eventhub:profile:%d", userId)

	if result, err := utils.GetFromRedis[dto.UserProfile](ctx, u.rc, key); err != nil {
		log.Println(err)
	} else {
		return result, nil
	}

	res, err := u.ur.GetProfileUser(ctx, userId)

	if err != nil {
		return dto.UserProfile{}, err
	}

	data := dto.UserProfile{
		Bio:           res.Bio,
		User_location: res.User_location,
		Position:      res.Position,
	}

	if err := utils.SetToRedis(ctx, u.rc, key, data); err != nil {
		log.Println(err)
	}

	return data, err
}

func (u *UserService) EditProfileUser(ctx context.Context, body dto.EditUserProfile, userId int, AvaPath *string) error {
	if err := u.ur.EditProfileUser(ctx, model.Account{
		Name:          body.Name,
		Bio:           body.Bio,
		User_location: body.User_location,
		Position:      body.Position,
		Avatar_url:    AvaPath,
	}, userId); err != nil {
		return err
	}

	key := fmt.Sprintf("eventhub:profile:%d", userId)

	mes := utils.DelFromRedis(ctx, u.rc, key)
	log.Println(mes)

	return nil
}

func (u *UserService) ChangePassword(ctx context.Context, userId int, body dto.ChangePassword) error {
	if len(body.CurrentPassword) == 0 && len(body.NewPassword) == 0 {
		return errorTemplate.ErrInvalidInputs
	}

	pass, err := u.ur.GetCurrPassword(ctx, userId)
	if err != nil {
		return errorTemplate.ErrEmailPasswordIncorrect
	}

	pkg.Compare(body.CurrentPassword, pass)

	hash := pkg.NewHashConfig().GenHash(body.NewPassword)

	if err := u.ur.ChangePassword(ctx, userId, hash); err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUserHeaderInformation(ctx context.Context, userId int) (dto.UserHeaderInfo, error) {
	res, err := u.ur.GetUserHeaderInformation(ctx, userId)
	if err != nil {
		return dto.UserHeaderInfo{}, nil
	}

	data := dto.UserHeaderInfo{
		Name:            res.Name,
		Email:           res.Email,
		Avatar_url:      res.Avatar_url,
		Role:            res.Role,
		EventJoin:       res.EventJoin,
		CommunityMember: res.CommunityMember,
	}

	return data, nil
}
