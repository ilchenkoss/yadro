package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"myapp/internal-web/core/domain"
	"net/http"
	"net/url"
)

type XkcdAPI struct {
	Client  http.Client
	BaseURL string
}

func NewXkcdAPI(baseURL string) *XkcdAPI {
	return &XkcdAPI{
		BaseURL: baseURL,
		Client:  http.Client{},
	}
}

func (a *XkcdAPI) UpdateDescription(id string, authToken string) error {
	type UpdateDescriptionResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	req, rErr := http.NewRequest("GET", fmt.Sprintf("%s/desc?description=%s", a.BaseURL, id), nil)
	if rErr != nil {
		return rErr
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", authToken))

	resp, respErr := a.Client.Do(req)
	if respErr != nil {
		return respErr
	}

	resBody, rErr := io.ReadAll(resp.Body)
	if rErr != nil {
		return rErr
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(resBody))
	}

	var searchResponse UpdateDescriptionResponse
	uErr := json.Unmarshal(resBody, &searchResponse)
	if uErr != nil {
		return uErr
	}

	return nil
}

func (a *XkcdAPI) GetComics(requestWords string, authToken string) ([]domain.Comics, error) {

	type SearchResponse struct {
		Success       bool            `json:"success"`
		Message       string          `json:"message"`
		FoundPictures []domain.Comics `json:"found_pictures"`
	}

	encodedRequestWords := url.QueryEscape(requestWords)

	req, rErr := http.NewRequest("GET", fmt.Sprintf("%s/pics?search=%s", a.BaseURL, encodedRequestWords), nil)
	if rErr != nil {
		return nil, rErr
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", authToken))

	resp, respErr := a.Client.Do(req)
	if respErr != nil {
		return nil, respErr
	}

	resBody, rErr := io.ReadAll(resp.Body)
	if rErr != nil {
		return nil, rErr
	}

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			return nil, domain.ErrToManyRequests
		case http.StatusUnauthorized:
			return nil, domain.ErrUnauthorized
		case http.StatusInternalServerError:
			return nil, domain.ErrAuthFailed
		}
		return nil, errors.New(string(resBody))
	}

	var searchResponse SearchResponse
	uErr := json.Unmarshal(resBody, &searchResponse)
	if uErr != nil {
		return nil, uErr
	}

	return searchResponse.FoundPictures, nil
}

func (a *XkcdAPI) Login(login string, password string) (string, error) {

	type requestLogin struct {
		Login    string
		Password string
	}

	type responseLogin struct {
		Success bool
		Message string
		Token   string
	}

	reqLogin := requestLogin{
		Login:    login,
		Password: password,
	}

	reqLoginJson, jErr := json.Marshal(reqLogin)
	if jErr != nil {
		return "", jErr
	}

	resLogin, rErr := http.Post(fmt.Sprintf("%s/login", a.BaseURL), "application/json", bytes.NewReader(reqLoginJson))
	if rErr != nil {
		return "", rErr
	}

	resBody, rErr := io.ReadAll(resLogin.Body)
	if rErr != nil {
		return "", rErr
	}

	if resLogin.StatusCode != http.StatusOK {
		switch resLogin.StatusCode {
		case http.StatusUnauthorized:
			return "", domain.ErrUnauthorized
		default:
			return "", errors.New(string(resBody))
		}
	}

	var response responseLogin
	uErr := json.Unmarshal(resBody, &response)
	if uErr != nil {
		return "", uErr
	}

	return response.Token, nil
}
