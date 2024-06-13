package port

import "myapp/internal-xkcd/core/domain"

type UserClient interface {
	UserRole(userID int64) (domain.UserRole, error)
	ChangeRole(reqUserID int64, reqRole domain.UserRole) error
}
