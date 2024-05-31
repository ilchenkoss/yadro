package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"myapp/internal/config"
	"myapp/internal/core/domain"
	"myapp/internal/core/port/mock"
	mock2 "myapp/internal/core/util/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
)

func TestNewScrapeHandler(t *testing.T) {

	tests := []struct {
		name  string
		mocks func(
			ssvc *mock.MockScrapeService,
			wsvc *mock.MockWeightService,
			cr *mock.MockComicsRepository,
			wr *mock.MockWeightRepository,
		)

		requestBody  interface{}
		expectedCode int
	}{
		{
			name: "Success",
			mocks: func(
				ssvc *mock.MockScrapeService,
				wsvc *mock.MockWeightService,
				cr *mock.MockComicsRepository,
				wr *mock.MockWeightRepository,
			) {
				cr.EXPECT().
					GetMissedIDs().Return(map[int]bool{2: true}, nil)
				cr.EXPECT().
					GetMaxID().Return(4, nil)
				ssvc.EXPECT().
					Scrape(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Comics{}, nil)
				cr.EXPECT().
					InsertComics(gomock.Any()).Return(10, nil)
				wsvc.EXPECT().
					WeightComics(gomock.Any()).Return(&[]domain.Weights{})
				wr.EXPECT().
					InsertWeights(gomock.Any()).Return(nil)
				cr.EXPECT().
					GetCountComics().Return(100, nil)
			},
			expectedCode: http.StatusOK,
		}, {
			name: "Error scrape",
			mocks: func(
				ssvc *mock.MockScrapeService,
				wsvc *mock.MockWeightService,
				cr *mock.MockComicsRepository,
				wr *mock.MockWeightRepository,
			) {
				cr.EXPECT().
					GetMissedIDs().Return(map[int]bool{2: true}, nil)
				cr.EXPECT().
					GetMaxID().Return(4, nil)
				ssvc.EXPECT().
					Scrape(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("scrape err"))
			},
			expectedCode: http.StatusInternalServerError,
		}, {
			name: "Error insert comics",
			mocks: func(
				ssvc *mock.MockScrapeService,
				wsvc *mock.MockWeightService,
				cr *mock.MockComicsRepository,
				wr *mock.MockWeightRepository,
			) {
				cr.EXPECT().
					GetMissedIDs().Return(map[int]bool{2: true}, nil)
				cr.EXPECT().
					GetMaxID().Return(4, nil)
				ssvc.EXPECT().
					Scrape(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Comics{}, nil)
				cr.EXPECT().
					InsertComics(gomock.Any()).Return(0, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		}, {
			name: "Error insert weights",
			mocks: func(
				ssvc *mock.MockScrapeService,
				wsvc *mock.MockWeightService,
				cr *mock.MockComicsRepository,
				wr *mock.MockWeightRepository,
			) {
				cr.EXPECT().
					GetMissedIDs().Return(map[int]bool{2: true}, nil)
				cr.EXPECT().
					GetMaxID().Return(4, nil)
				ssvc.EXPECT().
					Scrape(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Comics{}, nil)
				cr.EXPECT().
					InsertComics(gomock.Any()).Return(10, nil)
				wsvc.EXPECT().
					WeightComics(gomock.Any()).Return(&[]domain.Weights{})
				wr.EXPECT().
					InsertWeights(gomock.Any()).Return(errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		}, {
			name: "Error get count comics",
			mocks: func(
				ssvc *mock.MockScrapeService,
				wsvc *mock.MockWeightService,
				cr *mock.MockComicsRepository,
				wr *mock.MockWeightRepository,
			) {
				cr.EXPECT().
					GetMissedIDs().Return(map[int]bool{2: true}, nil)
				cr.EXPECT().
					GetMaxID().Return(4, nil)
				ssvc.EXPECT().
					Scrape(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Comics{}, nil)
				cr.EXPECT().
					InsertComics(gomock.Any()).Return(10, nil)
				wsvc.EXPECT().
					WeightComics(gomock.Any()).Return(&[]domain.Weights{})
				wr.EXPECT().
					InsertWeights(gomock.Any()).Return(nil)
				cr.EXPECT().
					GetCountComics().Return(0, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		}, {
			name: "Concurrent requests",
			mocks: func(
				ssvc *mock.MockScrapeService,
				wsvc *mock.MockWeightService,
				cr *mock.MockComicsRepository,
				wr *mock.MockWeightRepository,
			) {
				cr.EXPECT().
					GetMissedIDs().Return(map[int]bool{2: true}, nil)
				cr.EXPECT().
					GetMaxID().Return(4, nil)
				ssvc.EXPECT().
					Scrape(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Comics{}, nil)
				cr.EXPECT().
					InsertComics(gomock.Any()).Return(10, nil)
				wsvc.EXPECT().
					WeightComics(gomock.Any()).Return(&[]domain.Weights{})
				wr.EXPECT().
					InsertWeights(gomock.Any()).Return(nil)
				cr.EXPECT().
					GetCountComics().Return(100, nil)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockScrapeService := mock.NewMockScrapeService(ctrl)
			mockWeightService := mock.NewMockWeightService(ctrl)
			mockComicsRepository := mock.NewMockComicsRepository(ctrl)
			mockWeightRepository := mock.NewMockWeightRepository(ctrl)

			tt.mocks(mockScrapeService, mockWeightService, mockComicsRepository, mockWeightRepository)

			scrapeHandler := NewScrapeHandler(mockScrapeService, mockWeightService, mockComicsRepository, mockWeightRepository, context.Background(), &config.Config{}, &mock2.MockOSFileSystem{Mfs: fstest.MapFS{}})

			body, _ := json.Marshal(tt.requestBody)

			if !strings.Contains(tt.name, "Concurrent requests") {
				req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(scrapeHandler.Update)
				handler.ServeHTTP(rr, req)
				assert.Equal(t, tt.expectedCode, rr.Code)
				fmt.Println(rr.Code, rr.Body)
			} else {
				var wg sync.WaitGroup
				wg.Add(2)

				req1 := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
				rr1 := httptest.NewRecorder()
				req2 := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
				rr2 := httptest.NewRecorder()

				go func() {
					defer wg.Done()
					handler := http.HandlerFunc(scrapeHandler.Update)
					handler.ServeHTTP(rr1, req1)
				}()
				go func() {
					defer wg.Done()
					handler := http.HandlerFunc(scrapeHandler.Update)
					handler.ServeHTTP(rr2, req2)
				}()

				wg.Wait()

				assert.True(t, rr1.Code == http.StatusOK || rr2.Code == http.StatusOK)
				assert.True(t, rr1.Code == http.StatusInternalServerError || rr2.Code == http.StatusInternalServerError)

			}

		})
	}
}
