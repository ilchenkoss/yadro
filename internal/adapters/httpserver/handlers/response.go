package handlers

type updateResponse struct {
	Success     bool   `json:"success" example:"true"`
	Message     string `json:"message" example:"Success"`
	NewComics   int    `json:"new_comics" example:"Success"`
	TotalComics int    `json:"total_comics" example:"Success"`
}

func newUpdateResponse(success bool, message string, newComics int, totalComics int) updateResponse {
	return updateResponse{
		Success:     success,
		Message:     message,
		NewComics:   newComics,
		TotalComics: totalComics,
	}
}

type searchResponse struct {
	Success       bool     `json:"success" example:"true"`
	Message       string   `json:"message" example:"Success"`
	FoundPictures []string `json:"found_pictures"`
}

func newSearchResponse(success bool, message string, fundPictures []string) searchResponse {
	return searchResponse{
		Success:       success,
		Message:       message,
		FoundPictures: fundPictures,
	}
}
