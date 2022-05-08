package product_test

import (
	"context"
	"garagesale/internal/platform/database/databasetest"
	"garagesale/internal/product"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestProducts(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	t.Cleanup(teardown)
	ctx := context.Background()

	np := product.NewProduct{
		Name:     "test",
		Quantity: 1,
		Cost:     10,
	}

	now := time.Now()

	p, err := product.Create(ctx, db, np, now)
	if err != nil {
		t.Fatalf("could not create product %v", err)
	}

	saved, err := product.Retrieve(ctx, db, strconv.Itoa(p.ID))
	if err != nil {
		t.Fatalf("could not retrieve product, id: %v", p.ID)
	}

	if diff := cmp.Diff(p, saved); diff != "" {
		t.Fatalf("saved product did not match created: \n%s", diff)
	}
}

func TestProductList(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	t.Cleanup(teardown)
	ctx := context.Background()

	newProducts := [2]product.NewProduct{
		{
			Name:     "first product",
			Quantity: 1,
			Cost:     10,
		},
		{
			Name:     "second product",
			Quantity: 2,
			Cost:     12,
		},
	}
	now := time.Now()

	var savedProducts []product.Product

	for _, value := range newProducts {
		saved, err := product.Create(ctx, db, value, now)
		if err != nil {
			t.Fatalf("could not create product %v", err)
		}

		savedProducts = append(savedProducts, *saved)
	}

	if len(savedProducts) != len(newProducts) {
		t.Fatalf("length of created and saved products are not equal")
	}

	for _, value := range savedProducts {
		retrieve, err := product.Retrieve(ctx, db, strconv.Itoa(value.ID))
		if err != nil {
			t.Fatalf("could not retrieve product, ID: %d", value.ID)
		}

		if diff := cmp.Diff(*retrieve, value); diff != "" {
			t.Fatalf("saved product did not match created: \n%s", diff)
		}
	}
}
