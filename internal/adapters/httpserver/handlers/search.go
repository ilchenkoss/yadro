package handlers

import (
	"encoding/json"
	"log/slog"
	"myapp/internal/adapters/httpserver/handlers/utils"
	"myapp/internal/core/port"
	"net/http"
)

type SearchHandler struct {
	wr port.WeightRepository
	ws port.WeightService
}

func NewSearchHandler(wr port.WeightRepository, ws port.WeightService, l utils.Limiter) *SearchHandler {
	return &SearchHandler{
		wr,
		ws,
	}
}

func (s *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestString := r.URL.Query().Get("search")

	requestWeights := s.ws.WeightRequest(requestString)

	weights, getWeightsByWordsErr := s.wr.GetWeightsByWords(requestWeights)
	if getWeightsByWordsErr != nil {
		slog.Error("Error find relevant pictures :", getWeightsByWordsErr)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	pictures, findRelevantComicsErr := s.ws.FindRelevantPictures(requestWeights, weights)
	if findRelevantComicsErr != nil {
		slog.Error("Error find relevant pictures :", findRelevantComicsErr)
		http.Error(w, "Find pictures error", http.StatusInternalServerError)
		return
	}

	if len(pictures) > 10 {
		pictures = pictures[:10]
	}

	//w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.NewSearchResponse(true, "Success", pictures))
}
