package middleware

import (
	"banner-service/internal/utils/jwter"
	"banner-service/internal/utils/responser"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func Auth(log *logrus.Logger, onlyAdmin bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("AccessToken")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				log.Debug("token cookie not found", err)
				responser.WriteStatus(w, http.StatusUnauthorized)
				return
			default:
				log.Error("faild to get token cookie", err)
				responser.WriteStatus(w, http.StatusUnauthorized)
				return
			}
		}

		claims, err := jwter.TokenManagerSingleton.ParseJWT(tokenCookie.Value)
		if err != nil {
			log.Error("jws token is invalid auth ", err)
			responser.WriteStatus(w, http.StatusUnauthorized)
			return
		}

		if onlyAdmin && !claims["is_admin"].(bool) {
			responser.WriteStatus(w, http.StatusForbidden)
			return
		}

		tagIdStr := r.URL.Query().Get("tag_id")
		if tagIdStr != "" && !claims["is_admin"].(bool) {
			tagId, err := strconv.Atoi(tagIdStr)
			if err == nil {
				if tagId != int(claims["tag_id"].(float64)) {
					responser.WriteStatus(w, http.StatusForbidden)
					return
				}
			}
		}

		ctx := context.WithValue(r.Context(), "user_id", int(claims["user_id"].(float64)))
		ctx = context.WithValue(ctx, "is_admin", claims["is_admin"].(bool))
		ctx = context.WithValue(ctx, "tag_id", int(claims["tag_id"].(float64)))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
