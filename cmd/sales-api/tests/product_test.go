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
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestProducts(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	t.Cleanup(teardown)

	log := log.New(os.Stdout, "TEST", log.Flags())

	tests := ProductTest{
		app: handlers.API(log, db, nil),
	}

	t.Log("RUN PRODUCT TESTS")
	t.Run("CREATE", tests.Create)
	t.Run("LIST", tests.List)
	t.Run("RETRIEVE", tests.Retrieve)
	t.Run("UPDATE", tests.Update)
	t.Run("DELETE", tests.DeleteProduct)

	t.Log("RUN SALES TESTS")
	t.Run("ADD SALE", tests.AddSale)
	t.Run("SALE LIST", tests.ListSales)
}

// ProductTests holds methods for each product subtest
// These type allows passing dependencies for tests
type ProductTest struct {
	app      http.Handler
	products []map[string]interface{}
	sales    []map[string]interface{}
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

	want := map[string]interface{}{
		"name":         "1",
		"quantity":     float64(1),
		"cost":         float64(1),
		"id":           product["id"],
		"sold":         product["sold"],
		"revenue":      product["revenue"],
		"date_created": product["date_created"],
		"date_updated": product["date_updated"],
	}

	if diff := cmp.Diff(want, product); diff != "" {
		t.Fatalf("expected product diff: \n%v", diff)
	}

	p.products = append(p.products, want)
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

func (p *ProductTest) Update(t *testing.T) {
	product := p.products[0]

	updateName := "update name"
	updateCost := 10
	updateQuantity := 20
	update := strings.NewReader(`
		{
			"name": "` + updateName + `",
			"cost": ` + strconv.Itoa(updateCost) + `,
			"quantity": ` + strconv.Itoa(updateQuantity) + `
		}
	`)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/products/%v", product["id"]), update)
	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	var got map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decoding: %v", err)
	}

	dateUpdatedStr := got["date_updated"].(map[string]interface{})["Time"].(string)
	dateUpdated, err := time.Parse(time.RFC3339, dateUpdatedStr)
	if err != nil {
		t.Fatalf("could not parse date: %v", err)
	}
	dateCreatedStr := got["date_updated"].(map[string]interface{})["Time"].(string)
	dateCreated, err := time.Parse(time.RFC3339, dateCreatedStr)
	if err != nil {
		t.Fatalf("could not parse date: %v", err)
	}

	if dateCreated.Before(dateUpdated) {
		t.Fatalf("time is not valid: created - %v, updated - %v", got["date_created"], got["date_updated"])
	}

	want := map[string]interface{}{
		"name":         updateName,
		"quantity":     float64(updateQuantity),
		"cost":         float64(updateCost),
		"id":           got["id"],
		"sold":         got["sold"],
		"revenue":      got["revenue"],
		"date_created": got["date_created"],
		"date_updated": got["date_updated"],
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("expected products diff: \n%v", diff)
	}

	p.products[0] = want
}

func (p *ProductTest) DeleteProduct(t *testing.T) {
	// create
	np := strings.NewReader(`{"name": "1", "quantity": 1,  "cost": 1}`)

	cResp := httptest.NewRecorder()
	cReq := httptest.NewRequest(http.MethodPost, "/v1/products", np)
	p.app.ServeHTTP(cResp, cReq)

	var product map[string]interface{}
	if err := json.NewDecoder(cResp.Body).Decode(&product); err != nil {
		t.Fatalf("decoding: %v", err)
	}

	// list
	lReq := httptest.NewRequest(http.MethodGet, "/v1/products", nil)
	lResp := httptest.NewRecorder()

	p.app.ServeHTTP(lResp, lReq)

	var list []map[string]interface{}
	if err := json.NewDecoder(lResp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding: %v", err)
	}

	// delete
	dReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/products/%v", product["id"]), nil)
	dResp := httptest.NewRecorder()

	p.app.ServeHTTP(dResp, dReq)

	if dResp.Code != http.StatusNoContent {
		t.Fatalf("expected status code: %d, got: %d", http.StatusNoContent, dResp.Code)
	}

	// new list
	p.app.ServeHTTP(lResp, lReq)
	var newList []map[string]interface{}
	if err := json.NewDecoder(lResp.Body).Decode(&newList); err != nil {
		t.Fatalf("decoding: %v", err)
	}

	if len(list)-1 != len(newList) {
		t.Fatalf("expected list length: %v, gotted list length: %v", len(list)-1, len(newList))
	}
}

func (p *ProductTest) AddSale(t *testing.T) {
	product := p.products[0]
	saleJson := strings.NewReader(`
		{
			"quantity": 1,
			"paid": 20
		}
	`)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/products/%v/sales", product["id"]), saleJson)
	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	var newSale map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&newSale); err != nil {
		t.Fatalf("decoding: %v", err)
	}

	want := map[string]interface{}{
		"quantity":     float64(1),
		"paid":         float64(20),
		"date_created": newSale["date_created"],
		"id":           newSale["id"],
		"product_id":   product["id"],
	}

	if diff := cmp.Diff(want, newSale); diff != "" {
		t.Fatalf("expected sales diff: \n%v", diff)
	}

	p.sales = append(p.sales, want)
}

func (p *ProductTest) ListSales(t *testing.T) {
	product := p.products[0]

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/products/%v/sales", product["id"]), nil)
	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	var fetched []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
		t.Fatalf("decoding: %v", err)
	}

	if len(fetched) != len(p.sales) {
		t.Fatalf("expected lenght is not equal. Want: %v, Get: %v", len(p.sales), len(fetched))
	}

	if diff := cmp.Diff(fetched, p.sales); diff != "" {
		t.Fatalf("expected sales diff: \n%v", diff)
	}
}
