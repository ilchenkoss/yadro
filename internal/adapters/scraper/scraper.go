package scraper

import (
	"io"
	"net/http"
	"strings"
	"time"
)

type Scraper struct {
	Client *http.Client
}

func NewScraper(timeoutSeconds time.Duration) *Scraper {
	client := http.Client{Timeout: timeoutSeconds * time.Second}
	return &Scraper{
		Client: &client,
	}
}

func (s *Scraper) GetResponse(url string, retries int) ([]byte, int, error) {

	var statusCode int
	var err error

	for range retries {

		resp, rErr := s.Client.Get(url)

		if rErr != nil {
			err = rErr
			if strings.Contains(err.Error(), "no such host") {
				time.Sleep(1 * time.Second)
			}
			continue
		}

		if resp == nil {
			continue
		}

		statusCode = resp.StatusCode

		if resp.StatusCode == http.StatusOK {

			body, errRead := io.ReadAll(resp.Body)
			if errRead != nil {
				continue
			}

			resp.Body.Close()
			return body, statusCode, nil
		}

	}

	return nil, statusCode, err
}
