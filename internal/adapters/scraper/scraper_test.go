package scraper

import (
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewScraper(t *testing.T) {
	timeoutSeconds := int(1)
	scraper := NewScraper(1)
	assert.Equal(t, timeoutSeconds, int(scraper.Client.Timeout.Seconds()))
}

func TestScraper_GetResponse_Success(t *testing.T) {
	expectedBody := "Hello, World!"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedBody))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	scraper := NewScraper(1)

	body, statusCode, err := scraper.GetResponse(server.URL, 1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, expectedBody, string(body))
}

func TestScraper_GetResponse_Retry(t *testing.T) {
	expectedBody := "Hello, Retry!"
	var try int
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		try++
		if try < 3 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, Retry!"))
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	scraper := NewScraper(1)

	body, statusCode, err := scraper.GetResponse(server.URL, 3)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, expectedBody, string(body))
}

type CustomRoundTripper struct{}

func (c *CustomRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, &net.DNSError{Err: "no such host"}
}

func NewTestScraper(client *http.Client) *Scraper {
	return &Scraper{Client: client}
}

func TestScraper_GetResponse_NoSuchHost(t *testing.T) {

	retries := 1
	customClient := http.Client{Timeout: 1 * time.Second, Transport: &CustomRoundTripper{}}
	scraper := NewTestScraper(&customClient)
	start := time.Now()
	_, _, err := scraper.GetResponse("http://ya.ru", retries)
	duration := time.Since(start)

	*scraper.Client = customClient

	assert.ErrorContains(t, err, "no such host")
	assert.True(t, duration >= time.Duration(retries)*time.Second)
}
