package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"myapp/internal-api/adapters/httpserver/handlers/utils"
	"myapp/internal-api/config"
	"myapp/internal-api/core/port"
	"myapp/internal-api/core/util"
	"net/http"
	"sync"
)

type ScrapeHandler struct {
	ssvc port.ScrapeService
	wsvc port.WeightService
	mu   *sync.Mutex
	cr   port.ComicsRepository
	wr   port.WeightRepository
	ctx  context.Context
	cfg  *config.Config
	fs   util.FileSystem
}

func NewScrapeHandler(ssvc port.ScrapeService, wsvc port.WeightService, cr port.ComicsRepository, wr port.WeightRepository, ctx context.Context, cfg *config.Config, fs util.FileSystem) *ScrapeHandler {
	var mu sync.Mutex
	return &ScrapeHandler{
		ssvc,
		wsvc,
		&mu,
		cr,
		wr,
		ctx,
		cfg,
		fs,
	}
}

func (sc *ScrapeHandler) Update(w http.ResponseWriter, r *http.Request) {

	if !sc.mu.TryLock() {
		slog.Info("Trying to request an update that is already in progress")
		http.Error(w, "Update already processing", http.StatusInternalServerError)
		return
	}
	defer sc.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	// Get missed IDs
	missedIDs, getMissedIDsErr := sc.cr.GetMissedIDs()
	if getMissedIDsErr != nil {
		slog.Error("Error get missed IDs :", "error", getMissedIDsErr.Error())
	}

	// Get max ID
	maxID, getMaxIDerr := sc.cr.GetMaxID()
	if getMaxIDerr != nil {
		slog.Error("Error get max ID :", "error", getMaxIDerr.Error())
	}

	// Init temper
	temper, tempErr := util.NewTemper(&sc.cfg.Temp, sc.fs)
	if tempErr != nil {
		slog.Error("Error init util.temp :", "error", tempErr.Error())
	} else {
		// temp init ok
		if temper.MaxTempedID > maxID {
			maxID = temper.MaxTempedID
		}
	}

	// Start scrape
	newComics, scrapeErr := sc.ssvc.Scrape(missedIDs, maxID, temper)
	if scrapeErr != nil {
		slog.Error("Error scrape :", "error", scrapeErr.Error())
		http.Error(w, "Update error", http.StatusInternalServerError)
		return
	}

	// Write comics to db
	insertedCount, insertError := sc.cr.InsertComics(&newComics)
	if insertError != nil {
		slog.Error("Error insert database :", "error", insertError.Error())
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	tempRemoveErr := temper.RemoveTemp()
	if tempRemoveErr != nil {
		slog.Error("Error remove temp folder :", "error", tempRemoveErr.Error())
	}

	weights := sc.wsvc.WeightComics(newComics)
	insertWeightsErr := sc.wr.InsertWeights(weights)
	if insertWeightsErr != nil {
		slog.Error("Error insert weights: ", "error", insertWeightsErr.Error())
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Get comics count
	comicsCount, errGetCount := sc.cr.GetCountComics()
	if errGetCount != nil {
		slog.Error("Error get comics count :", "error", errGetCount.Error())
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(utils.NewUpdateResponse(true, "Success", insertedCount, comicsCount))
	if err != nil {
		//nothing
		return
	}
}
