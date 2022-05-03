package tests

import (
	"encoding/json"
	"fmt"
	"garagesale/cmd/sales-api/internal/handlers"
	"garagesale/internal/platform/database/databasetest"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProducts(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	t.Cleanup(teardown)

	log := log.New(os.Stdout, "TEST", log.Flags())

	tests := ProductTest{
		app: handlers.API(log, db),
	}

	t.Run("CREATE", tests.Create)
	t.Run("LIST", tests.List)
	t.Run("RETRIEVE", tests.Retrieve)
}

// ProductTests holds methods for each product subtest
// These type allows passing dependencies for tests
type ProductTest struct {
	app      http.Handler
	products []map[string]interface{}
}

func (p *ProductTest) Create(t *testing.T) {
	np := strings.NewReader(`{"name": "1", "quantity": 1,  "cost": 1}`)

	req := httptest.NewRequest(http.MethodPost, "/v1/products", np)
	resp := httptest.NewRecorder()
	p.app.ServeHTTP(resp, req)

	var product map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		t.Fatalf("decoding: %v", err)
	}

	p.products = append(p.products, product)
}

func (p *ProductTest) List(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/products", nil)
	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status code: %d, got: %d", http.StatusOK, resp.Code)
	}

	var list []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding: %v", err)
	}

	if diff := cmp.Diff(p.products, list); diff != "" {
		t.Fatalf("expected products diff: \n%v", diff)
	}
}

func (p *ProductTest) Retrieve(t *testing.T) {
	want := p.products[0]

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/products/%v", want["id"]), nil)
	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	var fetched map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
		t.Fatalf("decoding: %v", err)
	}

	if diff := cmp.Diff(want, fetched); diff != "" {
		t.Fatalf("expected products diff: \n%v", diff)
	}
}
