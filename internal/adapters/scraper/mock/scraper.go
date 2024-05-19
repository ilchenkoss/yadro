package mock

import (
	"fmt"
	"net/http"
	"strings"
)

type MockScraper struct{}

func (m *MockScraper) GetResponse(url string, retries int) ([]byte, int, error) {
	switch url {
	case "https://xkcd.com/404/info.0.json":
		return []byte{}, http.StatusNotFound, nil
	case "https://xkcd.com/500/info.0.json":
		return []byte{}, http.StatusNotFound, nil
	default:
		idStr := strings.TrimPrefix(url, "https://xkcd.com/")
		idStr = strings.TrimSuffix(idStr, "/info.0.json")
		return []byte(fmt.Sprintf(`{"num": %s, "img": "pic%s.jpg", "title": "title words", "alt": "alt words", "transcript": "transcript words"}`, idStr, idStr)), http.StatusOK, nil
	}
}
