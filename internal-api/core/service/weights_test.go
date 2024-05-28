package service

import (
	"github.com/stretchr/testify/assert"
	"myapp/internal-api/core/domain"
	"sort"
	"testing"
)

func TestNewWeightService(t *testing.T) {
	weightService := NewWeightService()
	assert.Equal(t, len(weightService.WeightsByPosition), 3)
}

func TestWeightService_WeightRequest(t *testing.T) {
	request := "A string that needs to be weight on wonderful scales. The weight of the word is measured in kg"
	expectedWordsByWeight := []string{"weight", "string", "need", "wonder", "scale", "word", "measur"}
	weightService := NewWeightService()
	weights := weightService.WeightRequest(request)

	var weightSlice []string
	for word := range weights {
		weightSlice = append(weightSlice, word)
	}

	sort.Slice(weightSlice, func(i, j int) bool {
		return weights[weightSlice[i]] > weights[weightSlice[j]]
	})

	assert.Equal(t, expectedWordsByWeight, weightSlice)
}

func TestWeightService_WeightComicsWordRepeat(t *testing.T) {
	comicWithRepeatWords := []domain.Comics{
		{Title: "wonderful wonderful"},
	}
	comicWithoutRepeatWords := []domain.Comics{
		{Title: "wonderful"},
	}

	weightService := NewWeightService()

	weightsWithRepeatWords := weightService.WeightComics(comicWithRepeatWords)
	assert.EqualValues(t, 1, len(*weightsWithRepeatWords))
	assert.EqualValues(t, "wonder", (*weightsWithRepeatWords)[0].Word.Word)

	weightsWithoutRepeatWords := weightService.WeightComics(comicWithoutRepeatWords)
	assert.EqualValues(t, 1, len(*weightsWithoutRepeatWords))
	assert.EqualValues(t, "wonder", (*weightsWithoutRepeatWords)[0].Word.Word)

	assert.True(t, (*weightsWithRepeatWords)[0].Weight > (*weightsWithoutRepeatWords)[0].Weight)
}

func TestWeightService_WeightComics(t *testing.T) {
	expectedWeights := []domain.Weights{{
		Word: &domain.Words{
			Word: "win",
		},
		Comic: &domain.Comics{
			ID: 2,
		},
		Position: &domain.Positions{
			Position: "title",
		},
		Weight: 1.0,
	}, {
		Word: &domain.Words{
			Word: "wonder",
		},
		Comic: &domain.Comics{
			ID: 1,
		},
		Position: &domain.Positions{
			Position: "title",
		},
		Weight: 2.0,
	}, {
		Word: &domain.Words{
			Word: "comic",
		},
		Comic: &domain.Comics{
			ID: 1,
		},
		Position: &domain.Positions{
			Position: "alt",
		},
		Weight: 3.0,
	}, {
		Word: &domain.Words{
			Word: "lost",
		},
		Comic: &domain.Comics{
			ID: 2,
		},
		Position: &domain.Positions{
			Position: "transcript",
		},
		Weight: 4.0,
	}, {
		Word: &domain.Words{
			Word: "appl",
		},
		Comic: &domain.Comics{
			ID: 1,
		},
		Position: &domain.Positions{
			Position: "transcript",
		},
		Weight: 5.0,
	}, {
		Word: &domain.Words{
			Word: "alt",
		},
		Comic: &domain.Comics{
			ID: 2,
		},
		Position: &domain.Positions{
			Position: "transcript",
		},
		Weight: 6.0,
	}}

	//clear weights
	for i := range expectedWeights {
		expectedWeights[i].Weight = 0.0
	}

	comics := []domain.Comics{
		{ID: 1, Picture: "1-Picture.jpg", Title: "wonderful", Alt: "comic", Transcript: "apple"},
		{ID: 2, Picture: "2-Picture.jpg", Title: "win", Alt: "", Transcript: "we lost alt"},
	}
	weightService := NewWeightService()
	weights := weightService.WeightComics(comics)

	//sort by weight
	sort.Slice(*weights, func(i, j int) bool {
		return (*weights)[i].Weight > (*weights)[j].Weight
	})

	//clear weights
	for i := range *weights {
		(*weights)[i].Weight = 0.0
	}
	assert.Equal(t, &expectedWeights, weights)

}

func TestFindRelevantPictures(t *testing.T) {
	request := "A string that needs to be weight on wonderful scales. The weight of the word is measured in kg"
	weightService := NewWeightService()
	weights := weightService.WeightRequest(request)

	weightsFromDB := []domain.Weights{{
		Word: &domain.Words{
			Word: "weight",
		},
		Comic: &domain.Comics{
			ID:      1,
			Picture: "1-Picture.jpg",
		},
		Position: &domain.Positions{
			Position: "title",
		},
		Weight: 2.0,
	}, {
		Word: &domain.Words{
			Word: "word",
		},
		Comic: &domain.Comics{
			ID:      1,
			Picture: "1-Picture.jpg",
		},
		Position: &domain.Positions{
			Position: "alt",
		},
		Weight: 3.0,
	}, {
		Word: &domain.Words{
			Word: "string",
		},
		Comic: &domain.Comics{
			ID:      2,
			Picture: "2-Picture.jpg",
		},
		Position: &domain.Positions{
			Position: "title",
		},
		Weight: 4.0,
	}, {
		Word: &domain.Words{
			Word: "need",
		},
		Comic: &domain.Comics{
			ID:      3,
			Picture: "3-picture.jpg",
		},
		Position: &domain.Positions{
			Position: "transcript",
		},
		Weight: 6.0,
	}}
	got, err := weightService.FindRelevantPictures(weights, &weightsFromDB)
	assert.NoError(t, err)
	assert.Equal(t, []string{"3-picture.jpg", "1-Picture.jpg", "2-Picture.jpg"}, got)
}
