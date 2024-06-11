package utils

import "myapp/internal-xkcd/core/domain"

type UpdateResponse struct {
	Success     bool   `json:"success" example:"true"`
	Message     string `json:"message" example:"Success"`
	NewComics   int    `json:"new_comics" example:"Success"`
	TotalComics int    `json:"total_comics" example:"Success"`
}

func NewUpdateResponse(success bool, message string, newComics int, totalComics int) UpdateResponse {
	return UpdateResponse{
		Success:     success,
		Message:     message,
		NewComics:   newComics,
		TotalComics: totalComics,
	}
}

type SearchResponse struct {
	Success       bool            `json:"success" example:"true"`
	Message       string          `json:"message" example:"Success"`
	FoundPictures []domain.Comics `json:"found_pictures"`
}

func NewSearchResponse(success bool, message string, fundPictures []domain.Comics) SearchResponse {
	return SearchResponse{
		Success:       success,
		Message:       message,
		FoundPictures: fundPictures,
	}
}

type LoginResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Token   string `json:"token" `
}

func NewLoginResponse(success bool, message string, token string) LoginResponse {
	return LoginResponse{
		Success: success,
		Message: message,
		Token:   token,
	}
}

type DescriptionResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
}

func NewDescriptionResponse(success bool, message string) LoginResponse {
	return LoginResponse{
		Success: success,
		Message: message,
	}
}
