package repository

import (
	"fmt"
	"myapp/internal-api/adapters/database"
	"myapp/internal-api/core/domain"
	"strings"
)

type WeightsRepository struct {
	db *database.DB
}

func NewWeightsRepository(db *database.DB) *WeightsRepository {
	return &WeightsRepository{
		db,
	}
}

func (wr *WeightsRepository) GetWeightsByWords(words map[string]float64) (*[]domain.Weights, error) {

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

	rows, qErr := wr.db.Query(query, args...)
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
			scanErr = fmt.Errorf("Error scanning row: %v", scanErr)
			return nil, scanErr
		}
		weights = append(weights, weight)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &weights, nil
}

func (wr *WeightsRepository) InsertWeights(weights *[]domain.Weights) error {

	//start transaction
	tx, txErr := wr.db.Begin()
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

func (wr *WeightsRepository) InsertPositions(positions *[]domain.Positions) error {

	placeholders := make([]string, 0, len(*positions))
	args := make([]interface{}, len(*positions)*2)

	for i, pos := range *positions {
		placeholders = append(placeholders, "(?, ?)")
		args[i*2] = pos.ID
		args[i*2+1] = pos.Position
	}

	query := fmt.Sprintf("INSERT OR IGNORE INTO positions (id, position) VALUES %s", strings.Join(placeholders, ", "))
	_, err := wr.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
