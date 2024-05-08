package handlers

import (
	"encoding/json"
	"log/slog"
	"myapp/internal/core/port"
	"net/http"
)

type SearchHandler struct {
	db port.Database
	w  port.WeightService
}

func NewSearchHandler(db port.Database, w port.WeightService) *SearchHandler {
	return &SearchHandler{
		db,
		w,
	}
}

func (s *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestString := r.URL.Query().Get("search")

	requestWeights := s.w.WeightRequest(requestString)

	weights, getWeightsByWordsErr := s.db.GetWeightsByWords(requestWeights)
	if getWeightsByWordsErr != nil {
		slog.Error("Error find relevant pictures :", getWeightsByWordsErr)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	pictures, findRelevantComicsErr := s.w.FindRelevantPictures(requestWeights, weights)
	if findRelevantComicsErr != nil {
		slog.Error("Error find relevant pictures :", findRelevantComicsErr)
		http.Error(w, "Find pictures error", http.StatusInternalServerError)
		return
	}

	if len(pictures) > 10 {
		pictures = pictures[:10]
	}

	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(newSearchResponse(true, "Success", pictures))
	w.Write(bytes)
}
