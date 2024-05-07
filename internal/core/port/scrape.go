package port

import (
	"myapp/internal/core/domain"
	"myapp/internal/core/util"
)

type Scraper interface {
	GetResponse(url string, retries int) ([]byte, int, error)
}

type ScrapeService interface {
	Scrape(missedIDs map[int]bool, maxID int, temper *util.Temper) ([]domain.Comics, error)
}
