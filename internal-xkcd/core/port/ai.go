package port

import (
	"myapp/internal-xkcd/core/domain"
)

type AiService interface {
	GetComicsDescription(comic domain.Comics) (string, error)
}
