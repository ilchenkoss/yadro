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

func LimiterMiddleware(id uint64, limiter *Limiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitErr := limiter.Add(id)
		if limitErr != nil {
			switch {
			case errors.Is(limitErr, domain.ErrRateLimitExceeded):
				http.Error(w, "Requests was exceeded", http.StatusTooManyRequests)
				return
			default:
				http.Error(w, limitErr.Error(), http.StatusInternalServerError)
				return
			}
		}
		next(w, r)
		defer limiter.Done()
	}
}

func AuthMiddleware(requiredRoles map[domain.UserRole]bool, ts port.TokenService, ur port.UserRepository, l *Limiter, next http.HandlerFunc) http.HandlerFunc {
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
		LimiterMiddleware(user.ID, l, next)(w, r)
	}
}

func OrdinaryMiddleware(f http.HandlerFunc, ts port.TokenService, ur port.UserRepository, l *Limiter) http.HandlerFunc {
	roles := map[domain.UserRole]bool{
		domain.Admin:      true,
		domain.Ordinary:   true,
		domain.SuperAdmin: true,
	}
	return AuthMiddleware(roles, ts, ur, l, f)
}

func AdminMiddleware(f http.HandlerFunc, ts port.TokenService, ur port.UserRepository, l *Limiter) http.HandlerFunc {
	roles := map[domain.UserRole]bool{
		domain.Admin:      true,
		domain.SuperAdmin: true,
	}
	return AuthMiddleware(roles, ts, ur, l, f)
}

func SuperAdminMiddleware(f http.HandlerFunc, ts port.TokenService, ur port.UserRepository, l *Limiter) http.HandlerFunc {
	roles := map[domain.UserRole]bool{
		domain.SuperAdmin: true,
	}
	return AuthMiddleware(roles, ts, ur, l, f)
}

func GuestMiddleware(f http.HandlerFunc, l *Limiter) http.HandlerFunc {
	return LimiterMiddleware(0, l, f)
}
