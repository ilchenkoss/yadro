package repository

import (
	"github.com/stretchr/testify/assert"
	"myapp/internal-api/adapters/database"
	"myapp/internal-api/config"
	"myapp/internal-api/core/domain"
	"testing"
)

func TestComicsRepository(t *testing.T) {
	cfg := &config.DatabaseConfig{
		DatabasePath: ":memory:",
	}

	db, err := database.NewConnection(cfg)
	assert.NoError(t, err)
	defer func(db *database.DB) {
		clConErr := db.CloseConnection()
		assert.NoError(t, clConErr)
	}(db)

	_, execErr := db.Exec(`CREATE TABLE IF NOT EXISTS comics (
    id INTEGER PRIMARY KEY,
    picture TEXT,
    title TEXT,
    alt TEXT,
    transcript TEXT,
    description TEXT
);`)
	assert.NoError(t, execErr)

	cRepo := NewComicsRepository(db)
	assert.NotNil(t, cRepo)

	comic1 := domain.Comics{
		ID:          1,
		Picture:     "comics1.jpg",
		Title:       "title words",
		Transcript:  "transcript words",
		Description: ""}

	Comics := []domain.Comics{comic1,
		{
			ID:          3,
			Picture:     "comics2.jpg",
			Title:       "title words",
			Alt:         "alt words",
			Transcript:  "transcript words",
			Description: "",
		}}

	//success
	insertedComics, icErr := cRepo.InsertComics(&Comics)
	assert.NoError(t, icErr)
	assert.Equal(t, 2, insertedComics)

	//success
	countComics, gccErr := cRepo.GetCountComics()
	assert.NoError(t, gccErr)
	assert.Equal(t, 2, countComics)

	//success
	comic, gcbIDErr := cRepo.GetComicsByID(comic1.ID)
	assert.NoError(t, gcbIDErr)
	assert.Equal(t, &comic1, comic)

	//success
	missedIDs, gmIDsErr := cRepo.GetMissedIDs()
	assert.NoError(t, gmIDsErr)
	assert.Equal(t, map[int]bool{2: true}, missedIDs)

	//success
	maxID, gmIDErr := cRepo.GetMaxID()
	assert.NoError(t, gmIDErr)
	assert.Equal(t, 3, maxID)

	//success
	ucdErr := cRepo.UpdateComicsDescriptionByID("3", "description words")
	assert.NoError(t, ucdErr)
	uComics, gucErr := cRepo.GetComicsByID(3)
	assert.NoError(t, gucErr)
	assert.Equal(t, &domain.Comics{
		ID:          3,
		Picture:     "comics2.jpg",
		Title:       "title words",
		Alt:         "alt words",
		Transcript:  "transcript words",
		Description: "description words",
	}, uComics)

}
