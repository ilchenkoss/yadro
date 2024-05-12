package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"myapp/internal/core/domain"
	"myapp/internal/core/port"
	"net/http"
)

type SearchHandler struct {
	wr port.WeightRepository
	ws port.WeightService
	l  *Limiter
}

func NewSearchHandler(wr port.WeightRepository, ws port.WeightService, l *Limiter) *SearchHandler {
	return &SearchHandler{
		wr,
		ws,
		l,
	}
}

func (s *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestString := r.URL.Query().Get("search")

	limitErr := s.l.Add(1)
	if limitErr != nil {
		switch {
		case errors.Is(limitErr, domain.ErrRateLimitExceeded):
			http.Error(w, "Requests was exceeded", http.StatusTooManyRequests)
			return
		default:
			http.Error(w, limitErr.Error(), http.StatusInternalServerError)
			return
		}
	}
	defer s.l.Done()

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
	json.NewEncoder(w).Encode(newSearchResponse(true, "Success", pictures))
}
