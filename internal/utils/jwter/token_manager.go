package jwter

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"

	"banner-service/internal/models"
)

type Claims struct {
	UserID  int  `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
	TagID   int  `json:"tag_id"`
	jwt.StandardClaims
}

type Manager struct {
	signingKey string
}

// TODO:избавиться от синглтона

var TokenManagerSingleton *Manager

func LoadSecret(signingKey string) error {
	if signingKey == "" {
		return errors.New("empty signing key")
	}

	TokenManagerSingleton = &Manager{signingKey: signingKey}
	return nil
}

func (m *Manager) GenerateJWT(user *models.User) (string, error) {
	claims := &Claims{
		UserID:  user.UserID,
		IsAdmin: user.IsAdmin,
		TagID:   user.TagID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 60).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) ParseJWT(accessToken string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (jwtKey interface{}, err error) {
		return []byte(m.signingKey), err
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token ")
	}

	return claims, nil
}
