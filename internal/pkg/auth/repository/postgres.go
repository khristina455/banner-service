package repository

import (
	"banner-service/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	createUser = `INSERT INTO "user" (login, password, is_admin, tag_id) VALUES ($1, $2, $3, $4) RETURNING id;`
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (ar *AuthRepo) CreateUser(ctx context.Context, user *models.User) (int, error) {
	var id int
	err := ar.db.QueryRow(ctx, createUser,
		user.Login, user.Password, user.IsAdmin, user.TagId).Scan(&id)

	if err != nil {
		err = fmt.Errorf("error happened in scan.Scan: %w", err)

		return 0, err
	}

	return id, nil
}

// TODO:сделать

func (ar *AuthRepo) GetUserByLogin(context.Context, string) (*models.User, error) {
	return &models.User{}, nil
}
