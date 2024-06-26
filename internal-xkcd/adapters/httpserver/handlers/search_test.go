package handlers

//
//func TestNewSearchHandler(t *testing.T) {
//
//	tests := []struct {
//		name  string
//		mocks func(
//			wr *mock.MockWeightRepository,
//			ws *mock.MockWeightService,
//		)
//
//		requestBody  interface{}
//		expectedCode int
//	}{
//		{
//			name: "Success",
//			mocks: func(
//				wr *mock.MockWeightRepository,
//				ws *mock.MockWeightService,
//			) {
//				ws.EXPECT().WeightRequest(gomock.Any()).Return(map[string]float64{"word": 1.1})
//				wr.EXPECT().GetWeightsByWords(gomock.Any()).Return(&[]domain.Weights{}, nil)
//				ws.EXPECT().FindRelevantPictures(gomock.Any(), gomock.Any()).Return([]domain.Comics{{Picture: "picture1.jpg"}, {Picture: "picture2.jpg"}}, nil)
//			},
//			expectedCode: http.StatusOK,
//		}, {
//			name: "Error database",
//			mocks: func(
//				wr *mock.MockWeightRepository,
//				ws *mock.MockWeightService,
//			) {
//				ws.EXPECT().WeightRequest(gomock.Any()).Return(map[string]float64{"word": 1.1})
//				wr.EXPECT().GetWeightsByWords(gomock.Any()).Return(nil, errors.New("db err"))
//			},
//			expectedCode: http.StatusInternalServerError,
//		}, {
//			name: "Error weights service",
//			mocks: func(
//				wr *mock.MockWeightRepository,
//				ws *mock.MockWeightService,
//			) {
//				ws.EXPECT().WeightRequest(gomock.Any()).Return(map[string]float64{"word": 1.1})
//				wr.EXPECT().GetWeightsByWords(gomock.Any()).Return(&[]domain.Weights{}, nil)
//				ws.EXPECT().FindRelevantPictures(gomock.Any(), gomock.Any()).Return(nil, errors.New("weights service err"))
//			},
//			expectedCode: http.StatusInternalServerError,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			mockWeightRepo := mock.NewMockWeightRepository(ctrl)
//			mockWeightService := mock.NewMockWeightService(ctrl)
//			mockComicsRepo := mock.NewMockComicsRepository(ctrl)
//			limiter := utils.NewLimiter(&config.HttpServerConfig{RateLimit: 1, ConcurrencyLimit: 1})
//			tt.mocks(mockWeightRepo, mockWeightService)
//
//			searchHandler := NewSearchHandler(mockWeightRepo, mockWeightService, mockComicsRepo, *limiter)
//
//			body, _ := json.Marshal(tt.requestBody)
//			req := httptest.NewRequest("GET", "/pics?search=binary,christmas,tree", bytes.NewBuffer(body))
//
//			rr := httptest.NewRecorder()
//
//			handler := http.HandlerFunc(searchHandler.Search)
//			handler.ServeHTTP(rr, req)
//
//			fmt.Println(rr.Code, rr.Body)
//
//			assert.Equal(t, tt.expectedCode, rr.Code)
//		})
//	}
//}
