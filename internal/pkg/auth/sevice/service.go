package sevice

import (
	"banner-service/internal/models"
	"banner-service/internal/pkg/auth"
	"context"
	"errors"
)

type AuthService struct {
	repo auth.Repository
}

func NewAuthService(repo auth.Repository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (as *AuthService) SignIn(ctx context.Context, user *models.User) error {
	u, err := as.repo.ReadUserByLogin(ctx, user.Login)
	if err != nil {
		return err
	}

	if u.Password == user.Password {
		user.UserID = u.UserID
		user.IsAdmin = u.IsAdmin
		user.TagID = u.TagID
		return nil
	}

	return errors.New("forbidden")
}

func (as *AuthService) SignUp(ctx context.Context, user *models.User) (int, error) {
	id, err := as.repo.CreateUser(ctx, user)
	return id, err
}
