package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MGT06/EventHub_Backend.git/internal/dto"
	errorTemplate "github.com/MGT06/EventHub_Backend.git/internal/error"
	"github.com/MGT06/EventHub_Backend.git/internal/model"
	"github.com/MGT06/EventHub_Backend.git/internal/repo"
	"github.com/MGT06/EventHub_Backend.git/internal/utils"
	"github.com/MGT06/EventHub_Backend.git/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	ar *repo.AuthRepo
	rc *redis.Client
	ns *NotificationService
}

func NewAuthService(ar *repo.AuthRepo, rc *redis.Client, ns *NotificationService) *AuthService {
	return &AuthService{
		ar: ar,
		rc: rc,
		ns: ns,
	}
}

func (a *AuthService) Register(ctx context.Context, body dto.Register) error {
	if len(body.FullName) == 0 || len(body.Email) == 0 || len(body.Password) == 0 {
		return errorTemplate.ErrInvalidInputs
	}

	_, err := a.ar.FindAccount(ctx, body.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	hash := pkg.NewHashConfig().GenHash(body.Password)

	id, err := a.ar.Register(ctx, model.Account{
		Name:     body.FullName,
		Email:    body.Email,
		Password: hash,
	})
	if err != nil {
		return err
	}

	a.ns.Notify(ctx, dto.CreateNotification{
		UserId:  id,
		Title:   "Registration confirmed",
		Message: "You're registered in this app",
		Type:    "registration",
	})

	return nil
}

func (a *AuthService) Login(ctx context.Context, body dto.Login) (string, error) {
	if len(body.Email) == 0 || len(body.Password) == 0 {
		return "", errorTemplate.ErrInvalidInputs
	}

	acc, err := a.ar.FindAccount(ctx, body.Email)
	if err != nil {
		return "", errorTemplate.ErrEmailPasswordIncorrect
	}

	if err := pkg.Compare(body.Password, acc.Password); err != nil {
		return "", errorTemplate.ErrEmailPasswordIncorrect
	}

	claims := pkg.NewJWTClaims(acc.UserId, acc.Role)
	return claims.GenToken()
}

func (a *AuthService) Logout(ctx context.Context, userId int, JTI string, Expired time.Duration) error {
	key := fmt.Sprintf("eventhub:tokenBlacklist:%s", JTI)

	if err := utils.SetToRedis(ctx, a.rc, key, userId, Expired); err != nil {
		return err
	}

	return nil
}

func (a *AuthService) ForgotPassword(ctx context.Context, body dto.RequestForgotPassword) (dto.ResponseForgotPassword, error) {
	acc, err := a.ar.FindAccount(ctx, body.Email)
	if err != nil {
		return dto.ResponseForgotPassword{}, err
	}

	data := dto.ResponseForgotPassword{
		UserId: acc.UserId,
	}

	return data, nil
}

func (a *AuthService) ResetPassword(ctx context.Context, body dto.ResetPassword) error {

	hash := pkg.NewHashConfig().GenHash(body.NewPassword)

	if err := a.ar.ResetPassword(ctx, hash, body.UserId); err != nil {
		return err
	}

	return nil
}
