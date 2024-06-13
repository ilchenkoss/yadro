package service

import (
	"context"
	"encoding/json"
	"fmt"
	"myapp/internal-xkcd/config"
	"myapp/internal-xkcd/core/domain"
	"myapp/internal-xkcd/core/port"
	"myapp/internal-xkcd/core/util"
	"net/http"
	"sync"
)

type ScrapeService struct {
	scraper port.Scraper
	ctx     context.Context
	scfg    config.ScrapeConfig
}

func NewScrapeService(ctx context.Context, scraper port.Scraper, scfg config.ScrapeConfig) *ScrapeService {
	return &ScrapeService{
		scraper,
		ctx,
		scfg,
	}
}

type Comics struct {
	ID         int    `json:"num"`
	Picture    string `json:"img"`
	Title      string `json:"title"`
	Alt        string `json:"alt"`
	Transcript string `json:"transcript"`
}

func (s *ScrapeService) Scrape(missedIDs map[int]bool, maxID int, temper *util.Temper) ([]domain.Comics, error) {

	// Create scrapeContext
	scrapeCtx, scrapeCtxCancel := context.WithCancel(s.ctx)

	// add mutex
	var mu sync.Mutex

	var result []domain.Comics

	// Create buffered channel
	IDsCh := make(chan int, 1)

	// Set scrape WaitGroup
	var wg sync.WaitGroup
	wg.Add(s.scfg.Parallel)

	// Append temped response
	go AppendTempedResponse(&wg, &mu, temper, &result)

	// Append IDs to scrape
	go AppendIDs(scrapeCtx, scrapeCtxCancel, IDsCh, s.scfg.ScrapePagesLimit, missedIDs, maxID)

	// Create scrape worker goroutines
	for w := 1; w <= s.scfg.Parallel; w++ {
		go ScrapeWorker(scrapeCtx, scrapeCtxCancel, &wg, &mu, IDsCh, &result, s.scfg.RequestRetries, s.scraper, temper)
	}

	wg.Wait()

	return result, nil
}

func ScrapeWorker(
	ctx context.Context,
	scrapeCtxCancel context.CancelFunc,

	wg *sync.WaitGroup,
	mu *sync.Mutex,

	IDsChan chan int,
	result *[]domain.Comics,

	retries int,

	scraper port.Scraper,
	temper *util.Temper,
) {

	for {
		select {
		case ID := <-IDsChan:

			url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", ID)
			response, statusCode, _ := scraper.GetResponse(url, retries)

			if statusCode == http.StatusOK {
				stdbIDerr := temper.SaveTempDataByID(response, ID)
				if stdbIDerr != nil {
					//nothing
					//literal escape
					nothing := 0
					nothing++
				}
				var comics Comics
				err := json.Unmarshal(response, &comics)
				if err != nil {
					continue
				}

				dComics := domain.Comics{
					ID:         comics.ID,
					Picture:    comics.Picture,
					Title:      comics.Title,
					Alt:        comics.Alt,
					Transcript: comics.Transcript,
				}

				mu.Lock()
				*result = append(*result, dComics)
				mu.Unlock()

			}

			if statusCode == http.StatusNotFound {
				scrapeCtxCancel()
			}

		case <-ctx.Done():
			wg.Done()
			return
		}
	}
}

func AppendIDs(scrapeCtx context.Context, scrapeCancel context.CancelFunc, IDsCh chan int, scrapeLimit int, missedIDs map[int]bool, lastID int) {

	ID := 1
	for {
		select {
		case <-scrapeCtx.Done():
			return
		default:
			//ignore 404 ID
			if ID == 404 {
				ID++
			}
			if scrapeLimit == 0 {
				scrapeCancel()
				return
			}
			if missedIDs[ID] {
				IDsCh <- ID
				scrapeLimit--
			}
			if ID > lastID {
				IDsCh <- ID
				scrapeLimit--
			}
			ID++
		}
	}
}

func AppendTempedResponse(wg *sync.WaitGroup, mu *sync.Mutex, temper *util.Temper, result *[]domain.Comics) {
	for tempFile := range temper.TempFiles {

		wg.Add(1)

		filePath := fmt.Sprintf("%s/%s", temper.TempDir, tempFile)
		tempData := temper.ReadTempFile(filePath)
		var comics Comics
		err := json.Unmarshal(tempData, &comics)
		if err != nil {
			wg.Done()
			continue
		}

		dComics := domain.Comics{
			ID:         comics.ID,
			Picture:    comics.Picture,
			Title:      comics.Title,
			Alt:        comics.Alt,
			Transcript: comics.Transcript,
		}

		mu.Lock()
		*result = append(*result, dComics)
		mu.Unlock()

		wg.Done()
	}

}
