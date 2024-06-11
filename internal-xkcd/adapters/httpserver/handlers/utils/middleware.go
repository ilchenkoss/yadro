package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"myapp/internal-xkcd/core/domain"
	"myapp/internal-xkcd/core/port"
	"net/http"
	"strconv"
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

func AuthMiddleware(requiredRoles map[domain.UserRole]bool, ac port.AuthClient, l *Limiter, next http.HandlerFunc) http.HandlerFunc {

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

		token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
		if err != nil {
			http.Error(w, "token is not valid", http.StatusBadGateway)
			return
		}

		userIDString, sErr := token.Claims.GetSubject()
		if sErr != nil {
			http.Error(w, "token is not valid", http.StatusBadGateway)
			return
		}

		userID, cErr := strconv.Atoi(userIDString)
		if cErr != nil {
			http.Error(w, "token is not valid", http.StatusBadGateway)
			return
		}

		//validate token and return role
		userRole, rErr := ac.UserRole(int64(userID))
		if rErr != nil {
			http.Error(w, "token is not valid", http.StatusBadGateway)
			return
		}

		//
		//userLogin, tErr := ts.GetUserByTokenString(accessToken)
		//if tErr != nil {
		//	switch {
		//	case errors.Is(tErr, domain.ErrTokenExpired):
		//		http.Error(w, "token expired", http.StatusUnauthorized)
		//		return
		//	case errors.Is(tErr, domain.ErrTokenNotValid):
		//		http.Error(w, "invalid token", http.StatusUnauthorized)
		//		return
		//	default:
		//		http.Error(w, "auth failed", http.StatusInternalServerError)
		//		return
		//	}
		//}
		//
		//user, rguErr := ur.GetUserByLogin(userLogin)
		//if rguErr != nil {
		//	switch {
		//	case errors.Is(rguErr, domain.ErrUserNotFound):
		//		slog.Error("Error attempt to log in using a non-existent login", slog.String("token", accessToken))
		//		http.Error(w, "invalid token", http.StatusUnauthorized)
		//		return
		//	default:
		//		http.Error(w, "auth failed", http.StatusInternalServerError)
		//		return
		//	}
		//}

		if !requiredRoles[userRole] {
			http.Error(w, "insufficient permissions", http.StatusForbidden)
			return
		}
		RateLimiterMiddleware(int64(userID), l, next)(w, r)
	}
}

func OrdinaryMiddleware(f http.HandlerFunc, ac port.AuthClient, l *Limiter) http.HandlerFunc {
	roles := map[domain.UserRole]bool{
		domain.Admin:     true,
		domain.Ordinary:  true,
		domain.SuperUser: true,
	}
	return AuthMiddleware(roles, ac, l, f)
}

func AdminMiddleware(f http.HandlerFunc, ac port.AuthClient, l *Limiter) http.HandlerFunc {
	roles := map[domain.UserRole]bool{
		domain.Admin:     true,
		domain.SuperUser: true,
	}
	return AuthMiddleware(roles, ac, l, f)
}

func SuperUserMiddleware(f http.HandlerFunc, ac port.AuthClient, l *Limiter) http.HandlerFunc {
	roles := map[domain.UserRole]bool{
		domain.SuperUser: true,
	}
	return AuthMiddleware(roles, ac, l, f)
}

func GuestMiddleware(f http.HandlerFunc, l *Limiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.cl.Add()
		defer l.cl.Done()
		f(w, r)
	}
}
