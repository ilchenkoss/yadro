package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"myapp/internal/config"
	"myapp/internal/core/port"
	"myapp/internal/core/util"
	"net/http"
	"sync"
)

type ScrapeHandler struct {
	ssvc port.ScrapeService
	wsvc port.WeightService
	mu   *sync.Mutex
	db   port.Database
	ctx  context.Context
	cfg  *config.Config
}

func NewScrapeHandler(ssvc port.ScrapeService, wsvc port.WeightService, db port.Database, ctx context.Context, cfg *config.Config) *ScrapeHandler {
	var mu sync.Mutex
	return &ScrapeHandler{
		ssvc,
		wsvc,
		&mu,
		db,
		ctx,
		cfg,
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
	missedIDs, getMissedIDsErr := sc.db.GetMissedIDs()
	if getMissedIDsErr != nil {
		slog.Error("Error get missed IDs :", getMissedIDsErr)
	}

	// Get max ID
	maxID, getMaxIDerr := sc.db.GetMaxID()
	if getMaxIDerr != nil {
		slog.Error("Error get max ID :", getMaxIDerr)
	}

	// Init temper
	temper, tempErr := util.NewTemper(&sc.cfg.Temp)
	if tempErr != nil {
		slog.Error("Error init util.temp :", tempErr)
	} else {
		// temp init ok
		if temper.MaxTempedID > maxID {
			maxID = temper.MaxTempedID
		}
	}

	// Start scrape
	newComics, scrapeErr := sc.ssvc.Scrape(missedIDs, maxID, temper)
	if scrapeErr != nil {
		slog.Error("Error scrape :", scrapeErr)
		http.Error(w, "Update error", http.StatusInternalServerError)
		return
	}

	// Write comics to db
	insertedCount, insertError := sc.db.InsertComics(&newComics)
	if insertError != nil {
		slog.Error("Error insert database :", insertError)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	tempRemoveErr := temper.RemoveTemp()
	if tempRemoveErr != nil {
		slog.Error("Error remove temp folder :", tempRemoveErr)
	}

	weights := sc.wsvc.WeightComics(newComics)
	insertWeightsErr := sc.db.InsertWeights(weights)
	if insertWeightsErr != nil {
		slog.Error("Error insert weights: ", insertWeightsErr)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Get comics count
	comicsCount, errGetCount := sc.db.GetCountComics()
	if errGetCount != nil {
		slog.Error("Error get comics count :", errGetCount)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newUpdateResponse(true, "Success", insertedCount, comicsCount))
	return
}
