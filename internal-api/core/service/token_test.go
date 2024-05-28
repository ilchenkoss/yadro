package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"myapp/internal-api/config"
	"myapp/internal-api/core/domain"
	"testing"
)

func TestTokenService(t *testing.T) {
	testCases := []struct {
		desc      string
		cfg       config.HttpServerConfig
		user      domain.User
		expectErr error
	}{
		{
			desc:      "Success",
			cfg:       config.HttpServerConfig{TokenMaxTime: 1},
			user:      domain.User{Login: "user"},
			expectErr: nil,
		}, {
			desc:      "Error token expired",
			cfg:       config.HttpServerConfig{TokenMaxTime: 0},
			user:      domain.User{Login: "user"},
			expectErr: domain.ErrTokenExpired,
		}, {
			desc:      "Error token not valid",
			cfg:       config.HttpServerConfig{TokenMaxTime: 1},
			user:      domain.User{Login: "user"},
			expectErr: domain.ErrTokenNotValid,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.desc, func(t *testing.T) {

			ts := NewTokenService(tc.cfg)

			tokenString, ctErr := ts.CreateToken(&tc.user)
			assert.NoError(t, ctErr)

			if errors.Is(tc.expectErr, domain.ErrTokenNotValid) {
				tokenString = tokenString[:1]
			}

			userLogin, guErr := ts.GetUserByTokenString(tokenString)
			assert.Equal(t, tc.expectErr, guErr)

			if guErr == nil {
				assert.Equal(t, tc.user.Login, userLogin)
			}

		})
	}
}
