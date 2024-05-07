package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"myapp/internal/core/domain"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"myapp/internal/config"
)

type DB struct {
	db *sql.DB
}

func NewConnection(cfg *config.DatabaseConfig) (*DB, error) {
	db, err := sql.Open("sqlite3", cfg.DatabasePath)

	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) GetComicsByID(ID int) (*domain.Comics, error) {

	query := "SELECT c.id, c.picture, c.title, c.alt, c.transcript FROM comics c WHERE id = ?"

	row := d.db.QueryRow(query, ID)

	var comics domain.Comics

	err := row.Scan(&comics.ID, &comics.Picture, &comics.Title, &comics.Alt, &comics.Transcript)
	if err != nil {
		return &comics, err
	}

	return &comics, nil
}

func (d *DB) GetWeightsByWords(words map[string]float64) (*[]domain.Weights, error) {

	//append args with words to query
	args := make([]interface{}, len(words))
	i := 0
	for word := range words {
		args[i] = word
		i++
	}

	placeholders := strings.Repeat("?,", len(args))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(`
		SELECT wd.word, w.comic_id,p.position, w.weight, c.picture
		FROM weights w
		JOIN positions p ON w.position_id = p.id
		JOIN words wd ON w.word_id = wd.id
		JOIN comics c ON w.comic_id = c.id
		WHERE wd.word IN (%s)
	`, placeholders)

	rows, qErr := d.db.Query(query, args...)
	if qErr != nil {
		return nil, qErr
	}
	defer rows.Close()

	var weights []domain.Weights
	for rows.Next() {
		weight := domain.Weights{
			Word:     &domain.Words{},
			Comic:    &domain.Comics{},
			Position: &domain.Positions{},
		}
		if scanErr := rows.Scan(&weight.Word.Word, &weight.Comic.ID, &weight.Position.Position, &weight.Weight, &weight.Comic.Picture); scanErr != nil {
			fmt.Errorf("Error scanning row: %v", scanErr)
			return nil, scanErr
		}
		weights = append(weights, weight)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println(len(weights))
	return &weights, nil
}

func (d *DB) GetCountComics() (int, error) {
	var comicsCount int
	queryErr := d.db.QueryRow("SELECT COUNT(*) FROM comics").Scan(&comicsCount)
	if queryErr != nil {
		return 0, queryErr
	}
	return comicsCount, nil
}

func (d *DB) CloseConnection() error {
	err := d.db.Close()
	return err
}

func (d *DB) Ping() error {
	err := d.db.Ping()
	return err
}

func (d *DB) GetMissedIDs() (map[int]bool, error) {

	missedIDs := make(map[int]bool)

	query, queryErr := d.db.Query(`
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
			fmt.Errorf("Error scanning row: %v", scanErr)
			return missedIDs, scanErr
		}
		missedIDs[missedID] = true
	}

	if queryErr = query.Err(); queryErr != nil {
		return missedIDs, queryErr
	}

	return missedIDs, nil
}

func (d *DB) GetMaxID() (int, error) {
	var maxID int

	err := d.db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM comics").Scan(&maxID)
	if err != nil {
		return 0, err
	}

	return maxID, nil

}

func (d *DB) InsertComics(comics *[]domain.Comics) (int, error) {

	//start transaction
	tx, txErr := d.db.Begin()
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

func (d *DB) InsertWeights(weights *[]domain.Weights) error {

	//start transaction
	tx, txErr := d.db.Begin()
	if txErr != nil {
		return txErr
	}

	stmtWords, pwdsErr := tx.Prepare(`
		INSERT INTO words(word) VALUES(?)
		ON CONFLICT(word) DO NOTHING;
	`)
	if pwdsErr != nil {
		return pwdsErr
	}

	stmtWeights, pwgtErr := tx.Prepare(`
	INSERT INTO weights(word_id, comic_id, position_id, weight) 
	SELECT wds.id, ?, pos.id, ? FROM words wds
	JOIN positions pos
	WHERE wds.word = ? AND pos.position = ?;
	`)
	if pwgtErr != nil {
		return pwgtErr
	}

	for _, weight := range *weights {
		//insert words to words
		_, ewdsErr := stmtWords.Exec(weight.Word.Word)
		if ewdsErr != nil {
			return ewdsErr
		}
		//insert weights
		_, ewgtErr := stmtWeights.Exec(weight.Comic.ID, weight.Weight, weight.Word.Word, weight.Position.Position)
		if ewgtErr != nil {
			return ewgtErr
		}
	}

	//close transaction
	errCommit := tx.Commit()
	if errCommit != nil {
		return errCommit
	}

	return nil
}

func (d *DB) InsertPositions(positions *[]domain.Positions) error {

	placeholders := make([]string, 0, len(*positions))
	args := make([]interface{}, len(*positions)*2)

	for i, pos := range *positions {
		placeholders = append(placeholders, "(?, ?)")
		args[i*2] = pos.ID
		args[i*2+1] = pos.Position
	}

	query := fmt.Sprintf("INSERT OR IGNORE INTO positions (id, position) VALUES %s", strings.Join(placeholders, ", "))
	_, err := d.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
