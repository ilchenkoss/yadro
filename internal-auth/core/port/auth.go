package port

import "myapp/internal-auth/core/domain"

type TokenService interface {
	CreateToken(userID int64) (string, error)
	GetUserID(token string) (int64, error)
}

type AuthService interface {
	Login(login string, password string) (string, error)
	Register(login string, password string, role domain.UserRole) (int64, error)
	UserRole(userID int64) (domain.UserRole, error)
	UserID(token string) (int64, error)
}
