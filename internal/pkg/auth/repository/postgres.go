package repository

import (
	"banner-service/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	createUser     = `INSERT INTO "user" (login, password, tag_id) VALUES ($1, $2, $3) RETURNING user_id;`
	getUserByLogin = `SELECT user_id, password, is_admin, tag_id FROM "user" WHERE login=$1`
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (ar *AuthRepository) CreateUser(ctx context.Context, user *models.User) (int, error) {
	var id int
	err := ar.db.QueryRow(ctx, createUser,
		user.Login, user.Password, user.TagId).Scan(&id)

	if err != nil {
		err = fmt.Errorf("error happened in scan.Scan: %w", err)

		return 0, err
	}

	return id, nil
}

// TODO:сделать

func (ar *AuthRepository) ReadUserByLogin(ctx context.Context, login string) (models.User, error) {
	u := models.User{}
	err := ar.db.QueryRow(ctx, getUserByLogin, login).Scan(&u.UserId, &u.Password, &u.IsAdmin, &u.TagId)

	if err != nil {
		err = fmt.Errorf("error happened in scan.Scan: %w", err)

		return models.User{}, err
	}

	return u, nil
}
