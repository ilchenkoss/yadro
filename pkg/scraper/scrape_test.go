package scraper

import (
	"myapp/pkg/database"
	"testing"
)

func TestMissID(t *testing.T) {

	scrapeIDs := 10
	dbData := database.ScrapeResult{
		Data:   map[int]database.ParsedData{},
		BadIDs: map[int]int{},
	}
	scrapeResult := MainScrape(dbData, scrapeIDs, 1)

	var IDs []int
	for goodID := range scrapeResult.Data {
		IDs = append(IDs, goodID)
	}

	for badID := range scrapeResult.BadIDs {
		IDs = append(IDs, badID)
	}

	if len(IDs) != scrapeIDs {

		t.Errorf("\nResult was incorrect. \n scrapes: %d, \n IDs: %d.", scrapeIDs, len(IDs))
	}
}
