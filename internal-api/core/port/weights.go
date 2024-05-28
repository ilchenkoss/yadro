package port

import "myapp/internal-api/core/domain"

type WeightRepository interface {
	InsertWeights(*[]domain.Weights) error
	InsertPositions(*[]domain.Positions) error
	GetWeightsByWords(words map[string]float64) (*[]domain.Weights, error)
}

type WeightService interface {
	FindRelevantPictures(requestWeights map[string]float64, weights *[]domain.Weights) ([]string, error)
	WeightComics(comics []domain.Comics) *[]domain.Weights
	WeightRequest(request string) map[string]float64
}
