package port

import "myapp/internal-web/core/domain"

type XkcdAPI interface {
	UpdateDescription(id string, authToken string) error
	GetComics(requestWords string, authToken string) ([]domain.Comics, error)
	Login(login string, password string) (string, error)
}
