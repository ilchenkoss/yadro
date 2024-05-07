package port

import "myapp/internal/core/domain"

type Database interface {
	GetMissedIDs() (map[int]bool, error)
	GetMaxID() (int, error)
	GetCountComics() (int, error)
	CloseConnection() error
	Ping() error
	InsertComics(*[]domain.Comics) (int, error)
	InsertWeights(*[]domain.Weights) error
	InsertPositions(*[]domain.Positions) error
	GetComicsByID(ID int) (*domain.Comics, error)
	GetWeightsByWords(words map[string]float64) (*[]domain.Weights, error)
}
