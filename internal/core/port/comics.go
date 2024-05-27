package port

import (
	"myapp/internal/core/domain"
)

type ComicsRepository interface {
	GetComicsByID(ID int) (*domain.Comics, error)
	GetMissedIDs() (map[int]bool, error)
	GetMaxID() (int, error)
	GetCountComics() (int, error)
	InsertComics(*[]domain.Comics) (int, error)
}
