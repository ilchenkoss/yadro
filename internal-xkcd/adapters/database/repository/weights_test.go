package repository

import (
	"github.com/stretchr/testify/assert"
	"myapp/internal-xkcd/adapters/database"
	"myapp/internal-xkcd/config"
	"myapp/internal-xkcd/core/domain"
	"testing"
)

func TestWeightsRepository(t *testing.T) {
	cfg := &config.DatabaseConfig{
		DatabasePath: ":memory:",
	}

	db, err := database.NewConnection(cfg)
	assert.NoError(t, err)
	defer func(db *database.DB) {
		clConErr := db.CloseConnection()
		assert.NoError(t, clConErr)
	}(db)

	_, execErr := db.Exec(`
		CREATE TABLE IF NOT EXISTS comics (
			id INTEGER PRIMARY KEY,
			picture TEXT,
			title TEXT,
			alt TEXT,
			transcript TEXT,
			description TEXT
		);
		INSERT INTO comics (id, picture) values (1, "pict1.jpg");
		INSERT INTO comics (id, picture) values (3, "pict2.jpg");
		CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY,
			word TEXT UNIQUE
		);
		CREATE TABLE IF NOT EXISTS positions (
			id INTEGER PRIMARY KEY,
			position TEXT
		);
		CREATE TABLE IF NOT EXISTS weights (
			word_id INTEGER,
			comic_id INTEGER,
			position_id INTEGER,
			weight REAL,
			FOREIGN KEY (comic_id) REFERENCES comics(id),
			FOREIGN KEY (word_id) REFERENCES words(id),
			FOREIGN KEY (position_id) REFERENCES positions(id)
		);`)
	assert.NoError(t, execErr)

	wRepo := NewWeightsRepository(db)
	assert.NotNil(t, wRepo)

	posTitle := domain.Positions{ID: 1, Position: "title"}
	posAlt := domain.Positions{ID: 2, Position: "alt"}
	posTranscript := domain.Positions{ID: 3, Position: "transcript"}

	positions := []domain.Positions{posTitle, posAlt, posTranscript}
	ipErr := wRepo.InsertPositions(&positions)
	assert.NoError(t, ipErr)

	weightsToInsert := []domain.Weights{
		{
			Word:     &domain.Words{Word: "word1"},
			Comic:    &domain.Comics{ID: 1},
			Position: &posTitle,
			Weight:   1.1,
		},
		{
			Word:     &domain.Words{Word: "word2"},
			Comic:    &domain.Comics{ID: 1},
			Position: &posAlt,
			Weight:   1.2,
		},
		{
			Word:     &domain.Words{Word: "word1"},
			Comic:    &domain.Comics{ID: 1},
			Position: &posTranscript,
			Weight:   1.3,
		},
		{
			Word:     &domain.Words{Word: "word1"},
			Comic:    &domain.Comics{ID: 2},
			Position: &posTitle,
			Weight:   1.1,
		},
	}
	iwErr := wRepo.InsertWeights(&weightsToInsert)
	assert.NoError(t, iwErr)

	weights, gwbwErr := wRepo.GetWeightsByWords(map[string]float64{"word1": 2.1, "word2": 2.2})
	assert.NoError(t, gwbwErr)
	assert.Equal(t, 3, len(*weights))
}
