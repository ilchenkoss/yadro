package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestOK(t *testing.T) {
	// create fake http server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)     //status code
		w.Write([]byte("Hello, world!")) //body

	}))

	defer server.Close()

	response := sendRequest(server.Client(), server.URL, 3, 100)

	if string(response) != "Hello, world!" {
		t.Errorf("unexpected response body: got %s, want %s", response, "Hello, world!")
	}
}

func TestRequestNotOK(t *testing.T) {
	// create fake http server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusInternalServerError) //status code
		w.Write([]byte("Hello, world!"))              //body

	}))

	defer server.Close()

	response := sendRequest(server.Client(), server.URL, 3, 100)

	if response != nil {
		t.Errorf("unexpected response body: got %s, want %s", response, "nil")
	}
}

func TestRequestChangeGlobalVar(t *testing.T) {
	// create fake http server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusNotFound) //status code
		w.Write([]byte("Hello, world!"))   //body

	}))

	defer server.Close()

	Condition = true

	response := sendRequest(server.Client(), server.URL, 3, 405)

	if Condition != false {
		t.Errorf("unexpected response body: got %s, want %s", response, "nil")
	}
}

func TestRequestFunny404(t *testing.T) {
	// create fake http server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusNotFound) //status code
		w.Write([]byte("Hello, world!"))   //body

	}))

	defer server.Close()

	Condition = true

	response := sendRequest(server.Client(), server.URL, 3, 404)

	if Condition != true {
		t.Errorf("unexpected response body: got %s, want %s", response, "nil")
	}
}
