package sevice

import (
	"banner-service/internal/models"
	"banner-service/internal/pkg/auth"
	"context"
	"fmt"
)

type AuthService struct {
	repo auth.AuthRepository
}

func NewAuthService(repo auth.AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (as *AuthService) SignIn(ctx context.Context, user *models.User) error {
	u, err := as.repo.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return err
	}

	if u.Password == user.Password {
		return nil
	}

	return fmt.Errorf("forbidden")
}

func (as *AuthService) SignUp(ctx context.Context, user *models.User) (int, error) {
	id, err := as.repo.CreateUser(ctx, user)
	return id, err
}
