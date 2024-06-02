package port

import (
	"myapp/internal-api/core/domain"
)

type ComicsRepository interface {
	UpdateComicsDescriptionByID(ID string, description string) error
	GetComicsByID(ID int) (*domain.Comics, error)
	GetMissedIDs() (map[int]bool, error)
	GetMaxID() (int, error)
	GetCountComics() (int, error)
	InsertComics(*[]domain.Comics) (int, error)
}
