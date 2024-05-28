package repository

import (
	"fmt"
	"myapp/internal-api/adapters/database"
	"myapp/internal-api/core/domain"
)

type ComicsRepository struct {
	db *database.DB
}

func NewComicsRepository(db *database.DB) *ComicsRepository {
	return &ComicsRepository{
		db,
	}
}

func (cr *ComicsRepository) GetComicsByID(ID int) (*domain.Comics, error) {

	query := "SELECT c.id, c.picture, c.title, c.alt, c.transcript FROM comics c WHERE id = ?"

	row := cr.db.QueryRow(query, ID)

	var comics domain.Comics

	err := row.Scan(&comics.ID, &comics.Picture, &comics.Title, &comics.Alt, &comics.Transcript)
	if err != nil {
		return &comics, err
	}

	return &comics, nil
}

func (cr *ComicsRepository) GetCountComics() (int, error) {
	var comicsCount int
	queryErr := cr.db.QueryRow("SELECT COUNT(*) FROM comics").Scan(&comicsCount)
	if queryErr != nil {
		return 0, queryErr
	}
	return comicsCount, nil
}

func (cr *ComicsRepository) GetMissedIDs() (map[int]bool, error) {

	missedIDs := make(map[int]bool)

	query, queryErr := cr.db.Query(`
		SELECT t1.id + 1 AS missing_id
		FROM comics AS t1
		LEFT JOIN comics AS t2 ON t1.id + 1 = t2.id
		WHERE t2.id IS NULL AND t1.id < (SELECT MAX(id) FROM comics)
	`)
	if queryErr != nil {
		return missedIDs, queryErr
	}
	defer query.Close()

	var missedID int
	for query.Next() {
		if scanErr := query.Scan(&missedID); scanErr != nil {
			scanErr = fmt.Errorf("Error scanning row: %v", scanErr)
			return missedIDs, scanErr
		}
		missedIDs[missedID] = true
	}

	if queryErr = query.Err(); queryErr != nil {
		return missedIDs, queryErr
	}

	return missedIDs, nil
}

func (cr *ComicsRepository) GetMaxID() (int, error) {
	var maxID int

	err := cr.db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM comics").Scan(&maxID)
	if err != nil {
		return 0, err
	}

	return maxID, nil

}

func (cr *ComicsRepository) InsertComics(comics *[]domain.Comics) (int, error) {

	//start transaction
	tx, txErr := cr.db.Begin()
	if txErr != nil {
		return 0, txErr
	}

	//prepare sql query
	stmt, stmtError := tx.Prepare("INSERT OR REPLACE INTO comics(id, picture, title, alt, transcript) VALUES(?, ?, ?, ?, ?)")
	if stmtError != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return 0, rollBackErr
		}
		return 0, stmtError
	}
	defer stmt.Close()

	//insert comics
	insertScore := 0
	for _, comic := range *comics {
		_, execErr := stmt.Exec(comic.ID, comic.Picture, comic.Title, comic.Alt, comic.Transcript) // _ -> result affected and last inserted ID
		if execErr != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				return insertScore, rollBackErr
			}
			return insertScore, execErr
		}
		insertScore++
	}

	//close transaction
	errCommit := tx.Commit()
	if errCommit != nil {
		return insertScore, errCommit
	}

	return insertScore, nil
}
