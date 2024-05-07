package port

import "myapp/internal/core/domain"

type WeightService interface {
	FindRelevantPictures(requestWeights map[string]float64, weights *[]domain.Weights) ([]string, error)
	WeightComics(comics []domain.Comics) *[]domain.Weights
	WeightRequest(request string) map[string]float64
}
