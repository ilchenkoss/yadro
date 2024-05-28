package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"myapp/internal-api/config"
	"myapp/internal-api/core/domain"
	"myapp/internal-api/core/port/mock"
	"myapp/internal-api/core/util"
	mockUtil "myapp/internal-api/core/util/mock"
	"slices"
	"sort"
	"sync"
	"testing"
	"testing/fstest"
	"time"
)

func TestScrapeService_Scrape(t *testing.T) {
	mockScraper := new(mock.MockScraper)
	ctx := context.Background()
	scrapeConfig := config.ScrapeConfig{
		Parallel:         5,
		ScrapePagesLimit: 10,
		RequestRetries:   3,
	}
	scrapeService := NewScrapeService(ctx, mockScraper, scrapeConfig)

	mockFS := mockUtil.MockOSFileSystem{
		Mfs: fstest.MapFS{
			"tempDir/pattern1-101": &fstest.MapFile{Data: []byte(`{"num": 1, "img": "pic1.jpg", "title": "title words", "alt": "alt words", "transcript": "transcript words"}`)},
		},
	}
	tempedIDs := len(mockFS.Mfs)
	tempCfg := config.TempConfig{
		TempDir:         "tempDir",
		TempFilePattern: "pattern",
	}
	temper, _ := util.NewTemper(&tempCfg, &mockFS)

	// Test successful scraping
	result, err := scrapeService.Scrape(map[int]bool{}, 403, temper)
	assert.NoError(t, err)
	assert.NotContains(t, result, 404)
	assert.Equal(t, scrapeConfig.ScrapePagesLimit+tempedIDs, len(result))
}

func TestAppendIDs(t *testing.T) {
	scrapeCtx, scrapeCancel := context.WithCancel(context.Background())

	IDsCh := make(chan int, 1)

	lastID := 403
	scrapeLimit := 3
	missedIDs := map[int]bool{1: true}
	go AppendIDs(scrapeCtx, scrapeCancel, IDsCh, scrapeLimit, missedIDs, lastID)

	var ids []int
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(scrapeLimit)

	go func() {
		for i := 0; i < 10000; i++ {
			id := <-IDsCh
			mu.Lock()
			ids = append(ids, id)
			mu.Unlock()
			assert.False(t, len(ids) > scrapeLimit)
			wg.Done()
		}
		//for {
		//	select {
		//	case id := <-IDsCh:
		//		mu.Lock()
		//		ids = append(ids, id)
		//		mu.Unlock()
		//		assert.False(t, len(ids) > scrapeLimit)
		//		wg.Done()
		//	}
		//}
	}()

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	select {
	case <-scrapeCtx.Done():
		//scrapeLimit
		wg.Wait()
		assert.True(t, len(ids) == scrapeLimit)
		//start after lastID, lastID do not include
		assert.NotContains(t, ids, 403)
		//start after lastID
		assert.Contains(t, ids, 405)
		//not add funny id
		assert.NotContains(t, ids, 404)
		//add missedIDs
		assert.Contains(t, ids, 1)
		return
	default:

		t.Error("context must be closed")
	}
}

func TestAppendIDs_ContextCancel(t *testing.T) {
	scrapeCtx, scrapeCancel := context.WithCancel(context.Background())

	IDsCh := make(chan int, 1)

	lastID := 0
	scrapeLimit := 5
	toCancel := 2
	missedIDs := map[int]bool{}
	go AppendIDs(scrapeCtx, scrapeCancel, IDsCh, scrapeLimit, missedIDs, lastID)

	var ids []int
	for {
		select {
		case <-IDsCh:
			toCancel--
			if toCancel <= 0 {
				scrapeCancel()
			}
		case <-scrapeCtx.Done():
			assert.True(t, toCancel <= len(ids) && len(ids) < scrapeLimit)
			return
		}
	}
}

func TestAppendTempedResponse(t *testing.T) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	mockFS := mockUtil.MockOSFileSystem{
		Mfs: fstest.MapFS{
			"tempDir/pattern1-101": &fstest.MapFile{Data: []byte(`{"num": 1, "img": "pic1.jpg", "title": "title words", "alt": "alt words", "transcript": "transcript words"}`)},
		},
	}
	tempCfg := config.TempConfig{
		TempDir:         "tempDir",
		TempFilePattern: "pattern",
	}
	temper, _ := util.NewTemper(&tempCfg, &mockFS)

	var result []domain.Comics
	AppendTempedResponse(&wg, &mu, temper, &result)
	assert.Equal(t, 1, len(result))

	expectedComics := domain.Comics{
		ID:         1,
		Picture:    "pic1.jpg",
		Title:      "title words",
		Alt:        "alt words",
		Transcript: "transcript words",
	}
	assert.Equal(t, expectedComics, result[0])
}

func TestScrapeWorker(t *testing.T) {
	tests := []struct {
		name         string
		id           []int
		expectComics *[]domain.Comics
		expectCancel bool
	}{
		{
			name: "Successful scrape",
			id:   []int{1, 2},
			expectComics: &[]domain.Comics{{
				ID:         1,
				Picture:    "pic1.jpg",
				Title:      "title words",
				Alt:        "alt words",
				Transcript: "transcript words",
			}, {
				ID:         2,
				Picture:    "pic2.jpg",
				Title:      "title words",
				Alt:        "alt words",
				Transcript: "transcript words",
			}},
			expectCancel: false,
		},
		{
			name: "Not found scrape stop",
			id:   []int{499, 500},
			expectComics: &[]domain.Comics{{
				ID:         499,
				Picture:    "pic499.jpg",
				Title:      "title words",
				Alt:        "alt words",
				Transcript: "transcript words",
			}},
			expectCancel: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			wg := &sync.WaitGroup{}

			mu := &sync.Mutex{}
			IDsChan := make(chan int, 1)
			result := &[]domain.Comics{}

			mockScraper := new(mock.MockScraper)
			mockOSFileSystem := mockUtil.MockOSFileSystem{
				Mfs: fstest.MapFS{},
			}
			mockTemper, _ := util.NewTemper(&config.TempConfig{}, &mockOSFileSystem)

			go func() {
				for _, id := range tt.id {
					IDsChan <- id
				}
			}()

			wg.Add(1)
			go ScrapeWorker(ctx, cancel, wg, mu, IDsChan, result, 0, mockScraper, mockTemper)

			if slices.Contains(tt.id, 500) {
				wg.Wait()
			} else {
				time.Sleep(200 * time.Millisecond)
			}

			mu.Lock()
			sort.Slice(*result, func(i, j int) bool {
				return (*result)[i].ID < (*result)[j].ID
			})

			if tt.expectComics != nil {
				sort.Slice(*tt.expectComics, func(i, j int) bool {
					return (*tt.expectComics)[i].ID < (*tt.expectComics)[j].ID
				})
			}
			mu.Unlock()

			assert.Equal(t, *tt.expectComics, *result)

			select {
			case <-ctx.Done():
				if tt.expectCancel {
					assert.NotNil(t, ctx.Err())
				} else {
					t.Error("Context was cancelled unexpectedly")
				}
			default:
				if tt.expectCancel {
					t.Error("Context was not cancelled as expected")
				} else {
					assert.Nil(t, ctx.Err())
				}
			}
		})
	}
}
