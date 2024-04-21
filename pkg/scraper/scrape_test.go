package scraper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(wantStatusCode int, ctxCancel context.CancelFunc, ID int) []byte {
	// create fake http server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(wantStatusCode)    //status code
		w.Write([]byte("Hello, world!")) //body

	}))

	response := sendRequest(server.Client(), server.URL, 3, ID, ctxCancel)
	return response
}

func TestRequestOK(t *testing.T) {

	ctx, ctxCancel := context.WithCancel(context.Background())
	response := testRequest(http.StatusOK, ctxCancel, 100)

	if ctx.Err() != nil {
		t.Errorf("Context closed, but comics dont end")
	}

	if response == nil {
		t.Errorf("Response do not return, but it shouldn't have")
	}
}

func TestRequestNotOK(t *testing.T) {

	ctx, ctxCancel := context.WithCancel(context.Background())
	response := testRequest(http.StatusInternalServerError, ctxCancel, 100)

	if ctx.Err() != nil {
		t.Errorf("Context closed, but comics dont end")
	}

	if response != nil {
		t.Errorf("Response return, but it shouldn't have")
	}
}

func TestResponseCLoseContext(t *testing.T) {

	ctx, ctxCancel := context.WithCancel(context.Background())
	testRequest(http.StatusNotFound, ctxCancel, 999)

	if ctx.Err() == nil {
		t.Errorf("Context not closed, but comics end")
	}
}

func TestResponseFunnyComics(t *testing.T) {

	ctx, ctxCancel := context.WithCancel(context.Background())
	testRequest(http.StatusNotFound, ctxCancel, 404)

	if ctx.Err() != nil {
		t.Errorf("Context closed, but comics not end")
	}
}
