package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"go.pixelfactory.io/needle/internal/api/handlers"
	"go.pixelfactory.io/needle/testdata"
)

func Test_DefaultHandler(t *testing.T) {
	is := require.New(t)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := handlers.NewDefaultHandler()
	handler.ServeHTTP(rr, req)

	is.Equal(rr.Code, http.StatusOK)
	is.Equal(rr.Body.Bytes(), handlers.Pixel)
}

func Test_CAHandler(t *testing.T) {
	is := require.New(t)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := handlers.NewCAHandler(testdata.Dir() + "/certs/root-ca.crt")
	handler.ServeHTTP(rr, req)

	ca, err := os.ReadFile(testdata.Dir() + "/certs/root-ca.crt")
	if err != nil {
		t.Fatal(err)
	}

	is.Equal(rr.Code, http.StatusOK)
	is.Equal(rr.Body.Bytes(), ca)
}
