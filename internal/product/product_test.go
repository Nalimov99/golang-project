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

func TestProductUpdate(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	t.Cleanup(teardown)
	ctx := context.Background()
	createdTime := time.Now()

	newProduct := product.NewProduct{
		Name:     "created name",
		Quantity: 10,
		Cost:     10,
	}
	createdProduct, err := product.Create(ctx, db, newProduct, createdTime)
	if err != nil {
		t.Fatal("could not create a product")
	}

	updatedTime := time.Now()
	updateName := "updated name"
	updateCost := 5
	updateQuantity := 5
	update := product.UpdateProduct{
		Name:     &updateName,
		Cost:     &updateCost,
		Quantity: &updateQuantity,
	}
	got, err := product.Update(ctx, db, strconv.Itoa(createdProduct.ID), update, updatedTime)
	if err != nil {
		t.Fatalf("could not update product: %v", err)
	}

	if createdTime.UTC().Format(time.ANSIC) != got.DateCreated.Time.Format(time.ANSIC) {
		t.Fatalf("dateCreated is not equal: \n%v, \n%v", createdTime, got.DateCreated.Time)
	}
	if updatedTime.UTC().Format(time.ANSIC) != got.DateUpdated.Time.Format(time.ANSIC) {
		t.Fatalf("dateUpdated is not equal: \n%v, \n%v", updatedTime, got.DateUpdated.Time)
	}

	want := product.Product{
		ID:          createdProduct.ID,
		Name:        *update.Name,
		Quantity:    *update.Quantity,
		Cost:        *update.Cost,
		Sold:        createdProduct.Sold,
		Revenue:     createdProduct.Revenue,
		DateCreated: got.DateCreated,
		DateUpdated: got.DateUpdated,
	}

	if diff := cmp.Diff(want, *got); diff != "" {
		t.Fatalf("expected product did not match: %v", diff)
	}
}

func TestProductDelete(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	t.Cleanup(teardown)
	ctx := context.Background()

	np := product.NewProduct{
		Name:     "2",
		Quantity: 1,
		Cost:     1,
	}

	created, err := product.Create(ctx, db, np, time.Now())
	if err != nil {
		t.Fatal("could not create")
	}

	id := strconv.Itoa(created.ID)

	if _, err := product.Retrieve(ctx, db, id); err != nil {
		t.Fatal("could not retrieve")
	}
	if err := product.Delete(ctx, db, id); err != nil {
		t.Fatalf("could not delete: %v", err)
	}

	_, err = product.Retrieve(ctx, db, id)
	if err != product.ErrNotFound {
		t.Fatal(err)
	}
}
