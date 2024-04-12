package auth

import (
	"banner-service/internal/models"
	"context"
)

type AuthRepository interface {
	CreateUser(context.Context, *models.User) (int, error)
	ReadUserByLogin(context.Context, string) (models.User, error)
}

type AuthService interface {
	SignIn(context.Context, *models.User) error
	SignUp(context.Context, *models.User) (int, error)
}
