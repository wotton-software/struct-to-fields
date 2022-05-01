package stf_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/wotton-software/struct-to-fields/pkg/stf"
	"testing"
)

type Order struct {
	ID       string   `json:"id"`
	Customer Customer `json:"customer"`
}

type Customer struct {
	Name  string
	Age   int
	Email string `stf:"-"`
	*Address
	Child *Customer
}

type Address struct {
	Line1    string
	Line2    *string
	Postcode string `json:"eirecode",stf:"eirecode"`
}

func TestExtractFields(t *testing.T) {
	child := &Customer{
		Name:    "Steve",
		Age:     9,
		Email:   "steve@job.com",
		Address: nil,
		Child:   nil,
	}

	bob := Customer{
		Name:  "Bob",
		Age:   25,
		Email: "bob@job.com",
		Address: &Address{
			Line1:    "SomeStreetLine",
			Line2:    nil,
			Postcode: "EIRE123",
		},
		Child: child,
	}

	child.Child = &bob //demonstrates cyclical order handling

	order := Order{
		ID:       "order-id-1234",
		Customer: bob,
	}

	e := stf.NewExtractor(stf.ExcludeNilsOption(true))

	actual, err := e.ExtractFields(order)
	if err != nil {
		panic(err)
	}

	expected := map[string]interface{}{
		"id":                        order.ID,
		"customer.Name":             order.Customer.Name,
		"customer.Age":              order.Customer.Age,
		"customer.Address.Line1":    order.Customer.Address.Line1,
		"customer.Address.eirecode": order.Customer.Address.Postcode,
		"customer.Child.Age":        order.Customer.Child.Age,
		"customer.Child.Child.Age":  order.Customer.Child.Child.Age,
		"customer.Child.Child.Name": order.Customer.Child.Child.Name,
		"customer.Child.Name":       order.Customer.Child.Name,
	}

	assert.Equal(t, expected, actual)
}
