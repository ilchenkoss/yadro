package utils

//
//func TestAuthMiddleware(t *testing.T) {
//
//	limiter := NewLimiter(&config.HttpServerConfig{ConcurrencyLimit: 1, RateLimit: 1})
//
//	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//		_, err := w.Write([]byte("Success"))
//		if err != nil {
//			return
//		}
//	})
//
//	tests := []struct {
//		name       string
//		authHeader string
//		mocks      func(
//			userRepo *mock.MockUserRepository,
//			tokenService *mock.MockTokenService,
//		)
//		requestRole  domain.UserRole
//		expectedCode int
//	}{
//		{
//			name:       "missing token",
//			authHeader: "",
//			mocks: func(
//				userRepo *mock.MockUserRepository,
//				tokenService *mock.MockTokenService,
//			) {
//			},
//			requestRole:  domain.Ordinary,
//			expectedCode: http.StatusUnauthorized,
//		},
//		{
//			name:       "invalid auth header",
//			authHeader: "invalid",
//
//			mocks: func(
//				userRepo *mock.MockUserRepository,
//				tokenService *mock.MockTokenService,
//			) {
//			},
//			requestRole:  domain.Ordinary,
//			expectedCode: http.StatusUnauthorized,
//		},
//		{
//			name:       "valid token, user not found",
//			authHeader: "Bearer validToken",
//			mocks: func(
//				userRepo *mock.MockUserRepository,
//				tokenService *mock.MockTokenService,
//			) {
//				tokenService.EXPECT().GetUserByTokenString("validToken").Return("userLogin", nil)
//				userRepo.EXPECT().GetUserByLogin("userLogin").Return(nil, domain.ErrUserNotFound)
//			},
//			requestRole:  domain.Ordinary,
//			expectedCode: http.StatusUnauthorized,
//		},
//		{
//			name:       "valid token, user found, no permissions",
//			authHeader: "Bearer validToken",
//			mocks: func(
//				userRepo *mock.MockUserRepository,
//				tokenService *mock.MockTokenService,
//			) {
//				tokenService.EXPECT().GetUserByTokenString("validToken").Return("userLogin", nil)
//				user := &domain.User{Role: domain.Ordinary}
//				userRepo.EXPECT().GetUserByLogin("userLogin").Return(user, nil)
//			},
//			requestRole:  domain.Admin,
//			expectedCode: http.StatusForbidden,
//		},
//		{
//			name:       "valid token, user found, have permissions",
//			authHeader: "Bearer validToken",
//			mocks: func(
//				userRepo *mock.MockUserRepository,
//				tokenService *mock.MockTokenService,
//			) {
//				tokenService.EXPECT().GetUserByTokenString("validToken").Return("userLogin", nil)
//				user := &domain.User{ID: 1, Login: "userLogin", Role: domain.SuperAdmin}
//				userRepo.EXPECT().GetUserByLogin("userLogin").Return(user, nil)
//			},
//			requestRole:  domain.Admin,
//			expectedCode: http.StatusOK,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			mockTokenService := mock.NewMockTokenService(ctrl)
//			mockUserRepository := mock.NewMockUserRepository(ctrl)
//
//			tt.mocks(mockUserRepository, mockTokenService)
//
//			req, err := http.NewRequest("GET", "/test", nil)
//			if err != nil {
//				t.Fatal(err)
//			}
//			if tt.authHeader != "" {
//				req.Header.Set("Authorization", tt.authHeader)
//			}
//
//			responseRecorder := httptest.NewRecorder()
//
//			handler := nextHandler
//			switch tt.requestRole {
//			case domain.Ordinary:
//				handler = OrdinaryMiddleware(nextHandler, mockTokenService, mockUserRepository, limiter)
//			case domain.Admin:
//				handler = AdminMiddleware(nextHandler, mockTokenService, mockUserRepository, limiter)
//			case domain.SuperAdmin:
//				handler = SuperAdminMiddleware(nextHandler, mockTokenService, mockUserRepository, limiter)
//			}
//
//			handler.ServeHTTP(responseRecorder, req)
//
//			assert.Equal(t, tt.expectedCode, responseRecorder.Code)
//		})
//	}
//}
