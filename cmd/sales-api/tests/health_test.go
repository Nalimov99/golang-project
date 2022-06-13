package tests

import (
	"encoding/json"
	"garagesale/cmd/sales-api/internal/handlers"
	"garagesale/internal/platform/database/databasetest"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestStatusCheck(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	t.Cleanup(teardown)

	log := log.New(os.Stdout, "TEST", log.Flags())

	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	resp := httptest.NewRecorder()

	app := handlers.API(log, db, nil)
	app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status code: %d, got: %d", http.StatusOK, resp.Code)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal("marshalling response body")
	}

	if body["status"] != "OK" {
		t.Fatalf("expect json body status: \"OK\", got: \"%v\"", body["status"])
	}
}
