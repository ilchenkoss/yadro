package weight

import (
	"math"
	"myapp/internal/core/domain"
	"myapp/pkg/words"
	"sort"
)

type WeightService struct {
	WeightsByPosition map[string]float64
}

func NewWeightService() *WeightService {
	comicsWordWeights := make(map[string]float64)
	comicsWordWeights["title"] = domain.WeightComicsWordPositionTitle
	comicsWordWeights["alt"] = domain.WeightComicsWordPositionAlt
	comicsWordWeights["transcript"] = domain.WeightComicsWordPositionTranscript
	return &WeightService{
		comicsWordWeights,
	}
}

func (w *WeightService) WeightRequest(request string) map[string]float64 {

	requestWords := words.StringNormalization(request)

	wordsWeight := make(map[string]float64)
	for word, wordInfo := range requestWords {
		var wordWeight float64

		// word repeat weight
		if wordInfo.Repeat > 1 {
			wordWeight += float64(wordInfo.Repeat) * domain.WeightRequestWordDuplicate
		}
		// word index weight
		wordWeight += domain.WeightRequestWordIndex / math.Log(float64(wordInfo.EntryIndex+2))
		wordsWeight[word] = wordWeight

	}

	return wordsWeight
}

func (w *WeightService) WeightComics(comics []domain.Comics) *[]domain.Weights {

	var weights []domain.Weights

	for _, comic := range comics {

		comicWordsByPosition := make(map[string]map[string]words.KeywordsInfo)
		comicWordsByPosition["title"] = words.StringNormalization(comic.Title)
		comicWordsByPosition["alt"] = words.StringNormalization(comic.Alt)
		comicWordsByPosition["transcript"] = words.StringNormalization(comic.Transcript)

		for position, comicWords := range comicWordsByPosition {

			for word, wordInfo := range comicWords {

				weight := domain.Weights{
					Comic:    &domain.Comics{ID: comic.ID},
					Word:     &domain.Words{Word: word},
					Weight:   float64(comic.ID) * domain.WeightComicsActual,
					Position: &domain.Positions{Position: position},
				}

				// word repeat weight
				if wordInfo.Repeat > 1 {
					weight.Weight += float64(wordInfo.Repeat) * domain.WeightComicsWordDuplicate
				}

				// word index weight
				weight.Weight += domain.WeightComicsWordIndex / math.Log(float64(wordInfo.EntryIndex+2))

				// word weight
				weight.Weight *= w.WeightsByPosition[position]

				weights = append(weights, weight)

			}
		}
	}
	return &weights
}

func (w *WeightService) FindRelevantPictures(requestWeights map[string]float64, weights *[]domain.Weights) ([]string, error) {

	type WeightsInfo struct {
		Alt        []string
		Title      []string
		Transcript []string
		Picture    string
		Weight     float64
	}

	weightsByID := make(map[int]WeightsInfo)

	for _, weight := range *weights {

		weightsInfo := weightsByID[weight.Comic.ID]

		switch weight.Position.Position {
		case "alt":
			weightsInfo.Alt = append(weightsInfo.Alt, weight.Word.Word)
		case "title":
			weightsInfo.Title = append(weightsInfo.Title, weight.Word.Word)
		case "transcript":
			weightsInfo.Transcript = append(weightsInfo.Transcript, weight.Word.Word)
		}

		weightsInfo.Picture = weight.Comic.Picture
		weightsInfo.Weight += weight.Weight
		weightsByID[weight.Comic.ID] = weightsInfo
	}

	var sliceWeights []WeightsInfo

	//coverage correction
	for id := range weightsByID {

		weightsInfo := weightsByID[id]
		//alt coverage
		weightsInfo.Weight += float64(len(weightsByID[id].Alt)/len(requestWeights)) * domain.WeightResponseRelevantCoverage * domain.WeightResponseRelevantCoverageAlt
		//title coverage
		weightsInfo.Weight += float64(len(weightsByID[id].Title)/len(requestWeights)) * domain.WeightResponseRelevantCoverage * domain.WeightResponseRelevantCoverageTitle
		//transcript coverage
		weightsInfo.Weight += float64(len(weightsByID[id].Transcript)/len(requestWeights)) * domain.WeightResponseRelevantCoverage * domain.WeightResponseRelevantCoverageTranscript

		weightsByID[id] = weightsInfo

		sliceWeights = append(sliceWeights, weightsInfo)
	}

	//sort by weight
	sort.Slice(sliceWeights, func(i, j int) bool {
		return sliceWeights[i].Weight >= sliceWeights[j].Weight
	})

	result := make([]string, 0, len(sliceWeights))

	//append result slice
	for _, value := range sliceWeights {
		result = append(result, value.Picture)
	}

	return result, nil
}
