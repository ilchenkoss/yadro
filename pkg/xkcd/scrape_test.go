package xkcd

import (
	"myapp/pkg/database"
	"testing"
)

func TestMissID(t *testing.T) {

	scrapeIDs := 10

	scrapeResult := MainScrape(database.ScrapeResult{}, scrapeIDs, 1)

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
