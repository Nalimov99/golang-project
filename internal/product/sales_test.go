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

func TestAddSale(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	t.Cleanup(teardown)
	ctx := context.Background()

	const quantity, paid = 1, 1

	ns := product.NewSale{
		Quantity: quantity,
		Paid:     paid,
	}
	now := time.Now()
	np := product.NewProduct{
		Name:     "test",
		Quantity: 1,
		Cost:     10,
	}

	p, err := product.Create(ctx, db, np, now)
	if err != nil {
		t.Fatalf("could not create product %v", err)
	}

	createdSale, err := product.AddSale(ctx, db, ns, strconv.Itoa(p.ID), now)
	if err != nil {
		t.Fatalf("could not create sale: %v", err)
	}

	if now.UTC().Format(time.ANSIC) != createdSale.DateCreated.Format(time.ANSIC) {
		t.Fatalf("dateCreated is not equal: \n%v, \n%v", now, createdSale.DateCreated)
	}

	resultSale := product.Sale{
		ID:          createdSale.ID,
		ProductID:   p.ID,
		Quantity:    quantity,
		Paid:        paid,
		DateCreated: createdSale.DateCreated,
	}

	if diff := cmp.Diff(resultSale, *createdSale); diff != "" {
		t.Fatalf("saved sale did not match created: \n%s", diff)
	}
}

func TestSalesList(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	t.Cleanup(teardown)
	ctx := context.Background()

	sales := [2]product.NewSale{
		{
			Quantity: 1,
			Paid:     1,
		},
		{
			Quantity: 2,
			Paid:     2,
		},
	}

	now := time.Now()
	np := product.NewProduct{
		Name:     "test",
		Quantity: 1,
		Cost:     10,
	}

	p, err := product.Create(ctx, db, np, now)
	if err != nil {
		t.Fatalf("could not create product %v", err)
	}

	for _, s := range sales {
		if _, err := product.AddSale(ctx, db, s, strconv.Itoa(p.ID), now); err != nil {
			t.Fatalf("could not create sale: %v", err)
		}
	}

	got, err := product.ListSales(ctx, db, strconv.Itoa(p.ID))
	if err != nil {
		t.Fatalf("could not get product list %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("length of created and expected does not equal: \n%v", got)
	}
}
