package utils

import (
	"errors"
	"myapp/internal-xkcd/core/domain"
	"myapp/internal-xkcd/core/port"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationType      = "bearer"
)

func RateLimiterMiddleware(id int64, l *Limiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitErr := l.rl.Add(id)
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
	}
}

func AuthMiddleware(requiredRoles map[domain.UserRole]bool, ac port.AuthClient, uc port.UserClient, l *Limiter, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		l.cl.Add()
		defer l.cl.Done()

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

		//validate token and return userID
		userID, uIDErr := ac.UserID(accessToken)
		if uIDErr != nil {
			switch {
			case errors.Is(uIDErr, domain.ErrTokenExpired):
				http.Error(w, "token expired", http.StatusUnauthorized)
			case errors.Is(uIDErr, domain.ErrTokenNotValid):
				http.Error(w, "token is not valid", http.StatusUnauthorized)
			default:
				http.Error(w, "", http.StatusUnauthorized)
			}
		}

		userRole, rErr := uc.UserRole(userID)
		if rErr != nil {
			//domain.ErrUserNotFound
			http.Error(w, "token is not valid", http.StatusUnauthorized)
			return
		}

		if !requiredRoles[userRole] {
			http.Error(w, "insufficient permissions", http.StatusForbidden)
			return
		}
		RateLimiterMiddleware(userID, l, next)(w, r)
	}
}

func OrdinaryMiddleware(f http.HandlerFunc, ac port.AuthClient, uc port.UserClient, l *Limiter) http.HandlerFunc {
	roles := map[domain.UserRole]bool{
		domain.Admin:     true,
		domain.Ordinary:  true,
		domain.SuperUser: true,
	}
	return AuthMiddleware(roles, ac, uc, l, f)
}

func AdminMiddleware(f http.HandlerFunc, ac port.AuthClient, uc port.UserClient, l *Limiter) http.HandlerFunc {
	roles := map[domain.UserRole]bool{
		domain.Admin:     true,
		domain.SuperUser: true,
	}
	return AuthMiddleware(roles, ac, uc, l, f)
}

func SuperUserMiddleware(f http.HandlerFunc, ac port.AuthClient, uc port.UserClient, l *Limiter) http.HandlerFunc {
	roles := map[domain.UserRole]bool{
		domain.SuperUser: true,
	}
	return AuthMiddleware(roles, ac, uc, l, f)
}

func GuestMiddleware(f http.HandlerFunc, l *Limiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.cl.Add()
		defer l.cl.Done()
		f(w, r)
	}
}
