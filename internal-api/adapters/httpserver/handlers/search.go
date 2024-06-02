package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"myapp/internal-api/adapters/httpserver/handlers/utils"
	"myapp/internal-api/core/port"
	"net/http"
)

type SearchHandler struct {
	wr port.WeightRepository
	cr port.ComicsRepository
	ws port.WeightService
}

func NewSearchHandler(wr port.WeightRepository, ws port.WeightService, cr port.ComicsRepository, l utils.Limiter) *SearchHandler {
	return &SearchHandler{
		wr,
		cr,
		ws,
	}
}

func (s *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestString := r.URL.Query().Get("search")

	requestWeights := s.ws.WeightRequest(requestString)
	if len(requestWeights) == 0 {
		http.Error(w, "Request empty", http.StatusBadRequest)
		return
	}

	weights, getWeightsByWordsErr := s.wr.GetWeightsByWords(requestWeights)
	if getWeightsByWordsErr != nil {
		slog.Error("Error find relevant pictures :", "error", getWeightsByWordsErr.Error())
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	pictures, findRelevantComicsErr := s.ws.FindRelevantPictures(requestWeights, weights)
	if findRelevantComicsErr != nil {
		slog.Error("Error find relevant pictures :", "error", findRelevantComicsErr.Error())
		http.Error(w, "Find pictures error", http.StatusInternalServerError)
		return
	}

	if len(pictures) > 10 {
		pictures = pictures[:10]
	}

	//w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(utils.NewSearchResponse(true, "Success", pictures))
	if err != nil {
		//nothing
		return
	}
}

func (s *SearchHandler) Description(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestComicID := r.URL.Query().Get("description")
	fmt.Println("requestID", requestComicID)
	result := fmt.Sprintf("test comics %s description", requestComicID)

	err1 := s.cr.UpdateComicsDescriptionByID(requestComicID, result)
	fmt.Println(err1)

	//w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(utils.NewDescriptionResponse(true, "Success"))
	if err != nil {
		//nothing
		return
	}
}
