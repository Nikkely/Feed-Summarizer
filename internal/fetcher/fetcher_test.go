package fetcher

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchHTML(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("<html><body>Test Page</body></html>")); err != nil {
			log.Fatalln("Failed to write response")
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	result, err := FetchHTML(ts.URL)
	assert.NoError(t, err, "fetchHTML returned an unexpected error")
	assert.Contains(t, result, "Test Page", "fetchHTML result mismatch")
}

func TestFetchHTML_Error(t *testing.T) {
	_, err := FetchHTML("http://invalid-url")
	assert.Error(t, err, "Expected an error but got nil")
}
