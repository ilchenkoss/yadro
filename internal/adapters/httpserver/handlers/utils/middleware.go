package utils

import (
	"errors"
	"log/slog"
	"myapp/internal/core/domain"
	"myapp/internal/core/port"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationType      = "bearer"
)

func AuthMiddleware(requiredRoles map[domain.UserRole]bool, ts port.TokenService, ur port.UserRepository, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authorizationHeaderKey)
		if len(authHeader) == 0 {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != authorizationType {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		accessToken := fields[1]
		token, tErr := ts.GetTokenByString(accessToken)
		if tErr != nil {
			switch {
			case errors.Is(tErr, domain.ErrTokenExpired):
				http.Error(w, "token expired", http.StatusUnauthorized)
				return
			default:
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
		}

		userLogin, lErr := ts.GetUserByToken(token)
		if lErr != nil {
			http.Error(w, "auth failed", http.StatusInternalServerError)
			return
		}

		user, rguErr := ur.GetUserByLogin(userLogin)
		if rguErr != nil {
			switch {
			case errors.Is(rguErr, domain.ErrUserNotFound):
				slog.Error("Error attempt to log in using a non-existent login. Token: ", fields[1])
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			default:
				http.Error(w, "auth failed", http.StatusInternalServerError)
				return
			}
		}

		if !requiredRoles[user.Role] {
			http.Error(w, "insufficient permissions", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func OrdinaryMiddleware(f http.HandlerFunc, ts port.TokenService, ur port.UserRepository) http.HandlerFunc {
	roles := map[domain.UserRole]bool{
		domain.Admin:    true,
		domain.Ordinary: true,
	}
	return AuthMiddleware(roles, ts, ur, f)
}

func AdminMiddleware(f http.HandlerFunc, ts port.TokenService, ur port.UserRepository) http.HandlerFunc {
	roles := map[domain.UserRole]bool{
		domain.Admin: true,
	}
	return AuthMiddleware(roles, ts, ur, f)
}
