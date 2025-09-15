package main

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Invoice struct {
	Id         uuid.UUID         `json:"id"`
	CustomerId uuid.UUID         `json:"customer_id"`
	Date       time.Time         `json:"invoice_date"`
	Positions  []InvoicePosition `json:"positions,omitempty"`
}

type InvoicePosition struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Desc   string    `json:"description"`
	Amount int       `json:"amount"`
	Price  float32   `json:"price"`
}
