package jwter

import (
	"banner-service/internal/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	UserId  int  `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
	TagId   int  `json:"tag_id"`
	jwt.StandardClaims
}

type Manager struct {
	signingKey string
}

var TokenManagerSingletone *Manager

func LoadSecret(signingKey string) error {
	if signingKey == "" {
		//return errors.New("empty signing key")
		signingKey = "kdkfkjelrug737gb"
	}

	TokenManagerSingletone = &Manager{signingKey: signingKey}
	return nil
}

func (m *Manager) GenerateJWT(user *models.User) (string, error) {
	claims := &Claims{
		UserId:  user.UserId,
		IsAdmin: user.IsAdmin,
		TagId:   user.TagId,
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
