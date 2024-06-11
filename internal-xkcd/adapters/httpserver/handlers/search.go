package handlers

import (
	"encoding/json"
	"log/slog"
	"myapp/internal-xkcd/adapters"
	"myapp/internal-xkcd/adapters/httpserver/handlers/utils"
	"myapp/internal-xkcd/core/port"
	"net/http"
	"strconv"
)

type SearchHandler struct {
	wr port.WeightRepository
	cr port.ComicsRepository
	ws port.WeightService
	ga *adapters.GptAPI
}

func NewSearchHandler(wr port.WeightRepository, ws port.WeightService, cr port.ComicsRepository, ga *adapters.GptAPI, l utils.Limiter) *SearchHandler {
	return &SearchHandler{
		wr,
		cr,
		ws,
		ga,
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

	comicID, atoiErr := strconv.Atoi(requestComicID)
	if atoiErr != nil {
		http.Error(w, "comic ID error", http.StatusInternalServerError)
	}
	comic, gcErr := s.cr.GetComicsByID(comicID)
	if gcErr != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
	}

	description, cdErr := s.ga.GetComicsDescription(*comic)
	if cdErr != nil {
		http.Error(w, "gpt error", http.StatusInternalServerError)
	}

	ucdErr := s.cr.UpdateComicsDescriptionByID(requestComicID, description)
	if ucdErr != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
	}

	err := json.NewEncoder(w).Encode(utils.NewDescriptionResponse(true, "Success"))
	if err != nil {
		//nothing
		return
	}
}
