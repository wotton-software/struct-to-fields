package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/wotton-software/struct-to-fields/pkg/stf"
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

func main() {
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

	child.Child = &bob //demonstrates cyclical data handling

	data := Order{
		ID:       "order-id-1234",
		Customer: bob,
	}

	extractor := stf.NewExtractor(stf.ExcludeNilsOption(true))

	tags, err := extractor.ExtractFields(data)
	if err != nil {
		panic(err)
	}

	log.WithFields(tags).Info("Some logging just happened")
}
