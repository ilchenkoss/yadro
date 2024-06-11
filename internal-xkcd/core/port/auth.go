package port

import "myapp/internal-xkcd/core/domain"

type AuthClient interface {
	Login(login string, password string) (string, error)
	Register(login string, password string, role domain.UserRole) (int64, error)
	UserRole(userID int64) (domain.UserRole, error)
}
