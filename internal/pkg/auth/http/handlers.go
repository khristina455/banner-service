package http

import (
	"banner-service/internal/models"
	"banner-service/internal/pkg/auth"
	"banner-service/internal/utils/jwter"
	"banner-service/internal/utils/responser"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type AuthHandler struct {
	service auth.AuthService
	logger  *logrus.Logger
}

func NewAuthHandler(s auth.AuthService, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{s, logger}
}

func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		responser.WriteStatus(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	u := &models.User{}
	err = json.Unmarshal(body, u)
	if err != nil {
		responser.WriteStatus(w, http.StatusInternalServerError)
		return
	}

	err = ah.service.SignIn(r.Context(), u)
	if err != nil {
		fmt.Println(err, " service")
		responser.WriteStatus(w, http.StatusInternalServerError)
		return
	}

	token, err := jwter.TokenManagerSingletone.GenerateJWT(u)
	if err != nil {
		responser.WriteStatus(w, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "AccessToken", Value: token})
	responser.WriteStatus(w, http.StatusOK)
}

func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signup")
	body, err := io.ReadAll(r.Body)

	if err != nil {
		responser.WriteStatus(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	u := &models.User{}
	err = json.Unmarshal(body, u)
	if err != nil {
		responser.WriteStatus(w, http.StatusInternalServerError)
		return
	}

	u.UserId, err = ah.service.SignUp(r.Context(), u)
	if err != nil {
		fmt.Println(err)
		responser.WriteStatus(w, http.StatusInternalServerError)
		return
	}

	token, err := jwter.TokenManagerSingletone.GenerateJWT(u)
	if err != nil {
		responser.WriteStatus(w, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "AccessToken", Value: token})
	responser.WriteStatus(w, http.StatusOK)
}